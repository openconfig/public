# Intergrated Circuit aggregated pipeline counters guide
## Introduction
This guide discusses semantics of different counters provided under the
`openconfig-platform/components/component/integrated-circuit/pipeline-counters` container.
The `INTEGRATED_CIRCUIT` or I-C, in this document refers to the OpenConfig [INTEGRATED_CIRCUIT](https://github.com/openconfig/public/blob/5d38d8531ef9c5b998262207eb6dbdae8968f9fe/release/models/platform/openconfig-platform-types.yang#L346) component type which is typically an ASIC or NPU (or combination of both) that provides packet processing capabilities.

## Per-block packets/octets counters
[TODO] more detailed description
## Drop packets/octets counters
The `/components/component/integrated-circuit/pipeline-counters/drop` container collects counters related to packets dropped by the `INTEGRATED_CIRCUIT`.
### Aggregated drop counters
These 4 counters should cover all packets dropped by the IC which are not already covered by the /interfaces tree.   For example, a packet which is dropped due to QoS policy for WRED should be counted only by the appropriate /interfaces path [dropped-pkts](https://github.com/openconfig/public/blob/5d38d8531ef9c5b998262207eb6dbdae8968f9fe/release/models/qos/openconfig-qos-interfaces.yang#L375).    

Aggregated drop counters are modeled as below:
```
module: openconfig-platform
  +--rw components
     +--rw component* [name]
        +--rw integrated-circuit
           +--ro oc-ppc:pipeline-counters
              +--ro oc-ppc:drop
                 +--ro oc-ppc:state
                    +--ro oc-ppc:adverse-aggregate?             oc-yang:counter64
                    +--ro oc-ppc:congestion-aggregate?          oc-yang:counter64
                    +--ro oc-ppc:packet-processing-aggregate?   oc-yang:counter64
                    +--ro oc-ppc:urpf-aggregate?                oc-yang:counter64
```
#### urpf-aggregate

##### Usability
The increments of this counter are typically signal of some form of attack with spoofed sourec address. Typically dDOS class.

#### packet-processing-aggregate

##### Usability
The increments of this counter are expected during convergence events as well as during stable operation. However rapid increase in drop rate **may** be a signal of network being unhealthy and typically requires further investigation. 
The further break down of this counter, if available as vendor extension under `/openconfig-platform:components/component/integrated-circuit/openconfig-platform-pipeline-counters:pipeline-counters/drop/vendor` container could help to further narrow-down cause of drops. 

If prolonged packet drops are found to be caused by lack of FIB entry for incomming packets, this suggest inconsistency between Network Control plane protocols (BGP, IGP, RSVP, gRIBI), FIB calculated by Controller Card and FIB programmed into given Integrated Circuit.

If implemetation supports `urpf-aggregate` counter, packets discarded due to uRPF should not be counted as `packet-processing-aggregate`. Else, uRPF discarded oacket should be counted against this counter.

#### congestion-aggregate


##### Usability
The increments of this counter are signal of given Integrated Circuit being overhelmed by incomming traffic and complexity of packet processing that is required. 

#### adverse-aggregate
##### Usability
The increments of this counter are generally a signal of a hardware defect (e.g. memory errors or signal integrity issues) or (micro)code software defects. 

#### Queue tail and AQM drops exeption discussion.
Drops associated with QoS queue tail or AQM are the result of egress interface congestion.   This is NOT the same as I-C congestion, and should be counted using /interfaces counters as it is expected state from the platform (router) point of view. It may be not expected state from a network design point of view but from the INTEGRATED_CIRCUIT, it is behaving according to design.   

The OpenConfig definition for [congestion-aggregate](https://github.com/openconfig/public/blob/5d38d8531ef9c5b998262207eb6dbdae8968f9fe/release/models/platform/openconfig-platform-pipeline-counters.yang#L1096-L1099) excludes "queue drop counters". It desirable to  not count QoS queue drops under this `congestion-aggregate` in order to maintain a clear signal of hitting I-C performance limitations, rather then blend it with basic, simple egress interface speed limitations.

### Per-Block drop copunters
[TODO] more detailed description for standard OpenConfig drop counters defined for Interface-, Lookup-, Queueing-, Fabric-  and Host-Interface- blocks. Also discuss relationship with Control plane traffic packets/octets counters.
### Vendor extensions
Please refer to [Vendor-Specific Augmentation for Pipeline Counter](vendor_counter_guide.md)
## Error counters
This counters **do not** counts **packets or bytes**.
They counte error events per block.

For example corruption of on chip, HBM or chip external memory buffers (soft-error) which also are not already counted as queue drops for interfaces.

[TODO] more detailed description
## Control plane traffic packets/octets counters
[TODO] more detailed description. Also discuss relationship with Host-Interface block counters.
### Standard OpenConfig counters
### Vendor extensions
