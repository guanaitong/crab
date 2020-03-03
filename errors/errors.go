package errors

import "fmt"

type Error interface {
	Code() int
	Msg() string
	Error() string
	Unwrap() error
}

// 系统异常,1=数据库 2=redis 3=remote service
type SystemError struct {
	err       error
	errorCode int
}

func (e *SystemError) Code() int { return e.errorCode }

func (e *SystemError) Msg() string {
	if e == nil {
		return "<nil>"
	}
	if e.errorCode == 1 {
		return "db error"
	} else if e.errorCode == 2 {
		return "redis error"
	} else if e.errorCode == 3 {
		return "remote service error"
	}

	return "system error "
}

func (e *SystemError) Unwrap() error { return e.err }

func (e *SystemError) Error() string {
	if e == nil {
		return "<nil>"
	}
	m := ""
	if e.err != nil {
		m = e.err.Error()
	}
	if e.errorCode == 0 {
		return "db error: " + m
	} else if e.errorCode == 1 {
		return "redis error: " + m
	}

	return "system error " + m
}

func NewSystemDbError(err error) *SystemError {
	return &SystemError{errorCode: 1, err: err}
}

func NewSystemRedisError(err error) *SystemError {
	return &SystemError{errorCode: 2, err: err}
}

func NewSystemRemoteServiceError(err error) *SystemError {
	return &SystemError{errorCode: 1, err: err}
}

// 业务异常
type BusinessError struct {
	code int
	msg  string
}

func (e *BusinessError) Code() int { return e.code }

func (e *BusinessError) Msg() string { return e.msg }

func (e *BusinessError) Unwrap() error { return nil }

func (e *BusinessError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("errCode is %d,errorMsg is %s", e.code, e.msg)
}

func NewBusinessError(code int, msg string) *BusinessError {
	return &BusinessError{code: code, msg: msg}
}

type Code struct {
	Code int
	Msg  string
}

// 远程服务的异常
type ApiError struct {
	code int
	msg  string
	err  error
}

func (e *ApiError) Code() int { return e.code }

func (e *ApiError) Msg() string { return e.msg }

func (e *ApiError) Unwrap() error { return e.err }

func (e *ApiError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.err != nil {
		return e.err.Error()
	}
	return fmt.Sprintf("errCode is %d,errorMsg is %s", e.code, e.msg)
}

func (e *ApiError) IsOK() bool {
	if e == nil {
		return true
	}
	return e.err == nil && e.code == 0
}

func NewApiError(code int, msg string, err error) *ApiError {
	return &ApiError{code: code, msg: msg, err: err}
}
