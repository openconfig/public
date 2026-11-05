
### **1. General Coding & Contribution Guidelines**

**Source:** `doc/openconfig_style_guide.md`

* All contributions should follow the [OpenConfig contributions guide](doc/openconfig_style_guide.md)

### **2. Comparison of OpenConfig implementation with industry implementations

**Source:** `doc/contributions_guide.md`

* Compare the yang changes with the documentation of implementations of network operating systems from Cisco IOS-XR, Nokia SR Linux, Juniper JunOS/JunOS Evolved, Arista EOS, NVIDIA Cumulus Linux and SONiC.    
* In the review output, Gemini should generate a table comparing the affected OpenConfig leaves with a description of a similar command line or group of commands for each network operating system that perform similar functions.  Because there is often  not be a 1:1 mapping between OpenConfig leaf and command line, it is sufficient to compare the changed OpenConfig model at a container level or even a top level and provide links to the relevant network operating system documentation.
* For each industry implementation if there is a match with openconfig the output should include a publicly accessible URL reference link, citing the source of the documentation where the relevant feature and commands are described.
* Verify that the URL is a valid page and returns a 200 OK.
* If there is no match or only a weak match
  * this should should be explicitly stated in the output.
  * There should not be any URL reference link.
* If an enum or identity is defined, compare all of them with the reference implementations.  
* If at least two of these implementations do to have a strong match, explicit feedback should be given, including suggestion of changes which may help make a match stronger.
* If two more implementations have a strong match with the yang changes, then the change meets the minimum OpenConfig guidelines for vendor neutral modeling.
