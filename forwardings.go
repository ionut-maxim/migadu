package migadu

import (
	"context"
	"fmt"
	"net/http"
)

type Forwarding struct {
	Address            string `json:"address"`
	BlockedAt          string `json:"blocked_at"`
	ConfirmationSentAt string `json:"confirmation_sent_at"`
	ConfirmedAt        string `json:"confirmed_at"`
	ExpiresOn          string `json:"expires_on"`
	IsActive           bool   `json:"is_active"`
	RemoveUponExpiry   bool   `json:"remove_upon_expiry"`
}

type CreateForwardingRequest struct {
	Address string `json:"address"`
}

type UpdateForwardingRequest struct {
	IsActive         *bool  `json:"is_active,omitempty"`
	ExpiresOn        string `json:"expires_on,omitempty"`
	RemoveUponExpiry *bool  `json:"remove_upon_expiry,omitempty"`
}

// ForwardingsService is the builder for collection-level forwarding operations.
type ForwardingsService struct {
	client     *Client
	domainName string
	localPart  string
}

// ForwardingService is the builder for single-forwarding operations.
type ForwardingService struct {
	client     *Client
	domainName string
	mailbox    string
	address    string
}

func (m *MailboxService) Forwardings() *ForwardingsService {
	return &ForwardingsService{client: m.client, domainName: m.domainName, localPart: m.localPart}
}

func (s *ForwardingsService) Forwarding(address string) *ForwardingService {
	return &ForwardingService{client: s.client, domainName: s.domainName, mailbox: s.localPart, address: address}
}

func (s *ForwardingsService) List(ctx context.Context) ([]Forwarding, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("domains/%s/mailboxes/%s/forwardings", s.domainName, s.localPart), nil)
	if err != nil {
		return nil, err
	}
	type response struct {
		Forwardings []Forwarding `json:"forwardings"`
	}
	resp, err := doAndDecode[response](s.client, req)
	if err != nil {
		return nil, err
	}
	return resp.Forwardings, nil
}

func (s *ForwardingsService) Create(ctx context.Context, request CreateForwardingRequest) (Forwarding, error) {
	r, err := encode(request)
	if err != nil {
		return Forwarding{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("domains/%s/mailboxes/%s/forwardings", s.domainName, s.localPart), r)
	if err != nil {
		return Forwarding{}, err
	}
	return doAndDecode[Forwarding](s.client, req)
}

func (s *ForwardingService) Get(ctx context.Context) (Forwarding, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("domains/%s/mailboxes/%s/forwardings/%s", s.domainName, s.mailbox, s.address), nil)
	if err != nil {
		return Forwarding{}, err
	}
	return doAndDecode[Forwarding](s.client, req)
}

func (s *ForwardingService) Update(ctx context.Context, request UpdateForwardingRequest) (Forwarding, error) {
	r, err := encode(request)
	if err != nil {
		return Forwarding{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("domains/%s/mailboxes/%s/forwardings/%s", s.domainName, s.mailbox, s.address), r)
	if err != nil {
		return Forwarding{}, err
	}
	return doAndDecode[Forwarding](s.client, req)
}

func (s *ForwardingService) Delete(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("domains/%s/mailboxes/%s/forwardings/%s", s.domainName, s.mailbox, s.address), nil)
	if err != nil {
		return err
	}
	_, err = doAndDecode[struct{}](s.client, req)
	return err
}
