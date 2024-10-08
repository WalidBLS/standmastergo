package errors

import "net/http"

type CustomError struct {
	Key string
	Err error
}

func (ce CustomError) Error() string {
	return ce.Err.Error()
}

func (ce CustomError) StatusCode() int {
	switch ce.Key {
	case BadRequest:
		return http.StatusBadRequest
	case Unauthorized:
	case InvalidCredentials:
	case InvalidCode:
	case ExpiredCode:
		return http.StatusUnauthorized
	case Forbidden:
		return http.StatusForbidden
	case NotFound:
		return http.StatusNotFound
	case MethodNotAllowed:
		return http.StatusMethodNotAllowed
	case Conflict:
	case EmailAlreadyExists:
		return http.StatusConflict
	case UnsupportedMediaType:
		return http.StatusUnsupportedMediaType
	case TooManyRequests:
		return http.StatusTooManyRequests
	case NotImplemented:
		return http.StatusNotImplemented
	case BadGateway:
		return http.StatusBadGateway
	case ServiceUnavailable:
		return http.StatusServiceUnavailable
	case GatewayTimeout:
		return http.StatusGatewayTimeout
	case InternalServerError:
	default:
		return http.StatusInternalServerError
	}

	return http.StatusInternalServerError
}
