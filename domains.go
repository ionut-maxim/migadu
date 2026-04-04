package migadu

import (
	"context"
	"fmt"
	"net/http"
)

type Domain struct {
	Name                    string   `json:"name"`
	State                   string   `json:"state"`
	Description             string   `json:"description"`
	Tags                    []string `json:"tags"`
	ActivatedAt             string   `json:"activated_at"`
	DeactivatedAt           *string  `json:"deactivated_at"`
	MxProxyEnabled          bool     `json:"mx_proxy_enabled"`
	CanAccess               bool     `json:"can_access"`
	CanReceive              bool     `json:"can_receive"`
	CanSend                 bool     `json:"can_send"`
	CatchallDestinations    []string `json:"catchall_destinations"`
	GreylistingEnabled      bool     `json:"greylisting_enabled"`
	JunkSubjectKeywordSpam  bool     `json:"junk_subject_keyword_spam"`
	SubjectRewritingEnabled bool     `json:"subject_rewriting_enabled"`
	HostedDns               bool     `json:"hosted_dns"`
	RecipientDenylist       []string `json:"recipient_denylist"`
	SenderAllowlist         []string `json:"sender_allowlist"`
	SenderDenylist          []string `json:"sender_denylist"`
	SpamAggressiveness      string   `json:"spam_aggressiveness"`
}

type DNSRecord struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type MXRecord struct {
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	Type     string `json:"type"`
	Value    string `json:"value"`
}

type DomainRecords struct {
	DomainName      string      `json:"domain_name"`
	DKIM            []DNSRecord `json:"dkim"`
	DMARC           DNSRecord   `json:"dmarc"`
	DNSVerification DNSRecord   `json:"dns_verification"`
	MXRecords       []MXRecord  `json:"mx_records"`
	SPF             DNSRecord   `json:"spf"`
}

type CreateDomainRequest struct {
	Name                    string   `json:"name"`
	HostedDNS               bool     `json:"hosted_dns"`
	CreateDefaultAddresses  bool     `json:"create_default_addresses,omitempty"`
	Description             string   `json:"description,omitempty"`
	Tags                    []string `json:"tags,omitempty"`
	SpamAggressiveness      string   `json:"spam_aggressiveness,omitempty"`
	GreylistingEnabled      *bool    `json:"greylisting_enabled,omitempty"`
	SubjectRewritingEnabled *bool    `json:"subject_rewriting_enabled,omitempty"`
	JunkSubjectKeywordSpam  *bool    `json:"junk_subject_keyword_spam,omitempty"`
	SenderDenylist          []string `json:"sender_denylist,omitempty"`
	SenderAllowlist         []string `json:"sender_allowlist,omitempty"`
	RecipientDenylist       []string `json:"recipient_denylist,omitempty"`
	CatchallDestinations    []string `json:"catchall_destinations,omitempty"`
	CanAccess               *bool    `json:"can_access,omitempty"`
	MxProxyEnabled          *bool    `json:"mx_proxy_enabled,omitempty"`
}

type UpdateDomainRequest struct {
	Description             string   `json:"description,omitempty"`
	Tags                    []string `json:"tags,omitempty"`
	SpamAggressiveness      string   `json:"spam_aggressiveness,omitempty"`
	GreylistingEnabled      *bool    `json:"greylisting_enabled,omitempty"`
	SubjectRewritingEnabled *bool    `json:"subject_rewriting_enabled,omitempty"`
	JunkSubjectKeywordSpam  *bool    `json:"junk_subject_keyword_spam,omitempty"`
	SenderDenylist          []string `json:"sender_denylist,omitempty"`
	SenderAllowlist         []string `json:"sender_allowlist,omitempty"`
	RecipientDenylist       []string `json:"recipient_denylist,omitempty"`
	CatchallDestinations    []string `json:"catchall_destinations,omitempty"`
	CanAccess               *bool    `json:"can_access,omitempty"`
	MxProxyEnabled          *bool    `json:"mx_proxy_enabled,omitempty"`
}

// DomainsService is the builder for collection-level domain operations.
type DomainsService struct {
	client *Client
}

// DomainService is the builder for single-domain operations.
type DomainService struct {
	client *Client
	name   string
}

func (c *Client) Domains() *DomainsService {
	return &DomainsService{client: c}
}

func (d *DomainsService) Domain(name string) *DomainService {
	return &DomainService{client: d.client, name: name}
}

func (d *DomainsService) List(ctx context.Context) ([]Domain, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "domains", nil)
	if err != nil {
		return nil, err
	}
	type response struct {
		Domains []Domain `json:"domains"`
	}
	resp, err := doAndDecode[response](d.client, req)
	if err != nil {
		return nil, err
	}
	return resp.Domains, nil
}

func (d *DomainsService) Create(ctx context.Context, request CreateDomainRequest) (Domain, error) {
	r, err := encode(request)
	if err != nil {
		return Domain{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "domains", r)
	if err != nil {
		return Domain{}, err
	}
	return doAndDecode[Domain](d.client, req)
}

func (d *DomainService) Get(ctx context.Context) (Domain, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("domains/%s", d.name), nil)
	if err != nil {
		return Domain{}, err
	}
	return doAndDecode[Domain](d.client, req)
}

func (d *DomainService) Update(ctx context.Context, request UpdateDomainRequest) (Domain, error) {
	r, err := encode(request)
	if err != nil {
		return Domain{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, fmt.Sprintf("domains/%s", d.name), r)
	if err != nil {
		return Domain{}, err
	}
	return doAndDecode[Domain](d.client, req)
}

func (d *DomainService) Records(ctx context.Context) (DomainRecords, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("domains/%s/records", d.name), nil)
	if err != nil {
		return DomainRecords{}, err
	}
	return doAndDecode[DomainRecords](d.client, req)
}

func (d *DomainService) Diagnostics(ctx context.Context) (Domain, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("domains/%s/diagnostics", d.name), nil)
	if err != nil {
		return Domain{}, err
	}
	return doAndDecode[Domain](d.client, req)
}

func (d *DomainService) Activate(ctx context.Context) (Domain, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("domains/%s/activate", d.name), nil)
	if err != nil {
		return Domain{}, err
	}
	return doAndDecode[Domain](d.client, req)
}
