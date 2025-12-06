package httpresponse

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testOption struct {
	response             *http.Response
	expectedStatusCode   int
	expectedContentType  string
	expectedBodylength   int
	expectedResponseBody string
	dontTestResponseBody bool
}

func testResponse(t *testing.T, option testOption) {
	t.Helper()
	body, err := ioutil.ReadAll(option.response.Body)
	if err != nil {
		t.Errorf("Cannot read body of response :%v", err)
	}

	if option.response.StatusCode != option.expectedStatusCode {
		t.Errorf("Expect status code %d but got %d", option.expectedStatusCode, option.response.StatusCode)
	}
	if len(option.expectedContentType) > 0 && option.response.Header.Get("Content-Type") != option.expectedContentType {
		t.Errorf("Expect body %s but got %s", option.expectedContentType, option.response.Header.Get("Content-Type"))
	}
	if option.expectedBodylength != 0 && option.expectedBodylength != len(body) {
		t.Errorf("Expect body of len %d but got %d", option.expectedBodylength, len(body))
	}
	if option.dontTestResponseBody == false && string(body) != option.expectedResponseBody {
		t.Errorf("Expect body %s but got %s", option.expectedResponseBody, string(body))
	}
}

func TestResponse(t *testing.T) {
	tests := map[string]struct {
		fn     func(http.ResponseWriter)
		option testOption
	}{
		"TestJson": {
			fn: func(w http.ResponseWriter) {
				Json(w, 200, map[string]string{"test": "test"})
			},
			option: testOption{
				expectedStatusCode:   200,
				expectedContentType:  "application/json",
				expectedResponseBody: `{"test":"test"}`,
			},
		},
		"TestJsonError": {
			fn: func(w http.ResponseWriter) {
				Json(w, 200, make(chan int))
			},
			option: testOption{
				expectedStatusCode:   500,
				expectedContentType:  "text/plain",
				expectedResponseBody: `json: unsupported type: chan int`,
			},
		},
		"TestError": {
			fn: func(w http.ResponseWriter) {
				Error(w, 200, "msg")
			},
			option: testOption{
				expectedStatusCode:   200,
				expectedContentType:  "application/json",
				expectedResponseBody: `{"error":"msg"}`,
			},
		},
		"TestOk": {
			fn: func(w http.ResponseWriter) {
				Ok(w, map[string]string{"test": "test"})
			},
			option: testOption{
				expectedStatusCode:   200,
				expectedContentType:  "application/json",
				expectedResponseBody: `{"test":"test"}`,
			},
		},
		"TestNotFound": {
			fn: func(w http.ResponseWriter) {
				NotFound(w)
			},
			option: testOption{
				expectedStatusCode:   404,
				expectedContentType:  "application/json",
				expectedResponseBody: `{"error":"entity not found"}`,
			},
		},
		"TestBadRequest": {
			fn: func(w http.ResponseWriter) {
				BadRequest(w, fmt.Errorf("test"))
			},
			option: testOption{
				expectedStatusCode:   400,
				expectedContentType:  "application/json",
				expectedResponseBody: `{"error":"test"}`,
			},
		},
		"TestCreated": {
			fn: func(w http.ResponseWriter) {
				Created(w, map[string]string{"test": "test"})
			},
			option: testOption{
				expectedStatusCode:   201,
				expectedContentType:  "application/json",
				expectedResponseBody: `{"test":"test"}`,
			},
		},
		"TestDeleted": {
			fn: func(w http.ResponseWriter) {
				Deleted(w)
			},
			option: testOption{
				expectedStatusCode:   204,
				dontTestResponseBody: true,
			},
		},
		"TestUnauthorized": {
			fn: func(w http.ResponseWriter) {
				Unauthorized(w)
			},
			option: testOption{
				expectedStatusCode:   401,
				expectedContentType:  "application/json",
				dontTestResponseBody: true,
			},
		},
		"TestConflicted": {
			fn: func(w http.ResponseWriter) {
				Conflicted(w, map[string]string{"test": "test"})
			},
			option: testOption{
				expectedStatusCode:   409,
				expectedContentType:  "application/json",
				expectedResponseBody: `{"reason":{"test":"test"}}`,
			},
		},
		"TestFile_RealPdf_LessThan512": {
			fn: func(w http.ResponseWriter) {
				File(w, ioutil.NopCloser(strings.NewReader(smallPdfContent)))
			},
			option: testOption{
				expectedStatusCode:   200,
				expectedContentType:  "application/pdf",
				expectedResponseBody: smallPdfContent,
				expectedBodylength:   len(smallPdfContent),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tt.fn(w)
			tt.option.response = w.Result()
			testResponse(t, tt.option)
		})
	}
}

const (
	smallPdfContent = `%PDF-1.2
9 0 obj
<<
>>
stream
BT/ 9 Tf(Test)' ET
endstream
endobj
4 0 obj
<<
/Type /Page
/Parent 5 0 R
/Contents 9 0 R
>>
endobj
5 0 obj
<<
/Kids [4 0 R ]
/Count 1
/Type /Pages
/MediaBox [ 0 0 99 9 ]
>>
endobj
3 0 obj
<<
/Pages 5 0 R
/Type /Catalog
>>
endobj
trailer
<<
/Root 3 0 R
>>
%%EOF`
)
