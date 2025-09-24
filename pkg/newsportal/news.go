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
	activeTx  *pg.Tx
}

func NewNewsService(dbo db.DB) *Service {
	repo := db.NewNewsRepo(dbo).WithEnabledOnly()
	validate := NewValidator()

	return &Service{
		db:        dbo,
		repo:      repo,
		validator: validate,
	}
}

func (s *Service) WithinLock(ctx context.Context, lockName string, fn func(*Service) error) error {
	return s.db.RunInLock(ctx, lockName, func(tx *pg.Tx) error {
		locked := NewNewsService(s.db)
		locked.repo = locked.repo.WithTransaction(tx)
		locked.activeTx = tx

		return fn(locked)
	})
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
	errs, _, err := s.validateSuggestion(ctx, req)

	return errs, err
}

func (s *Service) validateSuggestion(
	ctx context.Context,
	req NewsSuggestion,
) (ValidationErrors, *Category, error) {
	var res ValidationErrors

	catDTO, err := s.repo.CategoryByID(ctx, req.CategoryID)
	if err != nil {
		return nil, nil, err
	}
	if catDTO == nil {
		res = append(res, ValidationError{
			Field: "categoryId",
			Error: "category does not exist",
		})
	}
	category := NewCategory(catDTO)

	err = s.validator.StructCtx(ctx, req)
	if err == nil {
		return res, category, nil
	}

	var errs validator.ValidationErrors
	if !errors.As(err, &errs) {
		return nil, nil, err
	}

	res = append(res, NewValidationErrors(errs)...)

	return res, category, nil
}

func (s *Service) Suggest(ctx context.Context, suggestion NewsSuggestion) (*News, error) {
	var news *News

	err := s.WithinLock(ctx, "news.Suggest", func(s *Service) error {
		vErrs, category, err := s.validateSuggestion(ctx, suggestion)
		if err != nil {
			return err
		}
		if len(vErrs) > 0 {
			return ErrBadRequest
		}

		tags, err := s.txCreateNonExistentTags(ctx, suggestion.Tags)
		if err != nil {
			return err
		}

		dto, err := s.repo.AddNews(ctx, suggestion.ToDB(tags.IDs()...))
		if err != nil {
			return err
		}

		news = NewNews(dto)
		news.SetTags(tags)
		news.Category = category

		return nil
	})
	if err != nil {
		return nil, err
	}

	return news, nil
}

func (s *Service) txCreateNonExistentTags(ctx context.Context, names []string) (Tags, error) {
	if err := s.requireTx(); err != nil {
		return nil, err
	}

	if len(names) == 0 {
		return nil, nil
	}

	tags, err := s.repo.TagsByFilters(ctx, &db.TagSearch{NameIn: names}, db.PagerNoLimit)
	if err != nil && !errors.Is(err, pg.ErrNoRows) {
		return nil, err
	}

	index := NewTags(tags).IndexByName()
	res := make(Tags, 0, len(names))

	for _, name := range names {
		if tag, ok := index[name]; ok {
			res = append(res, tag)

			continue
		}

		dto, err := s.repo.AddTag(ctx, &db.Tag{
			Name:     name,
			StatusID: db.StatusEnabled,
		})
		if err != nil {
			return nil, err
		}

		res = append(res, *(NewTag(dto)))
	}

	return res, nil
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

func (s *Service) requireTx() error {
	if s.activeTx == nil {
		return errNotInTx
	}

	return nil
}
