# wifi-models
Development of configuration and operational state models for WiFi

Open Config is a collaboration between network operators to develop and standardize a set of vendor-neutral network configuration models. Models are written using the YANG data modeling language (IETF RFC 6020).

The configuration models in this repository are made available under the Apache 2.0 license (see the LICENSE file). Since the models are intended to be published in the IETF, participants must be willing to adhere to the IETF Note Well statement as well as BCP 78 and BCP 79.

# Working with YANG models
The recommended tool for "compiling" and manipulating YANG models is the open source pyang tool. For example to see a text-based tree view of a yang module, run the following command:

pyang -f tree --ietf --strict <module.yang>

Note that the --ietf and --strict flags should be used to ensure compatibility with IETF module guidelines, and strict YANG compliance, respectively. Additional OpenConfig validation can be achieved through using the OpenConfig linter.

# WiFi Model Flow
The following layout will be used to model the various components that make up a WiFi Network. In general if there are >2 leafs related to one feature, those are placed in their own container. See openconfig-wifi-mac.yang:wmm/dot11r/dot11v etc for example

## openconfig-wifi-mac.yang
 All MAC layer configuration/state.

## openconfig-wifi-phy.yang
 All PHY layer configuration/state.

## openconfig-system-wifi-ext.yang
Augmentation model of openconfig-system.yang
Used for non-802.11 config/state, for example:
 NTP, clock, hostname, Joined APs, etc.

## openconfig-wifi-types.yang
Types definition file for wifi-specific types.

In addition to the above models, it is expected that some WiFi networks utilize network elements (such as Wireless Lan Controllers) with components that share common config/state leafs with Routers/Switches. In such cases, vendors will be expected to support existing models, already published. 
For example:
SFP's and there related config/state, will utilize openconfig-platform-transceiver.yang.
