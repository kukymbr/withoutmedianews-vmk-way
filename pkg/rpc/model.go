package rpc

import (
	"time"

	"apisrv/pkg/newsportal"
)

//go:generate go tool colgen -imports=apisrv/pkg/newsportal
//colgen:News,Category,Tag,ValidationError
//colgen:News:MapP(newsportal.News)
//colgen:Category:MapP(newsportal.Category)
//colgen:Tag:MapP(newsportal.Tag)
//colgen:ValidationError:MapP(newsportal.ValidationError)

type NewsListReq struct {
	CategoryID int `json:"categoryId"`
	TagID      int `json:"tagId"`
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
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
	PublishedAt time.Time `json:"publishedAt"`

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
	Title      string
	Text       string
	ShortText  string
	CategoryID int
	Tags       []string
}

func (ns NewsSuggestion) ToDomain() newsportal.NewsSuggestion {
	return newsportal.NewsSuggestion{
		Title:      ns.Title,
		Text:       ns.Text,
		ShortText:  ns.ShortText,
		Tags:       ns.Tags,
		CategoryID: ns.CategoryID,
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
	StatusID int    `json:"statusId"`
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
	Field      string `json:"field"`
	Error      string `json:"error"`
	Constraint string `json:"constraint"`
}

func NewValidationError(in *newsportal.ValidationError) *ValidationError {
	if in == nil {
		return nil
	}

	return &ValidationError{
		Field:      in.Field,
		Error:      in.Error,
		Constraint: in.Constraint,
	}
}
