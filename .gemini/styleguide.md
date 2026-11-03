
### **1. General Coding & Contribution Guidelines**

**Source:** `doc/openconfig_style_guide.md`

* All contributions should follow the [OpenConfig contributions guide](doc/openconfig_style_guide.md)

### **2. Comparison of OpenConfig implementations with industry implementations

**Source:** `doc/contributions_guide.md`

* Gemini should compare the yang changes with the documentation of implementations of network operating systems from Cisco IOS-XR, Nokia SR Linux, Juniper JunOS/JunOS Evolved, Arista EOS, NVIDIA Cumulus Linux and SONiC.    
* In the review output, Gemini should generate a table comparing the affected OpenConfig leaves with a description of a similar command line or group of commands for each network operating system that perform similar functions.  Because there is often not be a 1:1 mapping between OpenConfig leaf and command line, it is sufficient to compare the changed OpenConfig model at a container level or even a top level and provide links to the relevant network operating system documentation.
* For each network operating system reference, find the source documentation which describes the features and commands that are similar to the OpenConfig implementation.  Provide a URL link to the source document.  Verify that the URL is accessible without special credentials required for access.
* If there is no match or only a weak match, this should should be explicitly stated in the output.  If there is no match, a URL does not need to be provided.
* If an enum or identity is defined, compare all of them with the reference implementations.

After the table output, include text for the following items:
* If at least two of these reference implementations do to have a strong match, suggestions of changes which may help make a match stronger may be added.
* If two more implementations have a strong match with the yang changes, then the change meets the minimum OpenConfig guidelines for vendor neutral modeling.
