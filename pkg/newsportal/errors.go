package newsportal

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
)
