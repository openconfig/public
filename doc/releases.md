# OpenConfig Release Versioning

## Background and Motivation

While each individual OpenConfig model can be tagged with a semantic version
(see [semver.md](semver.md)), models are often interdependent, or need to be
used together, for example when managing a full device. It is therefore useful
to define OpenConfig "releases" that contain a set of models that are designed
to work together. This also enables tracking breaking changes at the repository
level, as well as allowing public users to view and download tagged collections
of self-consistent models (see also the description of
[GitHub releases](https://help.github.com/articles/creating-releases/)).

In light of the above, this proposal introduces tagging
[semantic versions](https://semver.org/) to the set of all OpenConfig models as
a whole along with some OpenConfig-specific guidelines. Each release is
therefore the cumulative set of models committed to the master branch at a
certain point in time, and is tied to a specific commit in the OpenConfig
repository. YANG validators ensure that each release consists of collection of
published OpenConfig models that work together: that is, interdependencies
(e.g., imports, augments) and cross-references (e.g., leafrefs) are all
resolved.

As a side note, these releases are compatible with the notion of
[YANG release bundles](https://github.com/openconfig/public/blob/master/release/models/catalog/openconfig-module-catalog.yang).

## Policy

### Basic Guidelines

1.  A regular release of https://github.com/openconfig/public containing a set
    of compatible models (consisting of the entire set of models within the
    `openconfig/public` repo) is created roughly every **quarter**. The tag is
    named `vx.y.z` (e.g. `v1.2.0`) following
    [semantic versioning rules](https://semver.org/). A major, minor, or patch
    version increment is possible at each release, although non-backward
    compatible releases SHOULD be released at a less-frequent cadence.

    At the current time, releases are only expected to occur at the HEAD branch
    of the repository, meaning patch releases for non-HEAD model versions are
    not expected to be made.

2.  Non-backward compatible model changes affecting a feature that has
    reasonable functional test coverage (via
    [OpenConfig featureprofiles](https://github.com/openconfig/featureprofiles/))
    or implemented on a device SHOULD be made infrequently. The OpenConfig
    working group will create non-backward compatible releases periodically by
    considering both velocity and maintenance cost implications.

    e.g. It is November 2022, and the current latest release of OpenConfig
    models is `v2.3.1`. The OpenConfig community decides to change the default
    value of the leaf path `/interfaces/interface/config/enabled`. This is a
    breaking change since it would cause featureprofile tests that test for the
    behaviour of the default value without explicitly setting this leaf to begin
    to fail. As a result, the pull request for this change is not merged until
    the end of the quarter in December. In January 2023, a new release is
    created, versioned `v3.0.0` containing this update.

3.  Any non-backward compatible change for a pre-`v1.0.0` YANG module does NOT
    by itself necessitate a major revision change for the overall models
    repository.

4.  [Patch releases](https://semver.org/#spec-item-6) may be created at any time
    for backward compatible bug fixes, or equivalently, where only patch number
    increases occurred in models.

5.  [Pre-releases](https://semver.org/#spec-item-9) may be created at anytime to
    quickly introduce new changes to the models. These are not intended to be
    long-term, stable releases -- they should be replaced with the next regular
    release that encompasses these changes as soon as it becomes available.

6.  Wherever possible, it is RECOMMENDED to make backward compatible API changes
    (e.g. deprecating leaves via the
    [status statement](https://datatracker.ietf.org/doc/html/rfc6020#section-7.19.2))
    for at least one minor release prior to a non-backward compatible API change
    in order to ease the transition to the new API. These leaves are then
    expected to be removed or modified in the next major version release. NOTE:
    This guideline may change once OpenConfig operators gain more experience
    managing breaking changes.

7. [Deprecated](https://datatracker.ietf.org/doc/html/rfc6020#section-7.19.2)
   nodes must be supported until they are deleted from the models. The deprecated
   status serves as a hint that the node will be removed in some future release of
   OpenConfig data models and operators are advised to stop using the node.  The
   node description will recommend an alternative node or action.

9.  Release documentation should include the list of models and their version
    numbers contained in the corresponding release.

Each release `vx.y.z` (e.g. `v1.2.0`) MAY be given a name for easier human
identification, e.g. "September 2022". A client can thus assert they are
compatible with the "September 2022" release of OpenConfig. It's expected that
vendors will have some deviations and augments from the baseline; further, some
vendors may offer the ability to configure their NOS (network operating system)
to support different releases of https://github.com/openconfig/public.

Note that release version numbers need not be a function of individual model
release numbers (e.g., the max version number of all of the models in the
release).

### Corner-Case Guidelines

For non-backward compatible changes involving changing the type of a leaf, the
new leaf SHOULD have a different name than the previous leaf.
