# Versioning Individual OpenConfig models

<em>contributors</em>: Anees Shaikh, Josh George, Rob Shakir, Kristian Larsson<br>

## Background and Motivation

*For versioning the set of all OpenConfig models as a whole, see
[releases.md](releases.md).*

This document proposes to adopt [Semantic Versioning](http://semver.org/)
(semver) for published OpenConfig YANG models in the same way that software
projects use similar versioning schemes to indicate the maturity and
compatibility of software as it evolves. Semver bases its versioning on an API
contract with developers and users. The basic format of a semver number is
XX.YY.ZZ-release where XX is the major version number, YY is the minor version
number, and ZZ is a patch level. Additional information can be optionally
provided with a suffix. Detailed specification on the semver versioning rules
are available at the [link](http://semver.org/) above. Any non
backward-compatible change to the API requires incrementing the major version
number, while feature changes that do not break clients are indicated by
incrementing the minor version number. Non-feature patches that are backward
compatible are indicated with an increment to the patch number.

Semantic versioning is proposed as an addition to YANG revision statements for a
number of reasons:

*   YANG language rules state that the API never changes in a
    backward-incompatible way. From RFC 6020: “... changes are not allowed if
    they have any potential to cause interoperability problems between a client
    using an original specification and a server using an updated
    specification.”

This is simply not practical (and is largely motivated by SNMP MIB notions).
YANG models are not mature (less than 5 models have been made IETF RFCs and
these are not implemented by any major device platform). Server and client
implementations are only now being developed and deployed and significantly more
operational experience is needed before APIs can be frozen.

*   YANG revision statements consist of a date and some informational text. As
    such, they offer little information about what has changed from one revision
    to the next. This is perhaps not surprising when considering the rigid rules
    in YANG about guaranteed API compatibility.
*   YANG revision statements are meant for human consumption -- they are not
    very useful for any sort of programmatic dependency checking.

Semantic versioning has its
[own issues](https://gist.github.com/jashkenas/cbd2b088e20279ae2c8e) and it may
be that in OpenConfig we will have to adapt the specification somewhat based on
considerations for versioning YANG models. Also semver does not address the
problem of how to version groups of interdependent modules (e.g., a device model
composed of many constituent models).

Note that we would continue to use revision statements, e.g., with a date set to
the day a new semantic version is published. This allows consumers to continue
to use current YANG constructs such as import by revision.

## General guidelines for versioning OpenConfig YANG modules

An immediate question that arises when considering how to version YANG modules
is what criteria should be used to judge that a module is mature enough that an
API contract should be established with a version number.

According to the semver specification, software that is pre-release with major
version 0 may break clients as long as the major version number remains < 1.
That is, with major version 0, there should be no expectation of compatibility
from one release to another, even if only the minor version number is changing.

Based on these considerations, the following basic guidelines are proposed for
versioning OpenConfig modules:

*   All modules should start out with a 0 major number. The major number should
    remain 0 as long as the model is being reviewed and revised based on
    feedback from the OpenConfig operators and from vendors implementing the
    model.
*   Semver guidelines should be followed while the model is at major number 0,
    i.e., API or feature changes should increment the minor number, while minor
    fixes should increment the patch number.
*   Once a vendor implementation for a model is in progress, the major number
    should be changed to 1 to acknowledge that the API is being used by
    implementors with correspondingly more disruption likely when the model
    changes in incompatible ways. Deciding that vendor implementations are
    sufficiently in-progress to justify moving to major version 1 may be
    somewhat subjective and should be based on detailed discussions with
    implementors to understand what stage they are in their implementations.

## API changes in YANG modules

For the purposes of semver, the API presented by a YANG model consists of its
data nodes and corresponding paths. Other elements of the model may not,
strictly speaking, be considered part of the API, but still could have
significant impact on the use of the model by developers or clients. Such
elements include default values, configurability of a node, and behavior of a
given data node (as described by the description statement).

Since the API of the YANG module is a combination of these explicit and implicit
elements, the criteria for determining when a revision requires a major number
increment is not always straightforward. Below we list some general rules for
determining the API has changed, and consequently would increment the major
version number.

*   Any leaf, leaf-list, list, or container modifications that result in
    changing an existing data node name, or the path to a data node (location in
    the model)
*   Changing the target of a leafref
*   Removal of a data node (leaf, leaf-list, list, container)
*   Changing the type of a leaf or leaf-list
*   Changing a type definition such that data based on the existing typedef
    would be invalid (e.g., removing a value from an enumeration, changing the
    base type in a typedef, etc.)
*   Changing the key of a list (i.e., using a different data node as the list
    key)
*   Changing a conditional statement, such as when or must, to be more
    restrictive, or to be based on a different condition altogether
