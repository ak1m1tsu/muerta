package validator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	Name     string
	Value    string `validate:"notblank"`
	Expected bool
}

func Test_notBlank(t *testing.T) {
	testCases := []testCase{
		{
			Name:     "blank",
			Value:    "",
			Expected: false,
		},
		{
			Name:     "not blank",
			Value:    "test",
			Expected: true,
		},
		{
			Name:     "start with blank",
			Value:    "     test",
			Expected: false,
		},
		{
			Name:     "end with blank",
			Value:    "test ",
			Expected: false,
		},
		{
			Name:     "contains numbers",
			Value:    "test1",
			Expected: false,
		},
		{
			Name:     "with many spaces",
			Value:    "    ",
			Expected: false,
		},
		{
			Name:     "words with spaces",
			Value:    "test test",
			Expected: true,
		},
		{
			Name:     "words with many spaces",
			Value:    "test  test",
			Expected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actual := Validate(&tc)
			if tc.Expected {
				assert.Empty(t, actual.Error())
			} else {
				assert.Contains(t, actual.Error(), fmt.Sprintf("`value` with value `%s` doesn't satisfy the `notblank` constraint", tc.Value))
			}
		})
	}
}
