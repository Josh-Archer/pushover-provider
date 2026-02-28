// Copyright (c) Josh Archer
// SPDX-License-Identifier: MPL-2.0

// Package pushover provides a client for the Pushover REST API.
package pushover

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const defaultBaseURL = "https://api.pushover.net/1"

// Client is the Pushover API client.
type Client struct {
	token      string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Pushover API client.
func NewClient(token string) *Client {
	return &Client{
		token:      token,
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{},
	}
}

// NewClientWithBase creates a Pushover client that targets a custom base URL.
// This is exported for use in tests only.
func NewClientWithBase(token, base string, httpClient *http.Client) *Client {
	return &Client{
		token:      token,
		baseURL:    base,
		httpClient: httpClient,
	}
}

// APIResponse is the base Pushover API response.
type APIResponse struct {
	Status  int      `json:"status"`
	Request string   `json:"request"`
	Errors  []string `json:"errors,omitempty"`
}

// MessageRequest holds all fields for sending a Pushover message.
type MessageRequest struct {
	Token     string `json:"token"`
	User      string `json:"user"`
	Message   string `json:"message"`
	Title     string `json:"title,omitempty"`
	URL       string `json:"url,omitempty"`
	URLTitle  string `json:"url_title,omitempty"`
	Priority  int    `json:"priority"`
	Sound     string `json:"sound,omitempty"`
	Device    string `json:"device,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	HTML      int    `json:"html,omitempty"`
	Monospace int    `json:"monospace,omitempty"`
	TTL       int    `json:"ttl,omitempty"`
	// Emergency priority (priority=2) fields
	Retry    int    `json:"retry,omitempty"`
	Expire   int    `json:"expire,omitempty"`
	Callback string `json:"callback,omitempty"`
}

// MessageResponse is the response from sending a message.
type MessageResponse struct {
	APIResponse
	Receipt string `json:"receipt,omitempty"`
}

// ReceiptResponse is the response from polling an emergency receipt.
type ReceiptResponse struct {
	APIResponse
	Acknowledged        int    `json:"acknowledged"`
	AcknowledgedAt      int64  `json:"acknowledged_at"`
	AcknowledgedBy      string `json:"acknowledged_by"`
	AcknowledgedByDevice string `json:"acknowledged_by_device"`
	LastDeliveredAt     int64  `json:"last_delivered_at"`
	Expired             int    `json:"expired"`
	ExpiresAt           int64  `json:"expires_at"`
	CalledBack          int    `json:"called_back"`
	CalledBackAt        int64  `json:"called_back_at"`
}

// SoundsResponse is the response from listing sounds.
type SoundsResponse struct {
	APIResponse
	Sounds map[string]string `json:"sounds"`
}

// Sound represents a single Pushover sound.
type Sound struct {
	Key  string
	Name string
}

// ValidateRequest holds fields for validating a user/group key.
type ValidateRequest struct {
	Token  string `json:"token"`
	User   string `json:"user"`
	Device string `json:"device,omitempty"`
}

// ValidateResponse is the response from validating a user.
type ValidateResponse struct {
	APIResponse
	Group   int      `json:"group"`
	Devices []string `json:"devices"`
	Licenses []string `json:"licenses"`
}

// GroupResponse is the response from group API calls.
type GroupResponse struct {
	APIResponse
	Name  string        `json:"name"`
	Users []GroupMember `json:"users"`
}

// GroupMember represents a member in a Pushover group.
type GroupMember struct {
	User     string `json:"user"`
	Device   string `json:"device,omitempty"`
	Memo     string `json:"memo,omitempty"`
	Disabled bool   `json:"disabled"`
}

// SendMessage sends a notification via the Pushover API.
func (c *Client) SendMessage(ctx context.Context, req *MessageRequest) (*MessageResponse, error) {
	if req.Token == "" {
		req.Token = c.token
	}

	params := url.Values{}
	params.Set("token", req.Token)
	params.Set("user", req.User)
	params.Set("message", req.Message)
	if req.Title != "" {
		params.Set("title", req.Title)
	}
	if req.URL != "" {
		params.Set("url", req.URL)
	}
	if req.URLTitle != "" {
		params.Set("url_title", req.URLTitle)
	}
	params.Set("priority", strconv.Itoa(req.Priority))
	if req.Sound != "" {
		params.Set("sound", req.Sound)
	}
	if req.Device != "" {
		params.Set("device", req.Device)
	}
	if req.Timestamp != 0 {
		params.Set("timestamp", strconv.FormatInt(req.Timestamp, 10))
	}
	if req.HTML != 0 {
		params.Set("html", strconv.Itoa(req.HTML))
	}
	if req.Monospace != 0 {
		params.Set("monospace", strconv.Itoa(req.Monospace))
	}
	if req.TTL != 0 {
		params.Set("ttl", strconv.Itoa(req.TTL))
	}
	if req.Priority == 2 {
		params.Set("retry", strconv.Itoa(req.Retry))
		params.Set("expire", strconv.Itoa(req.Expire))
		if req.Callback != "" {
			params.Set("callback", req.Callback)
		}
	}

	var resp MessageResponse
	if err := c.doPost(ctx, "/messages.json", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetReceipt retrieves delivery status for an emergency message receipt.
func (c *Client) GetReceipt(ctx context.Context, receipt string) (*ReceiptResponse, error) {
	path := fmt.Sprintf("/receipts/%s.json?token=%s", receipt, url.QueryEscape(c.token))
	var resp ReceiptResponse
	if err := c.doGet(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelReceipt cancels an outstanding emergency notification.
func (c *Client) CancelReceipt(ctx context.Context, receipt string) (*APIResponse, error) {
	params := url.Values{}
	params.Set("token", c.token)
	var resp APIResponse
	if err := c.doPost(ctx, fmt.Sprintf("/receipts/%s/cancel.json", receipt), params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSounds returns the list of available Pushover sounds.
func (c *Client) GetSounds(ctx context.Context) ([]Sound, error) {
	path := fmt.Sprintf("/sounds.json?token=%s", url.QueryEscape(c.token))
	var resp SoundsResponse
	if err := c.doGet(ctx, path, &resp); err != nil {
		return nil, err
	}
	sounds := make([]Sound, 0, len(resp.Sounds))
	for k, v := range resp.Sounds {
		sounds = append(sounds, Sound{Key: k, Name: v})
	}
	return sounds, nil
}

// ValidateUser validates a Pushover user or group key.
func (c *Client) ValidateUser(ctx context.Context, req *ValidateRequest) (*ValidateResponse, error) {
	if req.Token == "" {
		req.Token = c.token
	}
	params := url.Values{}
	params.Set("token", req.Token)
	params.Set("user", req.User)
	if req.Device != "" {
		params.Set("device", req.Device)
	}
	var resp ValidateResponse
	if err := c.doPost(ctx, "/users/validate.json", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetGroup retrieves information about a Pushover delivery group.
func (c *Client) GetGroup(ctx context.Context, groupKey string) (*GroupResponse, error) {
	path := fmt.Sprintf("/groups/%s.json?token=%s", groupKey, url.QueryEscape(c.token))
	var resp GroupResponse
	if err := c.doGet(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RenameGroup renames a Pushover delivery group.
func (c *Client) RenameGroup(ctx context.Context, groupKey, name string) (*APIResponse, error) {
	params := url.Values{}
	params.Set("token", c.token)
	params.Set("name", name)
	var resp APIResponse
	if err := c.doPost(ctx, fmt.Sprintf("/groups/%s/rename.json", groupKey), params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AddGroupUser adds a user to a Pushover delivery group.
func (c *Client) AddGroupUser(ctx context.Context, groupKey, user, device, memo string) (*APIResponse, error) {
	params := url.Values{}
	params.Set("token", c.token)
	params.Set("user", user)
	if device != "" {
		params.Set("device", device)
	}
	if memo != "" {
		params.Set("memo", memo)
	}
	var resp APIResponse
	if err := c.doPost(ctx, fmt.Sprintf("/groups/%s/add_user.json", groupKey), params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RemoveGroupUser removes a user from a Pushover delivery group.
func (c *Client) RemoveGroupUser(ctx context.Context, groupKey, user, device string) (*APIResponse, error) {
	params := url.Values{}
	params.Set("token", c.token)
	params.Set("user", user)
	if device != "" {
		params.Set("device", device)
	}
	var resp APIResponse
	if err := c.doPost(ctx, fmt.Sprintf("/groups/%s/delete_user.json", groupKey), params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// EnableGroupUser re-enables a disabled user in a Pushover delivery group.
func (c *Client) EnableGroupUser(ctx context.Context, groupKey, user, device string) (*APIResponse, error) {
	params := url.Values{}
	params.Set("token", c.token)
	params.Set("user", user)
	if device != "" {
		params.Set("device", device)
	}
	var resp APIResponse
	if err := c.doPost(ctx, fmt.Sprintf("/groups/%s/enable_user.json", groupKey), params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DisableGroupUser disables a user in a Pushover delivery group.
func (c *Client) DisableGroupUser(ctx context.Context, groupKey, user, device string) (*APIResponse, error) {
	params := url.Values{}
	params.Set("token", c.token)
	params.Set("user", user)
	if device != "" {
		params.Set("device", device)
	}
	var resp APIResponse
	if err := c.doPost(ctx, fmt.Sprintf("/groups/%s/disable_user.json", groupKey), params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) doPost(ctx context.Context, path string, params url.Values, out interface{}) error {
	u := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, strings.NewReader(params.Encode()))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	// Check for API-level errors
	type statusChecker struct {
		Status int      `json:"status"`
		Errors []string `json:"errors"`
	}
	var sc statusChecker
	_ = json.Unmarshal(body, &sc)
	if sc.Status != 1 {
		return fmt.Errorf("pushover API error: %s", strings.Join(sc.Errors, "; "))
	}

	return nil
}

func (c *Client) doGet(ctx context.Context, path string, out interface{}) error {
	u := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	type statusChecker struct {
		Status int      `json:"status"`
		Errors []string `json:"errors"`
	}
	var sc statusChecker
	_ = json.Unmarshal(body, &sc)
	if sc.Status != 1 {
		return fmt.Errorf("pushover API error: %s", strings.Join(sc.Errors, "; "))
	}

	return nil
}

// IsGroupKey returns true if the validation response indicates the key is a group key.
func (v *ValidateResponse) IsGroupKey() bool {
	return v.Group == 1
}
