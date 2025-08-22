package types

type BizError struct {
	code    int
	message string
}

func (e *BizError) Code() int {
	return e.code
}

func (e *BizError) Error() string {
	return e.message
}

func NewBizError(code int, message string) *BizError {
	return &BizError{code, message}
}

func BizErr(message string) *BizError {
	return NewBizError(500, message)
}

var (
	ParamsError = NewBizError(400, "参数错误")
	SignError   = NewBizError(400, "参数错误2")
)
