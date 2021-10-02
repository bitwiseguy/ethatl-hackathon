package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/consensys/quorum/accounts/abi"
	"github.com/consensys/quorum/accounts/abi/bind"
	"github.com/consensys/quorum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var vrfContractAddress_rinkeby = "0xb3dCcb4Cf7a26f6cf6B120Cf5A73875B7BBc655B"
var conn *ethclient.Client

var vrfConsumerABI = `[
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
	}
]`

// bindToken binds a generic wrapper to an already deployed contract.
func bindABI(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(vrfConsumerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
}

// NewTokenTransactor creates a new write-only instance of Token, bound to a specific deployed contract.
func NewTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenTransactor, error) {
	contract, err := bindABI(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &TokenTransactor{contract: contract}, nil
}


func web3Connect() (error) {
	// Create an IPC based RPC connection to a remote node
	conn, err := ethclient.Dial("https://rinkeby.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	fmt.Println("Successful connection:", conn)
	return err
}

func fetchEntropy() {
	// TODO: the work... query ethereum, chainlink network, etc....
}

func deployVRFConsumerContract() {

}

func main() {
	err := web3Connect()
	deployVRFConsumerContract()
	fetchEntropy()
}
