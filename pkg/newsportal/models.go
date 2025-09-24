package newsportal

import (
	"time"

	"apisrv/pkg/db"
)

//go:generate go tool colgen -imports=apisrv/pkg/db
//colgen:News,Category,Tag
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
	Title     string   `validate:"required,min=3,max=255"`
	Text      string   `validate:"required,min=1024"`
	ShortText string   `validate:"required,max=255"`
	Tags      []string `validate:"required,dive,alphanumunicode"`
}
