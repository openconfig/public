
### **1. General Coding & Contribution Guidelines**

**Source:** `doc/openconfig_style_guide.md`

* All contributions should follow the [OpenConfig contributions guide](doc/openconfig_style_guide.md)

### **2. Comparison of OpenConfig implemention with industry implementions

**Source:** `doc/contributions_guide.md`

* Gemini should compare the yang changes with the documentation of implementions of network operating systems from Cisco IOX XR, Nokia SR Linux, Juniper JunOS, Arista EOS, NVIDIA Cumulus Linux and SONiC.  
* In the review output, Gemini should generate a table comparing the affected OpenConfig leaves with a description of a similar command line or group of commands that perform similar functions.  Because there is often  not be a 1:1 mapping between OpenConfig leaf and command line, it is sufficient to compare the changed model at a container level or even a top level and provide links to the relevant network operating system documentation.
* If there is no match or only a weak match, this should should be stated.  
* If an enum or identity is defined, compare all of them with the reference implementations.  
* If at least two of these implementations do to have a strong match, explicit feedback should be given, including suggestion of changes which may help make a match stronger.
