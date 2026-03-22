package enrollment

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComputePaymentState(t *testing.T) {
	tests := []struct {
		name           string
		finalFee       float64
		paid           float64
		expectedRemain float64
		expectedStatus string
	}{
		{name: "unpaid", finalFee: 100, paid: 0, expectedRemain: 100, expectedStatus: "unpaid"},
		{name: "partial", finalFee: 100, paid: 40, expectedRemain: 60, expectedStatus: "partial"},
		{name: "paid", finalFee: 100, paid: 100, expectedRemain: 0, expectedStatus: "paid"},
		{name: "over paid clamps", finalFee: 100, paid: 120, expectedRemain: 0, expectedStatus: "paid"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			remaining, status := computePaymentState(tc.finalFee, tc.paid)
			require.Equal(t, tc.expectedRemain, remaining)
			require.Equal(t, tc.expectedStatus, status)
		})
	}
}
