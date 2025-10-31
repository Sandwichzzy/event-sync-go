package ethclient

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
)

type EthClient struct {
	client *ethclient.Client
}

func NewEthClient(rpcUrl string) (*EthClient, error) {
	ethClient, err := ethclient.DialContext(context.Background(), rpcUrl)
	if err != nil {
		log.Error("new eth client fail", "err", err)
		return nil, err
	}
	return &EthClient{ethClient}, nil
}

// GetTxReceiptByHash eth_getTransactionReceipt
func (eth *EthClient) GetTxReceiptByHash(txHash string) (*types.Receipt, error) {
	r, err := eth.client.TransactionReceipt(context.Background(), common.HexToHash(txHash))
	if err != nil {
		log.Error("get tx receipt fail", "err", err)
		return nil, err
	}
	return r, nil
}

// GetLogs eth_getLogs
func (eth *EthClient) GetLogs(startBlock, endBlock *big.Int, contractAddressList []common.Address) ([]types.Log, error) {
	filterQueryParams := ethereum.FilterQuery{
		FromBlock: startBlock,
		ToBlock:   endBlock,
		Addresses: contractAddressList,
	}
	result, err := eth.client.FilterLogs(context.Background(), filterQueryParams)
	if err != nil {
		log.Error("filter logs fail", "err", err)
		return nil, err
	}
	return result, nil
}

// GetBlockReceipt eth_getBlockReceipts
func (eth *EthClient) GetBlockReceipt(blockNumber *rpc.BlockNumber, blockHash *common.Hash) ([]*types.Receipt, error) {
	ethBlockIf := rpc.BlockNumberOrHash{BlockNumber: blockNumber, BlockHash: blockHash, RequireCanonical: true}

	return eth.client.BlockReceipts(context.Background(), ethBlockIf)
}
