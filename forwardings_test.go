package migadu

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

const forwardingJSON = `{
	"address": "other@example.com",
	"blocked_at": "",
	"confirmation_sent_at": "",
	"confirmed_at": "",
	"expires_on": "",
	"is_active": true,
	"remove_upon_expiry": false
}`

const forwardingsJSON = `{"forwardings": [` + forwardingJSON + `]}`

func forwardingChain(transport http.RoundTripper) *ForwardingsService {
	return newTestClient(transport).Domains().Domain("example.com").Mailboxes().Mailbox("user").Forwardings()
}

func Test_Forwardings_List(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantCount int
		wantErr   bool
	}{
		{
			name:      "returns decoded forwardings",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, forwardingsJSON)},
			wantCount: 1,
		},
		{
			name:      "empty list",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, `{"forwardings":[]}`)},
			wantCount: 0,
		},
		{
			name:      "non-200 returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusUnauthorized, `{"error":"unauthorized","message":"bad credentials"}`)},
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
			forwardings, err := forwardingChain(tt.transport).List(context.Background())
			assertResult(t, err, tt.wantErr, nil)
			if err == nil && len(forwardings) != tt.wantCount {
				t.Fatalf("expected %d forwardings, got %d", tt.wantCount, len(forwardings))
			}
		})
	}
}

func Test_Forwardings_Create(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "creates forwarding",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, forwardingJSON)},
		},
		{
			name:      "non-200 returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusUnprocessableEntity, `{"error":"invalid","message":"invalid address"}`)},
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
			_, err := forwardingChain(tt.transport).Create(context.Background(), CreateForwardingRequest{
				Address: "other@example.com",
			})
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}

func Test_Forwarding_Get(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "returns forwarding",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, forwardingJSON)},
		},
		{
			name:      "not found returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusNotFound, `{"error":"not_found","message":"forwarding not found"}`)},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := forwardingChain(tt.transport).Forwarding("other@example.com").Get(context.Background())
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}

func Test_Forwarding_Update(t *testing.T) {
	active := false
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "updates forwarding",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, forwardingJSON)},
		},
		{
			name:      "non-200 returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusUnprocessableEntity, `{"error":"invalid","message":"bad field"}`)},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := forwardingChain(tt.transport).Forwarding("other@example.com").Update(context.Background(), UpdateForwardingRequest{
				IsActive: &active,
			})
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}

func Test_Forwarding_Delete(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "deletes forwarding",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, `{}`)},
		},
		{
			name:      "not found returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusNotFound, `{"error":"not_found","message":"forwarding not found"}`)},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := forwardingChain(tt.transport).Forwarding("other@example.com").Delete(context.Background())
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}
