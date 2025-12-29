package tio

import (
	"context"
)

// SessionAPI handles session operations.
type SessionAPI struct {
	client *Client
}

// SessionInfo represents the current session information.
type SessionInfo struct {
	ID            int    `json:"id"`
	UUID          string `json:"uuid"`
	Username      string `json:"username"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Type          string `json:"type"`
	ContainerID   int    `json:"container_id"`
	ContainerUUID string `json:"container_uuid"`
	ContainerName string `json:"container_name"`
	Permissions   int    `json:"permissions"`
	Groups        []int  `json:"groups"`
	Lockout       bool   `json:"lockout"`
	Enabled       bool   `json:"enabled"`
}

// Get retrieves the current session information.
func (s *SessionAPI) Get(ctx context.Context) (*SessionInfo, error) {
	var result SessionInfo
	_, err := s.client.Get(ctx, "session", &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Edit updates the current session.
func (s *SessionAPI) Edit(ctx context.Context, name, email string) (*SessionInfo, error) {
	payload := map[string]string{
		"name":  name,
		"email": email,
	}
	var result SessionInfo
	_, err := s.client.Put(ctx, "session", payload, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ChangePassword changes the current user's password.
func (s *SessionAPI) ChangePassword(ctx context.Context, currentPassword, newPassword string) error {
	payload := map[string]string{
		"password":     currentPassword,
		"new_password": newPassword,
	}
	_, err := s.client.Put(ctx, "session/chpasswd", payload, nil)
	return err
}

// GenerateAPIKeys generates new API keys for the current session.
func (s *SessionAPI) GenerateAPIKeys(ctx context.Context) (*APIKeys, error) {
	var result APIKeys
	_, err := s.client.Put(ctx, "session/keys", nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAPIKeys retrieves the API keys for the current session.
func (s *SessionAPI) GetAPIKeys(ctx context.Context) (*APIKeys, error) {
	var result APIKeys
	_, err := s.client.Get(ctx, "session/keys", &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// TwoFactorSetup represents two-factor authentication setup.
type TwoFactorSetup struct {
	Code string `json:"code,omitempty"`
}

// EnableTwoFactor enables two-factor authentication.
func (s *SessionAPI) EnableTwoFactor(ctx context.Context, phone string) error {
	payload := map[string]interface{}{
		"sms_enabled": true,
		"sms_phone":   phone,
	}
	_, err := s.client.Put(ctx, "session/two-factor", payload, nil)
	return err
}

// VerifyTwoFactor verifies the two-factor code.
func (s *SessionAPI) VerifyTwoFactor(ctx context.Context, code string) error {
	payload := map[string]string{"verification_code": code}
	_, err := s.client.Post(ctx, "session/two-factor/verify", payload, nil)
	return err
}

// DisableTwoFactor disables two-factor authentication.
func (s *SessionAPI) DisableTwoFactor(ctx context.Context) error {
	payload := map[string]interface{}{
		"sms_enabled":   false,
		"email_enabled": false,
	}
	_, err := s.client.Put(ctx, "session/two-factor", payload, nil)
	return err
}

// SendVerificationCode sends a verification code for two-factor auth.
func (s *SessionAPI) SendVerificationCode(ctx context.Context) error {
	_, err := s.client.Post(ctx, "session/two-factor/send-verification", nil, nil)
	return err
}
