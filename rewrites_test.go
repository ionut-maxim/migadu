package migadu

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

const rewriteJSON = `{
	"domain": "example.com",
	"name": "catch-support",
	"local_part_rule": "support+*",
	"order_num": 1,
	"destinations": ["user@example.com"]
}`

const rewritesJSON = `{"rewrites": [` + rewriteJSON + `]}`

func rewriteChain(transport http.RoundTripper) *RewritesService {
	return newTestClient(transport).Domains().Domain("example.com").Rewrites()
}

func Test_Rewrites_List(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantCount int
		wantErr   bool
	}{
		{
			name:      "returns decoded rewrites",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, rewritesJSON)},
			wantCount: 1,
		},
		{
			name:      "empty list",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, `{"rewrites":[]}`)},
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
			rewrites, err := rewriteChain(tt.transport).List(context.Background())
			assertResult(t, err, tt.wantErr, nil)
			if err == nil && len(rewrites) != tt.wantCount {
				t.Fatalf("expected %d rewrites, got %d", tt.wantCount, len(rewrites))
			}
		})
	}
}

func Test_Rewrites_Create(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "creates rewrite",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, rewriteJSON)},
		},
		{
			name:      "non-200 returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusUnprocessableEntity, `{"error":"invalid","message":"name taken"}`)},
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
			_, err := rewriteChain(tt.transport).Create(context.Background(), CreateRewriteRequest{
				Name:          "catch-support",
				LocalPartRule: "support+*",
				Destinations:  []string{"user@example.com"},
			})
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}

func Test_Rewrite_Get(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "returns rewrite",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, rewriteJSON)},
		},
		{
			name:      "not found returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusNotFound, `{"error":"not_found","message":"rewrite not found"}`)},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := rewriteChain(tt.transport).Rewrite("catch-support").Get(context.Background())
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}

func Test_Rewrite_Update(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "updates rewrite",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, rewriteJSON)},
		},
		{
			name:      "non-200 returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusUnprocessableEntity, `{"error":"invalid","message":"bad field"}`)},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := rewriteChain(tt.transport).Rewrite("catch-support").Update(context.Background(), UpdateRewriteRequest{
				Destinations: []string{"other@example.com"},
			})
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}

func Test_Rewrite_Delete(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "deletes rewrite",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, `{}`)},
		},
		{
			name:      "not found returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusNotFound, `{"error":"not_found","message":"rewrite not found"}`)},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rewriteChain(tt.transport).Rewrite("catch-support").Delete(context.Background())
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}
