package httpproxy

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Records struct {
	Responses  map[string][]Response `json:"responses"`
	At         time.Time             `json:"recorded_at"`
	TokenStack []string              `json:"tokens"`
	Env        map[string]string     `json:"env"`
}

func LoadZip(path string) (Records, error) {
	r := Records{}
	f, err := zip.OpenReader(path)
	if err != nil {
		return r, fmt.Errorf("cant open api zip replay : %w", err)
	}
	defer f.Close()

	file := f.File[0]
	zippedFile, err := file.Open()
	if err != nil {
		return r, err
	}
	defer zippedFile.Close()

	err = json.NewDecoder(zippedFile).Decode(&r)
	return r, err
}

func (r Records) Clone() (Records, error) {
	out := Records{}
	buf := &bytes.Buffer{}

	err := json.NewEncoder(buf).Encode(r)
	if err != nil {
		return out, err
	}
	err = json.NewDecoder(buf).Decode(&out)
	return out, err
}

func Save(p *Recorder, path string) error {
	archive, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cant save api response : %w", err)
	}
	defer archive.Close()
	zipWriter := zip.NewWriter(archive)

	f, err := zipWriter.Create("request.json")
	if err != nil {
		return err
	}
	r := Records{
		Responses:  p.Captured,
		At:         time.Now(),
		TokenStack: p.TokenStack,
		Env:        p.Env,
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	err = enc.Encode(r)
	if err != nil {
		return err
	}

	err = zipWriter.Close()
	return err
}
