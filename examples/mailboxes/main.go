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

	// List mailboxes
	mailboxes, err := domain.Mailboxes().List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range mailboxes {
		fmt.Printf("mailbox: %s\n", m.Address)
	}

	// Create a mailbox
	created, err := domain.Mailboxes().Create(ctx, migadu.CreateMailboxRequest{
		Name:      "John Doe",
		LocalPart: "john",
		Password:  "supersecret",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("created: %s\n", created.Address)

	// Get a mailbox
	mailbox, err := domain.Mailboxes().Mailbox("john").Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("got: %s (may_send=%v)\n", mailbox.Address, mailbox.MaySend)

	// Update a mailbox
	updated, err := domain.Mailboxes().Mailbox("john").Update(ctx, migadu.UpdateMailboxRequest{
		Name: "John D.",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("updated: %s\n", updated.Name)

	// Delete a mailbox
	if err = domain.Mailboxes().Mailbox("john").Delete(ctx); err != nil {
		if apiErr, ok := errors.AsType[*migadu.Error](err); ok {
			fmt.Printf("delete failed: %s\n", apiErr.Message)
			return
		}
		log.Fatal(err)
	}
	fmt.Println("deleted mailbox")
}
