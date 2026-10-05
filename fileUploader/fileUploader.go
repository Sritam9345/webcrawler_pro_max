package fileUploader

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"mime/multipart"
	"strconv"
)





func UploadURLs(urls []string) error {
	var body bytes.Buffer

	// Create multipart form
	writer := multipart.NewWriter(&body)

	// Create an in-memory file called urls.csv
	file, err := writer.CreateFormFile("file", "urls.csv")
	if err != nil {
		return err
	}

	// Write CSV directly into the multipart file
	csvWriter := csv.NewWriter(file)

	// Header
	if err := csvWriter.Write([]string{"index", "url"}); err != nil {
		return err
	}

	// URLs
	for i, url := range urls {
		if err := csvWriter.Write([]string{
			strconv.Itoa(i),
			url,
		}); err != nil {
			return err
		}
	}

	csvWriter.Flush()

	if err := csvWriter.Error(); err != nil {
		return err
	}

	// VERY IMPORTANT
	if err := writer.Close(); err != nil {
		return err
	}

	// Create HTTP request
	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/",
		&body,
	)
	if err != nil {
		return err
	}

	// Set multipart content type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("upload failed: %s", resp.Status)
	}

	return nil
}