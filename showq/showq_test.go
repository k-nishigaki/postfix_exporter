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

func getCount(t *testing.T, gauge *prometheus.GaugeVec) map[string]float64 {
	t.Helper()
	var values = make(map[string]float64)
	metrics := make(chan prometheus.Metric, 10)
	gauge.Collect(metrics)
	close(metrics)

	for m := range metrics {
		dtoMetric := dto.Metric{}
		if err := m.Write(&dtoMetric); err != nil {
			t.Fatalf("failed to write metric: %v", err)
		}
		if gaugeValue := dtoMetric.GetGauge(); gaugeValue != nil {
			for _, label := range dtoMetric.Label {
				if label.GetName() == "queue" {
					values[label.GetValue()] = gaugeValue.GetValue()
				}
			}
		}
	}
	return values
}

func TestCollectShowqFromReader(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                string
		file                string
		wantErr             bool
		expectedTotalCount  float64
		expectedActiveCount float64
		expectedHoldCount   float64
		expectedOtherCount  float64
	}{
		{
			name:                "basic test",
			file:                "../testdata/showq.txt",
			wantErr:             false,
			expectedTotalCount:  118702,
			expectedActiveCount: 16,
			expectedHoldCount:   0,
			expectedOtherCount:  8,
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
			counts := getCount(t, s.queueMessageGauge)
			assert.Equal(t, tt.expectedActiveCount, counts["active"], "Expected active count to match")
			assert.Equal(t, tt.expectedHoldCount, counts["hold"], "Expected hold count to match")
			assert.Equal(t, tt.expectedOtherCount, counts["other"], "Expected other count to match")
		})
	}
}
