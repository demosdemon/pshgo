package errors

import (
	"fmt"
	"net/http"
	"strings"
)

type HTTPError struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
	Err        error  `json:"error,omitempty"`
}

func (e HTTPError) Error() string {
	var b strings.Builder

	_, _ = fmt.Fprintf(&b, "HTTP [%03d] %s", e.StatusCode, e.Message)

	if e.Err != nil {
		_, _ = fmt.Fprintf(&b, ": %v", e.Err)
	}

	return b.String()
}

func InternalServerError(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusInternalServerError,
		Message:    msg,
		Err:        err,
	}
}
