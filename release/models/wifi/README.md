# wifi-models
Development of configuration and operational state models for WiFi.

Open Config is a collaboration between network operators to develop and 
standardize a set of vendor-neutral network configuration models. Models are 
written using the YANG data modeling language (IETF RFC 6020).

The configuration models in this repository are made available under the Apache
 2.0 license (see the LICENSE file). 

# Working with YANG models
The recommended tool for "compiling" and manipulating YANG models is the open 
source pyang tool. 

# WiFi Model Flow
The following models will be used for the various components that make up a WiFi
 Network.

## openconfig-wifi-mac.yang
 All MAC layer configuration/state.

## openconfig-wifi-phy.yang
 All PHY layer configuration/state.

## openconfig-access-points.yang
 Top level WiFi model for a list of Access Points. This mostly imports the MAC 
& PHY models, and includes additional State and Config container(s) which would 
not be appropriate to fall within a particular SSID or Radio, such as BSSID 
counters, System configuration, and assigned AP manager(s).

## openconfig-ap-manager.yang
 Top level configuration and state data for a system which manages Access 
Points. Note, all models are agnostic to which system is responsible for 
data-plane functions. See below for more info.

## openconfig-wifi-types.yang
Types definition file for wifi-specific types.


# OpenConfig WiFi Model Summary
In addition to the above models, it is expected that some WiFi networks utilize 
network elements (such as Wireless LAN Controllers) with components that share 
common config/state leafs with Routers/Switches. In such cases, vendors will be 
expected to support existing models, already published. 
For example, SFP's and there related config/state, will utilize 
openconfig-platform-transceiver.yang.

Furthermore, OpenConfig WiFi models are agnostic to which system is performing 
802.11 to 802.3 data-plane functions. These architectures are sometimes 
referred to as Controller Vs. Controller-less. The semantics of what a 
Controller is, remains out-of-scope of these models. Whether the network vendor 
utilizes a Wireless LAN Controller, Cloud controller, or simply an NMS 
(referring to themselves as 'controller-less') most WiFi vendors require some 
sort of management-plane and/or control plane system to aggregate configuration 
and state of the deployed APs.
Remaining architecture agnostic is driven from the point of view that network 
engineers deal with designing, deploying, and monitoring WiFi networks from the 
Access Point's perspective. AP name, which may or may not be FQDN, being the 
unique identifier makes the most sense to the most verticals.

An oversimplified workflow is as follows:
* Day 0: APs are shipped on-site or to a build room.
* Day 1: APs are provisioned, using openconfig-ap-manager.yang. The config 
container 'provision-aps' is used for initial assignment of the friendly 
ap-name to the factory issued MAC address. This is also where country-code is 
assigned.
* Day 2-n: The remaining configuration and state is done almost entirely 
through the openconfig-access-points.yang model.
