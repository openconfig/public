# Openconfig Terminal Device Manifest files guidelines

**Contributors:** Arturo Mayoral López de Lerma

This page documents the purpose and usability guidelines for the openconfig-terminal-device-properties.yang and openconfig-terminal-device-property-types.yang models included in openconfig/devices-manifest folder. These models enter in the category of devices manifest files, which is a new category of openconfig models which are intended to expose the static properties of a devices, which are traditionally documented by the device's manufacturers in the product datasheets. These properties are typically required by network planning teams to plan the network deployments.

# Motivation

Current optical transport networks are evolving into new degrees of openness and flexibility, where the optical terminal units become more and more independent from the rest of the DWDM line system, management and control. This means for instance that an optical terminal device, i.e., a transponder/muxponder device can be deployed as a standalone unit and connected to the line system of a third-party provider. In order, to make this type of integrations efficient, the optical channel signal characteristics can be exposed through the management interface (i.e., through OpenConfig models in this case), to allow the remote controller entity to ingest the required data for successful optical planning and impairment validation of the end-to-end transmission across the Optical Line System (OLS). This intends to remove ambiguity on how suppliers expose data required for network planning and physical impairment validation of end-to-end Optical Channels (OCh/OTSi) and foster openness across the optical transport industry.

Currently in OpenConfig, the optical channels characteristics are opaque to some extent and the model only offers an abstraction named 'operational-mode', which allows the manufactures to encode the different transmission modes supported by the device into a single integer value. While this may be sufficient in some cases (the manufacturer can offer the related mode details offline to its customers), and very practical for configuration purposes, it becomes cucumbersome for an network operator to encode this information into their network controller application, specially if the network contains many different vendor solutions which may expose their signal characteristics in different formats, adding or omitting some information in each case. Thus, this model intends to provide uniformity on the way operational-modes are characterized by including a set of static attributes associated to each mode.

This proposal was submitted by the Telecom Infra Project (TIP) Mandatory Use Cases for SDN Transport (MUST) sub-group. This is an operator centric initiative which intends to achieve a wider consensus about the SDN standards to use on the network transport segment.

# Model content
The model contains two main building blocks:
* **operational-mode-capabilities** this set of attributes contains all characteristic information of the signal (modulation format, FEC, bit rate...), relevant information for the physical impairment validation (OSNR Rx sensitivity, CD/PMD tolerance and penalties).
* **optical-channel-config-value-constrains** which contains the transmission configuration constrains/ranges of the optical-channel's attributes characterized by the operational-mode, i.e., the central frequency range, the frequency grid and the configurable transmitted power.


# Model logic

1. **Terminal device manufacturing – Operational modes hardcoding**
Terminal device manufacturer encodes supported transmission modes characteristic information into the terminal-device's manifest file which is hardcoded into the device firmware. This process shall takes into account the transmission modes supported third-party transceiver pluggable modules compatible with the terminal device.

2. **Terminal device w/o pluggable – Initial setup**
When the device is started, the operational modes list is complete and contains all the information about the compatible transceivers and their associated operational modes. The manifest file defined by the openconfig-terminal-device-properties.yang model is static and will remain available and invariant through the terminal device's management interface which exposes openconfig models.

3. **Terminal device with equipped pluggable – Running State**
Once the line-card module is successfully equipped with the fixed or pluggable transceiver and the module is discovered by the Terminal Device, the operational data store is updated with the actual modes available (see openconfig-terminal-device.yang module). The list of modes present in the terminal-device/operational-modes list represents the actual modes which are available in runtime by the device and which can be configured as part of the optical-channel configuration. Please note, that this information can be dynamically updated if the pluggable unit changes. For each operational-mode/mode-id present in the operational data store, there should be an operational-modes/mode-descriptor item with the same mode-id, included in the manifest file.

# Model implementation guidelines

Manifest files are a special OpenConfig model category since they do not represent operational data. Thus, this category of models will be included within the "models/tree/master/yang/devices-manifest" folder.

In order to keep separated them from the rest of operational models, the following openconfig extension is included in the model, to enrich the module metadata:
```
  oc-ext:origin "openconfig-properties";
```
# Extensions introduced in the v0.2.0 release of the model.

### Motivation

After the initial release of the [openconfig-terminal-device-properties.yang](https://github.com/openconfig/public/blob/master/release/models/devices-manifest/openconfig-terminal-device-properties.yang), there have been significant technical questions and discussions happening within the Telecom Infra Project (TIP) Open Optical & Packet Transport (OOPT) community between operators and vendors.

This issue summarizes the motivation and issues detected in the first release of the model and it will serve as an introduction and motivation of a new pull request (#911) with a new proposed comprehensive update of the model which will be accompanied by the relevant explanations on how the new model proposal will try to overcome the detected issues. It is worth mentioning that the current analysis and the new proposal are the outcomes of an extensive technical discussion within the OOPT community between vendors and service providers and that it consolidates an already discussed proposal starting from the issues and motivations explained here.

### Context

The current proposed terminal-device-properties model was designed with the objective of allowing the terminal devices' system vendors to expose the intrinsic properties (Modulation Format, FEC, Baud-rate) and performance characteristics (Rx-OSNR, CD/PMD limits) of the device's supported transmission modes.

The initial version of the model was designed as a flat list of mode properties, where each entry represents a mode supported by the terminal device and includes the list of characteristics of that mode. However, this initial version presents a significant list of limitations.

### Initial release limitations
-  **First issue**: The current model exposes a list of modes available in the device, however, the characteristics of a mode of transmission are affected by the HW transceiver supplying it. In other words, two different transceivers (supported by the same terminal device) might support the same mode, but their mode characteristics are different.

<img width="1062" alt="image" src="https://github.com/openconfig/public/assets/10529210/85d8f159-d753-434e-89e3-38a9b24cdd92">

- **Second issue**: There is not an interoperability matrix within the Terminal Device's which exposes the compatibility between Terminal Device's chassis, linecards, transceivers and modes. Right now there is no compatibility information available in the model, to allow the supplier to properly describe which modes are supported by each transceiver module available in the terminal device.

#### Operational challenges

- It is not clear how mode IDs will be assigned and who will assign them.
- Clarification of the bookended solution target by the model.

### New proposal scope and initial assumptions

Clarify the target of the next extension targets

1. Bookended solutions, and interoperability between terminal devices of the same system vendor.
2. Interoperability between different system vendors O-OTs through standard modes

### Solution proposed

This pull request covers a proposed solution to the issues described in #910.

The changes to the existing model **are not backward compatible.**

The summary of the changes proposed is the following:

**1. Operational-mode list:** 
  - The list of the exposed operational modes properties by the Terminal Device is augmented with the set of **CHARacteristic properties** of the operational mode. 
  - The mode-ids are the same used within the operational datastore of the terminal device (exposed by the [openconfig-terminal-device.yang](https://github.com/openconfig/public/blob/master/release/models/optical-transport/openconfig-terminal-device.yang) model) and have network-wide scope assuming they guarantee interoperability in bookended scenarios (two Terminal Devices of the same system vendor supporting the same mode).
  - The mode-ids are defined by the system vendor. 

**2. Mode-descriptors list:** 
 - The **design properties** of the modes, which are dependent on the transceiver component that implements the mode, are exposed as a nested list within the operational mode list. It is assumed that a single operational mode might be implemented by different transceivers which might have associated different design properties exposed by different mode descriptors.
 - The **mode-descriptor-id** is a local index that does not have interoperability meaning outside the specific Terminal device which reports it. In other words, the same mode descriptors might be exported by different Terminal Devices with different IDs.

**3. Interoperable mode list.**
 - A given proprietary operational mode might be capable to comply with a certain number of standards or elsewhere publicly defined operational modes defined by other organizations e.g., ITU-T or OpenROADM. 
 - This compatibility characteristic shall be assured by the system vendor and imply that the design properties of the implementations of that specific operational mode are a superset of all listed supported standard modes.
 - This model block will replace the previous "G.698.2" node tree with a more generic definition able to accommodate multiple standards or MSA-defined modes.

**4. Transceiver-descriptors list.**
 - The Terminal Device exposes the list of the transceiver components which are supported by the device.
 - Each transceiver exposes the list of compatible modes and their associated mode descriptor.
 - This new model block enhances the previous version by allowing to expose explicitly the compatibility matrix between the Terminal Device, the set of transceiver modules supported and the associated modes supported.

**5. Linecard-descriptors list.**
 - The Terminal Device exposes the list of linecard components which are supported by the device.
 - Each linecard component exposes the list of transceivers that are supported.
 - Each linecard constrains the list of modes that can be supported among the ones supported by the transceiver.
 - Each linecard constrains the optical-channel configuration, e.g., target-output-power and frequency range.
 - This new model block enhances the previous version by allowing to expose explicitly the compatibility matrix between the Terminal Device, the line cards, the set of transceiver module per line card and the associated modes supported.

Following the model tree with the 5 blocks described above. In green the new leaves/containers are added in this proposal; in black the non-modified leaves, even if they have been reallocated within the tree under different containers/lists.

<img width="1150" alt="image" src="https://github.com/openconfig/public/assets/10529210/a1810dd0-2e84-40c5-baf4-bfc43aad623c">

For more clarity on the above please check the following common definitions and assumptions defined during the design process of this proposal within the Telecom Infra Project (TIP) OOPT MUST project.

#### Common definitions
- System-vendor = the O-OT host platform provider (e.g. muxponder shelf, router, switch) and system integrator including the network operating system of the O-OT
- Manufacturer = Transceiver manufacturer (pluggable)
- Bookended scenario definition.
   - The System Vendor is the same in the two O-OTs.
   - The O-OTs of the same system vendor might host different Transceiver manufacturers.
   - Mode-ids are defined and maintained by the system vendor.
   - Interoperability shall be guaranteed by the system vendor if the same mode-id is configured in the line interfaces /optical-channels of the two O-OTs.

#### Assumptions
- The **openconfig-terminal-device-properties.yang** is a standalone model which represents static properties of a given terminal device, including:
  - **Operational-modes’ characteristic properties** on the transceiver configuration which characterize the mode (modulation-format, baud-rate, bit-rate, fec-format, filter…)
  - **Mode-descriptors** which describe the transmission design properties (Tx/Rx properties + CHARacteristic properties) of the implementation of the mode conditioned by the HW platform (transceiver + linecard) (Rx/Tx-OSNR, CD/PMD tolerances, penalties…)
  - **optical-channel’s configuration constraints** introduced by the HW implementation (min/max-central-frequency, min/max-output-power…)
- The openconfig-terminal-device-properties.yang is a standalone model representing a given mode’s static properties. We shall avoid any reference btw the manifest files and the operational openconfig trees which change dynamically. In other words, the references between mode-descriptors and operational modes, shall be through absolute identifiers.


