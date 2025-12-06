package httpcache

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Response struct {
	StatusCode int
	Body       []byte
	Header     map[string][]string
}

func saveResponse(rsp *http.Response) (Response, error) {
	defer rsp.Body.Close()
	buf := &bytes.Buffer{}
	_, err := io.Copy(buf, rsp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("cannot read server response : %w", err)
	}
	return Response{
		StatusCode: rsp.StatusCode,
		Body:       buf.Bytes(),
		Header:     rsp.Header,
	}, nil
}

func newResponseFrom(rsp Response) (*http.Response, error) {
	return &http.Response{
		StatusCode: rsp.StatusCode,
		Header:     rsp.Header,
		Body:       io.NopCloser(bytes.NewReader(rsp.Body)),
	}, nil
}
