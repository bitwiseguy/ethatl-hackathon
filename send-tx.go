package main

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

var randomNumberContractAddress = "0x18E07922265D22a4e71401534F1AD2e406a32C78"
var vrfConsumerABI = `[
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "bytes32",
				"name": "",
				"type": "bytes32"
			},
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
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "bytes32",
				"name": "",
				"type": "bytes32"
			}
		],
		"name": "NewRequestId",
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

func setupEventListeners(client ethclient.Client) {
  query := ethereum.FilterQuery{
    Addresses: []common.Address{common.HexToAddress(randomNumberContractAddress)},
  }

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
  if err != nil {
    log.Fatal(err)
  }

	newRandomNumberEventSig := crypto.Keccak256Hash([]byte("NewRandomNumber(bytes32,uint256)"))
	newRequestIdEventSig    := crypto.Keccak256Hash([]byte("NewRequestId(bytes32)"))
	
	contractAbi, err := abi.JSON(strings.NewReader(vrfConsumerABI))
	
	for {
		select {
		case err := <-sub.Err():
			log.Fatal("Error with event subscription:", err)
		case vLog := <-logs:
			log.Println("Event received. Parsing logs...")

			if vLog.Topics[0] == newRandomNumberEventSig {
			  log.Println("NewRandomNumber event detected.")
				event, err := contractAbi.Unpack("NewRandomNumber", vLog.Data)
        if err != nil {
          log.Fatal("Error unpacking event log:", err)
        }
			  log.Printf("Request ID: %+v\n", event[0])
			  log.Println("Random number result:", event[1])
			} else if vLog.Topics[0] == newRequestIdEventSig {
			  log.Println("NewRequestId event detected.")
			  event, err := contractAbi.Unpack("NewRequestId", vLog.Data)
        if err != nil {
          log.Fatal("Error unpacking event log:", err)
        }
			  log.Printf("Request ID: %+v\n", event[0])
		  }
		}
	}
}

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
	log.Println("Sender's public address:", fromAddress)
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
		log.Println("Error sending tx:", err);
	} else {
		log.Println("Sent tx:", signedTx.Hash());
	}
}
