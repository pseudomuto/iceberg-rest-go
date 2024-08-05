package tests

// All of these tests use VCR to record the HTTP interactions with the API server.
//
// This test suite is by no means exhaustive. It's meant to be used as a quick smoke test to ensure translations are
// happening appropriately and that we can make requests the API server. If you find a bug, open a PR with a regression
// test! Or if that's not your jam, an issue would be great. Don't forget to tag @pseudomuto on it.
//
// If you need to update them, follow these steps:
//
// * Delete the cassettes (in testdata/cassettes) corresponding with the test you want to rerecord.
// * Run `make catalog` (an ephemeral API server).
// * Run `make test` which will rerecord any tests without cassettes.

import (
	"os"
	"testing"

	. "github.com/pseudomuto/iceberg-rest-go"
	"github.com/seborama/govcr/v15"
	"github.com/stretchr/testify/require"
)

var host = "http://localhost:8181"

func TestMain(m *testing.M) {
	if v, ok := os.LookupEnv("CATALOG_HOST"); ok {
		host = v
	}

	os.Exit(m.Run())
}

func record(t *testing.T, fn func(ClientWithResponsesInterface)) {
	vcr := govcr.NewVCR(
		govcr.NewCassetteLoader("testdata/"+t.Name()+".json"),
		govcr.WithRequestMatchers(govcr.NewMethodURLRequestMatchers()...),
	)

	client, err := NewClientWithResponses(host, WithHTTPClient(vcr.HTTPClient()))
	require.NoError(t, err)

	fn(client)
}
