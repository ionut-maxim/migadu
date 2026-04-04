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

	// List rewrites
	rewrites, err := domain.Rewrites().List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range rewrites {
		fmt.Printf("rewrite: %s (%s) -> %v\n", r.Name, r.LocalPartRule, r.Destinations)
	}

	// Create a rewrite
	created, err := domain.Rewrites().Create(ctx, migadu.CreateRewriteRequest{
		Name:          "catch-support",
		LocalPartRule: "support+*",
		Destinations:  []string{"john@example.com"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("created: %s\n", created.Name)

	// Get a rewrite
	rewrite, err := domain.Rewrites().Rewrite("catch-support").Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("got: %s -> %v\n", rewrite.Name, rewrite.Destinations)

	// Update a rewrite
	updated, err := domain.Rewrites().Rewrite("catch-support").Update(ctx, migadu.UpdateRewriteRequest{
		Destinations: []string{"john@example.com", "jane@example.com"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("updated destinations: %v\n", updated.Destinations)

	// Delete a rewrite
	if err := domain.Rewrites().Rewrite("catch-support").Delete(ctx); err != nil {
		if apiErr, ok := errors.AsType[*migadu.Error](err); ok {
			fmt.Printf("delete failed: %s\n", apiErr.Message)
			return
		}
		log.Fatal(err)
	}
	fmt.Println("deleted rewrite")
}
