#!/bin/bash

go run $GOPATH/src/github.com/openconfig/ygot/generator/generator.go -path=../../../../.. -output_file=oc.go \
  -package_name=telemetry -generate_fakeroot -fakeroot_name=device -compress_paths=true \
  -shorten_enum_leaf_names \
  -trim_enum_openconfig_prefix \
  -typedef_enum_with_defmod \
  -enum_suffix_for_simple_union_enums \
  -exclude_modules=ietf-interfaces \
  -generate_rename \
  -generate_append \
  -generate_getters \
  -generate_leaf_getters \
  -generate_simple_unions \
  -annotations \
  -list_builder_key_threshold=3 \
  ../../../network-instance/openconfig-network-instance.yang \
  ../../../interfaces/openconfig-interfaces.yang \
  ../../../interfaces/openconfig-if-ip.yang \
  ../../../interfaces/openconfig-if-aggregate.yang \
  ../../../interfaces/openconfig-if-ethernet.yang \
  ../../../interfaces/openconfig-if-ip-ext.yang \
  ../../../qos/openconfig-qos.yang \
  ../../../platform/openconfig-platform-cpu.yang \
  ../../../lacp/openconfig-lacp.yang \
  ../../openconfig-ate-flow.yang \
  ../../openconfig-ate-intf.yang
gofmt -w -s oc.go
