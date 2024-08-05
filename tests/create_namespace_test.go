package tests

import (
	"context"
	"net/http"
	"testing"

	. "github.com/pseudomuto/iceberg-rest-go"
	"github.com/stretchr/testify/require"
)

func TestCreateNamespace(t *testing.T) {
	record(t, func(client ClientWithResponsesInterface) {
		req := CreateNamespaceJSONRequestBody{
			Namespaces: Namespaces{"testing", "newnamespace"},
			Properties: &map[string]string{
				"owner":    "testsuite",
				"location": "file:///tmp/testing/newnamespace",
			},
		}

		ctx := context.Background()
		resp, err := client.CreateNamespaceWithResponse(ctx, req)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode())
		require.Equal(t, req.Namespaces, resp.JSON200.Namespaces)
		require.Equal(t, req.Properties, resp.JSON200.Properties)

		// Ensure it exists.
		nsResp, err := client.LoadNamespaceMetadataWithResponse(ctx, NamespaceString(req.Namespaces))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, nsResp.StatusCode())

		// Handles Iceberg errors appropriately.
		resp, err = client.CreateNamespaceWithResponse(ctx, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, resp.StatusCode())
		require.Equal(t, "AlreadyExistsException", resp.JSON409.Error.Type)
		require.EqualError(t, resp.Error(), "Namespace already exists: testing.newnamespace")

		// Delete the namespace.
		delResp, err := client.DropNamespaceWithResponse(ctx, NamespaceString(req.Namespaces))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, delResp.StatusCode())

		// Ensure it's gone.
		nsResp, err = client.LoadNamespaceMetadataWithResponse(ctx, NamespaceString(req.Namespaces))
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, nsResp.StatusCode())
	})
}
