package migadu

import (
	"context"
	"fmt"
	"net/http"
)

type Rewrite struct {
	Domain        string   `json:"domain"`
	Name          string   `json:"name"`
	LocalPartRule string   `json:"local_part_rule"`
	OrderNum      int      `json:"order_num"`
	Destinations  []string `json:"destinations"`
}

type CreateRewriteRequest struct {
	Name          string   `json:"name"`
	LocalPartRule string   `json:"local_part_rule"`
	Destinations  []string `json:"destinations"`
}

type UpdateRewriteRequest struct {
	Name          string   `json:"name,omitempty"`
	LocalPartRule string   `json:"local_part_rule,omitempty"`
	Destinations  []string `json:"destinations,omitempty"`
}

// RewritesService is the builder for collection-level rewrite operations.
type RewritesService struct {
	client     *Client
	domainName string
}

// RewriteService is the builder for single-rewrite operations.
type RewriteService struct {
	client     *Client
	domainName string
	name       string
}

func (d *DomainService) Rewrites() *RewritesService {
	return &RewritesService{client: d.client, domainName: d.name}
}

func (s *RewritesService) Rewrite(name string) *RewriteService {
	return &RewriteService{client: s.client, domainName: s.domainName, name: name}
}

func (s *RewritesService) List(ctx context.Context) ([]Rewrite, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("domains/%s/rewrites", s.domainName), nil)
	if err != nil {
		return nil, err
	}
	type response struct {
		Rewrites []Rewrite `json:"rewrites"`
	}
	resp, err := doAndDecode[response](s.client, req)
	if err != nil {
		return nil, err
	}
	return resp.Rewrites, nil
}

func (s *RewritesService) Create(ctx context.Context, request CreateRewriteRequest) (Rewrite, error) {
	r, err := encode(request)
	if err != nil {
		return Rewrite{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("domains/%s/rewrites", s.domainName), r)
	if err != nil {
		return Rewrite{}, err
	}
	return doAndDecode[Rewrite](s.client, req)
}

func (s *RewriteService) Get(ctx context.Context) (Rewrite, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("domains/%s/rewrites/%s", s.domainName, s.name), nil)
	if err != nil {
		return Rewrite{}, err
	}
	return doAndDecode[Rewrite](s.client, req)
}

func (s *RewriteService) Update(ctx context.Context, request UpdateRewriteRequest) (Rewrite, error) {
	r, err := encode(request)
	if err != nil {
		return Rewrite{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("domains/%s/rewrites/%s", s.domainName, s.name), r)
	if err != nil {
		return Rewrite{}, err
	}
	return doAndDecode[Rewrite](s.client, req)
}

func (s *RewriteService) Delete(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("domains/%s/rewrites/%s", s.domainName, s.name), nil)
	if err != nil {
		return err
	}
	_, err = doAndDecode[struct{}](s.client, req)
	return err
}
