package errors

import "fmt"

const (
	OK = iota
	DbError
	RedisError
	RemoteServiceError
)

type Error interface {
	ErrorCode() int
	ErrorMsg() string
	Error() string
	Unwrap() error
}

// 系统异常:1=数据库,2=redis,3=remote service
type SystemError struct {
	Err  error
	Code int
}

func (e *SystemError) ErrorCode() int { return e.Code }

func (e *SystemError) ErrorMsg() string {
	if e == nil {
		return "<nil>"
	}
	if e.Code == DbError {
		return "db error"
	} else if e.Code == RedisError {
		return "redis error"
	} else if e.Code == RemoteServiceError {
		return "remote service error"
	}

	return "system error "
}

func (e *SystemError) Unwrap() error { return e.Err }

func (e *SystemError) Error() string {
	if e == nil {
		return "<nil>"
	}
	m := ""
	if e.Err != nil {
		m = e.Err.Error()
	}
	if e.Code == DbError {
		return "db error: " + m
	} else if e.Code == RedisError {
		return "redis error: " + m
	} else if e.Code == RemoteServiceError {
		return "remote service error: " + m
	}

	return "system error " + m
}

func NewDbError(err error) *SystemError {
	if err == nil {
		return nil
	}
	return &SystemError{Err: err, Code: DbError}
}

func NewRedisError(err error) *SystemError {
	if err == nil {
		return nil
	}
	return &SystemError{Err: err, Code: RedisError}
}

func NewRemoteServiceError(err error) *SystemError {
	if err == nil {
		return nil
	}
	return &SystemError{Err: err, Code: RemoteServiceError}
}

func Cast(err error) Error {
	if err == nil {
		return nil
	}
	return err.(Error)
}

// 业务异常
type BusinessError struct {
	Code int
	Msg  string
}

func (e *BusinessError) ErrorCode() int { return e.Code }

func (e *BusinessError) ErrorMsg() string { return e.Msg }

func (e *BusinessError) Unwrap() error { return nil }

func (e *BusinessError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("errCode is %d,errorMsg is %s", e.Code, e.Msg)
}

// 远程服务的异常
type ApiError struct {
	Code int
	Msg  string
	Err  error
}

func (e *ApiError) ErrorCode() int { return e.Code }

func (e *ApiError) ErrorMsg() string { return e.Msg }

func (e *ApiError) Unwrap() error { return e.Err }

func (e *ApiError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return fmt.Sprintf("errCode is %d,errorMsg is %s", e.Code, e.Msg)
}

func (e *ApiError) IsOK() bool {
	if e == nil {
		return true
	}
	return e.Err == nil && e.Code == 0
}
