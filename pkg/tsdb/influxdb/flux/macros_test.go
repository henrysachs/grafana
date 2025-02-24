package flux

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInterpolate(t *testing.T) {
	timeRange := backend.TimeRange{
		From: time.Unix(1632305571, 310985041),
		To:   time.Unix(1632309171, 310985042),
	}

	options := queryOptions{
		Organization:  "grafana1",
		Bucket:        "grafana2",
		DefaultBucket: "grafana3",
	}

	tests := []struct {
		name   string
		before string
		after  string
	}{
		{
			name:   "interpolate flux variables",
			before: `v.timeRangeStart, something.timeRangeStop, XYZ.bucket, uuUUu.defaultBucket, aBcDefG.organization, window.windowPeriod, a91{}.bucket`,
			after:  `2021-09-22T10:12:51.310985041Z, 2021-09-22T11:12:51.310985042Z, "grafana2", "grafana3", "grafana1", 1s, a91{}.bucket`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := queryModel{
				RawQuery:      tt.before,
				Options:       options,
				TimeRange:     timeRange,
				MaxDataPoints: 1,
				Interval:      1000 * 1000 * 1000,
			}
			interpolatedQuery, err := interpolate(query)
			require.NoError(t, err)
			diff := cmp.Diff(tt.after, interpolatedQuery)
			assert.Equal(t, "", diff)
		})
	}
}
