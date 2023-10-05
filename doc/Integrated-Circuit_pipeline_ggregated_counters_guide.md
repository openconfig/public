# Intergrated Circuit aggregated pipeline counters guide
## Introduction
This guide discusses semantics of different counters provided under the
`openconfig-platform/components/component/integrated-circuit/pipeline-counters` container.
The `INTEGRATED_CIRCUIT` or I-C, in this document refers to the OpenConfig [INTEGRATED_CIRCUIT](https://github.com/openconfig/public/blob/5d38d8531ef9c5b998262207eb6dbdae8968f9fe/release/models/platform/openconfig-platform-types.yang#L346) component type which is typically an ASIC or NPU (or combination of both) that provides packet processing capabilities.

## Per-block packets/octets counters
[TODO] more detailed description
## Drop packets/octets counters
The drop container collects counters related to packet dropped for varouus reasons and in varous places inside "Integrated Circuit".
### Aggregated drop counters
This 4 counters should cover all packets dropped inside I-C with one exeption - packet driopped due to QoS queue tail-drop or AQM (RED/WRED).  Aggregated drop couters are modeled as below:
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
> From OpenConfig definition:\
>This aggregation of counters represents the conditions in which packets are dropped due to failing uRPF lookup check. This counter and the packet-processing-aggregate counter should be incremented for each uRPF packet drop.

This counter counts packet discarded as resutlt of Unicast Reverse Path Forwarding verification. ([RFC2827](https://datatracker.ietf.org/doc/html/rfc2827), [RFC3704](https://datatracker.ietf.org/doc/html/rfc3704)).

##### Usability
The increments of this counter are typically signal of some form of attack with spoofed sourec address. Typically dDOS class.

#### packet-processing-aggregate
> From OpenConfig definition:\
> This aggregation of counters represents the conditions in which packets are dropped due to legitimate forwarding decisions (ACL drops, No Route etc.)

This counter counts packet discarded as resutlt of processing **non-corrupted packtet** against **legitimate, non-corrupted** state of I-C program (FIB content, ACL content, rate-limiting token-bucktes) which mandate packet drop. The examples of this class of discard are:
- dropping packets which destination address to no match any FIB entry
- dropping packets which destination address matches FIB entry pinting discard next-hop (e.g. route to null0)
- dropping packts due to ACL/packet filter decission
- dropping packets due to its TTL = 1
- dropping packts due to its size exceeds egress interface MTU and packet ca'nt be fragmented (IPv6 or Dont-Fragmemt bit is set)
- etc

Note: Form the I-C perspective it is doing exectly what it is told (programed) to do, and packet is parsable.

##### Usability
The increments of this counter are expected during convergence events as well as during stable operation. However rapid increase in drop rate **may** be a signal of network being unhealthy and typically requires further investigation. 
The further break down of this counter, if available as vendor extension under `/openconfig-platform:components/component/integrated-circuit/openconfig-platform-pipeline-counters:pipeline-counters/drop/vendor` container could help to further narrow-down cause of drops. 

If prolonged packet drops are found to be caused by lack of FIB entry for incomming packets, this suggest inconsistency between Network Control plane protocols (BGP, IGP, RSVP, gRIBI), FIB calculated by Controller Card and FIB programmed into given Integrated Circuit.

If implemetation supports `urpf-aggregate` counter, packets discarded due to uRPF should not be counted as `packet-processing-aggregate`. Else, uRPF discarded oacket should be counted against this counter.

#### congestion-aggregate
> From OpenConfig definition:\
>This tracks the aggregation of all counters where the expected conditions of packet drops due to internal congestion in some block of the hardware that **may not be visible** in through other congestion indicators like interface discards or **queue drop counters**.

This counter counts packet discarded as resutlt of exceedding performance limits of Integrated-Circuit, when it sprocess non-corruptec packets against to legitimate, non-corrupted progreamming state of I-C.

The typial example is overloading given IC with higher packet rate (pps) then given chip can handle. For exeple, let's assume chip X can process 3.6bps of incomming traffic and 2000 Mpps. However if averange incoming packet size is 150B, at full ingress rate this become 3000Mpps. Hence 1/3 of packets would be cropped and should be counted against `congestion-aggregate`.

Another example is the case when some I_C data bus is too narrow/slow for handling traffic. For example let's assume chip X needs to sent 3Tbps of it's ingress traffic to external buffer memory, which has only 2Tbps access I/O. It this case pactes would be discarded, because of internal congestion of memory I/O bus. Note, this packet are discarded even if queues are very little used, hence this are NOT QoS queue tail-drops nor WRED drops.

Yet another example is the case where extreemly large and long ACL/filter requires more cycles to process then NPU is bugeted for. 

##### Usability
The increments of this counter are signal of given Integrated Circuit being overhelmed by incomming traffic and complexity of packet processing that is required. 

#### adverse-aggregate
> From OpenConfig definition:\
> This captures the aggregation of all counters where the switch is **unexpectedly** dropping packets. Occurrence of these drops on a stable (no recent hardware or config changes) and otherwise healthy switch needs further investigation.

This counter counts packet discarded as resutlt of **corrupted** programming state in I-C or **corrupted** data structures of packet descriptors.

Note: corrupted packet recived on ingress interface should be counted separatly in `/interfaces/interface/state/counters/in-errors` and NOT counted as `adverse-aggregate`. This is because incomming corrupted packets are NOT a signal of adverse state of given I-C (but rather of upstream system). Therefore it is better not to count such drops as `adverse-aggregate` to keep it clean signal of I-C adverse state.


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
This counters do not counts packets of bytes. They counte error events per block.
[TODO] more detailed description
## Control plane traffic packets/octets counters
[TODO] more detailed description. Also discuss relationship with Host-Interface block counters.
### Standard OpenConfig counters
### Vendor extensions
