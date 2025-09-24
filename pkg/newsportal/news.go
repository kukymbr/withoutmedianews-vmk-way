package newsportal

import (
	"context"
	"fmt"

	"apisrv/pkg/db"

	"github.com/go-playground/validator/v10"
)

type Service struct {
	db        db.DB
	repo      db.NewsRepo
	validator *validator.Validate
}

func NewNewsService(dbo db.DB, validator *validator.Validate) *Service {
	return &Service{
		repo:      db.NewNewsRepo(dbo).WithEnabledOnly(),
		validator: validator,
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

func (s *Service) GetNews(ctx context.Context, id int) (News, error) {
	dto, err := s.repo.OneNews(
		ctx,
		&db.NewsSearch{ID: &id},
		db.EnabledOnly(),
		db.AlreadyPublished(),
		db.WithColumns(db.Columns.News.Category),
	)
	if err != nil {
		return News{}, fmt.Errorf("read news item: %w", err)
	}

	if dto == nil {
		return News{}, ErrNotFound
	}

	list, err := s.enrichNewsesWithTags(ctx, NewNewsList([]db.News{*dto}))
	if err != nil {
		return News{}, err
	}

	return list[0], nil
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

func (s *Service) ValidateSuggestion(ctx context.Context, req NewsSuggestion) error {
	return s.validator.StructCtx(ctx, req)
}

func (s *Service) Suggest(ctx context.Context, suggestion NewsSuggestion) error {
	if err := s.ValidateSuggestion(ctx, suggestion); err != nil {
		return err
	}

	if err := s.repo.CreateNonExistentTags(ctx, suggestion.Tags); err != nil {
		return err
	}

	tagDTOs, err := s.repo.TagsByFilters(ctx, &db.TagSearch{NameIn: suggestion.Tags}, db.PagerNoLimit)
	if err != nil {
		return err
	}

	dto := db.News{
		Title:     suggestion.Title,
		ShortText: suggestion.ShortText,
		Content:   &suggestion.Text,
		TagIDs:    NewTags(tagDTOs).IDs(),
		StatusID:  db.StatusDisabled,
	}

	_, err = s.repo.AddNews(ctx, &dto)

	return err
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
