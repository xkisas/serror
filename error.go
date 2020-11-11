package serror

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
)

func New(msg string, parent ErrorWrapper, code ...int) ErrorWrapper {
	if len(code) > 0 {
		return errorWrap{msg: msg, err: parent, code: code[0]}
	}
	return errorWrap{msg: msg, err: parent, code: parent.Code()}
}

func CheckError(err error) {
	if err != nil {
		if !errors.Is(err, RootError) {
			err = New(err.Error(), UnexpectedError)
		}
		throwWithCallerDepth(err, 2)
	}
}

func throwWithCallerDepth(err error, callerDepth int) {
	var fileLine string
	if _, file, line, ok := runtime.Caller(callerDepth); ok {
		fileLine = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}
	if e, ok := err.(ErrorWrapper); ok {
		err = errorWrap{msg: e.Error(), err: e, code: e.Code(), position: fileLine}
	} else {
		err = errorWrap{msg: err.Error(), err: UnexpectedError, code: UnexpectedError.code, position: fileLine}
	}
	panic(err)
}

func ThrowMsg(msg string, parent ErrorWrapper) {
	throwWithCallerDepth(New(msg, parent), 2)
}

func ThrowMsgWithCallerDepth(msg string, parent ErrorWrapper, callerDepth int) {
	throwWithCallerDepth(New(msg, parent), callerDepth)
}

func TryCatch(try func(), catch func(err error), errs ...error) {
	defer func() {
		if recv := recover(); recv != nil {
			if e, ok := recv.(error); ok {
				for _, err := range errs {
					if errors.Is(e, err) {
						catch(e)
						return
					}
				}
			}
			panic(recv)
		}
	}()
	try()
}
