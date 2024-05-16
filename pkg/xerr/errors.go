package xerr

import "github.com/zeromicro/x/errors"

func New(code int, msg string) error {
	return errors.New(code, msg)
}

func NewMsg(msg string) error {
	return errors.New(ServerCommonError, msg)
}

func NewDBErr() error {
	return errors.New(DbError, ErrMsg(DbError))
}

func NewInternalErr() error {
	return errors.New(ServerCommonError, ErrMsg(ServerCommonError))
}
