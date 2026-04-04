package migadu

import (
	"context"
	"fmt"
	"net/http"
)

type Mailbox struct {
	Address               string   `json:"address"`
	LocalPart             string   `json:"local_part"`
	DomainName            string   `json:"domain_name"`
	Name                  string   `json:"name"`
	IsInternal            bool     `json:"is_internal"`
	WildcardSender        bool     `json:"wildcard_sender"`
	MaySend               bool     `json:"may_send"`
	MayReceive            bool     `json:"may_receive"`
	MayAccessImap         bool     `json:"may_access_imap"`
	MayAccessPop3         bool     `json:"may_access_pop3"`
	MayAccessManagesieve  bool     `json:"may_access_managesieve"`
	PasswordRecoveryEmail string   `json:"password_recovery_email"`
	SpamAction            string   `json:"spam_action"`
	SpamAggressiveness    string   `json:"spam_aggressiveness"`
	SenderDenylist        []string `json:"sender_denylist"`
	SenderAllowlist       []string `json:"sender_allowlist"`
	RecipientDenylist     []string `json:"recipient_denylist"`
	FooterActive          bool     `json:"footer_active"`
	FooterPlainBody       string   `json:"footer_plain_body"`
	FooterHtmlBody        string   `json:"footer_html_body"`
}

type CreateMailboxRequest struct {
	Name                  string `json:"name"`
	LocalPart             string `json:"local_part"`
	Password              string `json:"password,omitempty"`
	PasswordMethod        string `json:"password_method,omitempty"`
	PasswordRecoveryEmail string `json:"password_recovery_email,omitempty"`
	ForwardingTo          string `json:"forwarding_to,omitempty"`
	IsInternal            *bool  `json:"is_internal,omitempty"`
}

type UpdateMailboxRequest struct {
	Name                  string   `json:"name,omitempty"`
	Password              string   `json:"password,omitempty"`
	PasswordRecoveryEmail string   `json:"password_recovery_email,omitempty"`
	ForwardingTo          string   `json:"forwarding_to,omitempty"`
	IsInternal            *bool    `json:"is_internal,omitempty"`
	MaySend               *bool    `json:"may_send,omitempty"`
	MayReceive            *bool    `json:"may_receive,omitempty"`
	MayAccessImap         *bool    `json:"may_access_imap,omitempty"`
	MayAccessPop3         *bool    `json:"may_access_pop3,omitempty"`
	MayAccessManagesieve  *bool    `json:"may_access_managesieve,omitempty"`
	SpamAction            string   `json:"spam_action,omitempty"`
	SpamAggressiveness    string   `json:"spam_aggressiveness,omitempty"`
	SenderDenylist        []string `json:"sender_denylist,omitempty"`
	SenderAllowlist       []string `json:"sender_allowlist,omitempty"`
	RecipientDenylist     []string `json:"recipient_denylist,omitempty"`
	FooterActive          *bool    `json:"footer_active,omitempty"`
	FooterPlainBody       string   `json:"footer_plain_body,omitempty"`
	FooterHtmlBody        string   `json:"footer_html_body,omitempty"`
}

// MailboxesService is the builder for collection-level mailbox operations.
type MailboxesService struct {
	client     *Client
	domainName string
}

// MailboxService is the builder for single-mailbox operations.
type MailboxService struct {
	client     *Client
	domainName string
	localPart  string
}

func (d *DomainService) Mailboxes() *MailboxesService {
	return &MailboxesService{client: d.client, domainName: d.name}
}

func (m *MailboxesService) Mailbox(localPart string) *MailboxService {
	return &MailboxService{client: m.client, domainName: m.domainName, localPart: localPart}
}

func (m *MailboxesService) List(ctx context.Context) ([]Mailbox, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("domains/%s/mailboxes", m.domainName), nil)
	if err != nil {
		return nil, err
	}
	type response struct {
		Mailboxes []Mailbox `json:"mailboxes"`
	}
	resp, err := doAndDecode[response](m.client, req)
	if err != nil {
		return nil, err
	}
	return resp.Mailboxes, nil
}

func (m *MailboxesService) Create(ctx context.Context, request CreateMailboxRequest) (Mailbox, error) {
	r, err := encode(request)
	if err != nil {
		return Mailbox{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("domains/%s/mailboxes", m.domainName), r)
	if err != nil {
		return Mailbox{}, err
	}
	return doAndDecode[Mailbox](m.client, req)
}

func (m *MailboxService) Get(ctx context.Context) (Mailbox, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("domains/%s/mailboxes/%s", m.domainName, m.localPart), nil)
	if err != nil {
		return Mailbox{}, err
	}
	return doAndDecode[Mailbox](m.client, req)
}

func (m *MailboxService) Update(ctx context.Context, request UpdateMailboxRequest) (Mailbox, error) {
	r, err := encode(request)
	if err != nil {
		return Mailbox{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("domains/%s/mailboxes/%s", m.domainName, m.localPart), r)
	if err != nil {
		return Mailbox{}, err
	}
	return doAndDecode[Mailbox](m.client, req)
}

func (m *MailboxService) Delete(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("domains/%s/mailboxes/%s", m.domainName, m.localPart), nil)
	if err != nil {
		return err
	}
	_, err = doAndDecode[struct{}](m.client, req)
	return err
}
