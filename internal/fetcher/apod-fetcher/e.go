package apodfetcher

import "fmt"

type RequestError struct {
	err error
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("%v", e.err)
}

func NewRequestError(e error) error {
	return &RequestError{
		err: fmt.Errorf("request error: %w", e),
	}
}

type DecodeJsonError struct {
	err error
}

func (e *DecodeJsonError) Error() string {
	return fmt.Sprintf("%v", e.err)
}

func NewDecodeJSONError(e error) error {
	return &DecodeJsonError{
		err: fmt.Errorf("decode JSON error: %w", e),
	}
}
