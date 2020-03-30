
# YANG authoring guidelines for OpenConfig models

**Contributors:** Anees Shaikh, Rob Shakir, Kristian Larsson<br>
**October 26, 2015**<br>
*Updated: June 2, 2019*


## Background
This document describes conventions adopted in the OpenConfig operator group
when writing YANG modules.  YANG is a domain-specific language for describing
configuration and operational state data for networking systems, protocols, and
software.  The [official language
specification](https://tools.ietf.org/html/rfc6020) is maintained by the [IETF
NETMOD](https://datatracker.ietf.org/wg/netmod/documents/) working group.  The
current version of the language is 1.0, with version 1.1 expected to be ratified
and released soon.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [General guidelines](#general-guidelines)
  - [IETF guidelines](#ietf-guidelines)
  - [Module compilation](#module-compilation)
  - [Line length](#line-length)
  - [Module template](#module-template)
  - [Modeling operational state](#modeling-operational-state)
  - [Top-level data nodes vs. groupings](#top-level-data-nodes-vs-groupings)
  - [Module version](#module-version)
- [YANG style conventions](#yang-style-conventions)
  - [Naming](#naming)
    - [Module naming](#module-naming)
    - [Submodule naming](#submodule-naming)
    - [Grouping naming](#grouping-naming)
    - [Prefix naming](#prefix-naming)
  - [Path references](#path-references)
    - [Intra-model paths](#intra-model-paths)
    - [Inter-model paths](#inter-model-paths)
  - [Capitalization](#capitalization)
    - [Enumerations](#enumerations)
    - [Identities](#identities)
- [YANG language usage](#yang-language-usage)
  - [`list`](#list)
  - [`presence`](#presence)
  - [`feature` and `if-feature`](#feature-and-if-feature)
  - [`choice`](#choice)
  - [XPath](#xpath)
  - [Regular expressions](#regular-expressions)
- [Appendix](#appendix)
  - [Example groupings for containers](#example-groupings-for-containers)
  - [OpenConfig YANG module template](#openconfig-yang-module-template)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## General guidelines

### IETF guidelines
OpenConfig YANG modules adopt some rules from IETF modeling standards, including
some of those defined in [RFC 6087](https://tools.ietf.org/html/rfc6087).

### Module compilation
All YANG modules should be validated / compiled with
[pyang](https://github.com/mbj4668/pyang) using the following flags:

`pyang --strict --lint <module>`[*](#pyang)

All errors and warnings should be corrected before submitting or
posting any modules.  Use of the `--lint` flag will cause pyang to check for
most of the guidelines mentioned in RFC 6087.  Adding the --ietf flag will also
check for conventions required to submit OpenConfig models to the IETF when
appropriate.

<a name="pyang"></a>*: The --lint flag was introduced in pyang v1.6.

### Line length

Per IETF formatting guidelines, lines should be no more than 70 characters (also
checked by the `--ietf` flag).

### Module template

OpenConfig has adapted the example template from RFC 6087 as a starting point
for writing new YANG modules (see the [Appendix](#appendix)).

### Modeling operational state

OpenConfig has adopted a structural convention for YANG models that emphasizes
the importance of modeling operational state (i.e., monitoring and telemetry
data), in addition to configuration data.  At a high level, this convention uses
specially named `config` and `state` containers, in every subtree to explicitly
indicate configuration and operational state data.  The rationale and details
for this convention are described in an [IETF
draft](https://tools.ietf.org/html/draft-openconfig-netmod-opstate-01).

These conventions are reflected in naming and structure of YANG groupings and
containers as described in more detail in the corresponding sections of this
document.

### Top-level data nodes vs. groupings

In general, directly defined data nodes should be avoided in all modules.
Instead, define the top-level container or other data nodes in a grouping, and
then instantiate it once in the module with a `uses` statement.

This allows maximum reuse of data definitions across models, and also makes it
easier to compose models using simple imports.

Modules should generally have a single `xxx-top` grouping that allows it to be
instantiated in other modules.  This top-level grouping should not have any
self-augmentations.

### Module version

Every module must have an `openconfig-version` statement indicating its
semantic version number.  This statement is a YANG extension defined in the
openconfig-extensions module.  The YANG revision statement should reference
semantic version.

```
oc-ext:openconfig-version "0.4.0";

  revision "2016-05-31" {
    description
      "Public release";
    reference "0.4.0";
  }
```

Individual YANG modules are versioned independently -- the
semantic version is generally incremented only when there is a
change in the corresponding file.  Submodules, however, must have
the same semantic version as their parent modules.  Further details on
versioning rules are available in the definition of the
`openconfig-version` extension in the `openconfig-extensions.yang`
module.


## YANG style conventions
Style conventions describe guidelines related to conventions used in writing
YANG modules.

### Naming

#### Module naming

YANG modules should have filenames of the form `openconfig-<function>.yang`.
The ‘openconfig’ prefix indicates that the module is originated by the
OpenConfig operator group.

Examples: `openconfig-bgp.yang`, `openconfig-mpls.yang`,
`openconfig-interfaces.yang`

#### Submodule naming

Related module and submodule filenames should be named
`openconfig-<function>-<subfunction>.yang`.

Examples: `openconfig-bgp-policy.yang`, `openconfig-mpls-te.yang`,
`openconfig-if-ethernet.yang`

#### Grouping naming

Grouping names should make it easy to quickly understand the nature of the
data within.  A suggested convention is `xxx-yyy[-config|state|top]`, where xxx
is the top-level module name (without the `openconfig` prefix), yyy is a string
which indicates the contents of the groupings.

For data that will be placed in a container, three groupings should be created:

* `xxx-yyy-config` -- configuration (read/write) leaves or leaf-lists
* `xxx-yyy-state` -- operational state (read-only) leaves or leaf-lists
* `xxx-yyy-top` -- a top-level grouping that defines the container structure,
  with the enclosing container, and the `config` and `state` containers within.

See the example in the [Appendix](#appendix).

#### Prefix naming

Each module requires a `prefix` statement with a prefix that other dependent
modules will use (also used in path references within the same module). Prefixes
should be short and clear, with abbreviations as appropriate.

Module prefixes should be of the form `oc-xxx[-yyy]`

Examples: `oc-types`, `oc-lldp`, `oc-if-ethernet`

### Path references

#### Intra-model paths

For leafrefs, XPaths, augments, etc. use relative paths when referencing nodes
in the same module.

#### Inter-model paths

For references external to the module (i.e., in another namespace), absolute
paths may be used.

### Capitalization

In most cases, identifiers in YANG modules, e.g., names of leaves, lists,
containers, etc. are **lower case with dashes between words**.  Further details
below.

#### Enumerations

`enum` values within an enumeration type should be UPPER_CASE_WITH_UNDERSCORES,
keeping with conventions used for enumerated types in many programming
languages. They MUST begin with an alphanumeric character (A-Z or 0-9),
optionally followed by a "_" or "." or additional alphanumeric characters
(A-Z or 0-9).

Example:
```
   type enumeration {
     enum ACCEPT_ROUTE {
       description "default policy to accept the route";
     }
     enum REJECT_ROUTE {
       description "default policy to reject the route";
     }
   }

```

#### Identities

YANG identities allow the definition of a "base" constant and additional values
that act as "derived" types -- `identity` values, including base identities,
should be UPPER_CASE_WITH_UNDERSCORES.

Since identities are most often implemented as enumerations in language
bindings, it is helpful to follow the same convention as with enumerations.
Identities should be upper case such that where an identityref is used in
preference to an enumeration, this is transparent to the entity interacting with
the model.

Example:
```
identity FIBER_CONNECTOR_TYPE {
  description
    "Type of optical fiber connector";
}

identity SC_CONNECTOR {
  base FIBER_CONNECTOR_TYPE;
  description
    "SC type fiber connector";
}

identity LC_CONNECTOR {
  base FIBER_CONNECTOR_TYPE;
  description
    "LC type fiber connector";
}
```


## YANG language usage
Language rules describe guidelines on use of specific YANG language statements,
including how modules should be structured and parsed.

### `list`

YANG list keys should be quoted:
```
list interfaces {
  key "name";
  ...
}

list servers {
  key "address port";
  ...
}
```

YANG requires leaf nodes that are list keys to be direct descendents of the `list`
statement.  Since key leaf nodes must also be members of the list data, they will
generally reside in a `config` or `state` container (see
[Modeling operational state](#modeling-operational-state)).  Hence, the list key leaf
nodes should be of type `leafref` with a `path` pointing to the corresponding "actual"
leaf in the config or state container.

```
grouping interfaces-config {

  leaf name {
    ...
  }
}

grouping interfaces-list-top
  list interface {
    key "name";

    leaf name {
      type leafref {
        path "../config/name";
      }
    }

    container config {

      uses interfaces-config;
    }

    ...
  }
}
```

Lists should have an enclosing container with no other data nodes inside
it.

```
container interfaces {

  list interface {
    ...
  }
}
```

### `presence`

Use of `presence` containers should be avoided.

Presence containers express implicit configuration semantics, which is more
difficult for management systems to interpret.  An alternative is to use an
explicit "enabled" leaf (or similar) to make activation of the corresponding
configuration explicit.  Presence containers are also incompatible with
hierarchical models in which lower levels inherit configuration from higher
levels.

Presence containers in YANG reflect CLIs which turn configuration on or off with
a single feature keyword, e.g., `signalling graceful-restart`, rather than
`signalling graceful-restart enable`.

### `feature` and `if-feature`

Use of `if-feature` should be avoided.

The `feature` and `if-feature` statements are to define an optional feature and
designate specific data as part of the optional feature.  OpenConfig models are
vendor-neutral and intended to express an operationally complete set of
features.  Non-compliance by implementors should be expressed by deviation files
rather than `if-feature`.

To add extensions or additional features to a model beyond the base OpenConfig
model, vendors and implementors should rather use YANG augmentations or
extension modules.

### `choice`

Use of `choice` statements should be avoided where possible.

YANG offers `choice` statements as an analog to case/switch statements in other
languages.  However, `choice` nodes do not appear in the actual data instances,
or in schema paths -- they are used primarily for validating instance data to
ensure that only one of the sets of data appears.

Example:

```
choice bandwidth {
  case explicit {
    leaf bw-value {
      type uint32;
    }
  }
  case auto {
    leaf min {
      type uint32;
    }
    leaf max {
      type uint32;
    }
  }
}
```

The corresponding path to the `bw-value` leaf in the example is `.../bw-value`
rather than `.../bandwidth/explicit/bw-value` which is much clearer.

Since very few nodes in the model generally need to be made mandatory, an
alternative approach is allow both options to appear in the data and rely on
separate semantic validation in the management system or device to flag an
invalid combination.

If a conditional set of values is really needed, a `when` statement could be
used to validate that certain data is allowed in the data instance.

Example:

```
leaf bandwidth-set {
  type enumeration {
    enum AUTO;
    enum EXPLICIT;
  }
}
container explicit {
  when "../bandwidth-set = EXPLICIT";
  leaf bw-value {
    ...
  }
}
container auto {
  when "../bandwidth-set = AUTO";
  leaf min {
    ...
  }
  leaf max {
    ...
  }
}
```

In this approach all nodes appear in schema paths, and the `when` statement
still allows the management system to validate instance data.

### XPath

Avoid complex XPath expressions. The goal is to keep it simple, both for the
sake of readability but also so that the OpenConfig models can be used in
environments that only support a basic set of XPath functions.

The following guidelines should be followed when using XPath expressions in
models:

* tests should use simple operators, like equality - avoid complex data
  manipulation in the XPath rules
* paths used in XPath should always be relative where possible
* only use the following YANG 1.0 functions:
 - and
 - current
 - not
 - or
* avoid using any of the new XPath types that are included in YANG 1.1

### Regular expressions

Use regular expressions available in the [POSIX Extended Regular
Expressions](http://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap09.html)
standard.

The YANG language specification lists the [W3C XML Schema
specification](http://www.w3.org/TR/2004/REC-xmlschema-2-20041028/) as its
reference for [regular
expressions](http://www.w3.org/TR/2004/REC-xmlschema-2-20041028/#regexs).
However, this is not a commonly used standard for implementors.


## Appendix

### Example groupings for containers

```
grouping rsvp-graceful-restart-config {
  description
    "Configuration data ";
}

grouping rsvp-graceful-restart-state {
  description
    "Operational state data ";
}

grouping rsvp-graceful-restart-top {
  description
    "Top-level grouping ";

  container graceful-restart {
    description
      "Top-level container ";

    container config {
      description
        "Configuration data ";

      uses rsvp-graceful-restart-config;
    }

    container state {

      config false;

      description
        "Operational state data ";

      uses rsvp-graceful-restart-config;
      uses rsvp-graceful-restart-state;
    }
  }
}
```

### OpenConfig YANG module template

```
module openconfig- {

  yang-version "1";

  // namespace
  namespace "http://openconfig.net/yang/";

  prefix "";

  // import some basic types
  import ietf-inet-types { prefix inet; }


  // meta
  organization "OpenConfig working group";

  contact
    "OpenConfig working group
    www.openconfig.net";

  description
    "This module ";

  revision "" {
    description
      "Initial revision";
    reference "TBD";
  }

  // extension statements

  // feature statements

  // identity statements

  // typedef statements

  // grouping statements

  // data definition statements

  // augment statements

  // rpc statements

  // notification statements

}
```
