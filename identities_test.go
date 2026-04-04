package migadu

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

const identityJSON = `{
	"local_part": "alias",
	"domain": "example.com",
	"address": "alias@example.com",
	"name": "Alias User",
	"may_send": true,
	"may_receive": true,
	"may_access_imap": true,
	"may_access_pop3": false,
	"may_access_managesieve": false,
	"footer_active": false,
	"footer_plain_body": "",
	"footer_html_body": ""
}`

const identitiesJSON = `{"identities": [` + identityJSON + `]}`

func identityChain(transport http.RoundTripper) *IdentitiesService {
	return newTestClient(transport).Domains().Domain("example.com").Mailboxes().Mailbox("user").Identities()
}

func Test_Identities_List(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantCount int
		wantErr   bool
	}{
		{
			name:      "returns decoded identities",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, identitiesJSON)},
			wantCount: 1,
		},
		{
			name:      "empty list",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, `{"identities":[]}`)},
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
			identities, err := identityChain(tt.transport).List(context.Background())
			assertResult(t, err, tt.wantErr, nil)
			if err == nil && len(identities) != tt.wantCount {
				t.Fatalf("expected %d identities, got %d", tt.wantCount, len(identities))
			}
		})
	}
}

func Test_Identities_Create(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "creates identity",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, identityJSON)},
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
			_, err := identityChain(tt.transport).Create(context.Background(), CreateIdentityRequest{
				Name:      "Alias User",
				LocalPart: "alias",
			})
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}

func Test_Identity_Get(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "returns identity",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, identityJSON)},
		},
		{
			name:      "not found returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusNotFound, `{"error":"not_found","message":"identity not found"}`)},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := identityChain(tt.transport).Identity("alias").Get(context.Background())
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}

func Test_Identity_Update(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "updates identity",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, identityJSON)},
		},
		{
			name:      "non-200 returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusUnprocessableEntity, `{"error":"invalid","message":"bad field"}`)},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := identityChain(tt.transport).Identity("alias").Update(context.Background(), UpdateIdentityRequest{
				Name: "Updated",
			})
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}

func Test_Identity_Delete(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "deletes identity",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, `{}`)},
		},
		{
			name:      "not found returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusNotFound, `{"error":"not_found","message":"identity not found"}`)},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := identityChain(tt.transport).Identity("alias").Delete(context.Background())
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}
