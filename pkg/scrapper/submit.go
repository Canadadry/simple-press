package scrapper

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type FormValue struct {
	String   string
	Filename string
	File     io.Reader
}

func (c *Client) Submit(name string, fields map[string]FormValue) error {
	if c.page == nil {
		return fmt.Errorf("cannot submit : the current scrapper page is nil")
	}

	form, err := GetForm(c.page, name)
	if err != nil {
		return fmt.Errorf("while trying to submit to %s : %w", name, err)
	}

	for fname := range fields {
		_, ok := form.Attribute[cleanFieldname(fname)]
		if !ok {
			return fmt.Errorf("field  %s dont exist on form %s", fname, name)
		}
	}

	body, bodyCtIn, err := createBodyRequest(fields)
	if err != nil {
		return fmt.Errorf("can't create body request of %s : %w", name, err)
	}

	previousURL := c.previousURL
	c.previousURL = c.currentURL
	currentURL := c.currentURL
	if form.Action != "" {
		currentURL = form.Action
	}

	req, err := http.NewRequest(strings.ToUpper(form.Method), c.endpoint+currentURL, body)
	if err != nil {
		return fmt.Errorf("cannot create form %s request %v", name, err)
	}
	req.Header.Add("Content-Type", bodyCtIn)
	req.Header.Add("Referer", previousURL)
	return c.makeRequest(req)
}

func createBodyRequest(f map[string]FormValue) (io.Reader, string, error) {
	for _, v := range f {
		if v.Filename != "" {
			return createMultiPartForm(f)
		}
	}
	return createURLEncodeForm(f)
}

func createURLEncodeForm(f map[string]FormValue) (io.Reader, string, error) {
	data := url.Values{}
	for k, v := range f {
		data.Set(k, v.String)
	}
	return strings.NewReader(data.Encode()), "application/x-www-form-urlencoded", nil
}

func createMultiPartForm(f map[string]FormValue) (io.Reader, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, fvalue := range f {
		if fvalue.String != "" {
			err := writer.WriteField(key, fvalue.String)
			if err != nil {
				return nil, "", fmt.Errorf("cant add field %s : %v", key, err)
			}
		} else if fvalue.File != nil {
			part, err := writer.CreateFormFile(cleanFieldname(key), fvalue.Filename)
			if err != nil {
				return nil, "", fmt.Errorf("cant add file  %s :  %w", key, err)
			}
			_, err = io.Copy(part, fvalue.File)
			if err != nil {
				return nil, "", fmt.Errorf("cant copy file %s into request body:  %w", key, err)
			}
		}
	}
	err := writer.Close()
	if err != nil {
		return nil, "", fmt.Errorf("cant close request body writer %w", err)
	}
	return body, writer.FormDataContentType(), nil
}

// TODO: Dirty way to manage the upload of several files
func cleanFieldname(fieldname string) string {
	r := regexp.MustCompile("\\+\\d+$")

	return r.ReplaceAllString(fieldname, "")
}
