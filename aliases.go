package migadu

import (
	"context"
	"fmt"
	"net/http"
)

type Alias struct {
	LocalPart    string   `json:"local_part"`
	Domain       string   `json:"domain"`
	Address      string   `json:"address"`
	IsInternal   bool     `json:"is_internal"`
	Destinations []string `json:"destinations"`
}

type CreateAliasRequest struct {
	LocalPart    string   `json:"local_part"`
	Destinations []string `json:"destinations"`
	IsInternal   *bool    `json:"is_internal,omitempty"`
}

type UpdateAliasRequest struct {
	Destinations []string `json:"destinations,omitempty"`
	IsInternal   *bool    `json:"is_internal,omitempty"`
}

// AliasesService is the builder for collection-level alias operations.
type AliasesService struct {
	client     *Client
	domainName string
}

// AliasService is the builder for single-alias operations.
type AliasService struct {
	client     *Client
	domainName string
	localPart  string
}

func (d *DomainService) Aliases() *AliasesService {
	return &AliasesService{client: d.client, domainName: d.name}
}

func (s *AliasesService) Alias(localPart string) *AliasService {
	return &AliasService{client: s.client, domainName: s.domainName, localPart: localPart}
}

func (s *AliasesService) List(ctx context.Context) ([]Alias, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("domains/%s/aliases", s.domainName), nil)
	if err != nil {
		return nil, err
	}
	type response struct {
		Aliases []Alias `json:"aliases"`
	}
	resp, err := doAndDecode[response](s.client, req)
	if err != nil {
		return nil, err
	}
	return resp.Aliases, nil
}

func (s *AliasesService) Create(ctx context.Context, request CreateAliasRequest) (Alias, error) {
	r, err := encode(request)
	if err != nil {
		return Alias{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("domains/%s/aliases", s.domainName), r)
	if err != nil {
		return Alias{}, err
	}
	return doAndDecode[Alias](s.client, req)
}

func (s *AliasService) Get(ctx context.Context) (Alias, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("domains/%s/aliases/%s", s.domainName, s.localPart), nil)
	if err != nil {
		return Alias{}, err
	}
	return doAndDecode[Alias](s.client, req)
}

func (s *AliasService) Update(ctx context.Context, request UpdateAliasRequest) (Alias, error) {
	r, err := encode(request)
	if err != nil {
		return Alias{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("domains/%s/aliases/%s", s.domainName, s.localPart), r)
	if err != nil {
		return Alias{}, err
	}
	return doAndDecode[Alias](s.client, req)
}

func (s *AliasService) Delete(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("domains/%s/aliases/%s", s.domainName, s.localPart), nil)
	if err != nil {
		return err
	}
	_, err = doAndDecode[struct{}](s.client, req)
	return err
}
