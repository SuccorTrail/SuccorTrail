package web3

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	defaultChainID = 1 // Mainnet
)

var succorTrailABI = `[CONTRACT_ABI]`

func main() {
	client, err := connectToEthereum()
	if err != nil {
		log.Fatalf("Error connecting to Ethereum: %v", err)
	}

	address := common.HexToAddress(getEnv("CONTRACT_ADDRESS", "0xYourContractAddress"))
	instance, err := bindContract(address, succorTrailABI, client)
	if err != nil {
		log.Fatalf("Error binding contract: %v", err)
	}

	if err := callReadFunction(instance); err != nil {
		log.Printf("Error calling read function: %v", err)
	}

	if err := sendTransaction(instance, client); err != nil {
		log.Printf("Error sending transaction: %v", err)
	}
}

func connectToEthereum() (*ethclient.Client, error) {
	infuraURL := getEnv("INFURA_URL", "https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID")
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Ethereum client: %w", err)
	}
	return client, nil
}

func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func bindContract(address common.Address, abiString string, client *ethclient.Client) (*bind.BoundContract, error) {
	contractABI, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %w", err)
	}
	return bind.NewBoundContract(address, contractABI, client, client, client), nil
}

func callReadFunction(instance *bind.BoundContract) error {
	var result []interface{}

	// Call the smart contract function and store the result in result
	err := instance.Call(&bind.CallOpts{}, &result, "message")
	if err != nil {
		return fmt.Errorf("failed to retrieve message: %w", err)
	}

	// Assuming the result is a string, convert the first element to string
	message := result[0].(string)
	fmt.Println("Message from smart contract:", message)
	return nil
}

func sendTransaction(instance *bind.BoundContract, client *ethclient.Client) error {
	privateKeyHex := getEnv("PRIVATE_KEY", "YOUR_PRIVATE_KEY")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to load private key: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(defaultChainID))
	if err != nil {
		return fmt.Errorf("failed to create authorized transactor: %w", err)
	}

	auth.GasLimit = uint64(300000)          // Set a gas limit
	auth.GasPrice = big.NewInt(20000000000) // Set a gas price (20 Gwei)

	tx, err := instance.Transact(auth, "setMessage", "Hello, Blockchain!")
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	fmt.Printf("Transaction sent: %s\n", tx.Hash().Hex())

	ctx := context.Background()
	txReceipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return fmt.Errorf("failed to retrieve transaction receipt: %w", err)
	}
	fmt.Printf("Transaction mined in block: %d\n", txReceipt.BlockNumber.Uint64())
	return nil
}
