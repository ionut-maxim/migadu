package migadu

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

const aliasJSON = `{
	"local_part": "info",
	"domain": "example.com",
	"address": "info@example.com",
	"is_internal": false,
	"destinations": ["user@example.com"]
}`

const aliasesJSON = `{"aliases": [` + aliasJSON + `]}`

func aliasChain(transport http.RoundTripper) *AliasesService {
	return newTestClient(transport).Domains().Domain("example.com").Aliases()
}

func Test_Aliases_List(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantCount int
		wantErr   bool
	}{
		{
			name:      "returns decoded aliases",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, aliasesJSON)},
			wantCount: 1,
		},
		{
			name:      "empty list",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, `{"aliases":[]}`)},
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
			aliases, err := aliasChain(tt.transport).List(context.Background())
			assertResult(t, err, tt.wantErr, nil)
			if err == nil && len(aliases) != tt.wantCount {
				t.Fatalf("expected %d aliases, got %d", tt.wantCount, len(aliases))
			}
		})
	}
}

func Test_Aliases_Create(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "creates alias",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, aliasJSON)},
		},
		{
			name:      "non-200 returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusUnprocessableEntity, `{"error":"invalid","message":"local_part taken"}`)},
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
			_, err := aliasChain(tt.transport).Create(context.Background(), CreateAliasRequest{
				LocalPart:    "info",
				Destinations: []string{"user@example.com"},
			})
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}

func Test_Alias_Get(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "returns alias",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, aliasJSON)},
		},
		{
			name:      "not found returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusNotFound, `{"error":"not_found","message":"alias not found"}`)},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := aliasChain(tt.transport).Alias("info").Get(context.Background())
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}

func Test_Alias_Update(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "updates alias",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, aliasJSON)},
		},
		{
			name:      "non-200 returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusUnprocessableEntity, `{"error":"invalid","message":"bad field"}`)},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := aliasChain(tt.transport).Alias("info").Update(context.Background(), UpdateAliasRequest{
				Destinations: []string{"other@example.com"},
			})
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}

func Test_Alias_Delete(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "deletes alias",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, `{}`)},
		},
		{
			name:      "not found returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusNotFound, `{"error":"not_found","message":"alias not found"}`)},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := aliasChain(tt.transport).Alias("info").Delete(context.Background())
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}
