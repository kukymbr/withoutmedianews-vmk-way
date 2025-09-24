package db

import (
	"context"
	"errors"

	"github.com/go-pg/pg/v10"
)

//go:generate go tool colgen
//colgen:Tag
//colgen:Tag:Index(Name)

func (nr NewsRepo) CreateNonExistentTags(ctx context.Context, names []string) (Tags, error) {
	if len(names) == 0 {
		return nil, nil
	}

	tags, err := nr.TagsByFilters(ctx, &TagSearch{NameIn: names}, PagerNoLimit)
	if err != nil && !errors.Is(err, pg.ErrNoRows) {
		return nil, err
	}

	index := Tags(tags).IndexByName()

	for _, name := range names {
		if _, ok := index[name]; ok {
			continue
		}

		dto := &Tag{
			Name:     name,
			StatusID: StatusEnabled,
		}

		tag, err := nr.AddTag(ctx, dto)
		if err != nil {
			return nil, err
		}

		tags = append(tags, *tag)
	}

	return sortTagsByNames(tags, names), nil
}
