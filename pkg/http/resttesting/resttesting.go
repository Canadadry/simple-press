package resttesting

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type TestOption struct {
	Handler              func(http.ResponseWriter, *http.Request)
	Method               string
	Url                  string
	Body                 string
	Headers              map[string]string
	ExpectedStatusCode   int
	ExpectedContentType  string
	ExpectedResponseBody func(*testing.T, string)
	Context              map[string]interface{}
}

func Match(ExpectedResponseBody string) func(*testing.T, string) {
	return func(t *testing.T, body string) {
		if string(body) != ExpectedResponseBody {
			t.Fatalf("Expect body \n%s\n but got \n%s", ExpectedResponseBody, string(body))
		}
	}
}

func TestRequest(t *testing.T, option TestOption) {

	bodyReader := strings.NewReader(option.Body)

	req, err := http.NewRequest(option.Method, option.Url, bodyReader)
	if err != nil {
		t.Fatalf("Cannot build query '%s' :%v", option.Url, err)
	}
	for k, v := range option.Headers {
		req.Header.Add(k, v)
	}
	for k, v := range option.Context {
		ctx := context.WithValue(req.Context(), k, v)
		req = req.WithContext(ctx)
	}
	w := httptest.NewRecorder()
	option.Handler(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Cannot read body of response :%v", err)
	}

	if resp.StatusCode != option.ExpectedStatusCode {
		t.Errorf("Expect status code %d but got %d \n %s", option.ExpectedStatusCode, resp.StatusCode, string(body))
	}
	if len(option.ExpectedContentType) > 0 && resp.Header.Get("Content-Type") != option.ExpectedContentType {
		t.Fatalf("Expect body %s but got %s", option.ExpectedContentType, resp.Header.Get("Content-Type"))
	}
	if option.ExpectedResponseBody != nil {
		option.ExpectedResponseBody(t, string(body))
	}
}
