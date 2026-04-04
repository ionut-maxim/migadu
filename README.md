# migadu

A Go client for the [Migadu API](https://migadu.com/api/).

## Install

```sh
go get github.com/ionut-maxim/migadu
```

## Usage

```go
client := migadu.New(os.Getenv("MIGADU_USER"), os.Getenv("MIGADU_API_KEY"))

domains, err := client.Domains().List(ctx)
mailboxes, err := client.Domains().Domain("example.com").Mailboxes().List(ctx)
aliases, err := client.Domains().Domain("example.com").Aliases().List(ctx)
```

## Resources

| Resource | Scope |
|---|---|
| Domains | `client.Domains()` |
| Aliases | `client.Domains().Domain(name).Aliases()` |
| Rewrites | `client.Domains().Domain(name).Rewrites()` |
| Mailboxes | `client.Domains().Domain(name).Mailboxes()` |
| Identities | `...Mailboxes().Mailbox(name).Identities()` |
| Forwardings | `...Mailboxes().Mailbox(name).Forwardings()` |

## Error handling

API errors come back as `*migadu.Error` with a status code and message.

```go
if apiErr, ok := errors.AsType[*migadu.Error](err); ok {
    fmt.Println(apiErr.StatusCode, apiErr.Message)
}
```

## Examples

See the [examples](./examples) folder for working code covering each resource.
