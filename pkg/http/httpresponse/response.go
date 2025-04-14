package httpresponse

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var PrettyPrint = false

type ResponseError interface {
	ContentType() string
	Write(io.Writer)
}

func GenericResponse(status int, payload interface{}) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		Json(w, status, payload)
	}
}

func Json(w http.ResponseWriter, status int, payload interface{}) error {
	marshal := json.Marshal
	if PrettyPrint {
		marshal = func(p interface{}) ([]byte, error) {
			return json.MarshalIndent(p, "", "    ")
		}
	}
	response, err := marshal(payload)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
	return nil
}

func Empty(w http.ResponseWriter, status int) error {
	w.Header().Set("Content-Length", "0")
	w.WriteHeader(status)
	return nil
}

func Error(w http.ResponseWriter, status int, message string) error {
	return Json(w, status, map[string]string{"error": message})
}

func ServerError(w http.ResponseWriter) error {
	return Error(w, http.StatusInternalServerError, "server error")
}

func Ok(w http.ResponseWriter, payload interface{}) error {
	return Json(w, http.StatusOK, payload)
}

func List[T any](w http.ResponseWriter, limit, offset, count int, entities []T) error {
	return Json(w, http.StatusOK, map[string]interface{}{
		"limit":    limit,
		"offset":   offset,
		"count":    count,
		"entities": entities,
	})
}

func NotFound(w http.ResponseWriter) error {
	return Error(w, http.StatusNotFound, "entity not found")
}

func RouteNotFound(w http.ResponseWriter, path string) error {
	return Error(w, http.StatusNotFound, fmt.Sprintf("route (%s) not found", path))
}

func BadRequest(w http.ResponseWriter, err error) error {
	respErr, ok := err.(ResponseError)
	if ok {
		w.Header().Set("Content-Type", respErr.ContentType())
		w.WriteHeader(http.StatusBadRequest)
		respErr.Write(w)
		return nil
	}
	return Error(w, http.StatusBadRequest, err.Error())
}

func Created(w http.ResponseWriter, payload interface{}) error {
	return Json(w, http.StatusCreated, payload)
}

func Deleted(w http.ResponseWriter) error {
	return Json(w, http.StatusNoContent, map[string]string{})
}

func MethodNotAllowed(w http.ResponseWriter) error {
	return Json(w, http.StatusMethodNotAllowed, map[string]string{})
}

func NoContent(w http.ResponseWriter) error {
	return Empty(w, http.StatusNoContent)
}

func Unauthorized(w http.ResponseWriter) error {
	return Error(w, http.StatusUnauthorized, "unauthorized")
}

func Forbidden(w http.ResponseWriter) error {
	return Json(w, http.StatusForbidden, "forbidden")
}

func Conflicted(w http.ResponseWriter, payload interface{}) error {
	respErr, ok := payload.(ResponseError)
	if ok {
		w.Header().Set("Content-Type", respErr.ContentType())
		w.WriteHeader(http.StatusConflict)
		respErr.Write(w)
		return nil
	}
	return Json(w, http.StatusConflict, map[string]interface{}{"reason": payload})
}

func File(w http.ResponseWriter, f io.ReadCloser) error {
	defer f.Close()
	filestart := make([]byte, 512)
	n, err := f.Read(filestart)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", http.DetectContentType(filestart))
	w.WriteHeader(http.StatusOK)
	w.Write(filestart[:n])
	_, err = io.Copy(w, f)
	return err
}

func Bytes(w http.ResponseWriter, b []byte) error {
	w.Header().Set("Content-Type", http.DetectContentType(b))
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(b)
	return err
}

func TooManyRequest(w http.ResponseWriter, err error) error {
	respErr, ok := err.(ResponseError)
	if ok {
		w.Header().Set("Content-Type", respErr.ContentType())
		w.WriteHeader(http.StatusTooManyRequests)
		respErr.Write(w)
		return nil
	}
	return Error(w, http.StatusTooManyRequests, err.Error())
}
