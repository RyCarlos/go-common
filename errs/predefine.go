package errs

const (
	NoneError = 0
	// 常规错误
	UnknownError          = 1000 // 未知错误
	ArgsError             = 1001 // 参数错误
	RouteNotFoundError    = 1002 // 路由不存在
	HttpMethodError       = 1003 // Http方法不存在
	PermissionDeniedError = 1004 // 权限不足
	TokenEmptyError       = 1005 // token为空
	TokenInvalidError     = 1006 // token无效
	TokenExpiredError     = 1007 // token过期
	DataNotFoundError     = 1008 // 数据不存在
	// 系统管理模块
	AdminNotExistError = 1500 // 管理员不存在
	AdminDisabledError = 1501 // 管理员被禁用
)

var (
	// 系统错误
	ErrUnknown          = NewErrorCode(UnknownError, "UnknownError")
	ErrArgs             = NewErrorCode(ArgsError, "ArgsError")
	ErrRoute            = NewErrorCode(RouteNotFoundError, "RouteNotFoundError")
	ErrHttpMethod       = NewErrorCode(HttpMethodError, "HttpMethodError")
	ErrPermissionDenied = NewErrorCode(PermissionDeniedError, "PermissionDeniedError")
	ErrTokenEmpty       = NewErrorCode(TokenEmptyError, "TokenEmptyError")
	ErrTokenInvalid     = NewErrorCode(TokenInvalidError, "TokenInvalidError")
	ErrTokenExpired     = NewErrorCode(TokenExpiredError, "TokenExpiredError")
	ErrDataNotFound     = NewErrorCode(DataNotFoundError, "DataNotFoundError")
	// 系统管理模块
	ErrAdminDisabled = NewErrorCode(AdminDisabledError, "AdminDisabledError")
	ErrAdminNotExist = NewErrorCode(AdminNotExistError, "AdminNotExistError")
)
