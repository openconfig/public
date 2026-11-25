# OpenConfig Model Review Style Guide & Guidelines

This guide defines the review standards and authoring conventions for OpenConfig YANG models in `openconfig/public`. It instructs automated reviewers (Gemini Code Assist for GitHub) and human reviewers to ensure vendor neutrality, schema consistency, semantic versioning rigor, and automated tooling compatibility (`ygot`, gNMI).

---

### **1. General Coding & Contribution Guidelines**

**Sources:**
- `doc/openconfig_style_guide.md`
- `doc/contributions-guide.md`
- `doc/semver.md`

All contributions must follow the [OpenConfig style guide](doc/openconfig_style_guide.md) and [contributions guide](doc/contributions-guide.md).

---

### **2. Multi-Vendor Implementation Feasibility (Mandatory for New Features)**

**Source:** `doc/contributions-guide.md`

OpenConfig models are operator-driven and strictly vendor-neutral. Every pull request adding or altering configuration/state models must demonstrate feasibility across industry network operating systems:
- **Target OSes:** Cisco IOS-XR, Nokia SR Linux, Juniper JunOS / JunOS Evolved, Arista EOS, Ciena SAOS, NVIDIA Cumulus Linux, and SONiC.
- **Review Requirement:**
  - Generate a comparison table in the review output comparing the affected OpenConfig leaves/containers with the **exact, concrete command-line (CLI) configuration syntax** in the target network operating systems.
  - **Hierarchy of Documentation Links:**
    1. **Primary: Public, Unauthenticated Configuration Guides & User Manuals:**
       - All links must be openly accessible without login credentials, paid subscriptions, or gated enterprise accounts.
       - **Arista:** Link to the public User Manual under `https://www.arista.com/en/um-eos/...` (e.g. `https://www.arista.com/en/um-eos/eos-bgp-pic-edge-for-evpn-vxlan-routes-for-remote-vtep-failures`). **STRICTLY FORBIDDEN:** Linking to Arista TOIs (`/support/toi/`) or TAC customer portals that prompt for a login or return "Access Denied".
       - **Juniper / HPE:** Link to public JunOS documentation under `https://networkingdocs.hpe.com/documentation/...` or `https://www.juniper.net/documentation/...` (e.g. `https://networkingdocs.hpe.com/documentation/us/en/software/junos/vpn-l3/topics/topic-map/l3-vpns-bgp-pic.html`).
       - **Cisco:** Link to public Command References or Configuration Guides under `https://www.cisco.com/c/en/us/td/docs/...` (never CCO login-gated tools).
       - **Nokia SR Linux:** Link to public documentation under `https://learn.nokia.com/srlinux` or official public doc portals.
       - **SONiC:** Link directly to public feature design docs or command guides under `https://github.com/sonic-net/SONiC/blob/master/doc/...` (never the `/wiki` homepage).
    2. **Fallback: Vendor Native YANG Models (`YangModels/yang`):**
       - If a public, non-paywalled CLI doc page cannot be linked or verified, fall back directly to the vendor's native YANG models in [github.com/YangModels/yang/tree/main/vendor](https://github.com/YangModels/yang/tree/main/vendor) (linking to the exact file and line number on GitHub).
    3. **Strict Prohibition on Generic Doc Roots & Paywalls:**
       - **Forbidden:** Linking to high-level doc homepages (e.g., `https://www.juniper.net/documentation/`, `https://learn.nokia.com`, `https://www.arista.com/support`, `https://www.cisco.com`).
       - If a feature is not supported by a platform, explicitly mark the column as `Unsupported / No Match` and do NOT provide a generic homepage URL fallback.
  - **Match Definition & Truthful Evaluation:**
    - A feature match requires that the target NOS provides explicit configuration knobs (CLI or native YANG) or operational telemetry matching the proposed OpenConfig schema semantics.
    - **No Forced Matches or Fabricated CLI:** If a target NOS does not expose a dedicated CLI configuration command or YANG container for the feature (e.g., SONiC implements BGP PIC as an underlying dataplane/SAI architecture without dedicated CLI configuration commands), the reviewer MUST explicitly mark it as `No Match / Unsupported via CLI` and state that no CLI configuration knob exists.
    - **Never invent or guess CLI syntax.** OpenConfig vendor-neutrality requires at least **two** independent implementations from different vendors (e.g., Cisco + Arista + Juniper is already 3 strong matches). Accurately recording `No Match` for platforms lacking CLI knobs is essential for high-signal reviews.
  - **Neutrality Threshold:** At least **two** separate implementations from different organizations must have a strong match with the YANG changes to satisfy OpenConfig vendor-neutrality criteria. If fewer than two have a match, provide explicit feedback and suggested modifications.

#### Required Table Format:
| Network OS | Exact CLI Command Syntax / Hierarchy | Exact Command Reference URL or Native YANG Path | Match Strength |
| :--- | :--- | :--- | :--- |
| **Cisco IOS-XR** | `telemetry model-driven \n destination-group <id> \n address ...` | `[Cisco telemetry cr](https://www.cisco.com/c/en/us/td/docs/routers/asr9000/software/asr9k-r7-3/telemetry/command/reference/b-telemetry-cr-asr9000-73x/telemetry-commands.html)` or `YangModels/yang:vendor/cisco/xr/...` | Strong |
| **Juniper JunOS** | `set protocols bgp family inet unicast extended-nexthop` | `[JunOS BGP PIC Guide](https://networkingdocs.hpe.com/documentation/us/en/software/junos/vpn-l3/topics/topic-map/l3-vpns-bgp-pic.html)` or `YangModels/yang:vendor/juniper/...` | Strong |
| **Arista EOS** | `bgp bestpath prefix-independent-convergence` | `[Arista User Manual](https://www.arista.com/en/um-eos/eos-bgp-pic-edge-for-evpn-vxlan-routes-for-remote-vtep-failures)` or `YangModels/yang:vendor/arista/...` | Strong |
| **Nokia SR Linux** | `routing-policy prefix-independent-convergence` | `[Nokia SR Linux Guide](https://learn.nokia.com/srlinux)` or `YangModels/yang:vendor/nokia/...` | Strong |
| **SONiC** | None (Dataplane/SAI architecture only; no dedicated CLI configuration knob) | `[SONiC BGP PIC Arch Doc](https://github.com/sonic-net/SONiC/blob/master/doc/pic/bgp_pic_arch_doc.md)` | No Match |

---

### **3. Semantic Versioning & Deprecation Policy**

**Source:** `doc/semver.md` and `openconfig-extensions.yang`

Every module and submodule must strictly adhere to Semantic Versioning (`oc-ext:openconfig-version`):
- **Deprecate Before Breaking:**
  - Avoid removing or breaking published leaves or containers in place.
  - Instead of replacing a leaf (e.g., changing scalar `leaf foo` into `leaf-list foo`), mark the existing leaf as `status deprecated;` with a descriptive deprecation note and introduce a new leaf/leaf-list alongside it (e.g., `leaf-list foo-types`).
  - If a leaf is deprecated, **retain** its existing default and type statements so existing controllers and clients do not break.
- **Version Bump Rules:**
  - **Major Bump (`X+1.0.0` or `0.Y+1.0` if `0.x`):** Backwards-incompatible / breaking changes (deletion of nodes, changing leaf types, narrowing value ranges, adding mandatory leaves, changing list keys).
  - **Minor Bump (`X.Y+1.0`):** Backwards-compatible additions (adding optional leaves, containers, lists, new identities, or deprecating existing nodes while adding replacements).
  - **Patch Bump (`X.Y.Z+1`):** Editorial changes, documentation clarifications, typo corrections, reference updates.
- **Submodule Version Parity:**
  - All submodules (`openconfig-*-<subfunction>.yang`) must share the exact same `openconfig-version` string as their parent module.
- **Revision Block:**
  - Add a new `revision "YYYY-MM-DD"` entry at the top of the revision list with a matching `reference "X.Y.Z";`.
  - The revision date must be the current date (never past or future).

---

### **4. Operational State & Grouping Architecture (`config` vs `state`)**

**Source:** `doc/openconfig_style_guide.md#modeling-operational-state`

- **Strict Container Separation:**
  - Every functional subtree must use explicit `container config` (read/write configuration) and `container state` (operational state, marked with `config false;`).
- **Grouping Structure:**
  - `xxx-yyy-config`: Defines configurable leaves.
  - `xxx-yyy-state`: Defines telemetry / read-only leaves (counters, timestamps, negotiated values, runtime status).
  - `xxx-yyy-top`: Top-level enclosing grouping that instantiates `container config` (`uses xxx-yyy-config;`) and `container state` (`uses xxx-yyy-config;` AND `uses xxx-yyy-state;`).
- **Pure Telemetry Isolation:**
  - Counters, packet statistics, transit errors, and dynamic state belong **exclusively** in `state` (and `xxx-yyy-state` groupings). They must **never** appear in `config`.
- **Top-Level Definitions:**
  - Avoid directly defined data nodes at the top level of a module. All nodes must be wrapped inside reusable groupings and instantiated once via `uses`.
- **No Defaults on State Leaves:**
  - Never specify a `default` statement on leaves that exist exclusively in `state`.

---

### **5. Schema Usability & Compression Rules (Critical for `ygot` / gNMI Tooling)**

**Source:** `doc/openconfig_style_guide.md#schema-usability`

To ensure programmatic schema compression and code generation (e.g. `ygot`, gNMI client bindings) function without identifier collisions:
- **No Sibling Container Collision:** A `list` node MUST NOT share its identifier with any sibling of its parent container (e.g. `/interfaces/interface` and `/interfaces/interface-common` is illegal).
- **No Grandparent Collision:** A leaf node MUST NOT share its identifier with any child of its grandparent node (e.g. `/a/config/leaf` and `/a/leaf` is illegal).
  *Exception:* List key leaves where the `key` leaf inside `config` or `state` shares the name with the direct child `key` leafref node.
- **List Enclosing Containers:** Every `list` must be enclosed within a surrounding plural container with no other sibling nodes in that container (e.g. `container interfaces { list interface { ... } }`).
- **List Keys:**
  - List keys must be direct children of `config` or `state` (path `../config/<key>` or `../state/<key>`). Referencing deeper descendant paths (e.g., `../state/counters/<key>`) is strictly illegal.
  - List keys in YANG must be quoted: `key "name";`.
  - Lists without keys are forbidden unless the enclosing container has the `oc-ext:telemetry-atomic;` extension.

---

### **6. YANG Language Statement Restrictions**

**Source:** `doc/openconfig_style_guide.md#yang-language-usage`

- **`default`:** Avoid `default` statements unless explicitly mandated by an authoritative RFC or IEEE specification. In vendor-neutral modeling, omit the default keyword and describe default system behavior in the `description` string instead.
- **`presence`:** Forbidden. Do not use presence containers. Use an explicit boolean leaf (e.g. `leaf enabled { type boolean; }`) to make activation explicit.
- **`choice` / `case`:** Avoid where possible because choice nodes do not appear in schema paths. Use explicit optional containers or `when` statements instead.
- **`if-feature`:** Forbidden. OpenConfig models express an operationally complete set of features; platform non-support must be handled via vendor deviation files, not `if-feature`.
- **XPath:** Use simple relative paths. Use only basic YANG 1.0 XPath functions (`and`, `current`, `not`, `or`). Do not use YANG 1.1 XPath extensions.
- **Regular Expressions:** Use POSIX Extended Regular Expressions (ERE), not W3C XML Schema regex.

---

### **7. Types, Units & Identities**

- **Mandatory Units:** Any leaf representing time, frequency, rates, bandwidth, memory, packet counts, or intervals MUST have an explicit `units` statement (e.g., `units "seconds";`, `units "bytes";`, `units "packets";`, `units "milliseconds";`).
- **Unsigned Counters:** Counters for packets, bytes, octets, errors, and discards must be unsigned (use `oc-yang:counter64` or `uint64`), **never** signed `int64`.
- **Identities vs Enums:** Prefer `identity` and `identityref` when a set of values is expected to be extensible across vendors or future revisions. Use `enumeration` only for strictly closed, standard sets.
- **Type Reuse:** Reuse standard OpenConfig types from `openconfig-types`, `openconfig-inet-types`, and `openconfig-yang-types` rather than creating redundant typedefs.

---

### **8. Naming, Formatting & Punctuation**

- **Case Conventions:**
  - Identifiers (leaves, containers, lists, groupings): `lower-case-with-hyphens` (e.g., `built-in-limits`, `downstream-on-demand`). No underscores, no camelCase.
  - Identities & Enums: `UPPER_CASE_WITH_UNDERSCORES` (e.g., `IPV4_ENCAP`, `DEFAULT_INSTANCE`).
  - Module filenames: `openconfig-<function>.yang` or `openconfig-<function>-<subfunction>.yang`.
  - Prefixes: `oc-<shortname>`.
- **No Redundant Naming:** Do not repeat the parent container's name in child leaf names (e.g., inside `container icmp`, name the leaf `rate-limit`, not `icmp-rate-limit`).
- **Line Length:** Maximum 70 characters per line per IETF/RFC formatting rules.
- **Description Formatting & Quality:**
  - Use multi-line indented style:
    ```yang
    description
      "First sentence summarizes purpose. Subsequent sentences detail
       operational semantics, units, corner-case behavior, and references.";
    ```
  - Descriptions must be complete sentences ending with a period.
  - Explain *operational semantics*, *calculation method*, behavior when unspecified, and cite relevant RFC or IEEE sections. Avoid trivial or empty descriptions.

---

### **9. Review Output Structure & Severity Rubric**

All automated and human reviews must follow this structured output:

1. **Executive Summary & Verdict:**
   - Summary of changes and overall architecture.
   - Verdict: `[APPROVE]`, `[CHANGES REQUESTED]`, or `[COMMENT / INQUIRY]`.
2. **Issues & Findings Table:**
   - `[Blocker]`: Breaking changes without major version bump; missing `openconfig-version`; presence containers; non-quoted list keys.
   - `[Major]`: Missing config/state separation; telemetry placed in config; signed counter types; missing units; camelCase or underscores in identifiers.
   - `[Nit/Style]`: Description formatting; line length > 70 characters; missing punctuation.
3. **Multi-Vendor Implementation Comparison Table:**
   - Table detailing equivalent commands for Cisco IOS-XR, Juniper JunOS, Arista EOS, Nokia SR Linux, Ciena SAOS, and SONiC with direct hyperlinks to the exact configuration reference page for each command.
4. **Suggested Code / Diff Blocks:**
   - Concrete before & after YANG snippets demonstrating exact required fixes.
