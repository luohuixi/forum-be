package err

var (
	// Common errors
	OK                  = &Errno{Code: 0, Message: "OK"}
	ErrDatabase         = &Errno{Code: 10001, Message: "Database error"}
	ErrBind             = &Errno{Code: 10002, Message: "Error occurred while binding the request body to the struct."}
	ErrBadRequest       = &Errno{Code: 10003, Message: "Request error"}
	ErrUserExisted      = &Errno{Code: 10004, Message: "User has existed"}
	ErrAuthToken        = &Errno{Code: 10005, Message: "Error occurred while handling the auth token"}
	ErrUserNotExisted   = &Errno{Code: 10006, Message: "User not existed"}
	ErrGetRedisList     = &Errno{Code: 10007, Message: "Get list from Redis out of expiration time"}
	ErrRewriteRedisList = &Errno{Code: 10008, Message: "rewrite list to Redis when cancel"}
	ErrRedis            = &Errno{Code: 10009, Message: "Redis error"}

	// oauth errors
	ErrRegister          = &Errno{Code: 20001, Message: "Error occurred while registering on auth-server"}
	ErrRemoteAccessToken = &Errno{Code: 20002, Message: "Error occurred while getting oauth access token from auth-server"}
	ErrLocalAccessToken  = &Errno{Code: 20003, Message: "Error occurred while getting oauth access token from local"}
	ErrGetUserInfo       = &Errno{Code: 20004, Message: "Error occurred while getting user info from oauth-server by access token"}
)

type Errno struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err Errno) Error() string {
	return err.Message
}

// // Err represents an error
// type Err struct {
// 	Code    int
// 	Message string
// 	Err     error
// }

// func New(errno *Errno, err error) *Err {
// 	return &Err{Code: errno.Code, Message: errno.Message, Err: err}
// }

// func (err *Err) Add(message string) error {
// 	err.Message += " " + message
// 	return err
// }

// func (err *Err) Addf(format string, args ...interface{}) error {
// 	err.Message += " " + fmt.Sprintf(format, args...)
// 	return err
// }

// func (err *Err) Error() string {
// 	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", err.Code, err.Message, err.Err)
// }

// // func IsErrUserNotFound(err error) bool {
// // 	code, _ := DecodeErr(err)
// // 	return code == ErrUserNotFound.Code
// // }

// func DecodeErr(err error) (int, string) {
// 	if err == nil {
// 		return OK.Code, OK.Message
// 	}

// 	switch typed := err.(type) {
// 	case *Err:
// 		return typed.Code, typed.Message
// 	case *Errno:
// 		return typed.Code, typed.Message
// 	default:
// 	}

// 	return InternalServerError.Code, err.Error()
// }
