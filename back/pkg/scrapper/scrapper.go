package scrapper

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type Jar struct {
	c map[string]string
}

func (j *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	for _, c := range cookies {
		j.c[c.Name] = c.Value
	}
}
func (j *Jar) Cookies(u *url.URL) []*http.Cookie {
	cookies := []*http.Cookie{}
	for name, value := range j.c {
		cookies = append(cookies, &http.Cookie{
			Name:  name,
			Value: value,
		})
	}

	return cookies
}

type Client struct {
	endpoint         string
	http             *http.Client
	ValidateResponse func(r *http.Response) error
	CheckRedirect    func(req *http.Request, via []*http.Request) error
	currentURL       string
	previousURL      string
	page             *Document
}

func New(api string) *Client {
	c := &Client{
		endpoint: api,
		http: &http.Client{
			Jar: &Jar{c: map[string]string{}},
		},
	}
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		c.previousURL = c.currentURL
		c.currentURL = req.URL.Path
		return nil
	}
	return c
}

func (c *Client) makeRequest(r *http.Request) error {


	c.http.CheckRedirect = c.CheckRedirect
	resp, err := c.http.Do(r)
	if err != nil {
		return fmt.Errorf("can't make request to %s : %w", r.URL.Path, err)
	}
	defer resp.Body.Close()
	if c.ValidateResponse != nil {
		err = c.ValidateResponse(resp)
		if err != nil {
			return fmt.Errorf("could not validate response : %w", err)
		}
	}

	doc, err := NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("could not parse response : %w", err)
	}
	c.page = doc
	return nil
}

func (c *Client) Get(path string) error {
	c.previousURL = c.currentURL
	c.currentURL = path

	req, err := http.NewRequest(http.MethodGet, c.endpoint+path, nil)
	if err != nil {
		return fmt.Errorf("can't create request to %s : %w", c.endpoint+path, err)
	}

	return c.makeRequest(req)
}

func (c *Client) Post(path string, form url.Values) error {
	c.previousURL = c.currentURL
	c.currentURL = path
	req, err := http.NewRequest(http.MethodPost, c.endpoint+path, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("can't create request to %s : %w", c.endpoint+path, err)
	}
	return c.makeRequest(req)
}

func (c *Client) GoBack() error {
	return c.Get(c.previousURL)
}

func (c *Client) Find(selector string) (*Selection, error) {
	return c.page.Find(selector)
}

func (c *Client) Render(w io.Writer) error {
	return html.Render(w, c.page.RootNode)
}

func (c *Client) String() string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "url: %s \n\n", c.currentURL)
	err := html.Render(buf, c.page.RootNode)
	if err != nil {
		return fmt.Sprintf("error: %v", err.Error())
	}
	return buf.String()
}

func (c *Client) GetFromReader(r io.Reader) error {
	doc, err := NewDocumentFromReader(r)
	if err != nil {
		return fmt.Errorf("could not parse response : %w", err)
	}
	c.page = doc
	return nil
}
