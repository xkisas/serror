package serror

type ErrorWrapper interface {
	Error() string
	Unwrap() error
	Code() int
	Position() string
	Throw()
}

var (
	RootError       = errorWrap{msg: "error", code: 500}
	UnexpectedError = errorWrap{msg: "unexpected error", err: RootError, code: 500}
)

type errorWrap struct {
	msg      string
	err      error
	code     int
	position string
}

func (e errorWrap) Error() string    { return e.msg }
func (e errorWrap) Unwrap() error    { return e.err }
func (e errorWrap) Code() int        { return e.code }
func (e errorWrap) Position() string { return e.position }
func (e errorWrap) Throw()           { throwWithCallerDepth(e, 2) }
