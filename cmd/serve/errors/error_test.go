package errors_test

import (
	. "github.com/demosdemon/pshgo/cmd/serve/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHTTPError_Error(t *testing.T) {
	t.Run("with err", func(t *testing.T) {
		err := HTTPError{
			StatusCode: 500,
			Message:    "Danger, Will Robinson!",
			Err:        assert.AnError,
		}
		assert.Equal(t, "HTTP [500] Danger, Will Robinson!: assert.AnError general error for testing", err.Error())
	})

	t.Run("without err", func(t *testing.T) {
		err := HTTPError{
			StatusCode: 404,
			Message:    "Whatcha talkin bout Willis?",
			Err:        nil,
		}
		assert.Equal(t, "HTTP [404] Whatcha talkin bout Willis?", err.Error())
	})
}

func TestBadRequest(t *testing.T) {
	err := BadRequest("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestUnauthorized(t *testing.T) {
	err := Unauthorized("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusUnauthorized, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestPaymentRequired(t *testing.T) {
	err := PaymentRequired("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusPaymentRequired, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestForbidden(t *testing.T) {
	err := Forbidden("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusForbidden, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestNotFound(t *testing.T) {
	err := NotFound("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestMethodNotAllowed(t *testing.T) {
	err := MethodNotAllowed("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusMethodNotAllowed, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestNotAcceptable(t *testing.T) {
	err := NotAcceptable("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotAcceptable, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestProxyAuthRequired(t *testing.T) {
	err := ProxyAuthRequired("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusProxyAuthRequired, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestRequestTimeout(t *testing.T) {
	err := RequestTimeout("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusRequestTimeout, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestConflict(t *testing.T) {
	err := Conflict("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusConflict, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestGone(t *testing.T) {
	err := Gone("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusGone, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestLengthRequired(t *testing.T) {
	err := LengthRequired("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusLengthRequired, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestPreconditionFailed(t *testing.T) {
	err := PreconditionFailed("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusPreconditionFailed, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestRequestEntityTooLarge(t *testing.T) {
	err := RequestEntityTooLarge("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusRequestEntityTooLarge, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestRequestURITooLong(t *testing.T) {
	err := RequestURITooLong("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusRequestURITooLong, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestUnsupportedMediaType(t *testing.T) {
	err := UnsupportedMediaType("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusUnsupportedMediaType, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestRequestedRangeNotSatisfiable(t *testing.T) {
	err := RequestedRangeNotSatisfiable("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusRequestedRangeNotSatisfiable, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestExpectationFailed(t *testing.T) {
	err := ExpectationFailed("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusExpectationFailed, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestTeapot(t *testing.T) {
	err := Teapot("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusTeapot, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestMisdirectedRequest(t *testing.T) {
	err := MisdirectedRequest("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusMisdirectedRequest, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestUnprocessableEntity(t *testing.T) {
	err := UnprocessableEntity("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusUnprocessableEntity, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestLocked(t *testing.T) {
	err := Locked("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusLocked, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestFailedDependency(t *testing.T) {
	err := FailedDependency("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusFailedDependency, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestTooEarly(t *testing.T) {
	err := TooEarly("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusTooEarly, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestUpgradeRequired(t *testing.T) {
	err := UpgradeRequired("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusUpgradeRequired, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestPreconditionRequired(t *testing.T) {
	err := PreconditionRequired("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusPreconditionRequired, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestTooManyRequests(t *testing.T) {
	err := TooManyRequests("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusTooManyRequests, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestRequestHeaderFieldsTooLarge(t *testing.T) {
	err := RequestHeaderFieldsTooLarge("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusRequestHeaderFieldsTooLarge, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestUnavailableForLegalReasons(t *testing.T) {
	err := UnavailableForLegalReasons("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusUnavailableForLegalReasons, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestInternalServerError(t *testing.T) {
	err := InternalServerError("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestNotImplemented(t *testing.T) {
	err := NotImplemented("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotImplemented, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestBadGateway(t *testing.T) {
	err := BadGateway("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadGateway, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestServiceUnavailable(t *testing.T) {
	err := ServiceUnavailable("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusServiceUnavailable, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestGatewayTimeout(t *testing.T) {
	err := GatewayTimeout("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusGatewayTimeout, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestHTTPVersionNotSupported(t *testing.T) {
	err := HTTPVersionNotSupported("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusHTTPVersionNotSupported, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestVariantAlsoNegotiates(t *testing.T) {
	err := VariantAlsoNegotiates("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusVariantAlsoNegotiates, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestInsufficientStorage(t *testing.T) {
	err := InsufficientStorage("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInsufficientStorage, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestLoopDetected(t *testing.T) {
	err := LoopDetected("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusLoopDetected, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestNotExtended(t *testing.T) {
	err := NotExtended("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotExtended, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}

func TestNetworkAuthenticationRequired(t *testing.T) {
	err := NetworkAuthenticationRequired("msg", nil)
	assert.Error(t, err)
	httpError, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNetworkAuthenticationRequired, httpError.StatusCode)
	assert.Equal(t, "msg", httpError.Message)
	assert.NoError(t, httpError.Err)
}
