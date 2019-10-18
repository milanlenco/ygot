#!/bin/bash

# Hack to ensure that if we are running on OS X with a homebrew installed
# GNU sed then we can still run sed.
runsed() {
  if hash gsed 2>/dev/null; then
    gsed "$@"
  else
    sed "$@"
  fi
}

git clone https://github.com/openconfig/public.git
mkdir deps
cp ../demo/getting_started/yang/{ietf,iana}* deps
go run ../ypathgen/generator/generator.go -path=public,deps -output_file=ocpath.go \
  -package_name=exampleoc -fakeroot_name=device \
  -exclude_modules=ietf-interfaces \
  -schema_struct_path=github.com/openconfig/ygot/exampleoc \
  public/release/models/network-instance/openconfig-network-instance.yang \
  public/release/models/optical-transport/openconfig-optical-amplifier.yang \
  public/release/models/optical-transport/openconfig-terminal-device.yang \
  public/release/models/optical-transport/openconfig-transport-line-protection.yang \
  public/release/models/platform/openconfig-platform.yang \
  public/release/models/policy/openconfig-routing-policy.yang \
  public/release/models/lacp/openconfig-lacp.yang \
  public/release/models/system/openconfig-system.yang \
  public/release/models/lldp/openconfig-lldp.yang \
  public/release/models/stp/openconfig-spanning-tree.yang \
  public/release/models/interfaces/openconfig-interfaces.yang \
  public/release/models/interfaces/openconfig-if-ip.yang \
  public/release/models/interfaces/openconfig-if-aggregate.yang \
  public/release/models/interfaces/openconfig-if-ethernet.yang \
  public/release/models/interfaces/openconfig-if-ip-ext.yang \
  public/release/models/relay-agent/openconfig-relay-agent.yang
runsed -i 's/This package was generated by.*/NOTE WELL: This is an example code file that is distributed with ygot.\nIt should not be used within your application, as it WILL change,\nwithout warning. Rather, you should generate structs directly from\nOpenConfig models using the ygot package.\n\nThis package was generated by github.com\/openconfig\/ygot/g' ocpath.go
gofmt -w -s ocpath.go
rm -rf deps public
