package migadu

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type mockTransport struct {
	resp *http.Response
	err  error
}

func (m *mockTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return m.resp, m.err
}

func newTestClient(transport http.RoundTripper) *Client {
	return &Client{
		httpClient: &http.Client{Transport: transport},
	}
}

func mockResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

const domainJSON = `{
	"name": "example.com",
	"state": "active",
	"can_send": true,
	"can_receive": true,
	"tags": [],
	"catchall_destinations": [],
	"recipient_denylist": [],
	"sender_allowlist": [],
	"sender_denylist": []
}`

const domainsJSON = `{"domains": [` + domainJSON + `]}`

func Test_Domains_List(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		ctx       func() context.Context
		wantCount int
		wantErr   bool
		wantErrIs error
	}{
		{
			name:      "returns decoded domains",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, domainsJSON)},
			ctx:       context.Background,
			wantCount: 1,
		},
		{
			name:      "empty list",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, `{"domains":[]}`)},
			ctx:       context.Background,
			wantCount: 0,
		},
		{
			name:      "non-200 returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusUnauthorized, `{"error":"unauthorized","message":"bad credentials"}`)},
			ctx:       context.Background,
			wantErr:   true,
		},
		{
			name:      "malformed json returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, `not json`)},
			ctx:       context.Background,
			wantErr:   true,
		},
		{
			name:      "network error",
			transport: &mockTransport{err: errors.New("connection refused")},
			ctx:       context.Background,
			wantErr:   true,
		},
		{
			name:      "cancelled context",
			transport: &mockTransport{err: context.Canceled},
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr:   true,
			wantErrIs: context.Canceled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domains, err := newTestClient(tt.transport).Domains().List(tt.ctx())
			assertResult(t, err, tt.wantErr, tt.wantErrIs)
			if err == nil && len(domains) != tt.wantCount {
				t.Fatalf("expected %d domains, got %d", tt.wantCount, len(domains))
			}
		})
	}
}

func Test_Domains_Create(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "creates domain",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, domainJSON)},
		},
		{
			name:      "non-200 returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusUnprocessableEntity, `{"error":"invalid","message":"name already taken"}`)},
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
			_, err := newTestClient(tt.transport).Domains().Create(context.Background(), CreateDomainRequest{
				Name:      "example.com",
				HostedDNS: false,
			})
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}

func Test_Domain_Get(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "returns domain",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, domainJSON)},
		},
		{
			name:      "not found returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusNotFound, `{"error":"not_found","message":"domain not found"}`)},
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
			_, err := newTestClient(tt.transport).Domains().Domain("example.com").Get(context.Background())
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}

func Test_Domain_Update(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "updates domain",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, domainJSON)},
		},
		{
			name:      "non-200 returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusUnprocessableEntity, `{"error":"invalid","message":"bad field"}`)},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc := "updated"
			_, err := newTestClient(tt.transport).Domains().Domain("example.com").Update(context.Background(), UpdateDomainRequest{
				Description: desc,
			})
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}

func Test_Domain_Records(t *testing.T) {
	const recordsJSON = `{
		"domain_name": "example.com",
		"dkim": [{"name": "dkim._domainkey", "type": "TXT", "value": "v=DKIM1"}],
		"dmarc": {"name": "_dmarc", "type": "TXT", "value": "v=DMARC1"},
		"dns_verification": {"name": "@", "type": "TXT", "value": "migadu-verify=abc"},
		"mx_records": [{"name": "@", "priority": 10, "type": "MX", "value": "aspmx1.migadu.com"}],
		"spf": {"name": "@", "type": "TXT", "value": "v=spf1 include:spf.migadu.com"}
	}`

	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "returns records",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, recordsJSON)},
		},
		{
			name:      "non-200 returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusNotFound, `{"error":"not_found","message":"domain not found"}`)},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newTestClient(tt.transport).Domains().Domain("example.com").Records(context.Background())
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}

func Test_Domain_Activate(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "activates domain",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, domainJSON)},
		},
		{
			name:      "dns check failed returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusUnprocessableEntity, `{"error":"dns_check_failed","message":"DNS checks failed"}`)},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newTestClient(tt.transport).Domains().Domain("example.com").Activate(context.Background())
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}

func Test_Domain_Diagnostics(t *testing.T) {
	tests := []struct {
		name      string
		transport *mockTransport
		wantErr   bool
	}{
		{
			name:      "returns diagnostics",
			transport: &mockTransport{resp: mockResponse(http.StatusOK, domainJSON)},
		},
		{
			name:      "non-200 returns error",
			transport: &mockTransport{resp: mockResponse(http.StatusNotFound, `{"error":"not_found","message":"domain not found"}`)},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newTestClient(tt.transport).Domains().Domain("example.com").Diagnostics(context.Background())
			assertResult(t, err, tt.wantErr, nil)
		})
	}
}

func assertResult(t *testing.T, err error, wantErr bool, wantErrIs error) {
	t.Helper()
	if wantErr {
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if wantErrIs != nil && !errors.Is(err, wantErrIs) {
			t.Fatalf("expected error %v, got %v", wantErrIs, err)
		}
		return
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
