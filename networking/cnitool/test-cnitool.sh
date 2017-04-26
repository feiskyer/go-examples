#!/bin/bash
set -x
set -o errexit
set -o nounset
set -o pipefail

export CNI_PATH=/opt/cni/bin
export NETCONFPATH=/etc/cni/net.d

CNITOOL_ROOT=$(dirname "${BASH_SOURCE}")
NETNS=test

cleanup() {
    ip netns del $NETNS
}

trap "cleanup" EXIT SIGINT

ip netns add $NETNS
cd $CNITOOL_ROOT
go build -o cnitool

./cnitool add netlist /var/run/netns/test
./cnitool del netlist /var/run/netns/test
