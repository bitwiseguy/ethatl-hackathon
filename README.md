# ethatl-hackathon

This is a submission for the EthAtl [Baseledger Beacon](https://docs.provide.services/hackathon/developers/bounties/baseledger-beacon) hackathon bounty.

If you run the steps in the `Quickstart`, you should be able to produce a random number produced by the Chainlink VRF oracle. There is a lot of overhead considering the output is a single random number. In some high-value situations it may be worth effort if you need guaranteed randomness. For example, when running a multiparty setup for a zero-knowledge circuit you may need a random number to use as a seed to help protect against fraudulent proofs.

As part of the development process, we deployed a VRF consumer on Rinkeby at 0x18E07922265D22a4e71401534F1AD2e406a32C78. This is the contract that gets targeted by the random number generation request. This contract then interacts with the VRFCoordinator contract, which the Chainlink nodes are subscribed to.

Our script sets up on-chain event listeners to the VRF consumer contract and prints out results for triggered events when a random number request is made and when the Chainlink VRF contract process the request by returning a new random number.

## Prerequisites

- Golang v1.16

## Quickstart

- `go mod tidy`
- `go run .`
