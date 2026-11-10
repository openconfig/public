# Next-hop Policy-tag Forwarding Policy

**Contributors**: henryhu@google.com, dloher@google.com
Date: Created Nov 10, 2026

## Summary
 
This document describes the proposal for adding next-hop-tag-based forwarding policy.

## Background

OpenConfig policy-forwarding matches and applies actions to packets. Policy-forwarding rules match on various properties of the traffic, such as destination IP, TTL, DSCP, protocol and apply actions, such as drop the packet, set TTL or set the next-hop.  These rules require inspecting fields inside the packet.

gRIBI can be used to inject routes and next-hops. In addition to the traditional next-hop attributes like destination IP address, additional attributes can be added such defining encapsulation headers.

In certain cases, it’s desirable to match gRIBI next-hops in policy-forwarding policies. This allows policy-forwarding actions to be applied to traffic sent to those next-hops. In this scenario, the policy-forwarding match rule is evaluated against next-hop table entries instead of matching on fields inside the packet.

## Next-hop Policy-tag Specification

This specification allows policy-tags to be attached to gRIBI next-hops, and enhances forwarding policies to match next-hops with such policy-tags.

When packets are forwarded using a route that resolves to a gRIBI-installed next-hop and that next-hop contains a policy-tag, such packets become subject to the actions of a forwarding policy with a match rule for the same policy-tag.

For traffic forwarded through a gRIBI next-hop, if it has a policy-tag attached, and if we have a forwarding policy which matches that policy-tag, then the forwarding actions specified in that forwarding policy are applied to the traffic.

### Details

To use this new feature, 2 steps need to be done:

  * First, attach policy-tags to next-hops programmed by gRIBI.
  * Second, configure policy-forwarding policies which match these gRIBI next-hop policy-tags, and specify actions for the policies.

For attaching policy-tags to next-hops, these tags are programmed together with next-hops through gRIBI. The OpenConfig model is extended with a policy-tag leaf in the next-hop container, and the policy-tag is set with this new leaf.

For configuring forwarding policies, the existing forwarding policies in OpenConfig are enhanced to match these next-hop policy-tags.  Only one policy may exist for a given policy-tag.  The maximum number of policy-tags and the number of type of policy-forwarding actions are platform dependent.

When the next-hops are programmed through gRIBI and the forwarding policies matching their policy-tags are configured, the traffic forwarded through such next-hops has the forwarding action specified in the forwarding policies applied to it. For example, a forwarding policy matches a next-hop policy-tag and specifies a traffic policer to be applied. All traffic forwarded through next-hops with that policy-tag is policed by that policer.

### Use Cases

In one use case, gRIBI routes are programmed on a device which directs traffic with different destination IPs to be encapsulated into different tunnels.  We want to apply a common bandwidth limit for the tunnels sharing the same destination, because such tunnels typically share the same path across the backbone network.

One option is adding forwarding policies which match specific destination prefixes of the traffic and apply policers. The problem with this approach is that it requires separate matching of destination prefixes for the forwarding policies, which needs additional hardware resources on the device. Specifically, the forwarding policies are often implemented with TCAM lookups, while the routes are often implemented with LPM lookups. Separate lookups mean additional hardware resource consumption. In addition, the scale achievable with TCAM is often lower than the scale achievable with LPM.

Using next-hop policy-tag, the lookup is done just once in LPM, and the policy-tag is part of the LPM lookup result. The policy-tag identifies a next-hop-tag-based forwarding policy, which then specifies a policer to be applied. This avoids the double lookups needed in the previous approach.

### Example

Suppose we have 2 interfaces. For interface I1, we have 3 gRIBI routes:

- Route 1: match prefix A, forward through next hop 1: send packets through tunnel 1, whose destination is in data center X
- Route 2: match prefix B, forward through next hop 2: send packers through tunnel 2, whose destination is in data center X
- Route 3: match prefix C, forward through next hop 3: send packets through tunnel 3, whose destination is in data center Y

For interface I2, we have 1 route:

- Route 4: match prefix D, forward through next hop 4: send packets through tunnel 4, whose destination is in data center Y

Goal: for bandwidth enforcement, we want to apply

- policer ‘i1-to-dc-X’ to traffic sent to data center X from interface 1
- policer ‘i1-to-dc-Y’ to traffic sent to data center Y from interface 1
- policer ‘i2-to-dc-Y’ to traffic sent to data center Y from interface 2

With the policy-tag approach, we

- first, mark next-hops 1 and 2 with policy-tag ‘DC-X’, and next-hops 3 and 4 with policy-tag ‘DC-Y’
- then, configure 3 policers, i1-to-dc-X, i1-to-dc-Y, and i2-to-dc-Y
- finally, configure forwarding policies which
  - matches interface I1 and policy-tag ‘DC-X’, apply policer ‘i1-to-dc-X’
  - matches interface I1 and policy-tag ‘DC-Y’, apply policer ‘i1-to-dc-Y’
  - matches interface I2 and policy-tag ‘DC-Y’, apply policer ‘i1-to-dc-Y’

This allows

- traffic resolved through next-hops 1 and 2 to be policed by policer ‘i1-to-cluster-X’
- traffic resolved through next-hop 3 to be policed by policer ‘i1-to-dc-Y’
- traffic resolved through next-hop 4 to be policed by policer ‘i2-to-dc-Y’

satisfying the requirements.

### Alternative Design Options

There are alternative ways to implement such forwarding policies. One idea is to have the policies match next-hops directly. With that approach, every change to the next-hops, whether an addition or a deletion, requires the policy to be updated to add or remove the next-hop, so that the policers are still applied to the correct next-hops. By having a level of indirection, the policy-tag allows the next-hops to be updated dynamically, while the policy can be configured in a static way.

Another idea is specifying the policer name directly in every next-hop. However, setting the policer through forwarding policies provides more flexibility. For example, the forwarding policy can match additional properties, such as ingressing interface, and applying different policers based on that. Also, forwarding policy can have other forwarding actions other than applying a policer.

## OpenConfig Changes

This proposal requires 2 parts of OpenConfig changes:

1. adding policy-tag to gRIBI next-hop
2. modeling next-hop-tag-based forwarding policies

### Add Policy-tag to gRIBI next-hop

A `policy-tag` leaf is added under `afts/next-hops/next-hop/state`. This allows a policy-tag to be specified for every next-hop. It’s an optional field, and if it’s not set, then no forwarding policy matching on policy-tag matches the next-hop.

#### OpenConfig AFT next-hops tree

The `policy-tag` leaf (`string`) is added under `/network-instances/network-instance/afts/next-hops/next-hop/state`:

```text
module: openconfig-network-instance
  +--rw network-instances
     +--rw network-instance* [name]
        +--ro afts
           +--ro next-hops
              +--ro next-hop* [index]
                 +--ro state
+                   +--ro policy-tag?   string
```

#### YANG definition

In [`openconfig-aft-common.yang`](http://google3/third_party/openconfig/vendor/nokia/yang/gx/v0_1/openconfig-aft-common.yang), `leaf policy-tag` is added to `grouping aft-common-entry-nexthop-state`:

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

### Modeling Next-hop Policy-tag Forwarding Policies

A `next-hop` container is added under `policy-forwarding/policies/policy/rules/rule`, which allows next-hop properties to be matched. A `policy-tag` leaf is added under `next-hop/config` and `next-hop/state`. This leaf allows a forwarding policy to match the policy-tags configured on next-hops. If it matches, then the forwarding action of this forwarding policy applies to the traffic resolved through that next-hop.

#### OpenConfig AFT next-hops tree

The `next-hop` container is added under `policy-forwarding/policies/policy/rules/rule`. The `policy-tag` leaf (`string`) is added under `next-hop/config` and `next-hop/state`.

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
+       |           +--rw next-hop
+       |           |  +--rw config
+       |           |  |  +--rw policy-tag?   string
+       |           |  +--ro state
+       |           |     +--ro policy-tag?   string
        |           +--rw action
        |              ...
        +--rw interfaces
        |  ...
        +--rw path-selection-groups
           ...
```

#### YANG definition

In `openconfig-pf-forwarding-policies.yang`, next-hop container is added, which contains a leaf matching policy-tag:

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
        have the policy-forwarding action will be applied.";
    }
  }
...
}
```

