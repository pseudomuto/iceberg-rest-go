package tests

import (
	"context"
	"net/http"
	"testing"

	. "github.com/pseudomuto/iceberg-rest-go"
	"github.com/stretchr/testify/require"
)

func TestCreateTable(t *testing.T) {
	record(t, func(client ClientWithResponsesInterface) {
		req := CreateTableJSONRequestBody{
			Name:     "my-table",
			Location: Ptr("/tmp/testing/tables/my-table"),
			Schema: Schema{
				Type: SchemaTypeStruct,
				Fields: []StructField{
					{
						Id:       1,
						Name:     "id",
						Required: true,
						Type:     Long,
						Doc:      Ptr("The unique id of the record"),
					},
					{
						Id:       2,
						Name:     "created_at",
						Required: true,
						Type:     Timestamp,
						Doc:      Ptr("When this record was created"),
					},
					{
						Id:   3,
						Name: "properties",
						Type: ToType(MapType{
							Type:    Map,
							Key:     String,
							Value:   String,
							KeyId:   4,
							ValueId: 5,
						}),
					},
					{
						Id:   6,
						Name: "thing_list",
						Type: ToType(ListType{
							Type:            List,
							Element:         String,
							ElementId:       7,
							ElementRequired: true,
						}),
					},
					{
						Id:   8,
						Name: "details",
						Type: ToType(StructType{
							Type: StructTypeTypeStruct,
							Fields: []StructField{
								{
									Id:   9,
									Name: "name",
									Type: String,
								},
								{
									Id:   10,
									Name: "age",
									Type: Int,
								},
							},
						}),
					},
				},
				IdentifierFieldIds: &[]int{1},
			},
			PartitionSpec: &PartitionSpec{
				SpecId: Ptr(1),
				Fields: []PartitionField{
					{
						Name:      "Hourly By Created At",
						SourceId:  2,
						FieldId:   Ptr(1),
						Transform: "hour",
					},
				},
			},
			WriteOrder: &SortOrder{
				OrderId: Ptr(1),
				Fields: []SortField{
					{
						SourceId:  2,
						Direction: Asc,
						NullOrder: NullsLast,
					},
				},
			},
			Properties: &map[string]string{
				"write.metadata.delete-after-commit.enabled": "true",
				"write.metadata.previous-versions-max":       "5",
				"write.object-store.enabled":                 "true",
				"write.parquet.compression-codec":            "zstd",
				"write.target-file-size-bytes":               "536870912",
			},
		}

		ctx := context.Background()

		// Create the namespace
		ns, err := client.CreateNamespaceWithResponse(ctx, CreateNamespaceRequest{
			Namespaces: Namespaces{"testing", "tables"},
			Properties: &map[string]string{
				"location": "/tmp/testing/tables",
			},
		})
		require.NoError(t, err)

		nsString := NamespaceString(ns.JSON200.Namespaces)

		resp, err := client.CreateTableWithResponse(ctx, nsString, nil, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode())

		require.Equal(t, req.Properties, resp.JSON200.Metadata.Properties)
		require.Equal(t, *req.Location, *resp.JSON200.Metadata.Location)
	})
}
