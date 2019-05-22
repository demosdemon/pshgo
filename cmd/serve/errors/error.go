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

func BadRequest(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusBadRequest,
		Message:    msg,
		Err:        err,
	}
}

func Unauthorized(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusUnauthorized,
		Message:    msg,
		Err:        err,
	}
}

func PaymentRequired(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusPaymentRequired,
		Message:    msg,
		Err:        err,
	}
}

func Forbidden(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusForbidden,
		Message:    msg,
		Err:        err,
	}
}

func NotFound(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusNotFound,
		Message:    msg,
		Err:        err,
	}
}

func MethodNotAllowed(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusMethodNotAllowed,
		Message:    msg,
		Err:        err,
	}
}

func NotAcceptable(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusNotAcceptable,
		Message:    msg,
		Err:        err,
	}
}

func ProxyAuthRequired(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusProxyAuthRequired,
		Message:    msg,
		Err:        err,
	}
}

func RequestTimeout(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusRequestTimeout,
		Message:    msg,
		Err:        err,
	}
}

func Conflict(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusConflict,
		Message:    msg,
		Err:        err,
	}
}

func Gone(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusGone,
		Message:    msg,
		Err:        err,
	}
}

func LengthRequired(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusLengthRequired,
		Message:    msg,
		Err:        err,
	}
}

func PreconditionFailed(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusPreconditionFailed,
		Message:    msg,
		Err:        err,
	}
}

func RequestEntityTooLarge(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusRequestEntityTooLarge,
		Message:    msg,
		Err:        err,
	}
}

func RequestURITooLong(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusRequestURITooLong,
		Message:    msg,
		Err:        err,
	}
}

func UnsupportedMediaType(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusUnsupportedMediaType,
		Message:    msg,
		Err:        err,
	}
}

func RequestedRangeNotSatisfiable(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusRequestedRangeNotSatisfiable,
		Message:    msg,
		Err:        err,
	}
}

func ExpectationFailed(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusExpectationFailed,
		Message:    msg,
		Err:        err,
	}
}

func Teapot(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusTeapot,
		Message:    msg,
		Err:        err,
	}
}

func MisdirectedRequest(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusMisdirectedRequest,
		Message:    msg,
		Err:        err,
	}
}

func UnprocessableEntity(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusUnprocessableEntity,
		Message:    msg,
		Err:        err,
	}
}

func Locked(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusLocked,
		Message:    msg,
		Err:        err,
	}
}

func FailedDependency(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusFailedDependency,
		Message:    msg,
		Err:        err,
	}
}

func TooEarly(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusTooEarly,
		Message:    msg,
		Err:        err,
	}
}

func UpgradeRequired(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusUpgradeRequired,
		Message:    msg,
		Err:        err,
	}
}

func PreconditionRequired(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusPreconditionRequired,
		Message:    msg,
		Err:        err,
	}
}

func TooManyRequests(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusTooManyRequests,
		Message:    msg,
		Err:        err,
	}
}

func RequestHeaderFieldsTooLarge(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusRequestHeaderFieldsTooLarge,
		Message:    msg,
		Err:        err,
	}
}

func UnavailableForLegalReasons(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusUnavailableForLegalReasons,
		Message:    msg,
		Err:        err,
	}
}

func InternalServerError(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusInternalServerError,
		Message:    msg,
		Err:        err,
	}
}

func NotImplemented(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusNotImplemented,
		Message:    msg,
		Err:        err,
	}
}

func BadGateway(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusBadGateway,
		Message:    msg,
		Err:        err,
	}
}

func ServiceUnavailable(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusServiceUnavailable,
		Message:    msg,
		Err:        err,
	}
}

func GatewayTimeout(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusGatewayTimeout,
		Message:    msg,
		Err:        err,
	}
}

func HTTPVersionNotSupported(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusHTTPVersionNotSupported,
		Message:    msg,
		Err:        err,
	}
}

func VariantAlsoNegotiates(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusVariantAlsoNegotiates,
		Message:    msg,
		Err:        err,
	}
}

func InsufficientStorage(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusInsufficientStorage,
		Message:    msg,
		Err:        err,
	}
}

func LoopDetected(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusLoopDetected,
		Message:    msg,
		Err:        err,
	}
}

func NotExtended(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusNotExtended,
		Message:    msg,
		Err:        err,
	}
}

func NetworkAuthenticationRequired(msg string, err error) error {
	return HTTPError{
		StatusCode: http.StatusNetworkAuthenticationRequired,
		Message:    msg,
		Err:        err,
	}
}
