#!/bin/bash
go build -o terraform-provider-cypher
export OS_ARCH="$(go env GOHOSTOS)_$(go env GOHOSTARCH)"
mkdir -p ~/.terraform.d/plugins/hashicorp.com/edu/cypher/0.2/$OS_ARCH
mv terraform-provider-cypher ~/.terraform.d/plugins/hashicorp.com/edu/cypher/0.2/$OS_ARCH