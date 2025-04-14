package stacktrace

import (
	"fmt"
	"runtime"
	"strings"
)

var defaultFrameCapacity = 20

func From(err error) *Error {
	if err == nil {
		return nil
	}
	e, ok := err.(*Error)
	if ok {
		return e
	}
	return trace(err, 2)
}

func Errorf(message string, args ...interface{}) *Error {
	err := fmt.Errorf(message, args...)
	for _, arg := range args {
		a, ok := arg.(*Error)
		if ok {
			return &Error{
				Err:    err,
				Frames: a.Frames,
			}
		}
	}
	return trace(fmt.Errorf(message, args...), 2)
}

func trace(err error, skip int) *Error {
	frames := make([]Frame, 0, defaultFrameCapacity)
	for {
		pc, path, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		frame := Frame{
			Func: fn.Name(),
			Line: line,
			Path: path,
		}
		frames = append(frames, frame)
		skip++
	}
	return &Error{
		Err:     err,
		Message: err.Error(),
		Frames:  frames,
	}
}

type Frame struct {
	Func string
	Line int
	Path string
}

func (f Frame) String() string {
	return fmt.Sprintf("%s:%d %s()", f.Path, f.Line, f.Func)
}

type Error struct {
	Err     error
	Message string
	Frames  []Frame
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) Trace() string {
	expectedRows := len(e.Frames) + 1
	rows := make([]string, 0, expectedRows)
	rows = append(rows, e.Error())
	for _, frame := range e.Frames {
		rows = append(rows, frame.String())
	}
	return strings.Join(rows, "\n")
}
