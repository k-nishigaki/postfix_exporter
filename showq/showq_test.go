package showq

import (
	"bufio"
	"os"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func getSum(t *testing.T, histogram *prometheus.HistogramVec) float64 {
	t.Helper()
	var total float64
	metrics := make(chan prometheus.Metric, 10)
	histogram.Collect(metrics)
	close(metrics)

	for m := range metrics {
		dtoMetric := dto.Metric{}
		if err := m.Write(&dtoMetric); err != nil {
			t.Fatalf("failed to write metric: %v", err)
		}
		if hist := dtoMetric.GetHistogram(); hist != nil {
			total += hist.GetSampleSum()
		}
	}
	return total
}

func TestCollectShowqFromReader(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name               string
		file               string
		wantErr            bool
		expectedTotalCount float64
	}{
		{
			name:               "basic test",
			file:               "../testdata/showq.txt",
			wantErr:            false,
			expectedTotalCount: 118702,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			file, err := os.Open(tt.file)
			if err != nil {
				t.Error(err)
			}
			fd := bufio.NewReader(file)
			s := NewShowq("")
			s.init(fd)
			defer file.Close()

			if err := s.collectTextualShowqFromScanner(fd); (err != nil) != tt.wantErr {
				t.Errorf("CollectShowqFromReader() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.expectedTotalCount, getSum(t, s.sizeHistogram), "Expected a lot more data.")
			assert.Less(t, 0.0, getSum(t, s.ageHistogram), "Age not greater than 0")
		})
	}
}
