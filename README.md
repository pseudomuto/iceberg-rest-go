# iceberg-rest-go

[![CI Status][ci-svg]][ci-url]

A Go client library for interacting with an [Iceberg Rest Catalog].

## Getting Started

```
go get https://github.com/pseudomuto/iceberg-rest-go

// Or simply
import "github.com/pseudomuto/iceberg-rest-go"
```

The test cases are the best place to look for examples, but just as a quick demo, here's how you can use this library.

**Getting a Namespace**

```go
func main() {
  client, err := NewClient("http://localhost:8181")
  if err != nil {
    log.Fatal(err)
  }

  ns := Namespaces{"data", "testing"}
  res, err := client.LoadNamespaceMetadataWithResponse(ctx, NamespaceString(ns))
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("Namespace properties:")
  for k, v := range *nsResp.JSON200.Properties {
    fmt.Printf("Key: %s, Value: %s\n", k, v)
  }
}
```

## Contributing

This library is (largely) auto-generated from the [OpenAPI spec] for the Rest catalog. Running `make update-api` will
fetch the latest version and perform a few updates to the downloaded file to make it work correctly (see the Makefile
for changes).

Once there's a new version, running `make generate` will generate the new Go code in _client.gen.go_.

`make help` is your friend.

### Testing

The tests all run against a running rest server (`make catalog`) and record their results in _tests/testdata/_ so future
test runs won't require a live server.

If you need to create a new test, create a new test file, run `make catalog` and then `make test` which will record the
interaction. To re-record, simply delete the JSON file and run again (ensuring the catalog is running). You may need to
restart the catalog server to ensure there's no existing data preventing the test from succeeding (it's ephemeral).

> Each test should be completely independent so that we can re-record one test without depending on ordering from any of
the others.

[Iceberg Rest Catalog]: https://iceberg.apache.org/concepts/catalog/#catalog-implementations
[OpenAPI spec]: https://github.com/apache/iceberg/raw/main/open-api/rest-catalog-open-api.yaml
[ci-svg]: https://github.com/pseudomuto/iceberg-rest-go/actions/workflows/ci.yaml/badge.svg?branch=main
[ci-url]: https://github.com/pseudomuto/iceberg-rest-go/actions/workflows/ci.yaml
