package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

func main() {
	input := flag.String("i", "", "Input file")
	output := flag.String("o", "", "Output file")
	help := flag.Bool("h", false, "Show usage information")
	flag.Parse()

	if *help || *input == "" || *output == "" {
		fmt.Println("Usage: go run main.go -i <input_file> -o <output_file>")
		fmt.Println("Generates a multipart/form-data request body with the file and a filename field.")
		return
	}

	info, err := os.Stat(*input)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Error: input file '%s' does not exist.\n", *input)
			os.Exit(1)
		} else {
			fmt.Printf("Error checking input file: %v\n", err)
			os.Exit(1)
		}
	}
	if info.IsDir() {
		fmt.Printf("Error: '%s' is a directory, not a file.\n", *input)
		os.Exit(1)
	}

	file, err := os.Open(*input)
	if err != nil {
		fmt.Printf("Error opening file '%s': %v\n", *input, err)
		os.Exit(1)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("content", *input)
	if err != nil {
		fmt.Printf("Error creating form file: %v\n", err)
		os.Exit(1)
	}
	if _, err := io.Copy(part, file); err != nil {
		fmt.Printf("Error copying file content: %v\n", err)
		os.Exit(1)
	}

	if err := writer.WriteField("name", info.Name()); err != nil {
		fmt.Printf("Error writing filename field: %v\n", err)
		os.Exit(1)
	}

	if err := writer.Close(); err != nil {
		fmt.Printf("Error finalizing multipart body: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(*output, body.Bytes(), 0644); err != nil {
		fmt.Printf("Error writing output file '%s': %v\n", *output, err)
		os.Exit(1)
	}

	fmt.Println("Content-Type:", writer.FormDataContentType())
	fmt.Printf("Multipart body written to '%s'\n", *output)
}
