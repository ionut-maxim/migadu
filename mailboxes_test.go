package migadu

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

const mailboxJSON = `{
	"address": "user@example.com",
	"local_part": "user",
	"domain_name": "example.com",
	"name": "Test User",
	"may_send": true,
	"may_receive": true,
	"may_access_imap": true,
	"may_access_pop3": false,
	"may_access_managesieve": false,
	"is_internal": false,
	"wildcard_sender": false,
	"spam_action": "folder",
	"spam_aggressiveness": "default",
	"sender_denylist": [],
	"sender_allowlist": [],
	"recipient_denylist": [],
	"footer_active": false,
	"footer_plain_body": "",
	"footer_html_body": ""
}`

const mailboxesJSON = `{"mailboxes": [` + mailboxJSON + `]}`

func Test_Mailboxes_List(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantCount int
		wantErr   bool
		wantErrIs error
	}{
		{
			name:      "returns decoded mailboxes",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, mailboxesJSON)},
			wantCount: 1,
		},
		{
			name:      "empty list",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, `{"mailboxes":[]}`)},
			wantCount: 0,
		},
		{
			name:      "non-200 returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusUnauthorized, `{"error":"unauthorized","message":"bad credentials"}`)},
			wantErr:   true,
		},
		{
			name:      "malformed json returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, `not json`)},
			wantErr:   true,
		},
		{
			name:      "network error",
			transport: &mockTransport{err: errors.New("connection refused")},
			wantErr:   true,
		},
		{
			name:      "cancelled context",
			transport: &mockTransport{err: context.Canceled},
			wantErr:   true,
			wantErrIs: context.Canceled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mailboxes, err := newTestClient(tt.transport).Domains().Domain("example.com").Mailboxes().List(context.Background())
			assertResult(t, err, tt.wantErr, tt.wantErrIs)
			if err == nil && len(mailboxes) != tt.wantCount {
				t.Fatalf("expected %d mailboxes, got %d", tt.wantCount, len(mailboxes))
			}
		})
	}
}

func Test_Mailboxes_Create(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "creates mailbox",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, mailboxJSON)},
		},
		{
			name:      "non-200 returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusUnprocessableEntity, `{"error":"invalid","message":"local_part already taken"}`)},
			wantErr:   true,
		},
		{
			name:      "network error",
			transport: &mockTransport{err: errors.New("connection refused")},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newTestClient(tt.transport).Domains().Domain("example.com").Mailboxes().Create(context.Background(), CreateMailboxRequest{
				Name:      "Test User",
				LocalPart: "user",
				Password:  "secret",
			})
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}

func Test_Mailbox_Get(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "returns mailbox",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, mailboxJSON)},
		},
		{
			name:      "not found returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusNotFound, `{"error":"not_found","message":"mailbox not found"}`)},
			wantErr:   true,
		},
		{
			name:      "network error",
			transport: &mockTransport{err: errors.New("connection refused")},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newTestClient(tt.transport).Domains().Domain("example.com").Mailboxes().Mailbox("user").Get(context.Background())
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}

func Test_Mailbox_Update(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "updates mailbox",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, mailboxJSON)},
		},
		{
			name:      "non-200 returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusUnprocessableEntity, `{"error":"invalid","message":"bad field"}`)},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newTestClient(tt.transport).Domains().Domain("example.com").Mailboxes().Mailbox("user").Update(context.Background(), UpdateMailboxRequest{
				Name: "Updated User",
			})
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}

func Test_Mailbox_Delete(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "deletes mailbox",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, `{}`)},
		},
		{
			name:      "not found returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusNotFound, `{"error":"not_found","message":"mailbox not found"}`)},
			wantErr:   true,
		},
		{
			name:      "network error",
			transport: &mockTransport{err: errors.New("connection refused")},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := newTestClient(tt.transport).Domains().Domain("example.com").Mailboxes().Mailbox("user").Delete(context.Background())
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}
