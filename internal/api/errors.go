package api

import "errors"

var (
	ErrCreateAPIRequest = errors.New("failed to create api request")
	ErrFailedAPIRequest = errors.New("failed to do api call")
	ErrReadResponseBody = errors.New("failed to copy downloaded file from body")
	ErrJsonMarshaling   = errors.New("failed to marshal request to json")
	ErrJsonUnmarshaling = errors.New("failed to unmarshal request to json")
)
