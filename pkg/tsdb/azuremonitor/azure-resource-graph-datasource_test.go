package azuremonitor

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/plugins"
	"github.com/stretchr/testify/require"
)

func TestBuildingAzureResourceGraphQueries(t *testing.T) {
	datasource := &AzureResourceGraphDatasource{}
	fromStart := time.Date(2018, 3, 15, 13, 0, 0, 0, time.UTC).In(time.Local)

	tests := []struct {
		name                      string
		queryModel                []backend.DataQuery
		timeRange                 plugins.DataTimeRange
		azureResourceGraphQueries []*AzureResourceGraphQuery
		Err                       require.ErrorAssertionFunc
	}{
		{
			name: "Query with macros should be interpolated",
			timeRange: plugins.DataTimeRange{
				From: fmt.Sprintf("%v", fromStart.Unix()*1000),
				To:   fmt.Sprintf("%v", fromStart.Add(34*time.Minute).Unix()*1000),
			},
			queryModel: []backend.DataQuery{
				{
					JSON: []byte(`{
						"queryType": "Azure Resource Graph",
						"azureResourceGraph": {
							"query":        "resources | where $__contains(name,'res1','res2')",
							"resultFormat": "table"
						}
					}`),
					RefID: "A",
				},
			},
			azureResourceGraphQueries: []*AzureResourceGraphQuery{
				{
					RefID:        "A",
					ResultFormat: "table",
					URL:          "",
					JSON: []byte(`{
						"queryType": "Azure Resource Graph",
						"azureResourceGraph": {
							"query":        "resources | where $__contains(name,'res1','res2')",
							"resultFormat": "table"
						}
					}`),
					InterpolatedQuery: "resources | where ['name'] in ('res1','res2')",
				},
			},
			Err: require.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queries, err := datasource.buildQueries(tt.queryModel, datasourceInfo{})
			tt.Err(t, err)
			if diff := cmp.Diff(tt.azureResourceGraphQueries, queries, cmpopts.IgnoreUnexported(simplejson.Json{})); diff != "" {
				t.Errorf("Result mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
