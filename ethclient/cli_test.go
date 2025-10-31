package ethclient

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const ConfirmDataStoreEventABI = "ConfirmDataStore(uint32,bytes32)"

var ConfirmDataStoreEventABIHash = crypto.Keccak256Hash([]byte(ConfirmDataStoreEventABI)) // 0xfbb7f4f1b0b9ad9e75d69d22c364e13089418d86fcb5106792a53046c0fb33aa

const DataLayerServiceManagerAddr = "0x5BD63a7ECc13b955C4F57e3F12A64c10263C14c1"

func TestEthClient_GetTxReceiptByHash(t *testing.T) {
	fmt.Println("test start for tx receipt")
	clint, err := NewEthClient("https://endpoints.omniatech.io/v1/eth/mainnet/public")
	if err != nil {
		fmt.Println("New eth client fail", err)
	}
	txReceipt, err := clint.GetTxReceiptByHash("0xfd26d40e17213bcafcf94bab9af92343302df9df970f20e1c9d515525e86e23e")
	fmt.Printf("Transaction Status: %d\n", txReceipt.Status)
	fmt.Printf("Total logs: %d\n", len(txReceipt.Logs))
	fmt.Printf("DataLayerServiceManagerAddr: %s\n", DataLayerServiceManagerAddr)

	if err != nil {
		fmt.Println("get tx receipt fail", err)
	}

	abiUint32, err := abi.NewType("uint32", "uint32", nil)
	if err != nil {
		fmt.Println("new uint32 abi type fail", err)
	}

	abiBytes32, err := abi.NewType("bytes32", "bytes32", nil)
	if err != nil {
		fmt.Println("new uint32 abi type fail", err)
	}
	confirmDataStoreArgs := abi.Arguments{
		{
			Name:    "dataStoreId",
			Type:    abiUint32,
			Indexed: false,
		}, {
			Name:    "headerHash",
			Type:    abiBytes32,
			Indexed: false,
		},
	}
	var dataStoreData = make(map[string]interface{})
	for _, rLog := range txReceipt.Logs {
		fmt.Println("address====", rLog.Address.String())
		if strings.ToLower(rLog.Address.String()) != strings.ToLower(DataLayerServiceManagerAddr) {
			continue
		}
		if rLog.Topics[0] != ConfirmDataStoreEventABIHash {
			continue
		}
		if len(rLog.Data) > 0 {
			err = confirmDataStoreArgs.UnpackIntoMap(dataStoreData, rLog.Data)
			if err != nil {
				fmt.Println("unpack data into mapping fail", err)
				continue
			}

			if dataStoreData != nil {
				fmt.Printf("dataStoreId: %v\n", dataStoreData["dataStoreId"])
				fmt.Printf("headerHash: %x\n", dataStoreData["headerHash"])
			}
		}
	}

	//found := false
	//for i, rLog := range txReceipt.Logs {
	//	fmt.Printf("\n--- Log %d ---\n", i)
	//	fmt.Printf("Address: %s\n", rLog.Address.String())
	//	fmt.Printf("Topics count: %d\n", len(rLog.Topics))
	//	fmt.Printf("ConfirmDataStoreEventABIHash: %x\n", ConfirmDataStoreEventABIHash)
	//	for j, topic := range rLog.Topics {
	//		fmt.Printf("Topic[%d]: %s\n", j, topic.Hex())
	//	}
	//
	//	addressMatch := strings.ToLower(rLog.Address.String()) == strings.ToLower(DataLayerServiceManagerAddr)
	//	topicMatch := len(rLog.Topics) > 0 && rLog.Topics[0] == ConfirmDataStoreEventABIHash
	//
	//	fmt.Printf("Address matches: %v\n", addressMatch)
	//	fmt.Printf("Topic matches: %v\n", topicMatch)
	//
	//	if !addressMatch {
	//		fmt.Println("Skipping - address doesn't match")
	//		continue
	//	}
	//	if !topicMatch {
	//		fmt.Println("Skipping - topic doesn't match")
	//		continue
	//	}
	//
	//	fmt.Println("*** MATCH FOUND! ***")
	//	if len(rLog.Data) > 0 {
	//		var dataStoreData = make(map[string]interface{})
	//		err = confirmDataStoreArgs.UnpackIntoMap(dataStoreData, rLog.Data)
	//		if err != nil {
	//			fmt.Println("unpack data into mapping fail:", err)
	//			continue
	//		}
	//
	//		if dataStoreData != nil {
	//			fmt.Printf("dataStoreId: %v\n", dataStoreData["dataStoreId"])
	//			fmt.Printf("headerHash: %x\n", dataStoreData["headerHash"])
	//			found = true
	//		}
	//	} else {
	//		fmt.Println("Log data is empty")
	//	}
	//}
	//
	//if !found {
	//	fmt.Println("\nNo matching logs found for the specified criteria")
	//}

	fmt.Println("=== Test End ===")
}

func TestEthClient_GetLogs(t *testing.T) {
	startBlock := big.NewInt(20487720)
	endBlock := big.NewInt(20487725)
	var contractList []common.Address
	addressCm := common.HexToAddress(DataLayerServiceManagerAddr)
	contractList = append(contractList, addressCm)
	clint, err := NewEthClient("https://eth.llamarpc.com")
	if err != nil {
		fmt.Println("connect ethereum fail", "err", err)
		return
	}
	logList, err := clint.GetLogs(startBlock, endBlock, contractList)
	if err != nil {
		fmt.Println("get log fail", err)
		return
	}
	abiUint32, err := abi.NewType("uint32", "uint32", nil)
	if err != nil {
		fmt.Println("Abi new uint32 type error", "err", err)
		return
	}
	abiBytes32, err := abi.NewType("bytes32", "bytes32", nil)
	if err != nil {
		fmt.Println("Abi new bytes32 type error", "err", err)
		return
	}
	confirmDataStoreArgs := abi.Arguments{
		{
			Name:    "dataStoreId",
			Type:    abiUint32,
			Indexed: false,
		}, {
			Name:    "headerHash",
			Type:    abiBytes32,
			Indexed: false,
		},
	}
	var dataStoreData = make(map[string]interface{})
	for _, rLog := range logList {
		fmt.Println(rLog.Address.String())
		if strings.ToLower(rLog.Address.String()) != strings.ToLower(DataLayerServiceManagerAddr) {
			continue

		}
		if rLog.Topics[0] != ConfirmDataStoreEventABIHash {
			continue
		}
		if len(rLog.Data) > 0 {
			err := confirmDataStoreArgs.UnpackIntoMap(dataStoreData, rLog.Data)
			if err != nil {
				fmt.Println("Unpack data into map fail", "err", err)
				continue
			}
			if dataStoreData != nil {
				fmt.Printf("dataStoreId: %v\n", dataStoreData["dataStoreId"])
				fmt.Printf("headerHash: %x\n", dataStoreData["headerHash"])
			}
			return
		}
	}
}
