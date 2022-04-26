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