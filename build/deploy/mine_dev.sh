#!/bin/bash

ABSOLUTE_MINERDIR=${1}

if [ -z ${ABSOLUTE_MINERDIR} ]; then
  echo "Missing path to mining data directory"
fi

source ../etc/meth/config.sh

if [ -z $BOOTNODE_ENODE_HASH ] || [ -z $BOOTNODE_IP ] || [ -z $BOOTNODE_PORT ]; then
  echo "Misconfigured config.sh, please check bootnode settings"
fi

if [ -z $MINER_PORT ] || [ -z $MINER_ETHERBASE ]; then
  echo "Misconfigured config.sh, please check miner settings"
fi

./gmeth --datadir ${ABSOLUTE_MINERDIR} \
  --port $MINER_PORT \
  --rpc \
  --rpcport $MINER_RPC_PORT \
  --rpcapi="db,eth,net,web3,personal" \
  --networkid 161027 \
  --mine --minerthreads 1 \
  --bootnodes enode://$BOOTNODE_ENODE_HASH@$BOOTNODE_IP:$BOOTNODE_PORT \
  --etherbase $MINER_ETHERBASE \
  --unlock $MINER_ETHERBASE \
  --password "../etc/meth/unlock.txt" \
  --rpcaddr "0.0.0.0" \
  --syncmode "full"
