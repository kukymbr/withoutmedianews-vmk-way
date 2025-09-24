package newsportal

import (
	"time"

	"apisrv/pkg/db"
	"github.com/go-playground/validator/v10"
)

//go:generate go tool colgen -imports=apisrv/pkg/db
//colgen:News,Category,Tag,ValidationError
//colgen:News:TagIDs,UniqueTagIDs,MapP(db.News)
//colgen:Category:MapP(db.Category)
//colgen:Tag:MapP(db.Tag)
//colgen:Tag:Index(Name)

type Tag struct {
	ID       int
	Name     string
	StatusID int
}

func NewTag(in *db.Tag) *Tag {
	if in == nil {
		return nil
	}

	return &Tag{
		ID:       in.ID,
		Name:     in.Name,
		StatusID: in.StatusID,
	}
}

func (tag *Tag) ToDB() *db.Tag {
	if tag == nil {
		return nil
	}

	return &db.Tag{
		ID:       tag.ID,
		Name:     tag.Name,
		StatusID: tag.StatusID,
	}
}

type Category struct {
	ID       int
	Title    string
	Sort     *int
	StatusID int
}

func NewCategory(in *db.Category) *Category {
	if in == nil {
		return nil
	}

	return &Category{
		ID:       in.ID,
		Title:    in.Title,
		Sort:     in.Sort,
		StatusID: in.StatusID,
	}
}

type News struct {
	ID          int
	Title       string
	ShortText   string
	Content     *string
	Author      *string
	CategoryID  int
	TagIDs      []int
	PublishedAt time.Time
	CreatedAt   time.Time
	StatusID    int
	Category    *Category
	Tags        Tags
}

func (news *News) SetTags(tags Tags) {
	if len(news.TagIDs) == 0 || len(tags) == 0 {
		return
	}

	index := tags.Index()

	news.Tags = make(Tags, 0, len(news.TagIDs))

	for _, id := range news.TagIDs {
		tag, ok := index[id]
		if ok {
			news.Tags = append(news.Tags, tag)
		}
	}
}

func NewNews(in *db.News) *News {
	if in == nil {
		return nil
	}

	return &News{
		ID:          in.ID,
		Title:       in.Title,
		ShortText:   in.ShortText,
		Content:     in.Content,
		Author:      in.Author,
		CategoryID:  in.CategoryID,
		TagIDs:      in.TagIDs,
		PublishedAt: in.PublishedAt,
		CreatedAt:   in.CreatedAt,
		StatusID:    in.StatusID,
		Category:    NewCategory(in.Category),
	}
}

func (list NewsList) SetTags(tags Tags) {
	index := tags.Index()

	for i, news := range list {
		itemTags := make([]Tag, 0, len(news.TagIDs))

		for _, id := range news.TagIDs {
			tag, ok := index[id]
			if ok {
				itemTags = append(itemTags, tag)
			}
		}

		list[i].Tags = itemTags
	}
}

type NewsSuggestion struct {
	Title      string   `validate:"required,min=3,max=255" json:"title"`
	Text       string   `validate:"required" json:"text"`
	ShortText  string   `validate:"required,max=255" json:"shortText"`
	Tags       []string `validate:"required,dive,alphanumunicode" json:"tags"`
	CategoryID int      `validate:"required" json:"categoryId"`
}

func (ns *NewsSuggestion) ToDB(tagIDs ...int) *db.News {
	if ns == nil {
		return nil
	}

	return &db.News{
		Title:      ns.Title,
		ShortText:  ns.ShortText,
		Content:    &ns.Text,
		CategoryID: ns.CategoryID,
		TagIDs:     tagIDs,
		StatusID:   db.StatusDisabled,
	}
}

type ValidationError struct {
	Field      string `json:"field"`
	Error      string `json:"error"`
	Constraint string `json:"constraint"`
}

func NewValidationError(fieldError validator.FieldError) ValidationError {
	return ValidationError{
		Field: fieldError.Field(),
		Error: fieldError.Error(),
	}
}

func NewValidationErrors(errs validator.ValidationErrors) ValidationErrors {
	vErrs := make(ValidationErrors, 0, len(errs))

	for _, err := range errs {
		vErrs = append(vErrs, NewValidationError(err))
	}

	return vErrs
}
