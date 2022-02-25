# gNMI Authorization Protocol

## Objective

This proto definition and the reference code(to be delivered seperately) serve
to describe an authorization framework for controlling which gNMI paths of a
network device users can access. The authorization policy is initially intended
to be deployed to a device, with the ability to define:

*   Policy rules - each rule defines a single authorization policy.
*   Groups of users - as a method to logically group users in the administrative
    domain, for instance: operators or administrators.
*   Users - individual referenced in rules or group definitions.

Authentication information is not included in this Authorization configuration.

Policy rules are matched based on the best match for the authorization request,
not the first match against a policy rule. Best match enables a configuration
which permits a user or group access to particular gNMI paths while
denying subordinate portions of the permitted paths, or the converse, without
regard to ordering of the rules in the configuration.

## Best Match

Authorization is performed for a singular user, gNMI path access, and access
methodology (READ/WRITE). The result of an Authorization evaluation is an
Action (Permit/Deny), policy version, and rule identifier.

A Best, or most specific, match is that which has the longest match to the
requested path and prefers:

*   a specific user over a group in the matching policy.
*   a defined KEY over a wildcard element in a keyed path.

Authorization rules must be defined such that a single best match is possible.
If the result of policy evaluation is more than one match, an error must be
raised.

Match rules permit a match against:

*   User or Group (not both)
*   an gNMI path
*   an access method (READ / WRITE / SUBSCRIBE)

An implicit deny is assumed, if there is no matching rule in the policy. Logging
may be specified on a per-policy-rule basis as well as a default for the whole
authorization policy.

As a request is evaluated against the configured policy, a READ / SUBSCRIBE
request for the configuration tree may traverse all of the tree and subtrees.
For portions of the tree for which the user has no access no data will be
returned. A WRITE request which attempts to write to a denied gNMI path or
element will return a "Permission Denied" error to the caller.

[gNMI paths](https://github.com/openconfig/reference/blob/master/rpc/gnmi/gnmi-specification.md#222-paths)
are hierarchical, and rooted at a defined "origin". gNMOpenConfigI may contain paths
such as:

```proto
    /a/b/c/d
```

Paths may also have attributes associated with the path elements such as:

```proto
    /a/b[key=foo]/c/d
```

Attributes may be wildcarded in the policy, such as:

```proto
    /a/b[key=*]/c/d
```

Wildcards are only valid in policy rules when applied to attributes.
Permitted use:

```proto
    /a/b[key=*]/c/d
```

Not permitted use:
```proto
    /a/b[key=foo]/*/d
```

An example of two rules with similar paths that differ only with respect
to attribute wildcarding:

```proto
   /a/b[key=FOO]/c/d
   /a/b[key=*]/c/d
```

The first rule specifies a single key, that rule is more specific than
the second rule, which specifies 'any value is accepted'.


## Bootstrap / Install Options

System bootstrap, or install, operations may include an authorization policy
delivered during bootstrap operations. It is suggested that the bootstrap
process include the complete authorization policy so all production tools
and services have immediate authorized access to finish installation and
move devices into production in a timely manner.

Using the Secure Zero Touch Provisioning (sZTP - RFC8572) process for
bootstrap/installation provides a method to accomplish this delivery, and
the delivery of all other bootstrap artifacts in a secure manner.

## Conflict Resolution

Policy evaluation should end with a single best match policy for the provided
user / path. If there is more than one 'best match', an error must be logged
and evaluation should return a failure.

## An Example Authorization Protobuf

```proto
# proto-file: //experimental/users/morrowc/gnXi/proto/authorization.proto
# proto-message: Policies
version: "UUID-1234-123123-123123"
created_on: 1234567890
# Define 2 groups.
group {
  name: "family-group"
  members { name: "stevie" }
  members { name: "brian" }
}
group {
  name: "test-group"
  members { name: "crusty" }
  members { name: "the-clown" }
}
# Action stevie to access /this/is/a/message_path in a READ manner.
policy {
  id: "one"
  log_level: LOG_NONE
  path {
    origin: "foo"
    elem { name: "this" }
    elem { name: "is" }
    elem { name: "a" }
    elem { name: "message_path" }
  }
  action: PERMIT
  level: READ
  user { name: "stevie" }
}
# Action members of family-group to access /this/is/a/different/message_path in
# both READ and WRITE methods.
policy {
  id: "two"
  log_level: LOG_BRIEF
  path {
    origin: "foo"
    elem { name: "this" }
    elem { name: "is" }
    elem { name: "a" }
    elem { name: "different" }
    elem { name: "message_path" }
  }
  action: PERMIT
  level: READ
  level: WRITE
  group { name: "family-group" }
}
# Demonstrate READ access to a key with an attribute defined.
policy {
  id: "key"
  log_level: LOG_BRIEF
  path {
    origin: "foo"
    elem { name: "this" }
    elem { name: "is" }
    elem { name: "a" }
    elem {
      name: "keyed"
      key{
        key: "Ethernet"
        value: "1/2/3"
      }
    }
    elem { name: "message_path" }
  }
  action: PERMIT
  level: READ
  group { name: "test-group" }
}
# Demonstrate READ access to a key with a wildcard attribute.
policy {
  id: "wyld"
  log_level: LOG_BRIEF
  path {
    origin: "foo"
    elem { name: "this" }
    elem { name: "is" }
    elem { name: "a" }
    elem {
      name: "keyed"
      key{
        key: "Ethernet"
        value: "*"
      }
    }
    elem { name: "message_path" }
  }
  action: PERMIT
  level: READ
  group { name: "family-group" }
}
# Demonstrate a key with a wildcard attribute and a user specific match.
# The previous policy matches all family-group users and permits a command
# path, the policy rule below specifically denies brian access to these paths.
policy {
  id: "wyld-stallions"
  log_level: LOG_BRIEF
  path {
    origin: "foo"
    elem { name: "this" }
    elem { name: "is" }
    elem { name: "a" }
    elem {
      name: "keyed"
      key{
        key: "Ethernet"
        value: "*"
      }
    }
    elem { name: "message_path" }
  }
  action: DENY
  level: READ
  group { name: "brian" }
}
# Add a final rule which is an explicit deny rule.
policy {
  id: "explicit-deny"
  log_level: LOG_FULL
  level: READ
  level: WRITE
  action: DENY
}
```

The example first policy rule:

```proto
# Action stevie to access /this/is/a/message_path for READ.
policy {
  id: "one"
  log_level: LOG_NONE
  path {
    origin: "foo"
    elem { name: "this" }
    elem { name: "is" }
    elem { name: "a" }
    elem { name: "message_path" }
  }
  action: PERMIT
  level: READ
  user { name: "stevie" }
}
```

permits the singular user "stevie" to access the path:

```shell
    /this/is/a/message_path
```

Additionally, "stevie" is permitted access to all paths below the defined path,
in a READ only mode, such as:

```shell
    /this/is/a/message_path/the
    /this/is/a/message_path/the/one
    /this/is/a/message_path/the/one/that
    /this/is/a/message_path/the/one/that/knocks
```

The second policy rule:

```proto
# Action members of family-group to run /this/is/a/different/message_path
policy {
  id: "two"
  log_level: LOG_BRIEF
  path {
    origin: "foo"
    elem { name: "this" }
    elem { name: "is" }
    elem { name: "a" }
    elem { name: "different" }
    elem { name: "message_path" }
  }
  action: PERMIT
  level: READ
  level: WRITE
  group { name: "family-group" }
}
```

example policy permits members or the family-group access to a single path, for
reading or writing:

```shell
    /this/is/a/different/message_path
```

and all path elements beyond "message_path":

```shell
    /this/is/a/different/message_path/foo
    /this/is/a/different/message_path/bar
    /this/is/a/different/message_path/foo/baz/bing/boop
```

The third policy rule:

```proto
# Demonstrate a key with an attribute defined.
policy {
  id: "key"
  log_level: LOG_BRIEF
  path {
    origin: "foo"
    elem { name: "this" }
    elem { name: "is" }
    elem { name: "a" }
    elem {
      name: "keyed"
      key{
        key: "name"
        value: "Ethernet1/2/3"
      }
    }
    elem { name: "message_path" }
  }
  action: PERMIT
  level: READ
  group { name: "test-group" }
}
```

Permits access by the "test-group" users to the keyed path, in a
read only manner:

```shell
    /this/is/a/keyed[name=Ethernet1/2/3]/message_path
```

and all path elementas as beyond "message_path". The final policy rule:

```proto
# Demonstrate a key with a wildcard attribute.
policy {
  id: "wyld"
  log_level: LOG_BRIEF
  path {
    origin: "foo"
    elem { name: "this" }
    elem { name: "is" }
    elem { name: "a" }
    elem {
      name: "keyed"
      key{
        key: "name"
        value: "*"
      }
    }
    elem { name: "message_path" }
  }
  action: PERMIT
  level: READ
  group { name: "family-group" }
}
```

permits access by the "family-group" users to the keyed path, with no
restrictions on the key values, but still as read-only:

```shell
    /this/is/a/keyed[name=Ethernet1/2/3]/message_path
    /this/is/a/keyed[name=POS3]/message_path
    /this/is/a/keyed[name=Serial4/1]/message_path
    /this/is/a/keyed[name=HSSI2]/message_path
```

Additionally, the path elements beyond "message_path" are available for access
to this group as well.

The wildcard character "*" (asterisk) is only able to be used as a value in
keyed elements, if the keys are missing in a keyed path a wildcard is assumed.
The wildcard is only used to mask out all possible values but not portions of
values, for instance:

```shell
    /this/is/a/keyed[name=*]/things - permitted usage of wildcard
    /this/is/a/keyed[name=Ethernet1/*/3]/things - NOT permitted usage of wildcard
```

The policy rule:

```proto
# Add a final rule which is an explicit deny rule.
policy {
  id: "explicit-deny"
  log_level: LOG_FULL
  level: READ
  level: WRITE
  action: DENY
}
```

provides an explcit deny for any request wich does not match any other policy
rule. This rule also requests that the result be logged in full fidelity.

