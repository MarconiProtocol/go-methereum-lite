.PHONY: build cluster test

default: build

build:
	go build -o bin/swarm    git.marconi.org/marconiprotocol/go-methereum-lite/cmd/swarm
	go build -o bin/geth     git.marconi.org/marconiprotocol/go-methereum-lite/cmd/geth
	go build -o bin/bootnode git.marconi.org/marconiprotocol/go-methereum-lite/cmd/bootnode

cluster: build
	scripts/boot-cluster.sh

test:
	go test -v git.marconi.org/marconiprotocol/go-methereum-lite/swarm/...
