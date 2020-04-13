package errors

import (
	"fmt"
	"github.com/guanaitong/crab/system"
)

var (
	OK                          = 0
	DbErrCodeDefault            = GenerateCode(1, 10)
	RedisErrCodeDefault         = GenerateCode(1, 20)
	RemoteServiceErrCodeDefault = GenerateCode(4, 20)
)

func GenerateCode(secondErrCode, thirdErrCode int) int {
	if secondErrCode >= 10 || secondErrCode < 0 {
		panic("second err code is invalid")
	}
	if thirdErrCode >= 1000 || thirdErrCode < 0 {
		panic("third err code is invalid")
	}
	return system.GetErrCodePrefix()*10000 + secondErrCode*1000 + thirdErrCode
}

type Error interface {
	ErrorCode() int
	ErrorMsg() string
	Error() string
	Unwrap() error
}

// 系统异常:1=数据库,2=redis,3=remote service
type systemError struct {
	Err  error
	Code int
}

func (e *systemError) ErrorCode() int { return e.Code }

func (e *systemError) ErrorMsg() string {
	if e == nil {
		return "<nil>"
	}
	if e.Code == DbErrCodeDefault {
		return "db error"
	} else if e.Code == RedisErrCodeDefault {
		return "redis error"
	} else if e.Code == RemoteServiceErrCodeDefault {
		return "remote service error"
	}

	return "system error "
}

func (e *systemError) Unwrap() error { return e.Err }

func (e *systemError) Error() string {
	if e == nil {
		return "<nil>"
	}
	m := ""
	if e.Err != nil {
		m = e.Err.Error()
	}
	if e.Code == DbErrCodeDefault {
		return "db error: " + m
	} else if e.Code == RedisErrCodeDefault {
		return "redis error: " + m
	} else if e.Code == RemoteServiceErrCodeDefault {
		return "remote service error: " + m
	}

	return "system error " + m
}

// 业务异常
type commonError struct {
	Code int
	Msg  string
}

func (e *commonError) ErrorCode() int { return e.Code }

func (e *commonError) ErrorMsg() string { return e.Msg }

func (e *commonError) Unwrap() error { return nil }

func (e *commonError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("errCode is %d,errorMsg is %s", e.Code, e.Msg)
}

// 远程服务的异常
type remoteServiceError struct {
	Code int
	Msg  string
	Err  error
}

func (e *remoteServiceError) ErrorCode() int { return e.Code }

func (e *remoteServiceError) ErrorMsg() string { return e.Msg }

func (e *remoteServiceError) Unwrap() error { return e.Err }

func (e *remoteServiceError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return fmt.Sprintf("errCode is %d,errorMsg is %s", e.Code, e.Msg)
}

func (e *remoteServiceError) IsOK() bool {
	if e == nil {
		return true
	}
	return e.Err == nil && e.Code == 0
}

func NewDbError(err error) Error {
	if err == nil {
		return nil
	}
	return &systemError{Err: err, Code: DbErrCodeDefault}
}

func NewRedisError(err error) Error {
	if err == nil {
		return nil
	}
	return &systemError{Err: err, Code: RedisErrCodeDefault}
}

func NewParamError(thirdErrCode int, msg string) Error {
	code := GenerateCode(2, thirdErrCode)
	return &commonError{Code: code, Msg: msg}
}

func NewBusinessError(thirdErrCode int, msg string) Error {
	code := GenerateCode(3, thirdErrCode)
	return &commonError{Code: code, Msg: msg}
}

func NewRemoteServiceError(err error) Error {
	if err == nil {
		return nil
	}
	return &systemError{Err: err, Code: RemoteServiceErrCodeDefault}
}

func Cast(err error) Error {
	if err == nil {
		return nil
	}
	return err.(Error)
}
