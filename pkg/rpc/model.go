package rpc

import (
	"time"

	"apisrv/pkg/newsportal"
)

//go:generate go tool colgen -imports=apisrv/pkg/newsportal
//colgen:News,Category,Tag
//colgen:News:MapP(newsportal.News)
//colgen:Category:MapP(newsportal.Category)
//colgen:Tag:MapP(newsportal.Tag)

type NewsListReq struct {
	CategoryID int `query:"category_id"`
	TagID      int `query:"tag_id"`
	Page       int `query:"page"`
	PerPage    int `query:"per_page"`
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
