package tt_pay

import (
	"fmt"
)

type Error struct {
	Code    string
	Msg     string
	SubCode string
	SubMsg  string
	Detail  string
}



const (
	ErrorPattern = `{"code": "%s", "msg": "%s", "sub_code": "%s", "sub_msg": "%s", "detail": "%s"}`
)

func (e *Error) Error() string {
	if e == nil {
		return "e is nil"
	}

	return fmt.Sprintf(ErrorPattern, e.Code, e.Msg, e.SubCode, e.SubMsg, e.Detail)
}
