#!/bin/bash

source ../etc/meth/config.sh

warnPeer() {
  echo "Misconfigured config.sh, please check peer settings"
}

if [ -z $PEER_DATADIR ]; then
  warnPeer
fi

HOME=`eval echo "~$USER"`
PEER_DATADIR=$HOME/$PEER_DATADIR

rm -rf $PEER_DATADIR
./gmeth --datadir $PEER_DATADIR init ../etc/meth/genesis.json
