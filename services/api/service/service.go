package service

import (
	"strconv"

	"github.com/Sandwichzzy/event-sync-go/database/worker"
	"github.com/Sandwichzzy/event-sync-go/services/api/models"
)

type Service interface {
	GetDepositTokensList(*models.QueryDTParams) (*models.DepositTokensResponse, error)

	QueryDTListParams(page string, pageSize string, order string) (*models.QueryDTParams, error)
}

type HandlerSvc struct {
	v                 *Validator
	depositTokensView worker.DepositTokensView
}

func (h HandlerSvc) GetDepositTokensList(params *models.QueryDTParams) (*models.DepositTokensResponse, error) {
	dtList, totalCount := h.depositTokensView.QueryDepositTokensList(params.Page, params.PageSize)
	return &models.DepositTokensResponse{
		Current: params.Page,
		Size:    params.PageSize,
		Total:   int64(totalCount),
		Result:  dtList,
	}, nil
}

func (h HandlerSvc) QueryDTListParams(page string, pageSize string, order string) (*models.QueryDTParams, error) {
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return nil, err
	}
	pageVal := h.v.ValidatePage(pageInt)

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		return nil, err
	}
	pageSizeVal := h.v.ValidatePageSize(pageSizeInt)
	orderBy := h.v.ValidateOrder(order)

	return &models.QueryDTParams{
		Page:     pageVal,
		PageSize: pageSizeVal,
		Order:    orderBy,
	}, nil
}

func New(v *Validator, dtv worker.DepositTokensView) Service {
	return &HandlerSvc{
		v:                 v,
		depositTokensView: dtv,
	}
}

func (h HandlerSvc) GetDepositList(params *models.QueryDTParams) (*models.DepositTokensResponse, error) {
	depositList, total := h.depositTokensView.QueryDepositTokensList(params.Page, params.PageSize)
	return &models.DepositTokensResponse{
		Current: params.Page,
		Size:    params.PageSize,
		Total:   int64(total),
		Result:  depositList,
	}, nil
}
