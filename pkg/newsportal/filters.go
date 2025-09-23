package newsportal

import "apisrv/pkg/db"

type NewsesFilter struct {
	CategoryID int
	TagID      int
}

func NewNewsFilter(categoryID int, tagID int) NewsesFilter {
	return NewsesFilter{
		CategoryID: categoryID,
		TagID:      tagID,
	}
}

func (f NewsesFilter) toDBSearch() *db.NewsSearch {
	search := &db.NewsSearch{}

	if f.CategoryID > 0 {
		search.CategoryID = &f.CategoryID
	}

	if f.TagID > 0 {
		search.TagID = &f.TagID
	}

	return search
}
