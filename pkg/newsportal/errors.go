package newsportal

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrBadRequest = errors.New("bad request")

	errNotInTx = errors.New("not in transaction")
)
