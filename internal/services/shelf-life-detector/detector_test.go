package shelflifedetector

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Detect(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		expected []string
	}{
		{
			name: "date start and date end",
			path: `test_data.webp`,
			expected: []string{
				"15.09.22",
				"24.09.22",
			},
		},
	}
	svc := DateDetectorService{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, _ := os.ReadFile(tc.path)
			dates, err := svc.Detect(data)
			assert.Nil(t, err)
			assert.NotNil(t, dates)
			assert.NotEmpty(t, dates)
			assert.Equal(t, tc.expected, dates)
		})
	}
}
