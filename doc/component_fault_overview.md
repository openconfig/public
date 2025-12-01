# Openconfig Component Fault design

**Contributors**: evangoodwin@google.com

## Problem Statement

`/system/alarms` provide a way to express that the device is in some kind of distress. However, alarms lack component attribution or a suggested response to the alarm. This requires operators to build up this knowledge over time. The ambition of this proposal is that since vendors already have this context (via FMEA) they can provide it along with the distress signal.

## Proposal

This proposal introduces a list of faults under a component's healthz container.

A fault indicates that a component has recently or is currently experiencing a problem as expressed by a symptom. These symptoms have underlying conditions that are platform-specific. The fault should not be deleted as soon as the underlying condition is no longer asserted to avoid excessive fault creation and deletion. Deletion of the fault should occur after a suitable period outside of the triggering condition

Nested under each fault is a list of actions recommended to remediate the fault condition.

```
        +--rw oc-platform-healthz:healthz
        |  +--rw faults <<<<<<<<<<<<<<<<<<<<<<<
        |     +--ro fault* [symptom]
        |        +--ro symptom         -> ../state/symptom
        |        +--ro state
        |        |  +--ro symptom?               identityref
        |        |  +--ro origin-time?           oc-types:timeticks64
        |        |  +--ro last-detection-time?   oc-types:timeticks64
        |        |  +--ro description?           string
        |        |  +--ro status?                enumeration
        |        |  +--ro counters
        |        |     +--ro occurrences?   oc-yang:counter64
        |        +--ro remediations
        |           +--ro remediation* [index]
        |              +--ro index    -> ../state/index
        |              +--ro state
        |                 +--ro index?    uint64
        |                 +--ro action?   identityref
        |                 +--ro target?   -> /oc-platform:components/component/name
```

## Example Use Cases
Repair Automation - By providing a list of remediations, when a fault occurs on a component, operators have a list of steps to inform repair strategies.This is the primary use case.

RMA “Contract” - The vendor can provide a list of steps they would like performed before an RMA is issued. Note the description of remediations indicates that they are not compulsory. So reception of these suggestions does not enforce they were performed before an RMA requested.


## Examples
A missing PSU component that can be recovered with a sequence of actions: reseat and/or device reboot.
```
"openconfig-platform:components": {
  "component": [
    {
      "name": "PSU 1",
      "openconfig-platform-healthz": {
        "state": {
          "status": "UNHEALTHY"
        },
        "openconfig-platform-healthz-fault:faults": {
          "fault": [
            {
              "symptom": "openconfig-platform-healthz-fault:SYMPTOM_MISSING_COMPONENT",
              "state": {
                "description": "PSU1::Missing",
                "last-detection-time": "1763937335936939878",
                "occurrences": "1",
                "origin-time": "1763764161843510422",
                "status": "ACTIVE",
                "symptom": "openconfig-platform-healthz-fault:SYMPTOM_MISSING_COMPONENT"
              },
              "remediations": {
                "remediation": [
                  {
                    "index": "0",
                    "state": {
                      "action": "openconfig-platform-healthz-fault:ACTION_RESEAT",
                      "target": "PSU 1"
                    }
                  },
                  {
                    "index": "1",
                    "state": {
                      "action": "openconfig-platform-healthz-fault:ACTION_WARM_REBOOT",
                      "target": "PSU 1"
                    }
                  }
                ]
              }
            }
          ]
        }
      }
    }
  ]
}
```

A hot temperature sensor that requires replacing a different component (“PSU 1”).
```
"openconfig-platform:components": {
  "component": [
    {
      "name": "temp-sensor",
      "openconfig-platform-healthz": {
        "state": {
          "status": "UNHEALTHY"
        },
        "openconfig-platform-healthz-fault:faults": {
          "fault": [
            {
              "symptom": "openconfig-platform-healthz-fault:SYMPTOM_OVER_THRESHOLD",
              "state": {
                "description": "temp too high; 100>90",
                "last-detection-time": "1763937335936939878",
                "occurrences": "1",
                "origin-time": "1763764161843510422",
                "status": "ACTIVE",
                "symptom": "openconfig-platform-healthz-fault:SYMPTOM_OVER_THRESHOLD"
              },
              "remediations": {
                "remediation": [
                  {
                    "index": "0",
                    "state": {
                      "action": "openconfig-platform-healthz-fault:ACTION_REPLACE",
                      "target": "PSU 1"
                    }
                  }
                ]
              }
            }
          ]
        }
      }
    }
  ]
}
```

An INACTIVE fault, which means that the underlying condition was not present during the last sampling. Note that some level of fault dampening is expected. Consider a sensor value where the instantaneous value oscillates along the boundary of its operating range. Once the fault is detected and reported, the expectation of this feature is that the underlying condition needs to be absent for a sufficient amount of time (to be defined by the vendor) before the fault is removed. Occurrences >1 indicates that some oscillation has occurred.
```
"openconfig-platform:components": {
  "component": [
    {
      "name": "temp-sensor",
      "openconfig-platform-healthz": {
        "state": {
          "status": "UNHEALTHY"
        },
        "openconfig-platform-healthz-fault:faults": {
          "fault": [
            {
              "symptom": "openconfig-platform-healthz-fault:SYMPTOM_OVER_THRESHOLD",
              "state": {
                "description": "temp too high; 100>90",
                "last-detection-time": "1763937335936939878",
                "occurrences": "12",
                "origin-time": "1763764161843510422",
                "status": "INACTIVE",
                "symptom": "openconfig-platform-healthz-fault:SYMPTOM_OVER_THRESHOLD"
              },
            }
          ]
        }
      }
    }
  ]
}
```



## Design Decisions/Questions
> The component model contains both hardware and software components. Are these faults limited to hardware components?

This proposal is motivated by hardware component failures, but there is no reason to prevent fault reporting on software components should they rise to the healthz “UNHEALTHY” standard; ie “the component is not performing the function expected of it.”

> Is the operator expected to ACK or clear remediation actions or faults?
There is no expectation that the operator should ack/clear the event. If the underlying fault condition is no longer asserted, after a vendor-defined cool down period (in the interest of fault dampening), the fault is expected to be removed by the device.

> Must the operator perform all remediations?

Not all actions in this list are necessarily required. Ultimately any fault can be mitigated by replacing the chassis, but that's expensive. Providing alternatives enables the operator to explore less expensive solutions first. The device is not expected to know whether any        remediations have previously been performed, so this list is not expected to change with progressive interventions. Additionally, the expense of executing a remediation depends on the deployment context and is ultimately defined by the operator. 

> How about fault categories or severities?

Severity was in the early proposal, but was ultimately removed. If a “component is no longer able to perform the function expected of it and requires remediation/intervention” that we’d always be in a “severe” place.

> What about historical and/or clear counts? e.g. the fault is not active now but there was prior record.

Many of the remediations (reboot, reimage, ect) would require this history to persist across these actions, which was deemed expensive to support. The prevailing thought was that the operator receiving these faults could monitor for recidivism and escalate as necessary.
