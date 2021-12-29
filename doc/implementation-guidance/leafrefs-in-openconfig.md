# Implementor Guidance: Leafrefs in OpenConfig/YANG

**Authors**: Aaron Beitch<sup>†</sup>, Justin Costa-Roberts<sup>†</sup>, Rob Shakir<sup>‡</sup>  
<small><sup>†</sup> Arista Networks, <sup>‡</sup> Google</small>  
**Published:** September 2020

## Contents

  - [Summary](#summary)
  - [Background](#background)
    - [Why use a `leafref`?](#why-use-a-leafref)
  - [Problems with Leafrefs](#problems-with-leafrefs)
    - [Expense of Validation](#expense-of-validation)
    - [Assumption of Complete and Consistent Data Tree](#assumption-of-complete-and-consistent-data-tree)
    - [Inability to Wholly Describe Validity](#inability-to-wholly-describe-validity)
  - [Relaxed Leafref Validation Requirements](#relaxed-leafref-validation-requirements)

## Summary
The validation of leafrefs poses many problems for YANG clients and servers
alike - however, they continue to have value as a schema construct to provide
common type conformance, and foreign keys within the schema. This document
defines the expected validation behaviour of an OpenConfig-compliant
implementation - particularly, allowing an implementation to relax validation
rules such that referenced entities within the schema are not required to exist
at the time of validation (i.e., always act as though `require-instance false`
were set for a leafref node).

## Background

As defined in [section 9.9 of
RFC6020](https://tools.ietf.org/html/rfc6020#section-9.9), a leafref is a type
of YANG leaf. Each leafref has a path defined in a limited XPath syntax. This
path must evaluate to a leaf in the schema. The values for this leafref are
constrained by the type of the leaf referred to by the leafref. In addition, a
leafref’s value is constrained to be in the nodeset returned from executing its
path as an XPath expression. In YANG 1.1, but not in the original YANG
specification, this second constraint can be disabled by the module author with
a statement of “require-instance false;”. See [section 9.9.3 of RFC
7950](https://tools.ietf.org/html/rfc7950#section-9.9.3). The OpenConfig YANG
models do not currently support YANG 1.1.

### Why use a `leafref`?

We observe that there are multiple reasons that a YANG schema author may choose
to use a `leafref` within their module, particularly:

1. To re-use a common type between two leaves. In this case, the author wishes
   to define that the type of the `leaf` that they are defining has exactly the
   same type as the type of the  remote `leaf` that is selected by the `path`
   specified.
   * A common use case for such leafrefs within the OpenConfig schema is for
     list keys, where there is an alias to the authoratitive source of a list
     key (under a `config` or `state` container) directly under the `list` structure
     itself, to comply with the YANG requirements.
2. To provide a foreign-key within one part of the schema to another schema
   element. The common use case for such leaves is to allow the schema to be
   normalized, and to allow programmatic means of traversing the schema (e.g.,
   automatically retrieving data that corresponds to such `leafref` references,
   without needing to write schema-path specific handling code).
   * For example, defining a `leafref` with the path
     `/interfaces/interface/config/name` provides a client a means to know that
     more definition of the referenced interface can be found at the corresponding
     schema path - and for a helper-library for a client querying such a leafref to
     automatically also retrieve the corresponding `/interfaces/interface` list
     entry.
3. To ensure consistency within a schema. In this case, the `leafref` is used
   to ensure that an entity that is referenced within a particular part of a
   schema actually exists within the schema.
   * For example, the list of interfaces with a particular protocol (say,
     IS-IS) enabled on them must reference an interface that is actually
     configured on the system. The leafref ensures that if IS-IS is enabled on
     `Ethernet423` such an interface actually exists on the local system.

Both 1 and 2 within the list are essentially schema-time definitions, they are
either convenying type information, or static information that relates to
characteristics (in this case, relationships between parts) of the schema tree.

The third case is somewhat different - it conveys characteristics of the data
tree through the leafref. We observe that implementing validation of such cases
requires additional consideration.

## Problems with Leafrefs

### Expense of Validation
Unlike most constraints in YANG, leafref constraint validation depends on more
than just the value in the node containing the leafref: it also depends on the
value of the nodes indicated by the path of the leafref. This means that a
modification to any node may invalidate a leafref node in any other part of the
tree. Validating a data tree after any node modification therefore requires
either re-examining the entire data tree or implementing complex bookkeeping
linking nodes to the leafref nodes that are constrained by them. 

In our implementation experience, detecting populated `leafref` leaves, and
performing data tree lookups to prove their validity is the dominant
computational expense during YANG tree validation.  This complexity also
affects the clients interacting with an OpenConfig implementation. Crafting a
valid Set operation requires knowledge of the leafrefs anywhere in the
OpenConfig models.

### Assumption of Complete and Consistent Data Tree

Leafrefs for validation assume a complete and consistent data tree is available
at the time of validation. In practice there are numerous deployments where
this constraint is not met:

 * During data tree synchronisation -- a natural way to synchronize two data
   trees (for instance, from a target device to a management system, or from a
   NOS’s native data into an OpenConfig tree) is to subscribe to updates from a
   source tree and apply them to the destination tree. If the destination must
   validate leafrefs, then either updates from the source tree must be composed
   into valid updates to apply to the target, or the source must only issue
   updates that satisfy the constraints imposed by leafrefs on any nodes contained
   in the update. In general, without support from the source, the former is not
   possible. If the source system does not provide such a capability, the
   destination tree must bypass leafref validation or implement complex systems
   for hiding the fact that some leafref leaves are temporarily invalid. In this
   model, assuming the source tree is valid at all times, the destination tree
   will also become valid, but may have periods where it is not.
 * Within modular configuration-generation systems -- in some cases, the
   configuration of a network element may not be generated by a single process.
   Individual handlers (e.g., a BGP session generator) may generate a part of the
   configuration, which has external references. At this point in time, the
   partial configuration is not available to be validated, however, it is still of
   interest whether the leaves that were defined are valid. In this case, the
   toolchain must provide means to be able to skip validation, as described above.
 * In mixed-schema configuration systems -- in today's network elements, there
   are typically multiple schemas available for configuration - legacy CLI, a
   vendor's native models, OpenConfig etc., with varying levels of support. We do
   not expect that leafrefs exist across these schemas, such that any entity that
   is handling a mixed-schema configuration must also provide means to bypass
   validation of leafrefs.
 * During pre-configuration within a NOS -- some network elements may wish to
   allow preconfiguration of a particular entity. For example, a client may be
   allowed to reference a particular target interface before it is configured. To
   support such behaviour, target implementations must provide some means to be
   able to disable leafref validation for a particular client operation if it is
   enabled by default.
 * In split RO vs. RW systems -- some system architectures split the collection
   of state data ('telemetry') from the storage of configuration data. Some
   systems may not therefore have visibility into both the 'intended'
   configuration as well as the 'applied' and 'derived' states of the system.
   In such systems, any leafref that is from a `config false` to a `config
   true` path in the data-tree will not resolve, and hence leafrefs between
   these views of the data tree cannot be validated. 
 
To this end, in our experience the choice of *when* the existence of a
leafref's target can be validated, and *which* leafrefs should be expected to
have valid targets is wholly dependent upon the context of the consuming system
-- rather than a generic property of a YANG schema.

### Inability to Wholly Describe Validity
Whilst the constraint that a particular referenced entity must exist on the
system is a useful validity constraint, it is generally insufficient to
describe the actual validity of the configuration. For example, if an interface
is to be included within a VRF, it is very likely to need to be defined, but it
likely also needs to be specified in a mode that means that it is compatible
with L3 routing - among a host of other constraints. In practice, therefore,
the mere existence of the interface is insufficient to describe validity.

To overcome this scenario `when` and `must` statements could be added to the
schema. However, there are multiple concerns with doing so:

 * It incurrs the same expense as is described in _Expense of Validation_ above
   on the server to ensure validity.
 * The choice of XPATH to describe these constraints means that toolchains must
   now handle a complex additional query language, and authors must craft such
   queries.
 * Both `when` and `must` statements have the same expectation of consistency
   that is described in _Assumption of a Complete and Consistent Data Tree_.

Use of `when`, `must` and `leafref` to indicate an instance is required are
somewhat outside the general scope of a data modelling language. For example,
protobuf, JSON schema, and thrift do not provide means to describe such
constraints. Use of these constraints in YANG implies an additional level of
complexity for the supporting toolchain, over and above that which is required
for other languages. 

In our implementation experience, in order to simplify the overall developer
experience and system architecture surrounding YANG modelled data, it is
significantly simpler to keep such constraints (if enforced) outside of the
data model specification.

## Relaxed Leafref Validation Requirements

By default, a server supporting OpenConfig:

 * SHOULD validate types for leaves of type `leafref`, along with any
   restrictions (e.g., `pattern` or `range` statemenets) described with those
   types.
 * SHOULD NOT validate the existence of the leaf that is pointed to by the
   leafref `path` by default. This behaviour MUST be explicitly enabled if
   required, and MAY be unsupported by a target. This guidance essentially treats
   all leafrefs within the schema as `require-instance false`.
 
Authors of OpenConfig models:

 * SHOULD continue to use `leafref` where references between parts of the
   schema are required.

System implementors:

 * MUST make their systems explicitly check references if they are required for
   their application.
 * MAY choose to implement mechanisms such as only validating the existence of
   leafref targets within a subset of the schema tree to overcome some of the
   issues described in this application note.
