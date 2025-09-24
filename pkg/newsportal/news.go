package newsportal

import (
	"context"
	"errors"
	"fmt"

	"apisrv/pkg/db"
	"github.com/go-pg/pg/v10"

	"github.com/go-playground/validator/v10"
)

type Service struct {
	db        db.DB
	repo      db.NewsRepo
	validator *validator.Validate
}

func NewNewsService(dbo db.DB) *Service {
	repo := db.NewNewsRepo(dbo).WithEnabledOnly()
	validate := NewValidator(repo)

	return &Service{
		db:        dbo,
		repo:      repo,
		validator: validate,
	}
}

func (s *Service) GetList(
	ctx context.Context,
	filter NewsesFilter,
	page, perPage int,
) ([]News, error) {
	items, err := s.repo.NewsByFilters(
		ctx,
		filter.toDBSearch(),
		db.NewPager(page, perPage),
		db.AlreadyPublished(),
		db.WithColumns(db.Columns.News.Category),
	)
	if err != nil {
		return nil, fmt.Errorf("read news list: %w", err)
	}

	return s.enrichNewsesWithTags(ctx, NewNewsList(items))
}

func (s *Service) GetNews(ctx context.Context, id int) (*News, error) {
	dto, err := s.repo.OneNews(
		ctx,
		&db.NewsSearch{ID: &id},
		db.EnabledOnly(),
		db.AlreadyPublished(),
		db.WithColumns(db.Columns.News.Category),
	)
	if err != nil {
		return nil, fmt.Errorf("read news item: %w", err)
	}

	if dto == nil {
		return nil, ErrNotFound
	}

	return s.enrichNewsWithTags(ctx, NewNews(dto))
}

func (s *Service) GetCount(ctx context.Context, filter NewsesFilter) (int, error) {
	count, err := s.repo.CountNews(ctx, filter.toDBSearch())
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Service) GetCategories(ctx context.Context) ([]Category, error) {
	categories, err := s.repo.CategoriesByFilters(ctx, nil, db.PagerNoLimit)
	if err != nil {
		return nil, fmt.Errorf("read categories from repo: %w", err)
	}

	return NewCategories(categories), nil
}

func (s *Service) GetTags(ctx context.Context) ([]Tag, error) {
	tags, err := s.repo.TagsByFilters(ctx, nil, db.PagerNoLimit)
	if err != nil {
		return nil, fmt.Errorf("read tags from repo: %w", err)
	}

	return NewTags(tags), nil
}

func (s *Service) ValidateSuggestion(ctx context.Context, req NewsSuggestion) (ValidationErrors, error) {
	err := s.validator.StructCtx(ctx, req)
	if err == nil {
		return nil, nil
	}

	var errs validator.ValidationErrors
	if !errors.As(err, &errs) {
		return nil, err
	}

	return NewValidationErrors(errs), nil
}

func (s *Service) Suggest(ctx context.Context, suggestion NewsSuggestion) (*News, error) {
	vErrs, err := s.ValidateSuggestion(ctx, suggestion)
	if err != nil {
		return nil, err
	}
	if len(vErrs) > 0 {
		return nil, ErrBadRequest
	}

	var news *News

	err = s.db.RunInLock(ctx, "news.Suggest", func(tx *pg.Tx) error {
		repo := s.repo.WithTransaction(tx)

		tagDTOs, err := repo.CreateNonExistentTags(ctx, suggestion.Tags)
		if err != nil {
			return err
		}

		tags := NewTags(tagDTOs)

		dto, err := repo.AddNews(ctx, suggestion.ToDB(tags.IDs()...))
		if err != nil {
			return err
		}

		news = NewNews(dto)
		news.SetTags(tags)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.enrichNewsWithTags(ctx, news)
}

func (s *Service) enrichNewsWithTags(ctx context.Context, news *News) (*News, error) {
	if news == nil {
		return nil, nil
	}

	list, err := s.enrichNewsesWithTags(ctx, NewsList{*news})
	if err != nil {
		return nil, err
	}

	return &list[0], nil
}

func (s *Service) enrichNewsesWithTags(ctx context.Context, newses NewsList) (NewsList, error) {
	tagIDs := newses.UniqueTagIDs()
	if len(tagIDs) == 0 {
		return newses, nil
	}

	dbTags, err := s.repo.TagsByFilters(ctx, &db.TagSearch{IDs: tagIDs}, db.PagerNoLimit)
	if err != nil {
		return nil, err
	}

	newses.SetTags(NewTags(dbTags))

	return newses, nil
}
