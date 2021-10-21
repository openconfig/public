# wifi-models
Development of configuration and operational state models for WiFi.

Open Config is a collaboration between network operators to develop
and standardize a set of vendor-neutral network configuration models.
Models are written using the YANG data modeling language
(IETF RFC 6020).

The configuration models in this repository are made available under
the Apache 2.0 license (see the LICENSE file).

# Working with YANG models
The recommended tool for "compiling" and manipulating YANG models is
the open source pyang tool.

# WiFi Model Flow
The following models will be used for the various components that make
up a WiFi Network.

## openconfig-wifi-mac.yang
All MAC layer configuration/state.

## openconfig-wifi-phy.yang
 All PHY layer configuration/state.

## openconfig-access-points.yang
Top level WiFi model for a list of Access Points. This mostly imports
the MAC, PHY, and other models which are applicable to aggregating
below the top-most network element. In the case of WiFi, this top-most
network element is the Access Point.

## openconfig-ap-manager.yang
 Top level configuration and state data for a system which manages
 Access Points. Note, all models are agnostic to which system is
 responsible for data-plane functions. See below for more info.

## openconfig-wifi-types.yang
Types definition file for wifi-specific types.


# OpenConfig WiFi Model Summary
In addition to the above models, it is expected that some WiFi
networks utilize network elements (such as Wireless LAN Controllers)
with components that share common config/state leafs with
Routers/Switches. In such cases, vendors will be expected to support
existing models, already published. For example, SFP's and there
related config/state, will utilize openconfig-platform-transceiver.yang.

Furthermore, OpenConfig WiFi models are agnostic to which system is
performing 802.11 to 802.3 data-plane functions. These architectures
are sometimes referred to as Controller Vs. Controller-less. The
semantics of what a Controller is, remains out-of-scope of these
models. Whether the network vendor utilizes a Wireless LAN Controller,
Cloud controller, or simply an NMS (referring to themselves as
  'controller-less') most WiFi vendors require some sort of
  management-plane and/or control plane system to aggregate
  configuration and state of the deployed APs. Remaining architecture
  agnostic is driven from the point of view that network engineers
  deal with designing, deploying, and monitoring WiFi networks from
  the Access Point's perspective. AP name, which may or may not be
  FQDN, being the unique identifier makes the most sense to the most
  verticals.

An oversimplified workflow is as follows:
* Day 0: APs are shipped on-site or to a build room.
* Day 1: APs are provisioned, using openconfig-ap-manager.yang. The
config container 'provision-aps' is used for initial assignment of the
friendly ap-name to the factory issued MAC address. This is also where
country-code is assigned.
* Day 2+: The remaining configuration and state is done almost entirely
through the openconfig-access-points.yang model.

## BSSID Telemetry
Since the SSID of a particular radio may or may not be dual-band,
and since the BSSIDs of a particular radio may not be known to
 the operator, there must be a method for 'discovery'. As such, the
 "bssids" container is a list with multiple keys. This allows the
 operator to utilize paths in their GetRequest's, providing the
 flexibility to GET/Subscribe to Telemetry for a particular BSSID,
 group of BSSID's, all BSSID's of a certain radio, or simply discover
 which BSSID's a radio is broadcasting. See the following examples,
 using gNMI:
1)
The following will return value of "num-associated-clients" for every
Radio (regardless of radio-id) which is broadcasting BSSID
00:11:22:33:44:55. If there were TWO radio's, which both have BSSID 00:11:22:33:44:55 on them, then the JSON being returned would include
the 'num-associated-clients' of BOTH.
```
GetRequest "/access-points:access-points/access-point/ssids/ssid[name="SSID-1"]/bssids/bssid[bssid=00:11:22:33:44:55]/state/num-associated-clients"
```
Note, this GetRequest does not include the key for radio-id, which is
the same as "*" (See [gNMI Specification](https://github.com/openconfig/reference/blob/master/rpc/gnmi/gnmi-path-conventions.md#paths-referencing-list-elements) for details).

2)
If we were ONLY interested in the telemetry of ONE of the BSSIDs on
radio-id 0 (even if it is broadcast exists on two radio-id's) we would
simply do:
```
GetRequest "/access-points:access-points/access-point/ssids/ssid[name="SSID-1"]/bssids/bssid[bssid=00:11:22:33:44:55]radio-id[radio-id=0]/state/num-associated-clients`
```
3)
If we wanted to "discover" what BSSID's are being broadcast by a
particular AP, for a particular SSID ("SSID-1") regardless of Radio ID
, we would do:
```
GetRequest "/access-points:access-points/access-point/ssids/ssid[name="SSID-1"]/bssids/bssid[bssid=*]radio-id[radio-id=*]/state/bssid`
```
This could be useful if you want to know the 5GHz BSSID of a dual-band SSID,
which you could then subsequently utilize to GET/Subscribe to only Telemetry
for that BSSID.

4)
Similar to the above, if we wanted to "discover" what BSSID's are
being broadcast but only for a particular Radio ("1"), we would do:
```
GetRequest "/access-points:access-points/access-point/ssids/ssid[name="SSID-1"]/bssids/bssid[bssid=*]radio-id[radio-id=1]/state/bssid`
```
This provides flexiility to GET/Subscribe to telemetry only for 5GHz
BSSID's, 2.4GHz BSSID's, or both, because we (the operator) know which
radio ID's we assign to which radio's.
## Radio list
The available radios on an access point are modeled as a list, keyed
by both an operator-assigned id and a frequency, in order to support
software-selectable and fixed-frequency radios. Radios on an AP are
typically associated with a physical slot id determined by the device.
Rather than having the operator "discover" which radio type is
supported by which slot on a given vendor platform, the device is
expected to maintain an internal mapping of the configured id to the
slot id.

Having a multi keyed list allows the implementer (vendor) to
internally map the OpenConfig modeled radio list entry to whichever
internal "Slot-ID" required, to make the configuration valid. Since
the operator can not change the operating-frequency of a radio,
without deleting/updating the list entries the implementor(vendor) can
release any lock on the hardware resource (eg Slot-ID) necessary to
make the desired configuration valid.

It's still possible for an operator to specify an invalid
configuration of radios, for example 2x 5GHz radios where only one
5GHz radio exists, and that is handled via standard error response
codes, like any other invalid configuration SetRequest.
