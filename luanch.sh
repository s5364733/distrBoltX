#!/bin/bash

set -e

trap 'killall main' SIGINT

cd $(dirname $0)

killall main || true

sleep 0.1

go install -v

main -db-location=shard0.db  -http-addr=127.0.0.1:8080 -config-location=sharding.toml -shard=shard0
main -db-location=shard1.db  -http-addr=127.0.0.1:8081 -config-location=sharding.toml -shard=shard1
main -db-location=shard2.db  -http-addr=127.0.0.1:8082 -config-location=sharding.toml -shard=shard2

wait