package migadu

import (
	"context"
	"fmt"
	"net/http"
)

type Identity struct {
	LocalPart            string `json:"local_part"`
	Domain               string `json:"domain"`
	Address              string `json:"address"`
	Name                 string `json:"name"`
	MaySend              bool   `json:"may_send"`
	MayReceive           bool   `json:"may_receive"`
	MayAccessImap        bool   `json:"may_access_imap"`
	MayAccessPop3        bool   `json:"may_access_pop3"`
	MayAccessManagesieve bool   `json:"may_access_managesieve"`
	FooterActive         bool   `json:"footer_active"`
	FooterPlainBody      string `json:"footer_plain_body"`
	FooterHtmlBody       string `json:"footer_html_body"`
}

type CreateIdentityRequest struct {
	Name      string `json:"name"`
	LocalPart string `json:"local_part"`
	Password  string `json:"password,omitempty"`
}

type UpdateIdentityRequest struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

// IdentitiesService is the builder for collection-level identity operations.
type IdentitiesService struct {
	client     *Client
	domainName string
	localPart  string
}

// IdentityService is the builder for single-identity operations.
type IdentityService struct {
	client     *Client
	domainName string
	mailbox    string
	localPart  string
}

func (m *MailboxService) Identities() *IdentitiesService {
	return &IdentitiesService{client: m.client, domainName: m.domainName, localPart: m.localPart}
}

func (s *IdentitiesService) Identity(localPart string) *IdentityService {
	return &IdentityService{client: s.client, domainName: s.domainName, mailbox: s.localPart, localPart: localPart}
}

func (s *IdentitiesService) List(ctx context.Context) ([]Identity, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("domains/%s/mailboxes/%s/identities", s.domainName, s.localPart), nil)
	if err != nil {
		return nil, err
	}
	type response struct {
		Identities []Identity `json:"identities"`
	}
	resp, err := doAndDecode[response](s.client, req)
	if err != nil {
		return nil, err
	}
	return resp.Identities, nil
}

func (s *IdentitiesService) Create(ctx context.Context, request CreateIdentityRequest) (Identity, error) {
	r, err := encode(request)
	if err != nil {
		return Identity{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("domains/%s/mailboxes/%s/identities", s.domainName, s.localPart), r)
	if err != nil {
		return Identity{}, err
	}
	return doAndDecode[Identity](s.client, req)
}

func (s *IdentityService) Get(ctx context.Context) (Identity, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("domains/%s/mailboxes/%s/identities/%s", s.domainName, s.mailbox, s.localPart), nil)
	if err != nil {
		return Identity{}, err
	}
	return doAndDecode[Identity](s.client, req)
}

func (s *IdentityService) Update(ctx context.Context, request UpdateIdentityRequest) (Identity, error) {
	r, err := encode(request)
	if err != nil {
		return Identity{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("domains/%s/mailboxes/%s/identities/%s", s.domainName, s.mailbox, s.localPart), r)
	if err != nil {
		return Identity{}, err
	}
	return doAndDecode[Identity](s.client, req)
}

func (s *IdentityService) Delete(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("domains/%s/mailboxes/%s/identities/%s", s.domainName, s.mailbox, s.localPart), nil)
	if err != nil {
		return err
	}
	_, err = doAndDecode[struct{}](s.client, req)
	return err
}
