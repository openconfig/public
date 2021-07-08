// Package stats provides functionality for gathering Ixia statistics.
package stats

import (
	"strconv"

	log "github.com/golang/glog"
	"github.com/pkg/errors"
)

// View is a statistics view.
type View string

const (
	// PortStatsView is the Port Statistics view.
	PortStatsView View = "Port Statistics"
	// PortCPUStatsView is the Port CPU Statistics view.
	PortCPUStatsView View = "Port CPU Statistics"
	// FlowStatsView is the Flow Statistics view.
	FlowStatsView View = "Flow Statistics"
	// EgressStatsView is the custom Egress Statistics view.
	EgressStatsView View = "EgressStatView"
)

// Row represents an individual row of statistics from an Ixia page.
// The key of the map is the column name, and the value of that stat.
// Ixia always provides the value of statistics as strings.
type Row map[string]string

func (r Row) parse(nameKey string, strKeys []string, intKeys []string, floatKeys []string) (map[string]*uint64, map[string]*float32, error) {
	name := r[nameKey]
	if name == "" {
		return nil, nil, errors.Errorf("row %v did not include value for required key %q", r, nameKey)
	}
	lookup := func(key string) string {
		v, ok := r[key]
		// Warn if the stat column is present but empty.
		if ok && v == "" {
			log.Warningf("got empty stat %q for key %q", key, name)
		}
		return v
	}

	for _, k := range strKeys {
		lookup(k) // log only
	}
	ints := make(map[string]*uint64)
	for _, k := range intKeys {
		v := lookup(k)
		if v == "" {
			continue
		}
		i, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "invalid value %q for int stat %q for key %q", v, k, name)
		}
		ints[k] = &i
	}
	floats := make(map[string]*float32)
	for _, k := range floatKeys {
		v := lookup(k)
		if v == "" {
			continue
		}
		f, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "invalid value %q for float stat %s for key %s", v, k, name)
		}
		f32 := float32(f)
		floats[k] = &f32
	}
	return ints, floats, nil
}

// Table is a set of Rows that combine to make a single statistics view
// returned from Ixia.
type Table []Row

// PortStats is the statistics in a Port Statistics view.
type PortStats struct {
	PortName, LineSpeed, LinkState                string
	FramesTx, FramesRx, BytesTx, BytesRx, CRCErrs *uint64
	TxRate, RxRate                                *float32
}

// ParsePortStats parses the table to port statistics.
func (t Table) ParsePortStats() ([]*PortStats, error) {
	const (
		nameKey      = "Port Name"
		lineSpeedKey = "Line Speed"
		linkStateKey = "Link State"
		framesTxKey  = "Frames Tx."
		framesRxKey  = "Valid Frames Rx."
		bytesTxKey   = "Bytes Tx."
		bytesRxKey   = "Bytes Rx."
		txRateKey    = "Tx. Rate (bps)"
		rxRateKey    = "Rx. Rate (bps)"
		crcErrsKey   = "CRC Errors"
	)
	strKeys := []string{lineSpeedKey, linkStateKey}
	intKeys := []string{framesTxKey, framesRxKey, bytesTxKey, bytesRxKey, crcErrsKey}
	floatKeys := []string{txRateKey, rxRateKey}

	var ps []*PortStats
	for _, row := range t {
		intVals, floatVals, err := row.parse(nameKey, strKeys, intKeys, floatKeys)
		if err != nil {
			return nil, err
		}
		ps = append(ps, &PortStats{
			PortName:  row[nameKey],
			LineSpeed: row[lineSpeedKey],
			LinkState: row[linkStateKey],
			FramesTx:  intVals[framesTxKey],
			FramesRx:  intVals[framesRxKey],
			BytesTx:   intVals[bytesTxKey],
			BytesRx:   intVals[bytesRxKey],
			CRCErrs:   intVals[crcErrsKey],
			TxRate:    floatVals[txRateKey],
			RxRate:    floatVals[rxRateKey],
		})
	}
	return ps, nil
}

// PortCPUStats is the statistics in a Port CPU Statistics view.
type PortCPUStats struct {
	PortName                         string
	TotalMemory, FreeMemory, CPULoad *uint64
}

// ParsePortCPUStats parses the table to port CPU statistics.
func (t Table) ParsePortCPUStats() ([]*PortCPUStats, error) {
	const (
		portNameKey = "Port Name"
		totalMemKey = "Total Memory(KB)"
		freeMemKey  = "Free Memory(KB)"
		cpuLoadKey  = "%CPU Load"
	)
	var strKeys []string
	intKeys := []string{totalMemKey, freeMemKey, cpuLoadKey}
	var floatKeys []string

	var pcs []*PortCPUStats
	for _, row := range t {
		intVals, _, err := row.parse(portNameKey, strKeys, intKeys, floatKeys)
		if err != nil {
			return nil, err
		}
		pcs = append(pcs, &PortCPUStats{
			PortName:    row[portNameKey],
			TotalMemory: intVals[totalMemKey],
			FreeMemory:  intVals[freeMemKey],
			CPULoad:     intVals[cpuLoadKey],
		})
	}
	return pcs, nil
}

// FlowStats is the statistics in a Flow Statistics view.
type FlowStats struct {
	TrafficItem                                       string
	RxBytes, TxFrames, RxFrames                       *uint64
	LossPct, TxRate, RxRate, TxFrameRate, RxFrameRate *float32
	// Optional ingress-tracking fields.
	RxPort, TxPort, SrcIPv4, DstIPv4, SrcIPv6, DstIPv6 string
	MPLSLabel                                          *uint64
}

// ParseFlowStats parses the table to flow statistics.
func (t Table) ParseFlowStats() ([]*FlowStats, error) {
	const (
		trafficItemKey = "Traffic Item"
		txFramesKey    = "Tx Frames"
		rxFramesKey    = "Rx Frames"
		lossPctKey     = "Loss %"
		txFrameRateKey = "Tx Frame Rate"
		rxFrameRateKey = "Rx Frame Rate"
		rxBytesKey     = "Rx Bytes"
		txRateKey      = "Tx Rate (bps)"
		rxRateKey      = "Rx Rate (bps)"
		// Optional ingress tracking fields.
		rxPortKey    = "Rx Port"
		txPortKey    = "Tx Port"
		mplsLabelKey = "MPLS:Label Value"
		srcIPv4Key   = "IPv4 :Source Address"
		dstIPv4Key   = "IPv4 :Destination Address"
		srcIPv6Key   = "IPv6 :Source Address"
		dstIPv6Key   = "IPv6 :Destination Address"
	)
	strKeys := []string{rxPortKey, txPortKey, srcIPv4Key, dstIPv4Key, srcIPv6Key, dstIPv6Key}
	intKeys := []string{rxBytesKey, txFramesKey, rxFramesKey, mplsLabelKey}
	floatKeys := []string{lossPctKey, txRateKey, rxRateKey, txFrameRateKey, rxFrameRateKey}

	var fs []*FlowStats
	for _, row := range t {
		intVals, floatVals, err := row.parse(trafficItemKey, strKeys, intKeys, floatKeys)
		if err != nil {
			return nil, err
		}
		fs = append(fs, &FlowStats{
			TrafficItem: row[trafficItemKey],
			RxBytes:     intVals[rxBytesKey],
			TxFrames:    intVals[txFramesKey],
			RxFrames:    intVals[rxFramesKey],
			LossPct:     floatVals[lossPctKey],
			TxRate:      floatVals[txRateKey],
			RxRate:      floatVals[rxRateKey],
			TxFrameRate: floatVals[txFrameRateKey],
			RxFrameRate: floatVals[rxFrameRateKey],
			// Optional ingress tracking fields.
			RxPort:    row[rxPortKey],
			TxPort:    row[txPortKey],
			MPLSLabel: intVals[mplsLabelKey],
			SrcIPv4:   row[srcIPv4Key],
			DstIPv4:   row[dstIPv4Key],
			SrcIPv6:   row[srcIPv6Key],
			DstIPv6:   row[dstIPv6Key],
		})
	}
	return fs, nil
}

// EgressStats is the statistics in a Egress Statistics view.
type EgressStats struct {
	*FlowStats
	Filter string
}

// ParseEgressStats parses the table to egress statistics.
func (t Table) ParseEgressStats() ([]*EgressStats, error) {
	const egressKey = "Egress Tracking"
	ps, err := t.ParseFlowStats()
	if err != nil {
		return nil, err
	}
	var es []*EgressStats
	for i, row := range t {
		es = append(es, &EgressStats{
			FlowStats: ps[i],
			Filter:    row[egressKey],
		})
	}
	return es, nil
}
