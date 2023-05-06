package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateHashFromPassword(t *testing.T) {
	testCase := struct {
		name     string
		password string
		salt     string
		expected string
	}{
		name:     "valid hash",
		password: "4*h0L28f#0198",
		salt:     "30612ede-2254-4708-9f4e-90afedbc33fb",
		expected: "77e31c1175a447652dec7a1665a0a6abff4933d11ccf044b6e95106e0fb28a5b",
	}
	t.Run(testCase.name, func(t *testing.T) {
		hash := GenerateHashFromPassword(testCase.password, testCase.salt)
		assert.NotEmpty(t, hash)
		assert.Equal(t, testCase.expected, hash)
	})
}

func Test_CompareHashAndPassword(t *testing.T) {
	testCase := struct {
		name     string
		password string
		salt     string
		hash     string
		expected bool
	}{
		name:     "valid compare",
		password: "4*h0L28f#0198",
		salt:     "30612ede-2254-4708-9f4e-90afedbc33fb",
		hash:     "77e31c1175a447652dec7a1665a0a6abff4933d11ccf044b6e95106e0fb28a5b",
		expected: true,
	}
	t.Run(testCase.name, func(t *testing.T) {
		isCompare := CompareHashAndPassword(testCase.password, testCase.salt, testCase.hash)
		assert.Equal(t, testCase.expected, isCompare)
	})
}
