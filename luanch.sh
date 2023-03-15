#!/bin/bash

set -e

trap 'killall main' SIGINT

cd $(dirname $0)

killall main || true

sleep 0.1

go install -v

main -db-location=shard0.db  -http-addr=127.0.0.2:8080 -config-file=sharding.toml -shard=shard0
main -db-location=shard0-r.db -http-addr=127.0.0.22:8080 -config-file=sharding.toml -shard=shard0 -replica=true &

main -db-location=shard1.db -http-addr=127.0.0.3:8080 -config-file=sharding.toml -shard=shard1 &
main -db-location=shard1-r.db -http-addr=127.0.0.33:8080 -config-file=sharding.toml -shard=shard1 -replica=true &

main -db-location=shard2.db -http-addr=127.0.0.4:8080 -config-file=sharding.toml -shard=shard2 &
main -db-location=shard2-r.db -http-addr=127.0.0.44:8080 -config-file=sharding.toml -shard=shard2 -replica=true &

main -db-location=shard3.db -http-addr=127.0.0.5:8080 -config-file=sharding.toml -shard=shard3 &
main -db-location=shard3-r.db -http-addr=127.0.0.55:8080 -config-file=sharding.toml -shard=shard3 -replica=true &

wait