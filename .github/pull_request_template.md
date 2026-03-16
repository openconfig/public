[Note: Please fill out the following template for your pull request. Replace
all the text in `[]` with your own content.]

[Note: Before this PR can be reviewed please agree to the CLA covering this
repo. Please also review the contribution guide -
https://github.com/openconfig/public/blob/master/doc/contributions-guide.md]

### Change Scope

* [Please briefly describe the change that is being made to the models.]
* [Please indicate whether this change is backwards compatible.]

### Platform Implementations

 * Implementation A: [link to documentation](http://foo.com) and/or
   implementation output.
 * Implementation B: [link to documentation](http://foo.com) and/or
   implementation output.

[Note: Please provide at least two references to implementations which are relevant to the model changes proposed.  Each implementation should be from separate organizations.]. 

[Note: If the feature being proposed is new - and something that is being
proposed as an enhancement to device functionality, it is sufficient to have
reviewers from the producers of two different implementations].

### Tree View

* [Please provide a view of the tree being modified.  It's preferred if a `diff` format is used for ease of review.]
* [Here are recommended steps to generate the tree view as a diff]

```
git checkout mychangebranch
pyang -f tree -p release/models/*/* > ~/new-tree.txt 
git checkout master
git pull
pyang -f tree -p release/models/*/* > ~/old-tree.txt
diff -bU 100 ~/old-tree.txt ~/new-tree.txt   | less
```

[Next, cut and paste the relevant portion of the tree with enough context for reviewers to quickly understand the change.]
```diff
 module: openconfig-interfaces
   +--rw interfaces
      +--rw interface* [name]
         +--ro state
         |  +--ro counters
         |  |  +--ro in-octets?               oc-yang:counter64
         |  |  +--ro in-pkts?                 oc-yang:counter64
         |  |  +--ro in-unicast-pkts?         oc-yang:counter64
         |  |  +--ro in-broadcast-pkts?       oc-yang:counter64
         |  |  +--ro in-multicast-pkts?       oc-yang:counter64
         |  |  +--ro in-errors?               oc-yang:counter64
         |  |  +--ro in-discards?             oc-yang:counter64
         |  |  +--ro out-octets?              oc-yang:counter64
         |  |  +--ro out-pkts?                oc-yang:counter64
         |  |  +--ro out-unicast-pkts?        oc-yang:counter64
         |  |  +--ro out-broadcast-pkts?      oc-yang:counter64
         |  |  +--ro out-multicast-pkts?      oc-yang:counter64
         |  |  +--ro out-discards?            oc-yang:counter64
         |  |  +--ro out-errors?              oc-yang:counter64
         |  |  +--ro last-clear?              oc-types:timeticks64
         |  |  +--ro in-unknown-protos?       oc-yang:counter64
         |  |  +--ro in-fcs-errors?           oc-yang:counter64
+        |  |  x--ro carrier-transitions?     oc-yang:counter64
-        |  |  +--ro carrier-transitions?     oc-yang:counter64
+        |  |  +--ro interface-transitions?   oc-yang:counter64
+        |  |  +--ro link-transitions?        oc-yang:counter64
         |  |  +--ro resets?                oc-yang:counter64
```
