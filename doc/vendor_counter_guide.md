# Vendor-Specific Augmentation for Pipeline Counter

**Contributors**: roland@arista.com

This document provides the guidelines for the vendor-specific portions of openconfig-platform. As implementations differ from vendor to vendor and platform to platform, a process of adding vendor-specific counters will be defined here.

## Usage: Vendor-specific pipeline drop counter

Each implementor should augment `/components/component/integrated-circuit/pipeline-counter/drop/vendor` with their own vendor and platform containers. The naming of the platform container may consist of the platform name, ASIC family, or a combination of both platform and ASIC family. Within the platform container, that container may use the utility grouping `oc-ppc:pipeline-vendor-drop-containers` that provides the adverse/congestion/packet-processing specific containers. For each set of adverse/congestion/packet-processing counters augmented into `oc-ppc:pipeline-vendor-drop-containers`, the sum of the counters should be included in the values of the aggregate leaves:

- Counters within `.../pipeline-counter/drop/vendor/<vendor>/<platform>/adverse/state` aggregate into `.../pipeline-counter/drop/state/adverse-aggregate`
- Counters within `.../pipeline-counter/drop/vendor/<vendor>/<platform>/congestion/state`  aggregate into `.../pipeline-counter/drop/state/congestion-aggregate`
- Counters within `.../pipeline-counter/drop/vendor/<vendor>/<platform>/packet-processing/state` aggregate into `.../pipeline-counter/drop/state/packet-processing-aggregate`

If these aggregate counters are implemented, the sum of the vendor-specific counters must match the aggregate counters.

If an integrated-circuit has a vendor-specific packet drop counter which cannot differentiate between packet-processing, congestion and adverse drops, then that counter should still be exposed as a vendor-specific packet-processing counter with an appropriate description in the vendor's augmentation.   The `packet-processing-aggregate` counter should be incremented in this scenario as expected above.  Such a counter is undesirable as it means this hardware cannot meet the goal for identifying adverse packet drops in the ASIC, but it is better not to ruin the fidelity of the `adverse-aggregate` drop counter with noise of intended packet drops.

## Example: Vendor-specific pipeline drop counter

This following is a sample augmentation file.

- Vendor: Acme
- Platform: AsicFamily

### Example YANG Augmentation

release/platform/acme-asicfamily-drop-augments.yang

```yang
grouping acme-asicfamily-adverse-drop-counters {
  leaf adverse-reason-counter-a {
    type oc-yang:counter64;
  }

  leaf adverse-reason-counter-b {
    type oc-yang:counter64;
  }

  leaf adverse-reason-counter-c {
    type oc-yang:counter64;
  }
}

grouping acme-asicfamily-congestion-drop-counters {
  leaf congestion-reason-counter-a {
    type oc-yang:counter64;
  }

  leaf congestion-reason-counter-b {
    type oc-yang:counter64;
  }

  leaf congestion-reason-counter-c {
    type oc-yang:counter64;
  }
}

grouping acme-asicfamily-packet-processing-drop-counters {
  leaf packet-processing-reason-counter-a {
    type oc-yang:counter64;
  }

  leaf packet-processing-reason-counter-b {
    type oc-yang:counter64;
  }

  leaf packet-processing-reason-counter-c {
    type oc-yang:counter64;
  }
}

augment "/components/component/integrated-circuit/pipeline-counter/drop/vendor" {
  container acme {
    container asic-family {
      uses oc-ppc:pipeline-vendor-drop-containers;
    }
  }
}

augment "/components/component/integrated-circuit/pipeline-counter/drop/vendor/acme/asic-family/adverse/state" {
  uses acme-asicfamily-adverse-drop-counters;
}

augment "/components/component/integrated-circuit/pipeline-counter/drop/vendor/acme/asic-family/congestion/state" {
  uses acme-asicfamily-congestion-drop-counters;
}

augment "/components/component/integrated-circuit/pipeline-counter/drop/vendor/acme/asic-family/adverse/state" {
  uses acme-asicfamily-packet-processing-drop-counters;
}
```

Note: Namespaces omitted from `augment <path>` for brevity

### Example pyang tree

```text
module: openconfig-platform
  +--rw components
     +--rw component* [name]
        +--rw integrated-circuit
           +--ro oc-ppc:pipeline-counters
              +--ro oc-ppc:drop
                 +--ro oc-ppc:vendor
                    +--ro acme-ppc:acme
                      +--ro acme-ppc:asic-family
                        +--ro oc-ppc:adverse
                           +--ro oc-ppc:state
                              +--ro acme-ppc:adverse-reason-counter-a?    oc-yang:counter64
                              +--ro acme-ppc:adverse-reason-counter-b?    oc-yang:counter64
                              +--ro acme-ppc:adverse-reason-counter-c?    oc-yang:counter64
                        +--ro oc-ppc:congestion
                           +--ro oc-ppc:state
                              +--ro acme-ppc:congestion-reason-counter-a?    oc-yang:counter64
                              +--ro acme-ppc:congestion-reason-counter-b?    oc-yang:counter64
                              +--ro acme-ppc:congestion-reason-counter-c?    oc-yang:counter64
                        +--ro oc-ppc:packet-processing
                           +--ro oc-ppc:state
                              +--ro acme-ppc:packet-processing-reason-counter-a?    oc-yang:counter64
                              +--ro acme-ppc:packet-processing-reason-counter-b?    oc-yang:counter64
                              +--ro acme-ppc:packet-processing-reason-counter-c?    oc-yang:counter64
```

## Usage: Vendor-specific control-plane traffic counter

Each implementor should augment `/components/component/integrated-circuit/control-plane-traffic/vendor` with their own `<vendor>/<platform>/state` containers. The naming of the platform container may consist of the platform name, ASIC family, or a combination of both platform and ASIC family. Within the state container, it should define vendor-specific counter containers. Each control-plane traffic counter may use the utility grouping `oc-ic:control-plane-traffic-vendor-counters` that provides the queued/dropped leaves. For each counters augmented into `<vendor>/<platform>/state`, the sum of the counters should be included in the values of the aggregate leaves:

- Counters within `.../control-plane-traffic/vendor/<vendor>/<platform>/state/<counter>/queued` aggregate into `.../control-plane-traffic/state/queued-aggregate`
- Counters within `.../control-plane-traffic/vendor/<vendor>/<platform>/state/<counter>/queued-bytes` aggregate into `.../control-plane-traffic/state/queued-bytes-aggregate`
- Counters within `.../control-plane-traffic/vendor/<vendor>/<platform>/state/<counter>/dropped` aggregate into `.../control-plane-traffic/state/dropped-aggregate`
- Counters within `.../control-plane-traffic/vendor/<vendor>/<platform>/state/<counter>/dropped-bytes` aggregate into `.../control-plane-traffic/state/dropped-bytes-aggregate`

If these aggregate counters are implemented, the sum of the vendor-specific counters must match the aggregate counters.

## Example: Vendor-specific pipeline drop counter

This following is a sample augmentation file.

- Vendor: Acme
- Platform: AsicFamily

### Example YANG Augmentation

release/platform/acme-asicfamily-control-plane-traffic-augments.yang

```yang
grouping acme-asicfamily-control-plane-traffic-counters {
  container queue-counter-a {
	 uses oc-ppc:control-plane-traffic-vendor-counters;
  }

  container queue-counter-b {
	 uses oc-ppc:control-plane-traffic-vendor-counters;
  }

  container queue-counter-c {
	 uses oc-ppc:control-plane-traffic-vendor-counters;
  }
}

augment "/components/component/integrated-circuit/pipeline-counters/control-plane-traffic/vendor" {
  container acme {
    container asic-family {
      container state {
        uses acme-asicfamily-control-plane-traffic-counters;
      }
    }
  }
}
```

Note: Namespaces omitted from `augment <path>` for brevity

### Example pyang tree

```text
module: openconfig-platform
  +--rw components
     +--rw component* [name]
        +--rw integrated-circuit
           +--ro oc-ppc:pipeline-counters
              +--ro oc-ppc:control-plane-traffic
                 +--ro oc-ppc:vendor
                    +--ro acme-ppc:acme
                       +--ro acme-ppc:asic-family
                          +--ro acme-ppc:state
                             +--ro acme-ppc:queue-counter-a
                             +--ro acme-ppc:queue-counter-b
                             +--ro acme-ppc:queue-counter-c
```
