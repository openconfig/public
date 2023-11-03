# Route Redistribution in OpenConfig Network Instance
**Contributors:**: robjs@google.com


One difference between vendor devices that OpenConfig seeks to abstract is the expected manner in which the various device RIBs are populated.  For some vendors, the model essentially reflects a per-protocol RIB (i.e., the routes of protocol A are isolated from those of protocol B), while for others the routes are expected to populate the same RIB, such that routes installed into the RIB by protocol A are directly accessible by protocol B.

OpenConfig models must therefore define a vendor-neutral way to express how the routes installed by one protocol are visible to another protocol.  Such a convention must also be readily mapped into the underlying configuration of both styles of implementations.

## RIB Modeling Approach used in OpenConfig

OpenConfig provides a set of `tables` within a `network instance`. The network instance represents an entire virtual routing and forwarding instance (`VRF` or `routing-instance`). The `tables` then represent various per-protocol RIBs that can exist in a network instance.

In implementations that provide a per-protocol RIB construct, such tables are expected to be supported natively, with little mapping complexity. For those that provide a single RIB that multiple protocols install routes into, then the tables in the network instance are essentially virtual RIBs formed by filtering the entries that are in the common multi-protocol RIB.

## Interconnection of Protocol RIBs

In numerous use cases an operator wishes to take routes that exist in one protocol and advertise them within another protocol. Common examples include:

* redistributing aggregate prefixes that are locally generated into BGP
* redistributing static routes into BGP or an IGP
* redistributing IGP routes into BGP for inter-domain advertisement.

In those implementations that maintain a (virtual-) RIB per protocol, the operator must configure explicit connections between these tables, usually alongside a policy to allow such inter-protocol advertisement to occur. While no such configuration is required in those implementations that utilise a single RIB, it is notable that even in these implementations a protocol will not advertise another protocol's routes from the RIB by default (since such a setup would mean that for instance, the BGP DFZ would be re-advertised into the IGP at peering edge devices, which is clearly not desirable). Instead, an operator must create explicit configuration that matches routes installed into the RIB by a particular protocol, and then "import" these routes into the protocol that is expected to advertise them.

In this latter case, the use of an import policy within a particular protocol instance is equivalent to the former - with the only difference being the context in which the policy is defined.   With a table per protocol, an explicit configuration of redistributing routes from protocol A to protocol B is specified.  In the latter case, the target protocol (protocol B) is specified, and a policy is used to select the source protocols from which routes may be imported for advertisement. 
In both cases, it is important that the destination protocol (protocol B) in this case only place the redistributed routes from the source protocol (protocol A in this case) in its "advertisement RIB view", but do not use those routes in the active/installed/forwarding route selection process and,consequently, do not place them in the FIB (forwarding information base).

## OpenConfig Approach to Protocol RIB Interconnection

OpenConfig models interconnection of RIBs through an explicit `table-connection` list within the `network-instance` model. In order to redistribute routes within protocol A to protocol B, an explicit connection between the tables corresponding to protocol A and protocol B is created.

An OpenConfig `routing-policy` is specified along with the connection - allowing an operator to apply policies to routes being imported into the destination protocol's RIB. Such policies are evaluated such that:

 * Routes that are rejected by the protocol are not imported into the destination RIB.
 * Modifications to attributes of matching routes result in the imported routes in the destination protocol's RIB being modified, but do not modify the source protocol's RIB.

It is expected that protocol-specific attributes (e.g., BGP communities) are set by such an import policy, allowing routes that are redistributed to carry information relating to their source (e.g., an IGP route may be tagged with a specific community using policy to indicate its provenance).

In the absense of an import-policy for table-connections, default-import-policy should take effect. In the absence of both, no routes should be allowed to be redistributed.

## Examples of OpenConfig Network instance

### Example: Aggregate Routes Redistributed to IS-IS Level 2

Consider the case that a router generates an aggregate default route (`0.0.0.0/0`), and advertises it into IS-IS Level 2.  Within an OpenConfig data instance, this configuration breaks down into a number of elements.

#### Protocol Configuration

Within a `network-instance` an instance of each protocol is created. In this example, the protocols required are `AGGREGATE` and `ISIS`. OpenConfig `network-instance` supports multiple instances of each protocol:

```json
"openconfig-network-instance:network-instances": {
    "network-instance": [
        {
            "config": {
                "name": "DEFAULT"
            },
            "name": "DEFAULT",
            "protocols": {
                "protocol": [
                    {
                        "local-aggregates": {
                            "aggregate": [
                                {
                                    "prefix": "0.0.0.0/0",
                                    "config": {
                                        "set-tag": 42,
                                        "prefix": "0.0.0.0/0"
                                    }
                                }
                            ]
                        },
                        "identifier": "LOCAL_AGGREGATE",
                        "config": {
                            "identifier": "openconfig-policy-types:LOCAL_AGGREGATE",
                            "name": "0"
                        },
                        "name": "0"
                    },
                    {
                        "identifier": "ISIS",
                        "config": {
                            "identifier": "openconfig-policy-types:ISIS",
                            "name": "15169"
                        },
                        "name": "15169"
                    }
                ]
            },
          }
    ]
}
```

The configuration above creates an instance of each protocol, and within the `LOCAL_AGGREGATE` protocol instance, generates the default route.

#### Instantiated tables

Since an instance of IS-IS and a local aggregate instance exist, tables for IPv4 routing within these protocols are created in the `network-instance/tables/table` list:

```json
"tables": {
    "table": [
        {
            "config": {
                "address-family": "openconfig-types:IPV4",
                "protocol": "ISIS"
            },
            "address-family": "IPV4",
            "protocol": "ISIS"
        },
        {
            "config": {
                "address-family": "openconfig-types:IPV4",
                "protocol": "LOCAL_AGGREGATE"
            },
            "address-family": "IPV4",
            "protocol": "LOCAL_AGGREGATE"
        }
    ]
}
```

These tables simply establish the RIBs that are required for IPv4. A device is expected to automatically create these tables (such that the NMS does not explicitly need to send this configuration).

#### Inter-table Connection

The connection between the `LOCAL_AGGREGATE` and `ISIS` table consists of two elements. The first is the policy, which informs the router which prefixes are to matched, and sets any relevant attributes that are required within IS-IS for the imported routes. The policy below permits only the default route (matched by a `prefix-set`) and sets the IS-IS level so that the prefix is imported to Level 2 only.

```json
"openconfig-routing-policy:routing-policy": {
    "policy-definitions": {
        "policy-definition": [
            {
                "statements": {
                    "statement": [
                        {
                            "conditions": {
                                "match-prefix-set": {
                                    "config": {
                                        "prefix-set": "DEFAULT"
                                    }
                                }
                            },
                            "config": {
                                "name": "TERM0"
                            },
                            "name": "TERM0",
                            "actions": {
                                "openconfig-isis-policy:isis-actions": {
                                    "config": {
                                        "set-level": 2
                                    }
                                }
                            }
                        }
                    ]
                },
                "config": {
                    "name": "DEFAULT-TO-ISIS-LEVEL2"
                },
                "name": "DEFAULT-TO-ISIS-LEVEL2"
            }
        ]
    },
    "defined-sets": {
        "prefix-sets": {
            "prefix-set": [
                {
                    "prefixes": {
                        "prefix": [
                            {
                                "masklength-range": "exact",
                                "ip-prefix": "0.0.0.0/0",
                                "config": {
                                    "masklength-range": "exact",
                                    "ip-prefix": "0.0.0.0/0"
                                }
                            }
                        ]
                    },
                    "prefix-set-name": "DEFAULT",
                    "config": {
                        "prefix-set-name": "DEFAULT"
                    }
                }
            ]
        }
    }
},
```

In order to create the connection, an explicit connection is created between the `LOCAL_AGGREGATE` `IPV4` table and the `ISIS` `IPV4` table:

```json
"table-connections": {
    "table-connection": [
        {
            "dst-protocol": "ISIS",
            "address-family": "IPV4",
            "config": {
                "dst-protocol": "ISIS",
                "address-family": "IPV4",
                "import-policy": [
                    "DEFAULT-TO-ISIS-LEVEL2"
                ],
                "src-protocol": "LOCAL_AGGREGATE"
            },
            "src-protocol": "LOCAL_AGGREGATE"
        }
    ]
}
```

In this case, the `DEFAULT-TO-ISIS-LEVEL2` policy is used for the connection - which performs two functions:

* Matches the default route.
* Sets the IS-IS level that `0.0.0.0/0` is imported into to level-2.

#### Mapping to Vendor Behaviors

Below, we show how this OpenConfig configuration can be mapped into several vendor configurations:

**IOS XE**

```
ip prefix-list DEFAULT-ONLY seq 5 permit 0.0.0.0/0
!
route-map DEFAULT-TO-ISIS-LEVEL2 permit 5
  match ip address prefix-list DEFAULT-ONLY
exit
!
router isis
  redistribute static level-2 route-map DEFAULT-TO-ISIS-LEVEL2
exit
```

**JUNOS**

In JUNOS, since the `export` policy within IS-IS doesn't have a context of the
static table that is being exported, the target protocol from the OpenConfig
`table-connection` is used to generate a wrapper policy that calls the original
(OpenConfig) policy, with an additional criteria specifying the source protocol.

```
policy-options {
  prefix-list DEFAULT-ONLY {
    0.0.0.0/0;
  }

  policy-statement DEFAULT-TO-ISIS-LEVEL2 {
    term TERM0 {
      from {
        prefix-list DEFAULT-ONLY;
      }
      to level2;
      then accept;
    }
  }

  policy-statement REDIST-STATIC-ISIS {
    term TERM0 {
      from {
        protocol aggregate;
        policy DEFAULT-TO-ISIS-LEVEL2;
      }
      then accept;
    }
  }

  isis {
    set export [REDIST-STATIC-ISIS];
  }
}
```

**IOS XR**
```
prefix-set DEFAULT-ONLY
 0.0.0.0/0
end-set

route-policy DEFAULT-TO-ISIS-LEVEL2
 if destination in DEFAULT-ONLY then
   pass
 else
   drop
 endif
end-policy

router isis
 redistribute static level-2 route-policy DEFAULT-TO-ISIS-LEVEL2
exit
```
