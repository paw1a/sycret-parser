package handler

import "errors"

var (
	ErrInvalidBody   = errors.New("invalid request body json")
	ErrFailedAPICall = errors.New("failed to make sycret API call")
	ErrCreateFile    = errors.New("failed to create temp file")
	ErrCopyDocFile   = errors.New("failed to copy downloaded doc file")
)
