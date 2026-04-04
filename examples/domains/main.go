package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ionut-maxim/migadu"
)

func main() {
	client := migadu.New(
		os.Getenv("MIGADU_USER"),
		os.Getenv("MIGADU_API_KEY"),
	)
	ctx := context.Background()

	// List all domains
	domains, err := client.Domains().List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, d := range domains {
		fmt.Printf("domain: %s (%s)\n", d.Name, d.State)
	}

	// Get a single domain
	domain, err := client.Domains().Domain("example.com").Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("got: %s\n", domain.Name)

	// Create a domain
	created, err := client.Domains().Create(ctx, migadu.CreateDomainRequest{
		Name:      "newdomain.com",
		HostedDNS: false,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("created: %s\n", created.Name)

	// Update a domain
	updated, err := client.Domains().Domain("example.com").Update(ctx, migadu.UpdateDomainRequest{
		Description: "My main domain",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("updated: %s\n", updated.Name)

	// Get DNS records
	records, err := client.Domains().Domain("example.com").Records(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("SPF: %s\n", records.SPF.Value)
	for _, mx := range records.MXRecords {
		fmt.Printf("MX: %s (priority %d)\n", mx.Value, mx.Priority)
	}

	// Run diagnostics
	_, err = client.Domains().Domain("example.com").Diagnostics(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("diagnostics passed")

	// Activate a domain
	activated, err := client.Domains().Domain("example.com").Activate(ctx)
	if err != nil {
		if apiErr, ok := errors.AsType[*migadu.Error](err); ok && apiErr.StatusCode == 422 {
			fmt.Println("DNS checks not ready yet")
			return
		}
		log.Fatal(err)
	}
	fmt.Printf("activated: %s\n", activated.Name)
}
