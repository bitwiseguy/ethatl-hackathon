package main

import (
	"log"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func web3Connect() (*ethclient.Client, error) {
	godotenv.Load(".env")
	ethClient_ws := os.Getenv("ETH_CLIENT_WS")
	
	// Create an IPC based RPC connection to a remote node
	conn, err := ethclient.Dial(ethClient_ws)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	log.Println("Successful connection to eth client!")
	return conn, err
}

func main() {
	conn, err := web3Connect()
	if err != nil {
		log.Println("Failed to connect to Eth Client")
	}

	go setupEventListeners(*conn)
	fetchEntropy(*conn, randomNumberContractAddress)
	for true {}
}
