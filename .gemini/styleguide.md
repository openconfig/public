
### **1. General Coding & Contribution Guidelines**

**Source:** `doc/openconfig_style_guide.md`

* All contributions should follow the [OpenConfig contributions guide](doc/openconfig_style_guide.md)

### **2. Comparison of OpenConfig implementations with industry implementations

**Source:** `doc/contributions_guide.md`

* Gemini should compare the OpenConfig YANG changes with
  * The documentation of implementations of network operating systems (NOS) from Cisco IOS-XR, Nokia SR Linux, Juniper JunOS/JunOS Evolved, Arista EOS, NVIDIA Cumulus Linux and SONiC.
  * The industry YANG models exposed at https://github.com/YangModels/yang
* In the review output, generate a table comparing the affected OpenConfig leaves with a description of a similar command line or group of commands for each NOS that perform similar features.  Because there is often not be a 1:1 mapping between OpenConfig leaf and command line, it is sufficient to compare the changed OpenConfig model at a container level or even a top level and provide links to the relevant NOS documentation.
* A similar features in the NOS should have configuration or telemetry items that match those that are being proposed.  An exact match in the configuration item indicates a strong match.  Similar, but not exactly matching items may be considered a partial match or weak match. 
* If an enum or identity is defined, compare all of them with the reference implementations.
* For each NOS reference, find the source documentation which describes the features and commands that are similar to the OpenConfig implementation.  Provide a URL link to the source document.  Verify that the URL is accessible without special credentials required for access.  If there is no match, no documentation URL needs to be provided.

After the table output, include text for the following items:
* If at least two of these reference implementations do to have a strong match, suggestions of changes which may help make a match stronger may be added.
* If two more implementations have a strong match with the yang changes, then the change meets the minimum OpenConfig guidelines for vendor neutral modeling.
