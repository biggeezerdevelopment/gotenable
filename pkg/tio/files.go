package tio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
)

// FilesAPI handles file upload operations.
type FilesAPI struct {
	client *Client
}

// Upload uploads a file to Tenable.io.
func (f *FilesAPI) Upload(ctx context.Context, filename string, data io.Reader, encrypted bool) (string, error) {
	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("Filedata", filepath.Base(filename))
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(part, data); err != nil {
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	// Create the request
	path := "file/upload"
	if encrypted {
		path = "file/upload?no_enc=1"
	}

	var result struct {
		Fileuploaded string `json:"fileuploaded"`
	}

	req := f.client.Request(ctx).
		SetHeader("Content-Type", writer.FormDataContentType()).
		SetBody(body.Bytes()).
		SetResult(&result)

	resp, err := req.Post(path)
	if err != nil {
		return "", err
	}

	if resp.StatusCode() >= 400 {
		return "", fmt.Errorf("upload failed with status %d", resp.StatusCode())
	}

	return result.Fileuploaded, nil
}

