package middleware

import (
	"app/pkg/clock"
	"app/pkg/router"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var PrettyPrint = false

type ctxSQLLoggerKey struct{}
type ctxSQLLoggerHideParamKey struct{}

func LogSQLQuery(ctx context.Context, mode, query string, time time.Duration, params ...interface{}) {
	value := ctx.Value(ctxSQLLoggerKey{})
	fn, ok := value.(func(queryLogLine))
	if !ok {
		return
	}
	hide := ctx.Value(ctxSQLLoggerHideParamKey{})
	if hide != nil {
		fn(queryLogLine{mode, query, fmt.Sprintf("%v", time), nil})
		return
	}
	fn(queryLogLine{mode, query, fmt.Sprintf("%v", time), params})
}

func HideSQLQueryParam(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxSQLLoggerHideParamKey{}, true)
}

type queryLogLine struct {
	Mode     string
	Query    string
	Duration string
	Params   []interface{}
}

func Logger(w io.Writer, clock clock.Clock, exclude func(string) bool) func(next router.HandlerFunc) router.HandlerFunc {
	marshal := json.Marshal
	if PrettyPrint {
		marshal = func(p interface{}) ([]byte, error) {
			return json.MarshalIndent(p, "", "    ")
		}
	}
	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) error {
			queries := []queryLogLine{}
			ww := &wrapWriter{w: rw, statusCode: http.StatusOK}
			ctx := context.WithValue(r.Context(), ctxSQLLoggerKey{}, func(q queryLogLine) {
				queries = append(queries, q)
			})
			startedAt := clock.Now()
			err := next(ww, r.WithContext(ctx))
			if err == nil && exclude != nil && exclude(r.URL.String()) {
				return nil
			}
			endedAt := clock.Now()
			logLine := struct {
				At           time.Time
				Level        string
				Method       string
				URL          string
				Pattern      string
				Status       int
				Bytes        int
				Elapsed      string
				Error        error
				ErrorMessage string
				SQL          []queryLogLine
			}{
				At:      startedAt,
				Level:   "trace",
				Method:  r.Method,
				URL:     r.URL.String(),
				Pattern: router.GetPattern(r),
				Status:  ww.statusCode,
				Bytes:   ww.byteWritten,
				Elapsed: fmt.Sprintf("%v", endedAt.Sub(startedAt)),
				Error:   err,
				SQL:     queries,
			}

			if err != nil {
				logLine.Status = 500
				logLine.ErrorMessage = err.Error()
			}

			if logLine.Status >= 500 {
				logLine.Level = "error"
			}
			out, _ := marshal(logLine)
			fmt.Fprintf(w, "%s\n", string(out))
			return err
		}
	}
}

type wrapWriter struct {
	w           http.ResponseWriter
	statusCode  int
	byteWritten int
}

func (ww *wrapWriter) Header() http.Header {
	return ww.w.Header()
}
func (ww *wrapWriter) Write(b []byte) (int, error) {
	n, err := ww.w.Write(b)
	ww.byteWritten += n
	return n, err
}
func (ww *wrapWriter) WriteHeader(statusCode int) {
	ww.statusCode = statusCode
	ww.w.WriteHeader(statusCode)
}
