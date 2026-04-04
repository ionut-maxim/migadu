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
	mailbox := client.Domains().Domain("example.com").Mailboxes().Mailbox("john")

	// List forwardings
	forwardings, err := mailbox.Forwardings().List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range forwardings {
		fmt.Printf("forwarding: %s (active=%v)\n", f.Address, f.IsActive)
	}

	// Create a forwarding
	created, err := mailbox.Forwardings().Create(ctx, migadu.CreateForwardingRequest{
		Address: "backup@gmail.com",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("created: %s\n", created.Address)

	// Get a forwarding
	forwarding, err := mailbox.Forwardings().Forwarding("backup@gmail.com").Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("got: %s (active=%v)\n", forwarding.Address, forwarding.IsActive)

	// Update a forwarding
	active := false
	updated, err := mailbox.Forwardings().Forwarding("backup@gmail.com").Update(ctx, migadu.UpdateForwardingRequest{
		IsActive: &active,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("updated active: %v\n", updated.IsActive)

	// Delete a forwarding
	if err := mailbox.Forwardings().Forwarding("backup@gmail.com").Delete(ctx); err != nil {
		if apiErr, ok := errors.AsType[*migadu.Error](err); ok {
			fmt.Printf("delete failed: %s\n", apiErr.Message)
			return
		}
		log.Fatal(err)
	}
	fmt.Println("deleted forwarding")
}
