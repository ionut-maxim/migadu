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

	// List identities
	identities, err := mailbox.Identities().List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range identities {
		fmt.Printf("identity: %s\n", i.Address)
	}

	// Create an identity
	created, err := mailbox.Identities().Create(ctx, migadu.CreateIdentityRequest{
		Name:      "John Work",
		LocalPart: "john.work",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("created: %s\n", created.Address)

	// Get an identity
	identity, err := mailbox.Identities().Identity("john.work").Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("got: %s\n", identity.Address)

	// Update an identity
	updated, err := mailbox.Identities().Identity("john.work").Update(ctx, migadu.UpdateIdentityRequest{
		Name: "John (Work)",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("updated: %s\n", updated.Name)

	// Delete an identity
	if err = mailbox.Identities().Identity("john.work").Delete(ctx); err != nil {
		if apiErr, ok := errors.AsType[*migadu.Error](err); ok {
			fmt.Printf("delete failed: %s\n", apiErr.Message)
			return
		}
		log.Fatal(err)
	}
	fmt.Println("deleted identity")
}
