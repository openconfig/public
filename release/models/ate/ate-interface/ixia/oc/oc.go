/*
Copyright 2020 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package oc translates Ixia statistics tables to an OpenConfig schema.
package oc

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	log "github.com/golang/glog"
	"github.com/openconfig/public/release/models/ate/ate-interface/ixia/stats"
	"github.com/openconfig/public/release/models/ate/ate-interface/telemetry"
	"github.com/openconfig/ygot/ygot"
	"github.com/pkg/errors"
)

var (
	translateFunctions = map[stats.View]func(stats.Table, []string) (*telemetry.Device, error){
		stats.PortStatsView:    translatePortStats,
		stats.PortCPUStatsView: translatePortCPUStats,
		stats.FlowStatsView:    translateFlowStats,
		stats.EgressStatsView:  translateEgressStats,
	}
	// debugLog enables logging of debug information in the library. It is expensive
	// so should only be enabled for development purposes.
	debugLog = false
)

// Translate converts Ixia statistics tables to an OpenConfig schema.
// The form of the input table is assumed to be:
//
//	{
//		StatisticsViewName1: [
//			{
//				Caption1: Value1,
//				Caption2: Value2
//			},
//			{
//				Caption1: Value1,
//				Caption2: Value2,
//			},
//		],
//		StatisticsViewName2: [
//			{
//				...
//			}
//		],
//	}
//
// where each item in the slice under a statistics view represents a row of data in that view.
//
// The itFlows parameter indicates which flows have ingress tracking enabled.
func Translate(stats map[stats.View]stats.Table, itFlows []string) (ygot.ValidatedGoStruct, error) {
	root := &telemetry.Device{}
	for k, v := range stats {
		if fn, ok := translateFunctions[k]; ok {
			d, err := fn(v, itFlows)
			if err != nil {
				log.Infof("Got error: %v while translating %s", err, k)
				continue
			}
			n, err := ygot.MergeStructs(root, d)
			if err != nil {
				log.Infof("Got error: %v while merging %s", err, k)
				continue
			}
			root = n.(*telemetry.Device)
		}
	}

	if debugLog {
		log.Infof("merged processed device, %s", jsonDebug(root))
	}

	return root, nil
}

// translatePortCPUStats maps the Ixia Port CPU Statistics to an OpenConfig
// device structure.
func translatePortCPUStats(in stats.Table, _ []string) (*telemetry.Device, error) {
	pcs, err := in.ParsePortCPUStats()
	if err != nil {
		return nil, err
	}

	d := &telemetry.Device{}
	for _, row := range pcs {
		portName, err := shortPortName(row.PortName)
		if err != nil {
			return nil, err
		}
		port := d.GetOrCreateComponent(portName)
		cpu := d.GetOrCreateComponent(fmt.Sprintf("%s_CPU", portName))

		port.Type = telemetry.PlatformTypes_OPENCONFIG_HARDWARE_COMPONENT_PORT
		port.NewSubcomponent(*cpu.Name)

		cpu.Type = telemetry.PlatformTypes_OPENCONFIG_HARDWARE_COMPONENT_CPU
		cpu.Parent = port.Name

		if row.TotalMemory != nil && row.FreeMemory != nil {
			port.Memory = &telemetry.Component_Memory{
				Available: ygot.Uint64(*row.FreeMemory),
				Utilized:  ygot.Uint64(*row.TotalMemory - *row.FreeMemory),
			}
		}

		if row.CPULoad != nil {
			cpu.GetOrCreateCpu().GetOrCreateUtilization().Instant = ygot.Uint8(uint8(*row.CPULoad))
		}
	}

	if debugLog {
		log.Infof("port CPU device, %s", jsonDebug(d))
	}

	return d, nil
}

// translatePortStats translates the Ixia "Port Stats" view to OpenConfig, returning
// the translated statistics as a ygot-generated struct.
func translatePortStats(in stats.Table, _ []string) (*telemetry.Device, error) {
	ps, err := in.ParsePortStats()
	if err != nil {
		return nil, err
	}

	d := &telemetry.Device{}
	for _, row := range ps {
		portName, err := shortPortName(row.PortName)
		if err != nil {
			return nil, err
		}
		i := d.GetOrCreateInterface(portName)
		i.Counters = &telemetry.Interface_Counters{
			InOctets:  row.BytesRx,
			InPkts:    row.FramesRx,
			OutOctets: row.BytesTx,
			OutPkts:   row.FramesTx,
		}
		i.GetOrCreateEthernet().GetOrCreateCounters().InCrcErrors = row.CRCErrs
		i.InRate = pfloat32Bytes(row.RxRate)
		i.OutRate = pfloat32Bytes(row.TxRate)

		i.Type = telemetry.IETFInterfaces_InterfaceType_ethernetCsmacd
		switch row.LinkState {
		case "":
			// No error when empty.
		case "Link Up":
			i.OperStatus = telemetry.Interface_OperStatus_UP
		case "Link Down", "No PCS Lock":
			i.OperStatus = telemetry.Interface_OperStatus_DOWN
		default:
			return nil, errors.Errorf("statistics row %q has an unmappable port link state %q", row.PortName, row.LinkState)
		}

		// TODO(team): Map different speed interfaces - need to determine what the possible Ixia values are.
		switch row.LineSpeed {
		case "":
			// No error when empty.
		case "10GE LAN":
			i.GetOrCreateEthernet().PortSpeed = telemetry.IfEthernet_ETHERNET_SPEED_SPEED_10GB
		case "100GE":
			i.GetOrCreateEthernet().PortSpeed = telemetry.IfEthernet_ETHERNET_SPEED_SPEED_100GB
		default:
			return nil, errors.Errorf("statistics row %q has an unmappable port link speed %q", row.PortName, row.LineSpeed)
		}
	}

	return d, nil
}

// translateFlowStats translates the Flow Statistics view in the
// supplied stats.Table to a telemetry.Device ygot-generated object.
func translateFlowStats(in stats.Table, itFlows []string) (*telemetry.Device, error) {
	fs, err := in.ParseFlowStats()
	if err != nil {
		return nil, err
	}

	d := &telemetry.Device{}
	for _, row := range fs {
		f := d.GetOrCreateFlow(row.TrafficItem)
		if isIngressTracked(row.TrafficItem, itFlows) {
			it := ingressTrackingFromFlowStats(f, row)
			it.Counters = &telemetry.Flow_IngressTracking_Counters{
				InOctets: row.RxBytes,
				InPkts:   row.RxFrames,
				OutPkts:  row.TxFrames,
			}
			it.LossPct = pfloat32Bytes(row.LossPct)
			it.InRate = pfloat32Bytes(row.RxRate)
			it.InFrameRate = pfloat32Bytes(row.RxFrameRate)
			it.OutRate = pfloat32Bytes(row.TxRate)
			it.OutFrameRate = pfloat32Bytes(row.TxFrameRate)
		} else {
			f.Counters = &telemetry.Flow_Counters{
				InOctets: row.RxBytes,
				InPkts:   row.RxFrames,
				OutPkts:  row.TxFrames,
			}
			f.LossPct = pfloat32Bytes(row.LossPct)
			f.InRate = pfloat32Bytes(row.RxRate)
			f.InFrameRate = pfloat32Bytes(row.RxFrameRate)
			f.OutRate = pfloat32Bytes(row.TxRate)
			f.OutFrameRate = pfloat32Bytes(row.TxFrameRate)
		}
	}
	return d, nil
}

// translateEgressStats translates the Custom Egress Stats view in the
// supplied stats.Table to a telemetry.Device ygot-generated object.
func translateEgressStats(in stats.Table, itFlows []string) (*telemetry.Device, error) {
	es, err := in.ParseEgressStats()
	if err != nil {
		return nil, err
	}

	d := &telemetry.Device{}
	var f *telemetry.Flow
	var it *telemetry.Flow_IngressTracking
	for _, row := range es {
		// Setting the traffic item key, if present.
		if row.TrafficItem != "" {
			f = d.GetOrCreateFlow(row.TrafficItem)
			if isIngressTracked(row.TrafficItem, itFlows) {
				it = ingressTrackingFromFlowStats(f, row.FlowStats)
				it.Filter = &row.Filter
			} else {
				f.Filter = &row.Filter
			}
			continue
		}

		if it == nil {
			ef := f.GetOrCreateEgressTracking(row.Filter)
			ef.Counters = &telemetry.Flow_EgressTracking_Counters{
				InOctets: row.RxBytes,
				InPkts:   row.RxFrames,
				OutPkts:  row.TxFrames,
			}
			ef.LossPct = pfloat32Bytes(row.LossPct)
			ef.InRate = pfloat32Bytes(row.RxRate)
			ef.InFrameRate = pfloat32Bytes(row.RxFrameRate)
			ef.OutRate = pfloat32Bytes(row.TxRate)
			ef.OutFrameRate = pfloat32Bytes(row.TxFrameRate)
		} else {
			ef := it.GetOrCreateEgressTracking(row.Filter)
			ef.Counters = &telemetry.Flow_IngressTracking_EgressTracking_Counters{
				InOctets: row.RxBytes,
				InPkts:   row.RxFrames,
				OutPkts:  row.TxFrames,
			}
			ef.LossPct = pfloat32Bytes(row.LossPct)
			ef.InRate = pfloat32Bytes(row.RxRate)
			ef.InFrameRate = pfloat32Bytes(row.RxFrameRate)
			ef.OutRate = pfloat32Bytes(row.TxRate)
			ef.OutFrameRate = pfloat32Bytes(row.TxFrameRate)
		}
	}

	return d, nil
}

func ingressTrackingFromFlowStats(flow *telemetry.Flow, row *stats.FlowStats) *telemetry.Flow_IngressTracking {
	return flow.GetOrCreateIngressTracking(
		row.RxPort,
		row.TxPort,
		mplsLabelFromUint(row.MPLSLabel),
		row.SrcIPv4,
		row.DstIPv4,
		row.SrcIPv6,
		row.DstIPv6,
	)
}

func mplsLabelFromUint(label *uint64) telemetry.Flow_IngressTracking_MplsLabel_Union {
	if label == nil {
		return telemetry.MplsTypes_MplsLabel_Enum_NO_LABEL
	}
	labelNum := *label
	switch labelNum {
	case 0:
		return telemetry.MplsTypes_MplsLabel_Enum_IPV4_EXPLICIT_NULL
	case 1:
		return telemetry.MplsTypes_MplsLabel_Enum_ROUTER_ALERT
	case 2:
		return telemetry.MplsTypes_MplsLabel_Enum_IPV6_EXPLICIT_NULL
	case 3:
		return telemetry.MplsTypes_MplsLabel_Enum_IMPLICIT_NULL
	case 7:
		return telemetry.MplsTypes_MplsLabel_Enum_ENTROPY_LABEL_INDICATOR
	}
	const maxUint32 = uint64(^uint32(0))
	if labelNum > maxUint32 {
		return telemetry.MplsTypes_MplsLabel_Enum_UNSET
	}
	return telemetry.UnionUint32(labelNum)
}

// Strips off the Ixia name from the port name.
func shortPortName(fullPortName string) (string, error) {
	parts := strings.SplitAfterN(fullPortName, "/", 2)
	if len(parts) < 2 || len(parts[1]) == 0 {
		return "", errors.Errorf("invalid port name: got %q, want [ixia_name]/[port_name]", fullPortName)
	}
	return parts[1], nil
}

func pfloat32Bytes(f *float32) telemetry.Binary {
	if f == nil {
		return nil
	}
	return float32Bytes(*f)
}

func float32Bytes(f float32) telemetry.Binary {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, math.Float32bits(f))
	return telemetry.Binary(b)
}

// jsonDebug renders a device struct to JSON such that it can be logged.
// Since it is used only in debugging, it logs an error if the device struct
// cannot be rendered to JSON, and returns an error string to the user.
func jsonDebug(d *telemetry.Device) string {
	js, err := ygot.ConstructIETFJSON(d, nil)
	if err != nil {
		log.Errorf("cannot render device to JSON during debugging, %v", err)
		return "(unrenderable)"
	}

	j, err := json.MarshalIndent(js, "", "  ")
	if err != nil {
		log.Errorf("cannot marshal JSON for device, %v", err)
		return "(unmarshallable)"
	}
	return string(j)
}

func isIngressTracked(flow string, itFlows []string) bool {
	for _, f := range itFlows {
		if f == flow {
			return true
		}
	}
	return false
}
