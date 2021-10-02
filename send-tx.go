package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

var vrfConsumerABI = `[
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"name": "NewRandomNumber",
		"type": "event"
	},
	{
		"inputs": [],
		"name": "getRandomNumber",
		"outputs": [
			{
				"internalType": "bytes32",
				"name": "requestId",
				"type": "bytes32"
			}
		],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "bytes32",
				"name": "requestId",
				"type": "bytes32"
			},
			{
				"internalType": "uint256",
				"name": "randomness",
				"type": "uint256"
			}
		],
		"name": "rawFulfillRandomness",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"inputs": [],
		"name": "randomResult",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	}
]`

func fetchEntropy(client ethclient.Client, contractAddress string) {
	godotenv.Load(".env")
	privKey := os.Getenv("WALLET_PRIVATE_KEY")
  
	d := time.Now().Add(5000 * time.Millisecond)
  ctx, cancel := context.WithDeadline(context.Background(), d)
  defer cancel()

  contractAbi, err := abi.JSON(strings.NewReader(vrfConsumerABI))

	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
    	log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
    	log.Fatal("error casting public key to ECDSA")
	}
    
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("Sender's public address:", fromAddress)
  nonce, err := client.NonceAt(ctx, fromAddress, nil)
  if err != nil {
    log.Fatal(err)
  }

  gasPrice, err := client.SuggestGasPrice(context.Background())
  if err != nil {
      log.Fatal(err)
  }

	bytesData, _ := contractAbi.Pack("getRandomNumber")
	tx := types.NewTransaction(nonce, common.HexToAddress(contractAddress), nil, 10000000, gasPrice, bytesData)
	signedTx, _ := types.SignTx(tx, types.LatestSignerForChainID(big.NewInt(4)), privateKey)
	err = client.SendTransaction(ctx, signedTx)

	if err != nil {
		fmt.Println("Error sending tx:", err);
	} else {
		fmt.Println("Sent tx:", signedTx.Hash());
	}
}
