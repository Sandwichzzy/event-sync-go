package node

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/Sandwichzzy/event-sync-go/common/bigint"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	ErrHeaderTraversalAheadOfProvider            = errors.New("the HeaderTraversal's internal state is ahead of the provider")
	ErrHeaderTraversalAndProviderMismatchedState = errors.New("the HeaderTraversal and provider have diverged in state")
)

type HeaderTraversal struct {
	ethClient EthClient
	chainId   uint

	latestHeader        *types.Header
	lastTraversedHeader *types.Header // 最后遍历的区块头

	blockConfirmationDepth *big.Int // 区块确认深度
}

func NewHeaderTraversal(ethClient EthClient, fromHeader *types.Header, confDepth *big.Int, chainId uint) *HeaderTraversal {
	return &HeaderTraversal{
		ethClient:              ethClient,
		lastTraversedHeader:    fromHeader,
		blockConfirmationDepth: confDepth,
		chainId:                chainId,
	}
}

func (f *HeaderTraversal) LatestHeader() *types.Header {
	return f.latestHeader
}

func (f *HeaderTraversal) LastTraversedHeader() *types.Header {
	return f.lastTraversedHeader
}

func (f *HeaderTraversal) NextHeaders(maxSize uint64) ([]types.Header, error) {
	latestHeader, err := f.ethClient.BlockHeaderByNumber(nil) //nil 是取最新的区块头
	if err != nil {
		return nil, fmt.Errorf("unable to query latest block: %w", err)
	} else if latestHeader == nil {
		return nil, fmt.Errorf("latest header unreported")
	} else {
		f.latestHeader = latestHeader
	}
	// 计算安全高度：最新高度 - 确认深度
	endHeight := new(big.Int).Sub(latestHeader.Number, f.blockConfirmationDepth)
	if endHeight.Sign() < 0 {
		// No blocks with the provided confirmation depth available
		return nil, nil
	}
	//检查是否已经遍历到安全高度
	if f.lastTraversedHeader != nil {
		cmp := f.lastTraversedHeader.Number.Cmp(endHeight)
		if cmp == 0 { // 已经同步到最新安全区块
			return nil, nil
		} else if cmp > 0 { // 本地状态超前于提供者
			return nil, ErrHeaderTraversalAheadOfProvider
		}
	}
	//计算下一个要获取的高度
	nextHeight := bigint.Zero
	if f.lastTraversedHeader != nil {
		nextHeight = new(big.Int).Add(f.lastTraversedHeader.Number, bigint.One)
	}
	// 限制获取范围（考虑maxSize 步长） clamp(start,end, size)
	endHeight = bigint.Clamp(nextHeight, endHeight, maxSize) //`end - start` <= size.
	//获取区块头范围
	headers, err := f.ethClient.BlockHeadersByRange(nextHeight, endHeight, f.chainId)
	if err != nil {
		return nil, fmt.Errorf("error querying blocks by range: %w", err)
	}

	numHeaders := len(headers)
	if numHeaders == 0 {
		return nil, nil
	} else if f.lastTraversedHeader != nil && headers[0].ParentHash != f.lastTraversedHeader.Hash() {
		fmt.Println(f.lastTraversedHeader.Number)
		fmt.Println(headers[0].Number)
		fmt.Println(len(headers))
		return nil, ErrHeaderTraversalAndProviderMismatchedState //验证链连续性（防止重组）
	}
	//更新最后遍历的区块头
	f.lastTraversedHeader = &headers[numHeaders-1]
	return headers, nil
}
