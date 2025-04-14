package router

import (
	"app/pkg/http/httpresponse"
	"fmt"
	"net/http"
)

func redirectToIndexOpen(fs http.FileSystem, index, name string) (http.File, error) {
	f, err := fs.Open(name)
	if err != nil {
		return fs.Open(index)
	}
	d, err := f.Stat()
	if err != nil {
		return fs.Open(index)
	}
	if d.IsDir() {
		f.Close()
		return fs.Open(index)
	}
	return f, nil
}

func FileServer(root http.FileSystem, index string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		url := GetField(r, 0)
		fmt.Printf("%#v\n", url)
		f, err := redirectToIndexOpen(root, index, url)
		if err != nil {
			return err
		}
		defer f.Close()
		return httpresponse.File(w, f)
	}
}
