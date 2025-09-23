//nolint:dupl
package vt

import (
	"time"

	"apisrv/pkg/db"
)

type Category struct {
	ID       int    `json:"id"`
	Title    string `json:"title" validate:"required,max=255"`
	Sort     *int   `json:"sort"`
	StatusID int    `json:"statusId" validate:"required,status"`

	Status *Status `json:"status"`
}

func (c *Category) ToDB() *db.Category {
	if c == nil {
		return nil
	}

	category := &db.Category{
		ID:       c.ID,
		Title:    c.Title,
		Sort:     c.Sort,
		StatusID: c.StatusID,
	}

	return category
}

type CategorySearch struct {
	ID       *int    `json:"id"`
	Title    *string `json:"title"`
	Sort     *int    `json:"sort"`
	StatusID *int    `json:"statusId"`
	IDs      []int   `json:"ids"`
}

func (cs *CategorySearch) ToDB() *db.CategorySearch {
	if cs == nil {
		return nil
	}

	return &db.CategorySearch{
		ID:         cs.ID,
		TitleILike: cs.Title,
		Sort:       cs.Sort,
		StatusID:   cs.StatusID,
		IDs:        cs.IDs,
	}
}

type CategorySummary struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Sort  *int   `json:"sort"`

	Status *Status `json:"status"`
}

type News struct {
	ID          int       `json:"id"`
	Title       string    `json:"title" validate:"required,max=255"`
	ShortText   string    `json:"shortText" validate:"required,max=1024"`
	Content     *string   `json:"content"`
	Author      *string   `json:"author" validate:"omitempty,max=255"`
	CategoryID  int       `json:"categoryId" validate:"required"`
	TagIDs      []int     `json:"tagIds" validate:"required"`
	PublishedAt time.Time `json:"publishedAt" validate:"required"`
	CreatedAt   time.Time `json:"createdAt"`
	StatusID    int       `json:"statusId" validate:"required,status"`

	Category *CategorySummary `json:"category"`
	Status   *Status          `json:"status"`
}

func (n *News) ToDB() *db.News {
	if n == nil {
		return nil
	}

	news := &db.News{
		ID:          n.ID,
		Title:       n.Title,
		ShortText:   n.ShortText,
		Content:     n.Content,
		Author:      n.Author,
		CategoryID:  n.CategoryID,
		TagIDs:      n.TagIDs,
		PublishedAt: n.PublishedAt,
		CreatedAt:   n.CreatedAt,
		StatusID:    n.StatusID,
	}

	return news
}

type NewsSearch struct {
	ID              *int       `json:"id"`
	Title           *string    `json:"title"`
	ShortText       *string    `json:"shortText"`
	Content         *string    `json:"content"`
	Author          *string    `json:"author"`
	CategoryID      *int       `json:"categoryId"`
	PublishedAt     *time.Time `json:"publishedAt"`
	CreatedAt       *time.Time `json:"createdAt"`
	StatusID        *int       `json:"statusId"`
	IDs             []int      `json:"ids"`
	TagID           *int       `json:"tagId"`
	PublishedBefore *time.Time `json:"publishedBefore"`
}

func (ns *NewsSearch) ToDB() *db.NewsSearch {
	if ns == nil {
		return nil
	}

	return &db.NewsSearch{
		ID:              ns.ID,
		TitleILike:      ns.Title,
		ShortTextILike:  ns.ShortText,
		ContentILike:    ns.Content,
		AuthorILike:     ns.Author,
		CategoryID:      ns.CategoryID,
		PublishedAt:     ns.PublishedAt,
		CreatedAt:       ns.CreatedAt,
		StatusID:        ns.StatusID,
		IDs:             ns.IDs,
		TagID:           ns.TagID,
		PublishedBefore: ns.PublishedBefore,
	}
}

type NewsSummary struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	ShortText   string    `json:"shortText"`
	Content     *string   `json:"content"`
	Author      *string   `json:"author"`
	CategoryID  int       `json:"categoryId"`
	PublishedAt time.Time `json:"publishedAt"`
	CreatedAt   time.Time `json:"createdAt"`

	Category *CategorySummary `json:"category"`
	Status   *Status          `json:"status"`
}

type Tag struct {
	ID       int    `json:"id"`
	Name     string `json:"name" validate:"required,max=64"`
	StatusID int    `json:"statusId" validate:"required,status"`

	Status *Status `json:"status"`
}

func (t *Tag) ToDB() *db.Tag {
	if t == nil {
		return nil
	}

	tag := &db.Tag{
		ID:       t.ID,
		Name:     t.Name,
		StatusID: t.StatusID,
	}

	return tag
}

type TagSearch struct {
	ID       *int    `json:"id"`
	Name     *string `json:"name"`
	StatusID *int    `json:"statusId"`
	IDs      []int   `json:"ids"`
}

func (ts *TagSearch) ToDB() *db.TagSearch {
	if ts == nil {
		return nil
	}

	return &db.TagSearch{
		ID:        ts.ID,
		NameILike: ts.Name,
		StatusID:  ts.StatusID,
		IDs:       ts.IDs,
	}
}

type TagSummary struct {
	ID   int    `json:"id"`
	Name string `json:"name"`

	Status *Status `json:"status"`
}
