package rpc

import (
	"context"
	"errors"

	"apisrv/pkg/newsportal"
	"github.com/vmkteam/zenrpc/v2"
)

//go:generate go tool zenrpc

type NewsService struct {
	zenrpc.Service

	service *newsportal.Service
}

func NewNewsService(service *newsportal.Service) *NewsService {
	return &NewsService{
		service: service,
	}
}

func (ctrl NewsService) Get(ctx context.Context, req NewsListReq) ([]News, error) {
	items, err := ctrl.service.GetList(
		ctx,
		newsportal.NewNewsFilter(req.CategoryID, req.TagID),
		req.Page,
		req.PerPage,
	)
	if err != nil {
		return nil, err
	}

	// TODO: switch to summary
	resp := NewNewsList(items)

	return resp, nil
}

func (ctrl NewsService) GetByID(ctx context.Context, id int) (*News, error) {
	item, err := ctrl.service.GetNews(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := NewNews(item)

	return resp, nil
}

func (ctrl NewsService) Count(ctx context.Context, req NewsListReq) (int, error) {
	count, err := ctrl.service.GetCount(
		ctx,
		newsportal.NewNewsFilter(req.CategoryID, req.TagID),
	)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (ctrl NewsService) Categories(ctx context.Context) ([]Category, error) {
	categories, err := ctrl.service.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	resp := NewCategories(categories)

	return resp, nil
}

func (ctrl NewsService) Tags(ctx context.Context) ([]Tag, error) {
	tags, err := ctrl.service.GetTags(ctx)
	if err != nil {
		return nil, err
	}

	resp := NewTags(tags)

	return resp, nil
}

func (ctrl NewsService) ValidateSuggestion(ctx context.Context, req NewsSuggestion) (ValidationErrors, error) {
	dtos, err := ctrl.service.ValidateSuggestion(ctx, req.ToDomain())
	if err != nil {
		return nil, err
	}
	if len(dtos) == 0 {
		return nil, nil
	}

	return NewValidationErrors(dtos), nil
}

func (ctrl NewsService) Suggest(ctx context.Context, req NewsSuggestion) (*News, error) {
	dto, err := ctrl.service.Suggest(ctx, req.ToDomain())
	if err != nil {
		return nil, err
	}

	switch {
	case errors.Is(err, newsportal.ErrBadRequest):
		return nil, newBadRequestError(err)
	case err != nil:
		return nil, newInternalError(err)
	}

	return NewNews(dto), nil
}
