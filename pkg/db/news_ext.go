package db

import (
	"context"

	"github.com/go-pg/pg/v10"
)

//go:generate go tool colgen
//colgen:Tag
//colgen:Tag:Index(Name)

func (nr NewsRepo) CreateNonExistentTags(ctx context.Context, names []string) error {
	if len(names) == 0 {
		return nil
	}

	return nr.db.RunInLock(ctx, "CreateNonExistentTags", func(tx *pg.Tx) error {
		repo := nr.WithTransaction(tx)

		tags, err := repo.TagsByFilters(ctx, &TagSearch{NameIn: names}, PagerNoLimit)
		if err != nil {
			return err
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

			// TODO: batches?
			if _, err := repo.AddTag(ctx, dto); err != nil {
				return err
			}
		}

		return nil
	})
}
