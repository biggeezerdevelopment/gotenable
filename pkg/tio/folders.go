package tio

import (
	"context"
	"fmt"
)

// FoldersAPI handles folder operations.
type FoldersAPI struct {
	client *Client
}

// Folder represents a scan folder.
type Folder struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	DefaultTag       int    `json:"default_tag"`
	Custom           int    `json:"custom"`
	UnreadCount      int    `json:"unread_count,omitempty"`
}

// List retrieves all folders.
func (f *FoldersAPI) List(ctx context.Context) ([]Folder, error) {
	var result struct {
		Folders []Folder `json:"folders"`
	}

	_, err := f.client.Get(ctx, "folders", &result)
	if err != nil {
		return nil, err
	}

	return result.Folders, nil
}

// Create creates a new folder.
func (f *FoldersAPI) Create(ctx context.Context, name string) (*Folder, error) {
	payload := map[string]string{"name": name}
	var result Folder
	_, err := f.client.Post(ctx, "folders", payload, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete removes a folder.
func (f *FoldersAPI) Delete(ctx context.Context, folderID int) error {
	_, err := f.client.Delete(ctx, fmt.Sprintf("folders/%d", folderID))
	return err
}

// Edit updates a folder's name.
func (f *FoldersAPI) Edit(ctx context.Context, folderID int, name string) error {
	payload := map[string]string{"name": name}
	_, err := f.client.Put(ctx, fmt.Sprintf("folders/%d", folderID), payload, nil)
	return err
}

