package grpc

import (
	"context"

	"github.com/Sandwichzzy/event-sync-go/services/grpc/eventpb"
)

func (rs *RpcService) GetDepositTokenList(ctx context.Context, request *eventpb.DepositTokenListReq) (*eventpb.DepositTokenListRep, error) {
	dtList, totalCount := rs.db.DepositTokens.QueryDepositTokensList(int(request.Page), int(request.PageSize))
	if totalCount == 0 {
		return &eventpb.DepositTokenListRep{
			Code:         eventpb.ReturnCode_SUCCESS,
			Message:      "No data in database",
			DepositToken: nil,
		}, nil
	}
	var depositTokenList []*eventpb.DepositToken
	for _, dt := range dtList {
		dtItem := &eventpb.DepositToken{
			Guid:         dt.GUID.String(),
			BlockNumber:  dt.BlockNumber.Uint64(),
			TokenAddress: dt.TokenAddress.String(),
			Sender:       dt.Sender.String(),
			Amount:       dt.Amount.Uint64(),
			Timestamp:    dt.Timestamp,
		}
		depositTokenList = append(depositTokenList, dtItem)
	}
	return &eventpb.DepositTokenListRep{
		Code:         eventpb.ReturnCode_SUCCESS,
		Message:      "get data success",
		DepositToken: depositTokenList,
	}, nil
}

func (rs *RpcService) GetDepositTokenDetail(ctx context.Context, request *eventpb.DepositTokenDetailReq) (*eventpb.DepositTokenDetailRep, error) {
	dt, err := rs.db.DepositTokens.QueryDepositTokensById(request.Guid)
	if err != nil {
		return &eventpb.DepositTokenDetailRep{
			Code:    eventpb.ReturnCode_ERROR,
			Message: "query data fail",
		}, nil
	}
	return &eventpb.DepositTokenDetailRep{
		Code:         eventpb.ReturnCode_SUCCESS,
		Message:      "get data success",
		Guid:         dt.GUID.String(),
		BlockNumber:  dt.BlockNumber.Uint64(),
		TokenAddress: dt.TokenAddress.String(),
		Sender:       dt.Sender.String(),
		Amount:       dt.Amount.Uint64(),
		Timestamp:    dt.Timestamp,
	}, nil
}
