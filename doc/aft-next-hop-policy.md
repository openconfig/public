# Next-hop Policy-tag Forwarding Policy

**Contributors**: henryhu@google.com, dloher@google.com, robjs@google.com
Date: Created Nov 11, 2026

## Summary
 
This document specifies next-hop policy-tag based policy-forwarding. A policy-tag is an attribute of an AFT next-hop, and a policy-forwarding rule can match on that policy-tag to apply actions to packets forwarded through the tagged next-hops.

## Background

OpenConfig policy-forwarding matches and applies actions to packets. Policy-forwarding rules match on various properties of the packets, such as destination IP, TTL, DSCP, protocol and apply actions, such as drop the packet, set TTL or set the next-hop.  These rules require inspecting fields inside the packet.

gRIBI can be used to inject routes and next-hops. In addition to the traditional next-hop attributes like destination IP address, additional attributes can be added such as encapsulation headers.

In certain cases, it’s desirable to match gRIBI next-hops in policy-forwarding policies. This allows policy-forwarding actions to be applied to packets sent to those next-hops. In this scenario, the policy-forwarding match rule is evaluated against next-hop table entries instead of matching on fields inside the packet.

## Next-hop Policy-tag Specification

This specification defines a policy-tag attribute for gRIBI next-hops, and a policy-forwarding rule match criterion for next-hops with such policy-tags.

When packets are forwarded using a route that resolves to a next-hop containing a policy-tag, such packets become subject to the actions of a policy-forwarding policy with a match rule for the same policy-tag.

### Details

Next-hop policy-tag forwarding is enabled in 2 steps:

  * First, policy-forwarding policies are configured which match gRIBI next-hop policy-tags and specify actions.
  * Second, routes and next-hops with policy-tags are injected using gRIBI.

The OpenConfig AFT model includes a policy-tag leaf in the next-hop container.  gRIBI is used to program next-hops which include the policy-tag leaf.


Packets forwarded to the matching next-hops have the policy-forwarding action applied to them. For example, a policy-forwarding policy matches a next-hop policy-tag and specifies a policer action. All packets forwarded through next-hops with that policy-tag are policed by that policer.  If a next-hop has a policy-tag with no matching policy-forwarding policy, no actions are applied.

### Use Cases

In one use case, gRIBI routes are programmed on a device which directs packets with different destination IPs to be encapsulated into different tunnels.  The operator applies a common bandwidth limit to the tunnels sharing the same destination, because such tunnels typically share the same path across the backbone network.

Without policy-tags, this requires policy-forwarding policies which match specific destination prefixes of the packets and apply policers. Such policies require matching of all destination prefixes which may be more numerous than next-hops, consuming additional hardware resources on the device. In addition, policy-forwarding policies are often implemented with TCAM which typically has more limited capacity than LPM used for next-hops.  Using next-hop policy-tag, the lookup is done just once in LPM, and the policy-tag is part of the LPM lookup result. The policy-tag identifies a policy-forwarding policy which then specifies a policer action to be applied. 

### Example

Consider a device with 2 interfaces. Interface I1 has 3 gRIBI routes:

- Route 1: match prefix A, forward through next hop 1: send packets through tunnel 1, whose destination is in data center X
- Route 2: match prefix B, forward through next hop 2: send packets through tunnel 2, whose destination is in data center X
- Route 3: match prefix C, forward through next hop 3: send packets through tunnel 3, whose destination is in data center Y

Interface I2 has 1 route:

- Route 4: match prefix D, forward through next hop 4: send packets through tunnel 4, whose destination is in data center Y

To limit bandwidth, the following policers are applied:

- policer ‘i1-to-dc-X’ to packets sent to data center X from interface 1
- policer ‘i1-to-dc-Y’ to packets sent to data center Y from interface 1
- policer ‘i2-to-dc-Y’ to packets sent to data center Y from interface 2

This is achieved with policy-tags as follows:

- Next-hops 1 and 2 are marked with policy-tag ‘DC-X’, and next-hops 3 and 4 with policy-tag ‘DC-Y’.
- 3 policers are configured: ‘i1-to-dc-X’, ‘i1-to-dc-Y’, and ‘i2-to-dc-Y’.
- Policy-forwarding policies are configured and applied to interfaces:

policy_if1 applied to interface i1
if tag DC-X, action ‘policer-i1-to-dc-X’
if tag DC-Y, action ‘policer-i1-to-dc-Y’

policy_if2 applied to interface i2
if tag DC-Y, action ‘policer-i2-to-dc-Y’

As a result:

- packets resolved through next-hops 1 and 2 are policed by policer ‘i1-to-dc-X’
- packets resolved through next-hop 3 are policed by policer ‘i1-to-dc-Y’
- packets resolved through next-hop 4 are policed by policer ‘i2-to-dc-Y’

### Design Rationale

Policy-forwarding rules match a policy-tag rather than matching next-hops directly. If rules matched next-hops directly, every addition or deletion of a next-hop would require a policy update so that the actions remain applied to the correct next-hops. The policy-tag provides a level of indirection which allows next-hops to be updated dynamically while the policy configuration remains static.

Actions are applied through policy-forwarding policies rather than being specified directly in each next-hop (for example, as a policer name). This provides more flexibility: a policy-forwarding rule can match additional properties, such as the ingress interface, and apply different actions based on them. A policy-forwarding policy can also specify actions other than applying a policer.

## OpenConfig Model

This specification defines 2 OpenConfig model elements:

1. a policy-tag leaf on gRIBI (AFT) next-hops
2. a next-hop policy-tag match criterion in policy-forwarding rules

### AFT Next-hop Policy-tag

A `policy-tag` leaf is defined under `afts/next-hops/next-hop/state`. This allows a policy-tag to be specified for every next-hop. It’s an optional field, and if it’s not set, then no forwarding policy matching on policy-tag matches the next-hop.

#### OpenConfig AFT next-hops tree

The `policy-tag` leaf (`string`) is located under `/network-instances/network-instance/afts/next-hops/next-hop/state`:

```text
module: openconfig-network-instance
  +--rw network-instances
     +--rw network-instance* [name]
        +--ro afts
           +--ro next-hops
              +--ro next-hop* [index]
                 +--ro state
                    +--ro policy-tag?   string
```

#### YANG definition

In `openconfig-aft-common.yang`, `leaf policy-tag` is defined in `grouping aft-common-entry-nexthop-state`:

```yang
grouping aft-common-entry-nexthop-state {
...
    leaf policy-tag {
      type string;
      description
        "A policy tag applied to the next-hop. This tag can be used to
        group next-hops together, so that policies can be applied at
        the group level. In one use case, this allows for a group of next-hops
        to be associated with a policy-forwarding policy. The next-hops may
        change over time, but the policy definition can remain static.";
    }
...
}
```

#### Next-hop example

Here is canonical OpenConfig JSON for an example AFT `next-hop` entry with `policy-tag`:

```json
{
  "network-instances": {
    "network-instance": [
      {
        "name": "DEFAULT",
        "afts": {
          "next-hops": {
            "next-hop": [
              {
                "index": 1001,
                "state": {
                  "index": 1001,
                  "policy-tag": "clusterA",
                  // other fields specifying other aspects, such as the forwarding behavior, of the next-hop
                }
              }
            ]
          }
        }
      }
    ]
  }
}
```

### Policy-forwarding Next-hop Policy-tag Match

A `next-hop` container is defined under `policy-forwarding/policies/policy/rules/rule`, which allows next-hop properties to be matched. A `policy-tag` leaf is defined under `next-hop/config` and `next-hop/state`. This leaf allows a policy-forwarding rule to match the policy-tags of next-hops. If it matches, then the action of the policy-forwarding rule applies to the packets resolved through that next-hop.

#### OpenConfig policy-forwarding tree

The `next-hop` container is located under `policy-forwarding/policies/policy/rules/rule`. The `policy-tag` leaf (`string`) is located under `next-hop/config` and `next-hop/state`.

```text
module: openconfig-policy-forwarding
   grouping policy-forwarding-top:
     +--rw policy-forwarding
        +--rw config
        |  +--rw global-decap-policy?   string
        +--ro state
        |  +--ro global-decap-policy?   string
        +--rw policies
        |  +--rw policy* [policy-id]
        |     +--rw policy-id          -> ../config/policy-id
        |     +--rw policy-counters
        |     |  ...
        |     +--rw config
        |     |  ...
        |     +--ro state
        |     |  ...
        |     +--rw rules
        |        +--rw rule* [sequence-id]
        |           +--rw sequence-id    -> ../config/sequence-id
        |           +--rw config
        |           |  +--rw sequence-id?      uint32
        |           |  +--rw address-family?   identityref
        |           |  +--rw description?      string
        |           +--ro state
        |           |  +--ro sequence-id?      uint32
        |           |  +--ro address-family?   identityref
        |           |  +--ro description?      string
        |           |  +--ro matched-pkts?     oc-yang:counter64
        |           |  +--ro matched-octets?   oc-yang:counter64
        |           +--rw l2
        |           |  ...
        |           +--rw ipv4
        |           |  ...
        |           +--rw ipv6
        |           |  ...
        |           +--rw transport
        |           |  ...
        |           +--rw next-hop
        |           |  +--rw config
        |           |  |  +--rw policy-tag?   string
        |           |  +--ro state
        |           |     +--ro policy-tag?   string
        |           +--rw action
        |              ...
        +--rw interfaces
        |  ...
        +--rw path-selection-groups
           ...
```

#### YANG definition

In `openconfig-pf-forwarding-policies.yang`, the `next-hop` container contains a leaf matching policy-tag:

```yang
submodule openconfig-pf-forwarding-policies {
...
 container rules {
   list rule {
...
     uses pf-forwarding-policy-rule-next-hop;
...
   }
 }
...
 grouping pf-forwarding-policy-rule-next-hop {
    description
      "Forwarding policy rule match parameters related to next-hops.";
    container next-hop {
      description
        "Container for forwarding policy match criteria related to
        next-hops.";
      container config {
        description
          "Configuration parameters for next-hop match criteria.";
        uses pf-forwarding-policy-rule-next-hop-config;
      }
      container state {
        config false;
        description
          "Operational state parameters for next-hop match criteria.";
        uses pf-forwarding-policy-rule-next-hop-config;
      }
    }
  }
  grouping pf-forwarding-policy-rule-next-hop-config {
    description
      "Configuration parameters for next-hop match criteria.";
    leaf policy-tag {
      type string;
      description
        "Match an AFT next-hop policy-tag
        (afts/next-hops/next-hop/state/policy-tag) in the same
        network-instance. All packets with next-hops with a matching policy-tag
        have the policy-forwarding action applied.";
    }
  }
...
}
```

## TODO

We need to create more examples of policies and consider possible issues
regarding order of operations, rules matching and conflicts.  (For context,
see doc/ocpf-match-action.md). Discussion is required to resolve how 
configuration of next-hop policy-tag fits into policy-forwarding overall.
