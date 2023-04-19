-include .env
export $(shell sed 's/=.*//' .env)

indexer:
	cd cmd/indexer && go run . -c ../../build/dipdup.yml

tester:
	cd cmd/tester && go run . -c ../../build/dipdup.yml

starknet_id:
	cd cmd/rpc_tester && go run . -c ../../cmd/rpc_tester/starknet_id.yml

loot_survivor:
	cd cmd/rpc_tester && go run . -c ../../cmd/rpc_tester/loot_survivor.yml

build-proto:
	protoc \
		-I=${GOPATH}/src \
		--doc_out=${GOPATH}/src/github.com/dipdup-io/starknet-indexer/pkg/grpc \
		--doc_opt=markdown,README.md \
		--go-grpc_out=${GOPATH}/src \
		--go_out=${GOPATH}/src \
		${GOPATH}/src/github.com/dipdup-io/starknet-indexer/pkg/grpc/proto/*.proto

build:
	docker-compose up -d -- build

lint:
	golangci-lint run

test:
	go test ./...