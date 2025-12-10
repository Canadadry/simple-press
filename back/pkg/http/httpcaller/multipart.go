package httpcaller

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
)

type FormValue struct {
	String   string
	Filename string
	File     io.Reader
}

func CreateMultiPartForm(f map[string][]FormValue) (io.Reader, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, farray := range f {
		for _, fvalue := range farray {
			if fvalue.String != "" {
				err := writer.WriteField(key, fvalue.String)
				if err != nil {
					return nil, "", fmt.Errorf("cant add field %s : %v", key, err)
				}
			} else if fvalue.File != nil {
				part, err := writer.CreateFormFile(key, fvalue.Filename)
				if err != nil {
					return nil, "", fmt.Errorf("cant add file  %s :  %w", key, err)
				}
				_, err = io.Copy(part, fvalue.File)
				if err != nil {
					return nil, "", fmt.Errorf("cant copy file %s into request body:  %w", key, err)
				}
			}
		}
	}
	err := writer.Close()
	if err != nil {
		return nil, "", fmt.Errorf("cant close request body writer %w", err)
	}
	return body, writer.FormDataContentType(), nil
}
