package tio

import (
	"context"
	"fmt"
)

// UsersAPI handles user-related operations.
type UsersAPI struct {
	client *Client
}

// User represents a Tenable.io user.
type User struct {
	ID              int        `json:"id"`
	UUID            string     `json:"uuid"`
	Username        string     `json:"username"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	Type            string     `json:"type"`
	ContainerUUID   string     `json:"container_uuid"`
	Permissions     int        `json:"permissions"`
	LoginFailCount  int        `json:"login_fail_count"`
	LoginFailTotal  int        `json:"login_fail_total"`
	Enabled         bool       `json:"enabled"`
	TwoFactor       *TwoFactor `json:"two_factor,omitempty"`
	LastLogin       int64      `json:"lastlogin,omitempty"`
	UUIDHash        string     `json:"uuid_id,omitempty"`
	PermissionsType string     `json:"user_type,omitempty"`
}

// TwoFactor contains two-factor authentication settings.
type TwoFactor struct {
	Verified     bool   `json:"verified"`
	SMSEnabled   bool   `json:"sms_enabled"`
	SMSPhone     string `json:"sms_phone,omitempty"`
	EmailEnabled bool   `json:"email_enabled"`
}

// UserCreateRequest represents a request to create a user.
type UserCreateRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Permissions int    `json:"permissions"`
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
	Type        string `json:"type,omitempty"`
}

// UserEditRequest represents a request to edit a user.
type UserEditRequest struct {
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
	Permissions int    `json:"permissions,omitempty"`
	Enabled     *bool  `json:"enabled,omitempty"`
}

// List retrieves all users.
func (u *UsersAPI) List(ctx context.Context) ([]User, error) {
	var result struct {
		Users []User `json:"users"`
	}

	_, err := u.client.Get(ctx, "users", &result)
	if err != nil {
		return nil, err
	}

	return result.Users, nil
}

// Create creates a new user.
func (u *UsersAPI) Create(ctx context.Context, req *UserCreateRequest) (*User, error) {
	var result User
	_, err := u.client.Post(ctx, "users", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves a specific user.
func (u *UsersAPI) Get(ctx context.Context, userID int) (*User, error) {
	var result User
	_, err := u.client.Get(ctx, fmt.Sprintf("users/%d", userID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Edit updates a user.
func (u *UsersAPI) Edit(ctx context.Context, userID int, req *UserEditRequest) (*User, error) {
	var result User
	_, err := u.client.Put(ctx, fmt.Sprintf("users/%d", userID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete removes a user.
func (u *UsersAPI) Delete(ctx context.Context, userID int) error {
	_, err := u.client.Delete(ctx, fmt.Sprintf("users/%d", userID))
	return err
}

// Enable enables a user account.
func (u *UsersAPI) Enable(ctx context.Context, userID int) (*User, error) {
	enabled := true
	return u.Edit(ctx, userID, &UserEditRequest{Enabled: &enabled})
}

// Disable disables a user account.
func (u *UsersAPI) Disable(ctx context.Context, userID int) (*User, error) {
	enabled := false
	return u.Edit(ctx, userID, &UserEditRequest{Enabled: &enabled})
}

// ChangePassword changes a user's password.
func (u *UsersAPI) ChangePassword(ctx context.Context, userID int, currentPassword, newPassword string) error {
	payload := map[string]string{
		"current_password": currentPassword,
		"password":         newPassword,
	}
	_, err := u.client.Put(ctx, fmt.Sprintf("users/%d/chpasswd", userID), payload, nil)
	return err
}

// GenerateAPIKeys generates new API keys for a user.
func (u *UsersAPI) GenerateAPIKeys(ctx context.Context, userID int) (*APIKeys, error) {
	var result APIKeys
	_, err := u.client.Put(ctx, fmt.Sprintf("users/%d/keys", userID), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// APIKeys contains user API key information.
type APIKeys struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

// GetAPIKeys retrieves a user's API keys.
func (u *UsersAPI) GetAPIKeys(ctx context.Context, userID int) (*APIKeys, error) {
	var result APIKeys
	_, err := u.client.Get(ctx, fmt.Sprintf("users/%d/keys", userID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteAPIKeys deletes a user's API keys.
func (u *UsersAPI) DeleteAPIKeys(ctx context.Context, userID int) error {
	_, err := u.client.Delete(ctx, fmt.Sprintf("users/%d/keys", userID))
	return err
}

// ImpersonateRequest represents a request to impersonate a user.
type ImpersonateRequest struct {
	UserID int `json:"user_id"`
}

// Impersonate gets impersonation token for a user.
func (u *UsersAPI) Impersonate(ctx context.Context, userID int) (string, error) {
	var result struct {
		Token string `json:"token"`
	}
	_, err := u.client.Post(ctx, "users/impersonate", &ImpersonateRequest{UserID: userID}, &result)
	if err != nil {
		return "", err
	}
	return result.Token, nil
}

// GroupsAPI handles user group operations.
type GroupsAPI struct {
	client *Client
}

// Group represents a user group.
type Group struct {
	ID            int    `json:"id"`
	UUID          string `json:"uuid"`
	Name          string `json:"name"`
	Permissions   int    `json:"permissions"`
	UserCount     int    `json:"user_count"`
	ContainerUUID string `json:"container_uuid"`
}

// List retrieves all groups.
func (g *GroupsAPI) List(ctx context.Context) ([]Group, error) {
	var result struct {
		Groups []Group `json:"groups"`
	}

	_, err := g.client.Get(ctx, "groups", &result)
	if err != nil {
		return nil, err
	}

	return result.Groups, nil
}

// Create creates a new group.
func (g *GroupsAPI) Create(ctx context.Context, name string) (*Group, error) {
	payload := map[string]string{"name": name}
	var result Group
	_, err := g.client.Post(ctx, "groups", payload, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete removes a group.
func (g *GroupsAPI) Delete(ctx context.Context, groupID int) error {
	_, err := g.client.Delete(ctx, fmt.Sprintf("groups/%d", groupID))
	return err
}

// Edit updates a group.
func (g *GroupsAPI) Edit(ctx context.Context, groupID int, name string) (*Group, error) {
	payload := map[string]string{"name": name}
	var result Group
	_, err := g.client.Put(ctx, fmt.Sprintf("groups/%d", groupID), payload, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// AddUser adds a user to a group.
func (g *GroupsAPI) AddUser(ctx context.Context, groupID, userID int) error {
	_, err := g.client.Post(ctx, fmt.Sprintf("groups/%d/users/%d", groupID, userID), nil, nil)
	return err
}

// RemoveUser removes a user from a group.
func (g *GroupsAPI) RemoveUser(ctx context.Context, groupID, userID int) error {
	_, err := g.client.Delete(ctx, fmt.Sprintf("groups/%d/users/%d", groupID, userID))
	return err
}

// ListUsers lists users in a group.
func (g *GroupsAPI) ListUsers(ctx context.Context, groupID int) ([]User, error) {
	var result struct {
		Users []User `json:"users"`
	}
	_, err := g.client.Get(ctx, fmt.Sprintf("groups/%d/users", groupID), &result)
	if err != nil {
		return nil, err
	}
	return result.Users, nil
}
