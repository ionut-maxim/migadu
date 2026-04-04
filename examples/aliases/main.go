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
	domain := client.Domains().Domain("example.com")

	// List aliases
	aliases, err := domain.Aliases().List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, a := range aliases {
		fmt.Printf("alias: %s -> %v\n", a.Address, a.Destinations)
	}

	// Create an alias
	created, err := domain.Aliases().Create(ctx, migadu.CreateAliasRequest{
		LocalPart:    "info",
		Destinations: []string{"john@example.com"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("created: %s\n", created.Address)

	// Get an alias
	alias, err := domain.Aliases().Alias("info").Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("got: %s -> %v\n", alias.Address, alias.Destinations)

	// Update an alias
	updated, err := domain.Aliases().Alias("info").Update(ctx, migadu.UpdateAliasRequest{
		Destinations: []string{"john@example.com", "jane@example.com"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("updated destinations: %v\n", updated.Destinations)

	// Delete an alias
	if err := domain.Aliases().Alias("info").Delete(ctx); err != nil {
		if apiErr, ok := errors.AsType[*migadu.Error](err); ok {
			fmt.Printf("delete failed: %s\n", apiErr.Message)
			return
		}
		log.Fatal(err)
	}
	fmt.Println("deleted alias")
}
