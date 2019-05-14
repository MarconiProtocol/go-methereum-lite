**This repository is forked from https://github.com/ethereum/go-ethereum**
# Go Marconi Ethereum - Lite
This is a fork of [go-ethereum](https://github.com/ethereum/go-ethereum) which we've modified to suit our needs.
Our repo is up to date with go-ethereum at commit [c9427004](https://github.com/ethereum/go-ethereum/commits/c9427004).  
Our goal is to write our version of the Ethereum Virtual Machine (EVM) in the future and release it as a seperate project (go-marconi). 

A high level overview of the biggest changes is as follows:
* Switched the hash function from **ethhash** to **CryptoNightR**
* Updated default chain configuration to target a **512kb block size** with a block time of **30s**.
* Removed mining code, see our rationale in this [blog post](https://medium.com/marconiprotocol/how-were-bootstrapping-marconi-50fb2edaad6e) 


## Quick Links
- [Compiling the Source](/README_geth.md#building-the-source)
- [Running Geth](/README_geth.md#running-geth)

## Explorer
The blockchain data can be viewed using our explorer sites:  
Main Net: [explorer.marconi.org](https://explorer.marconi.org)  
Test Net: [explorer.testnet.marconi.org](https://explorer.testnet.marconi.org/)  

## Related Repos
- [explorer](https://git.marconi.org/marconiprotocol/explorer)
- [xmr-stak](https://git.marconi.org/marconiprotocol/xmr-stak)
- [marconi-cryptonight](https://git.marconi.org/marconiprotocol/marconi-cryptonight)


