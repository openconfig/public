# Openconfig EVPN Support

**Contributors:** Oscar Gonzalez de Dios, Samier Barguil, Mike Wiebe

This page documents how the Openconfig Yang data models can be used to
implement different Ethernet VPN (EVPN) solutions. Current
implementation covers the following use cases: 
*  BGP MPLS-Based
Ethernet VPNs [RFC 7432](https://datatracker.ietf.org/doc/html/rfc7432)
with VLAN based service. 
*  Provider Backbone Bridging Combined with
Ethernet VPN (PBB-EVPN) [RFC
7263](https://datatracker.ietf.org/doc/html/rfc7263) with VLAN based
service. 
* Network Virtualization Overlay (NVO) EVPN [RFC
8365](https://datatracker.ietf.org/doc/html/rfc8365) with VLAN based
service and symmetric IRB.

# EVPN concepts in Openconfig 
## MAC VRF
A MAC-VRF (A Virtual Routing
and Forwarding table for Media Access Control (MAC) addresses) on a
Provider Edge (PE) device with EVPN support is instantiated by creating
a network instance with the 'type = oc-ni-types:L2VSI' in combination
with the parameters under the '/network-instances/network-instance/evpn'
container.

The presence of the evpn container during the netconf/gRPC request is
the indication for the device of the EVPN support in that VRF.

## EVPN Instance
A PE device can participate in one or more EVPN
Instances (EVI). An EVPN Instance spans several PE devices. For each EVI
in which the PE participates, a new entry of the evpn-instance container
within evpn needs to be created. This enables multi-domain and Data
center interconnect use cases. 

For each EVPN Instance the following items can be configured: 
* **Unique EVPN Instance Identifier** (EVI-ID).  
* **Data plane encapsulation** (MPLS or NVO (VXLAN) EVPN): 
  * In each EVPN instance, the type of encapsulation of
  the traffic between PEs can be selected. This is performed via the
  ‘network-instance/evpn/evpn-instances/evpn-instance/config/encapsulation-type’
  container. Note that  Currently “oc-ni-types:MPLS” and
  “oc-ni-types:VXLAN” can be selected.  
  * Note that, in the EVPN case, the
  ´network-instance/config/encapsulation-type’ **has to be ignored**.
  Hence, it allows the creation of two EVPN Instances one EVI with VXLAN
  and another one with MPLS for a data center interconnect use case. 
* **Service type**: EVPN defines three different Service Types (VLAN BASED,
VLAN BUNDLE and VLAN AWARE). In the current version of EVPN Model in
Openconfig, **VLAN_BASED is supported**. 
* **Multicast group and multicast mask**:
Multicast group address and mask for BUM traffic. 
* **Replication mode**.
The PE may use ingress replication for flooding BUM (Broadcast, Unicast
or Multicast) traffic.  
* **Route Distinguisher (RD) associated to the
evpn-instance**. An RD MUST be assigned for a given evpn-instance on a PE.
This RD MUST be unique across all evpn-instances on a PE. The RD is
configured using the parameter 'route-distinguisher' at this level. **It
overrides the route-distinguisher value defined under
'network-instance/config'**. 
* **Import and Export Route Target (RT)**. In
order to participate in an EVPN Instance, it is neede to indicate, in
the import-export-policy container, the values of the Route Target for
import and export.  
* **VXLAN**. In case VXLAN is used as the data plane
encapsulation, use the 'vxlan' container to configure the VNI and
overlay end point reference. 
* **PBB**. In case Provider Backbone Bridging
(PBB) is combined with Ethernet VPN, use the pbb container to configure
the b (backbone) and i components.

Example of configuration of an EVPN container for a VLAN based MPLS EVPN
with explicit RD assignment. 

``` 
{
  "evpn-instances": {
    "evpn-instance": {
      "config": {
        "evi": "1001",
        "encapsulation-type": "oc-ni-types:MPLS",
        "service-type": "oc-evpn-types:VLAN_BASED",
        "route-distinguisher": "65000:2"
      }
    }
  }
}
``` 

## Ethernet Segment
A customer site is connected to one or more PEs via
a set of Ethernet links. That set of links is referred to as an
'Ethernet segment' (see RFC 7432 section 5 ) . The Ethernet segments
that are associated to the VRF are listed in the 
‘/network-instances/network-instance/evpn/ethernet-segments’ container.
Each ethernet segment needs the following information: 
* **Name**: Uniquely
identifies the ethernet segment. It has only local meaning. 
* **ESI type**.
[RFC 7432](https://datatracker.ietf.org/doc/html/rfc7432#section-5)
defines several types of identifiers, with explicit rules for the ESI
value assignment. Use this field to select the desired type (see table
bellow for the values) 
* **ESI**: The ethernet segment identifier can be
auto assigned or explicitly configured. In case of  autoasign, choose
the autotype, otherwise, indicate the identifier.
 
  * The following table summarized the ESI types and possible values:

| [ESI Type](https://datatracker.ietf.org/doc/html/rfc7432#section-5) | Typedef value| ESI|
|----------|------------------------------|-----------------------------------------------|
| Type 0   | TYPE_0_OPERATOR_CONFIGURED   | Directly configured by the operator           |
| Type 1   | TYPE_1_LACP_BASED            | AUTO enum must be used                        |
| Type 2   | TYPE_2_BRIDGE_PROTOCOL_BASED | AUTO enum must be used                        |
| Type 3   | TYPE_3_MAC_BASED             | Directly configured or AUTO enum must be used |
| Type 4   | TYPE_4_ROUTER_ID_BASED       | Directly configured or AUTO enum must be used |
| Type 5   | TYPE_5_AS_BASED              | Directly configured or AUTO enum must be used |


* **Redundancy mode**. For load balancing purposes, in case of multi-homing,
the allowed possibilities are Single Active or All Active (see
https://datatracker.ietf.org/doc/html/rfc7432#section-14.1) 
* **Designated Forwarder (DF) election**. In the multi-homing cases and to manage BUM
(Broadcast, Unicast or Multicast) traffic, it is possible to configure
the method and algorithm to choose if the node is the designated
forwarded in an EVPN Instance. 
  * **df-election-method**. The possible
  methods, according to [RFC
  8584](https://datatracker.ietf.org/doc/html/rfc8584) are DEFAULT,
  HIGHEST RANDOM and PREFERENCE. 
  * **preference**. Only applies in case
  PREFERENCE Method is selected. Use this field to indicate the value of
  the preference.  
  * **revertive**. Only applies in case PREFERENCE Method is
  selected. This boolean indicates if there is revertion after the
  recovery of the ethernet segment. By default is TRUE, meaning it is
  revertive. 
  * **election-wait time**. Only applies in case PREFERENCE Method
  is selected. for the configured ethernet segment, when the DF timer with
  this time value expires, the device selects the DF with the highest
  preference value.

 Example of Ethernet segment with directly configured ESI by the
 operator of type 0

```  
{
  "ethernet-segments": {
    "ethernet-segment": {
      "name": "esi-1",
      "config": {
        "name": "esi-1",
        "esi-type ": "oc-evpn-types:TYPE_0_OPERATOR_CONFIGURED ",
        "esi": "01:00:00:00:00:71:00:00:00:01",
        "redundancy-mode": "oc-evpn-types:SINGLE_ACTIVE"
      }
    }
  }
} 
```

Example of Ethernet segment with auto-generated ESI using LACP 
``` 
{
  "ethernet-segments": {
    "ethernet-segment": {
      "name": "esi-1",
      "config": {
        "name": "esi-2",
        "esi-type ": "oc-evpn-types:TYPE_1_LACP_BASED ",
        "esi": "AUTO",
        "redundancy-mode": "oc-evpn-types:SINGLE_ACTIVE"
      }
    }
  }
}
``` 

Example of Ethernet segment with directly configured ESI by the
operator of type 0 and DF election method preference

```  
{
  "ethernet-segments": {
    "ethernet-segment": {
      "name": "esi-1",
      "config": {
        "name": "esi-1",
        "esi-type ": "oc-evpn-types:TYPE_0_OPERATOR_CONFIGURED",
        "esi": "01:00:00:00:00:71:00:00:00:01",
        "redundancy-mode": "oc-evpn-types:SINGLE_ACTIVE"
      },
      "df-election": {
        "config": {
          "df-election": "PREFERENCE",
          "preference": "1",
          "revertive": "true",
          "election-wait-time": "3"
        }
      }
    }
  }
}
```
