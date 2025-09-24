package rpc

import (
	"time"

	"apisrv/pkg/newsportal"

	"github.com/go-playground/validator/v10"
)

//go:generate go tool colgen -imports=apisrv/pkg/newsportal,github.com/go-playground/validator/v10
//colgen:News,Category,Tag,ValidationError
//colgen:News:MapP(newsportal.News)
//colgen:Category:MapP(newsportal.Category)
//colgen:Tag:MapP(newsportal.Tag)
//colgen:ValidationError:Map(validator.FieldError)

type NewsListReq struct {
	CategoryID int `json:"category_id"`
	TagID      int `json:"tag_id"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
}

type NewsCountResponse struct {
	Count int `json:"count"`
}

type APIError struct {
	Message string `json:"message"`
}

type News struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	ShortText   string    `json:"short_text"`
	Content     *string   `json:"content"`
	Author      *string   `json:"author"`
	PublishedAt time.Time `json:"published_at"`

	Category *Category `json:"category"`
	Tags     Tags      `json:"tags"`
}

func NewNews(in *newsportal.News) *News {
	if in == nil {
		return nil
	}

	return &News{
		ID:          in.ID,
		Title:       in.Title,
		ShortText:   in.ShortText,
		Content:     in.Content,
		Author:      in.Author,
		PublishedAt: in.PublishedAt,
		Category:    NewCategory(in.Category),
		Tags:        NewTags(in.Tags),
	}
}

type NewsSuggestion struct {
	Title     string
	Text      string
	ShortText string
	Tags      []string
}

func (ns NewsSuggestion) ToDomain() newsportal.NewsSuggestion {
	return newsportal.NewsSuggestion{
		Title:     ns.Title,
		Text:      ns.Text,
		ShortText: ns.ShortText,
		Tags:      ns.Tags,
	}
}

type Category struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func NewCategory(in *newsportal.Category) *Category {
	if in == nil {
		return nil
	}

	return &Category{
		ID:    in.ID,
		Title: in.Title,
	}
}

type Tag struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	StatusID int    `json:"status_id"`
}

func NewTag(in *newsportal.Tag) *Tag {
	if in == nil {
		return nil
	}

	return &Tag{
		ID:       in.ID,
		Name:     in.Name,
		StatusID: in.StatusID,
	}
}

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func NewValidationError(fieldError validator.FieldError) ValidationError {
	return ValidationError{
		Field: fieldError.Field(),
		Error: fieldError.Error(),
	}
}
