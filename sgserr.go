package sgs

import (
	"errors"
)

var ErrNotFound = errors.New("no such client")
var ErrConflict = errors.New("client already exists")
