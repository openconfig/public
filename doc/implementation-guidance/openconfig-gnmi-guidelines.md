# OpenConfig use of gNMI
**Authors**: robjs@google.com  
**Published**: October 2020

## Aim
This document provides details of the use of gNMI and OpenConfig in conjunction
with one another. gNMI is specifically designed to be independent of:

- The modelling language that is used to describe the tree structure
- The models that are used within the transport

This flexibility allows use cases such as [Mixed
Schema](https://github.com/openconfig/reference/blob/master/rpc/gnmi/mixed-schema.md)
where vendor command line interface input  is represented using the `cli`
origin, as well as transport for data trees that are modelled in YANG as well
as protobuf.

OpenConfig and gNMI are designed to be used together, and some simplifying
assumptions are made such that they work well together. This document describes
specific implementation considerations of using the two together.

## Implementation Guidelines

### Path Specifications and the OpenConfig `origin`
The gNMI specification defines a path as consisting of a node name, and
attributes associated with the name. In the context of YANG, this means the
schema tree node name, and keys/values that are used within YANG `list` entries
which are stored in the data tree. YANG allows for complex use-cases, such as
two schema tree nodes with the same name being defined by the same module -
i.e., it is valid to have the following set of modules:

```yang
module a {
  ...
  container root { }
}
```

```yang
module b {
	...
	import a { prefix a; }
	augment /a:root { leaf foo { type string; } } }
}
```

```yang
module c {
	...
	import a { prefix a; }
	augment /a:root { leaf foo { type string; } } }
}
```

In this case, there are two definitions of `/root/foo`. To support this, YANG
prefixes the name of an element with the module that defines it - i.e., the two
`foo` leaves above become `/a:root/b:foo` and `/a:root/c:foo`. 

In practice, we find that this causes unnecessary complexity in an
implementation. Since paths must now be specified using the namespace, or
module name, within which they are defined: 

* a user now must care about the structure of the YANG modules that form the
  tree that they are considering. This means that the re-use of the module
  definitions now has an impact on the end user who may not even know that a YANG
  model is used to define the tree.
* there are user-facing changes required as the lifecycle of schema nodes is
  progressed through - for example, an implementation specific module may be
  used before upstreaming a change to a core module set. When this happens, often
  no node names will change, but the defining module may.

Therefore within OpenConfig -- i.e., when the `origin` YANG extension is set to
`openconfig`, and in gNMI when the `origin` field is set to `openconfig` or an
empty string:

* All paths MUST be unique when considering only the node name, and **NOT**
  prefixing this name with a module namespace.
* Base OpenConfig modules and augmentations MUST comply with this requirement.
  * We expect that a simple means to achieve this is to create vendor-specific
    node names for augmentations - e.g., `vendorFoo/config/new-config-leaf`.
* RFC7951-formatted JSON MAY omit the module name prefixes, as described in the
  [specification](https://tools.ietf.org/html/rfc7951#section-4).
 
This set of simplifying rules allow for separation of users from the exact
module structure, whilst still supporting extensibility of the OpenConfig tree.
We observe that asserting thse constraints avoids alternate approaches, such as
defining a single huge module for the entire tree, as some implementations do -
which causes significant difficulty in consuming the model for specific use
cases which make use of only a subset of nodes. This single module approach is
taken by some implementations, e.g.,
[`nokia-conf`](https://github.com/nokia/7x50_YangModels/blob/master/latest_sros_20.7/nokia-conf.yang),
and [`configuration`](https://github.com/Juniper/yang/blob/master/17.1/17.1R1/configuration.yang)
module. In these cases, there will be user-facing impact of breaking these
modules down into more consumable chunks, which the approach described above
avoids.

It should be noted that where other `origin` values are specified, the same
simplifying assumptions MUST NOT be assumed to be true, and rather each
specified origin should define, along with the schema used for the origin, any
specific constraints that can be made. Thus, outside of the `openconfig`
origin, gNMI paths MAY contain the YANG module name.
