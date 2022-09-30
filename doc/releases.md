# OpenConfig Releases

## Rationale

As the OpenConfig project matures and is adopted by more implementors and
network operators, model changes are becoming increasingly frequent. This has
made identifying compatible and well-known sets of models more difficult.
Furthermore, uncontrolled and frequent breaking changes can create excessive
burden on implementors and operators alike.

This revision proposal aims to address these problems by introducing
[semantic versioning](https://semver.org/) to the set of OpenConfig models as a
whole along with some OpenConfig-specific guidelines. As a side note, these
releases are compatible with the notion of
[YANG release bundles](https://github.com/openconfig/public/blob/master/release/models/catalog/openconfig-module-catalog.yang).

## Policy

### Basic Guidelines

1.  A regular release of https://github.com/openconfig/public containing a set
    of compatible models (consisting of the entire set of models within the
    `openconfig/public` repo) is created roughly every quarter. The tag is named
    vx.y.z (e.g. v1.2.0) following
    [semantic versioning rules](https://semver.org/). A major, minor, or patch
    version increment is possible at each release.
2.  Any backwards-incompatible change for a pre-1.0.0 model does NOT on its own
    necessitate a major revision change for the overall models repository.
3.  [Patch releases](https://semver.org/#spec-item-6) may be created at any time
    for backwards-compatible bug fixes.
4.  [Pre-releases](https://semver.org/#spec-item-9) may be created at anytime to
    quickly introduce new changes to the models. These are not intended to be
    long-term, stable releases -- they should be replaced with the next regular
    release that encompasses these changes as soon as it becomes available.
5.  Backwards-incompatible model changes affecting a feature that has reasonable
    functional test coverage via
    [OpenConfig featureprofiles](https://github.com/openconfig/featureprofiles/))
    or implemented on a device SHOULD be made infrequently. The OpenConfig
    working group will create backwards-incompatible releases periodically by
    considering both velocity and maintenance cost implications.
6.  Wherever possible, it is RECOMMENDED to make backwards compatible API
    changes (e.g. deprecating leaves via the
    [status statement](https://www.rfc-editor.org/rfc/rfc7950#section-7.21.2))
    for at least one minor release prior to a backwards-incompatible API change
    in order to ease the transition to the new API. These leaves are then
    expected to be removed or modified in the next major version release. NOTE:
    This guideline may change once OpenConfig operators gain more experience
    managing breaking changes.

Each release vx.y.z (e.g. v1.2.0) MAY be given a name for easier human
identification, e.g. "September 2022". A client can thus assert they are
compatible with the "September 2022" release of OpenConfig. It's expected that
vendors will have some deviations and augments from the baseline; further, some
vendors may offer the ability to configure their NOS (network operating system)
to support different releases of OC.

### Corner-Case Guidelines

For backwards-incompatible changes involving changing the type of a leaf, the
new leaf SHOULD have a different name than the previous leaf.
