package db

// TODO: do we have generator for this?
func sortTagsByNames(tags Tags, names []string) Tags {
	if len(tags) != len(names) {
		return tags
	}

	index := tags.IndexByName()
	sorted := make(Tags, 0, len(tags))

	for _, name := range names {
		sorted = append(sorted, index[name])
	}

	return sorted
}
