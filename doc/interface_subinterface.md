# Representing Parent and Logical "Default" Child Interfaces in OpenConfig
**Contributors**: robjs<sup>†</sup>, ayhchan<sup>†</sup>  
<sup>†</sup> @google.com
**April 2019**

## Problem Statement
The OpenConfig interfaces model provides two constructs for interfaces:

* Interfaces - which represent individual physical or logical interfaces which
  may or may not be subdivided on a system. A physical interface in OpenConfig
  *cannot* have a Layer 3 address (i.e., IPv4, or IPv6). For example, the
  following are all interfaces:
   - an Ethernet port,
   - a logical interface which a device, which is switching a VLAN, uses to
     act as a host within the same VLAN (typically referred to as an SVI).
   - a loopback.

* Subinterfaces are logical subdivisions of an OpenConfig interface - which may
  or may not represent a difference in encapsulation. For example, the following
  are subinterfaces:
   - A logical 'unit' that has IPv4 or IPv6 address configuration on a physical
     interface.
   - A logical construct on physical Ethernet port that matches a particular
     802.1Q VLAN to configure an IPv4 or IPv6 address for the local system.
   - A logical construct on a physical Ethernet port that matches one or more
     VLANs to be treated as within a particular L2 service (VSI, VPLS etc.).

Whilst some vendor systems share this same division between the physical port
and logical subdivisions of it (e.g., SROS has a `port` and further `SAP`
entities associated with the port, JUNOS has a `interface` which has a series
of `unit` entities), other systems (e.g., IOS, IOS XR) do not have such a
division. This note seeks to clarify the expected behaviour for both sets of
systems.

## Expected Behaviour

### Systems with Explicit Subinterfaces

For systems with explicit subinterfaces, each OpenConfig subinterface should
map to an underlying subinterface. Configuration applied to the OpenConfig
`interface` should be applied to the entity representing the port - whilst
`subinterface` configuration should be mapped to the subinterface referenced by
the OpenConfig input.

Operational state for the parent 'interface' should:

- Represent the sum of all subinterface counters -- i.e., be the total traffic
  on a particular interface across all subinterfaces.
- Represent the state of the physical or logical interface - whether or not it
  is subdivided.

### Systems with a 'Default' Subinterface

For systems that allow configuration directly onto a parent interface that
cannot be configured directly in OpenConfig (e.g., IP address information) -
this configuration should be mapped to a subinterface with index 0. This
subinterface should have the encapsulation of the 'default' subinterface
specified if this is relevant - in most cases, the encapsulation is
expected to be null.

In this case, operational state of the parent `interface` should:

 - Represent the sum of all subinterfaces -- i.e., again be the total traffic
   on a parent interface across all subinterfaces.
 - Represent the state of the physical or logical interface.

The operational state for subinterface 0 should:
 
 - Represent the total of traffic matched by that subinterface. In the case
   that there are no subinterfaces other than index 0 specified, then this will
   represent the total traffic on the interface. In the case that there are other
   subinterfaces, it should match ONLY the traffic that is matched by the default
   subinterface. For example, if there is a default subinterface that has no VLAN
   encapsulation, but a subsequent subinterface with encapsulation of VLAN 42,
   then subinterface 0 should only show the traffic that matched VLANs other than
   42.
 - Represent state of the logical subinterface only, i.e., for example there
   may be separate `oper-status` between subinterfaces.

For systems that have this default mapping of subinterface 0 to a parent
interface, the configuration of separate IPv4 and IPv6 subinterfaces with no
specified encapsulations is invalid - since in this case there is ambiguity as
to which OpenConfig `subinterface` entities map to the "default" subinterface's
config on the device. (i.e., OpenConfig `subinterface` 4 and 6 configured on
`GigabitEthernet42` would potentially both map into configuration and state
directly under `GigabitEthernet` in the native schema). Conversely, a single
subinterface that specifies both IPv4 and IPv6 addresses is valid, since this
unambiguously maps to the "default" subinterface on these systems.
