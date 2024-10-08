package errors

const (
	BadRequest           = "BAD_REQUEST"
	Unauthorized         = "UNAUTHORIZED"
	Forbidden            = "FORBIDDEN"
	NotFound             = "NOT_FOUND"
	MethodNotAllowed     = "METHOD_NOT_ALLOWED"
	Conflict             = "CONFLICT"
	UnsupportedMediaType = "UNSUPPORTED_MEDIA_TYPE"
	TooManyRequests      = "TOO_MANY_REQUESTS"
	NotImplemented       = "NOT_IMPLEMENTED"
	BadGateway           = "BAD_GATEWAY"
	ServiceUnavailable   = "SERVICE_UNAVAILABLE"
	GatewayTimeout       = "GATEWAY_TIMEOUT"
	InternalServerError  = "INTERNAL_SERVER_ERROR"

	EmailAlreadyExists = "EMAIL_ALREADY_EXISTS"
	InvalidCredentials = "INVALID_CREDENTIALS"
	InvalidCode        = "INVALID_CODE"
	ExpiredCode        = "EXPIRED_CODE"
)
