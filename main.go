package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

var randomNumberContractAddress = "0xe7D85ad235B9C9E86E9904426a7E1d4F303c3aB3"

func web3Connect() (*ethclient.Client, error) {
	// Create an IPC based RPC connection to a remote node
	conn, err := ethclient.Dial("https://rinkeby.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	fmt.Println("Successful connection to eth client!")
	return conn, err
}

func setupEventListeners() {

}

func main() {
	conn, err := web3Connect()
	if err != nil {
		fmt.Println("Failed to connect to Eth Client")
	}

	setupEventListeners()
	fetchEntropy(*conn, randomNumberContractAddress)
}
