package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/stretchr/testify/assert"
)

func Test_CreateToken(t *testing.T) {
	testCases := []struct {
		name    string
		details *dto.TokenDetails
	}{
		{
			name: "valid details",
			details: &dto.TokenDetails{
				User: &dto.TokenPayload{
					UserID:   1,
					Username: "username",
					Roles:    []string{"user"},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rng := rand.Reader
			pk, _ := rsa.GenerateKey(rng, 4096)
			keyPem := pem.EncodeToMemory(&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(pk),
			})
			actual, err := CreateToken(tc.details.User, time.Minute*15, keyPem)
			assert.Nil(t, err)
			assert.NotNil(t, actual)
			assert.NotEmpty(t, actual.Token)
			assert.NotEmpty(t, actual.UUID)
			assert.Equal(t, tc.details.User.UserID, actual.User.UserID)
			assert.Equal(t, tc.details.User.Username, actual.User.Username)
			assert.Equal(t, tc.details.User.Roles, actual.User.Roles)
		})
	}
}
