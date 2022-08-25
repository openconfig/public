# Contributions to OpenConfig
*Contributors*: robjs<sup>†</sup>, aashaikh<sup>†</sup>, chris_luke<sup>⸸</sup>  
† @google.com, ⸸ @comcast.com  
Created: May 2018  
Updated: August 2022

## Rationale
As the OpenConfig project matures and is adopted by implementors and network
operators, an increasing number of vendors and open source projects have
augmented the models with additional fields that have not fallen into the
initial "minimum operationally viable" approach that OpenConfig has taken.
However, this is not to say that these features are not used, or are not of
interest for including in the models. It is desirable that these augmentations
- where they meet the style, and support requirements for including in the
models - be included within the standard set of models, rather than have many
different augmentations covering the same feature set. To this end, this
document outlines a process for external contributions to OpenConfig.

OpenConfig considers schema consistency paramount -- since it greatly impacts
the usability of the models for consumers be they machine or human. As external
contributions are accepted into the OpenConfig working group, the intent is to
maintain the vendor-neutral, operationally-focused modelling approach taken
thus far. In order that these requirements are met, external model changes will
still be reviewed and approved by the core operator group before being merged
into the schema. It should be noted that it is explicitly required that two
vendors have implemented a particular feature before a feature is considered
for inclusion in OpenConfig -- features that do not meet this requirement
should continue to use vendor-specific augmentations to the schema.

The process for making a contribution is outlined below.

## Contributing to OpenConfig

OpenConfig prefers code (i.e., YANG) contributions, rather than feature
requests. If you wish to discuss the suitability or approach for a change, or
addition to the models, this can be done with an issue in the [OpenConfig
public GitHub](https://github.com/openconfig/public/issues).

All contributions to OpenConfig MUST be Apache 2.0 licensed. A contributor
license agreement (CLA), namely the [Google
CLA](https://cla.developers.google.com/), MUST be signed for any contribution
to be accepted.

The CLA is used to ensure that the rights to use the contribution are well
understood by the OpenConfig working group, and consumers of the OpenConfig
models. Since copyright over each contribution is assigned to its authors, code
comments or the `description` field of a YANG model should reflect the
contribution made, and the copyright holder. No code will be reviewed if the
license is not explicitly stated, or the CLA has not been signed.

Note that we use the Google CLA because the OpenConfig project is [administered
and maintained by Google](https://opensource.google.com/docs/cla/#why), not to
ascribe any specific rights to a single OpenConfig member.

To make a contribution to OpenConfig:

1. Open a pull request in the
 [openconfig/public](https://github.com/openconfig/public) repo. The pull
  request template for the repository details the information that is expected,
  please fill it out, along with any additional information that is useful for
  reviewers. In addition:
    * Pull requests should be kept small. An ideal change is less than 500 lines
     of YANG. Small changes allow detailed discussions of the additions that are
     being made to the model, whilst also ensuring that course-corrections can be
     made early in the process. In some cases, changes larger than 500 lines may
     be unavoidable - these should be rare, and generally only be the case when
     entirely new modules are being added to the model. In this case, it is very
     likely an issue should have been created to discuss the addition prior to
     code review.
    * When the pull request adds a new feature that is supported across vendors,
     the author must include links to public-facing documentation showing
     the implementation of the feature within the change description. This
     simplifies the process of reviewing differences and the chosen abstractions
     (if any are used).

1. The pull request should include both the suggested YANG additions, as well
 as any relevant changes to the `.spec.yml` files that are included within the
 repo. In general, where new content is being added to existing files - no
 change to the build specification should be required. In the case that new
 files are being added, new files should be added to the relevant stanza of an
 existing rule, or entirely new stanzas should be added.

1. The automated CI running against each pull request will lint the pull
 request against the OpenConfig [style guide
 rules](https://github.com/openconfig/public/blob/master/doc/openconfig_style_guide.md).
 An example of the output of the linter can be found
 [here](https://gist.github.com/OpenConfigBot/139f5263ec20957124c7d05edc2c79ff).
 The list of style rules that are enforced are found in this [pyang
 plugin](https://github.com/openconfig/oc-pyang/blob/master/openconfig_pyang/plugins/openconfig.py).
 The linter can be run locally by installing the `openconfig_pyang` package from
 PyPi. Automated CI also runs a small number of integration tests with publicly
 available YANG toolchains, in order to detect regression issues that may occur
 due to OpenConfig model changes.

1. Discussion of the PR is carried out in the `openconfig/public` repository -
 in order to ensure that different viewpoints can be considered from the
 community. Real-time discussions (either scheduled or ad-hoc) can be arranged
 where needed.

1. When the model changes are approved. The pull request will be
 merged in the public repository.

The aim of this process is not to restrict contributions to OpenConfig, but
simply to maintain the model quality and approach that the working group has
strived for since its inception in 2014. Questions prior to making submissions
are welcome, please use the [netopenconfig Google
group](mailto:netopenconfig@googlegroups.com), or the [public repository
issues](https://github.com/openconfig/public/issues).
