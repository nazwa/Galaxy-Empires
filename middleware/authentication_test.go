package middleware

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

var (
	TEST_AUTH_KEY    []byte = []byte("This is a test key @@")
	TEST_PRIVATE_KEY []byte = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA4f5wg5l2hKsTeNem/V41fGnJm6gOdrj8ym3rFkEU/wT8RDtn
SgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7mCpz9Er5qLaMXJwZxzHzAahlfA0i
cqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBpHssPnpYGIn20ZZuNlX2BrClciHhC
PUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2XrHhR+1DcKJzQBSTAGnpYVaqpsAR
ap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3bODIRe1AuTyHceAbewn8b462yEWKA
Rdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy7wIDAQABAoIBAQCwia1k7+2oZ2d3
n6agCAbqIE1QXfCmh41ZqJHbOY3oRQG3X1wpcGH4Gk+O+zDVTV2JszdcOt7E5dAy
MaomETAhRxB7hlIOnEN7WKm+dGNrKRvV0wDU5ReFMRHg31/Lnu8c+5BvGjZX+ky9
POIhFFYJqwCRlopGSUIxmVj5rSgtzk3iWOQXr+ah1bjEXvlxDOWkHN6YfpV5ThdE
KdBIPGEVqa63r9n2h+qazKrtiRqJqGnOrHzOECYbRFYhexsNFz7YT02xdfSHn7gM
IvabDDP/Qp0PjE1jdouiMaFHYnLBbgvlnZW9yuVf/rpXTUq/njxIXMmvmEyyvSDn
FcFikB8pAoGBAPF77hK4m3/rdGT7X8a/gwvZ2R121aBcdPwEaUhvj/36dx596zvY
mEOjrWfZhF083/nYWE2kVquj2wjs+otCLfifEEgXcVPTnEOPO9Zg3uNSL0nNQghj
FuD3iGLTUBCtM66oTe0jLSslHe8gLGEQqyMzHOzYxNqibxcOZIe8Qt0NAoGBAO+U
I5+XWjWEgDmvyC3TrOSf/KCGjtu0TSv30ipv27bDLMrpvPmD/5lpptTFwcxvVhCs
2b+chCjlghFSWFbBULBrfci2FtliClOVMYrlNBdUSJhf3aYSG2Doe6Bgt1n2CpNn
/iu37Y3NfemZBJA7hNl4dYe+f+uzM87cdQ214+jrAoGAXA0XxX8ll2+ToOLJsaNT
OvNB9h9Uc5qK5X5w+7G7O998BN2PC/MWp8H+2fVqpXgNENpNXttkRm1hk1dych86
EunfdPuqsX+as44oCyJGFHVBnWpm33eWQw9YqANRI+pCJzP08I5WK3osnPiwshd+
hR54yjgfYhBFNI7B95PmEQkCgYBzFSz7h1+s34Ycr8SvxsOBWxymG5zaCsUbPsL0
4aCgLScCHb9J+E86aVbbVFdglYa5Id7DPTL61ixhl7WZjujspeXZGSbmq0Kcnckb
mDgqkLECiOJW2NHP/j0McAkDLL4tysF8TLDO8gvuvzNC+WQ6drO2ThrypLVZQ+ry
eBIPmwKBgEZxhqa0gVvHQG/7Od69KWj4eJP28kq13RhKay8JOoN0vPmspXJo1HY3
CKuHRG+AP579dncdUnOMvfXOtkdM4vk0+hWASBQzM9xzVcztCa+koAugjVaLS9A+
9uQoqEeVNTckxx0S2bYevRy7hGQmUJTyQm3j1zEUR5jpdbL83Fbq
-----END RSA PRIVATE KEY-----`)
)

func mockLoginToken(id string, expires time.Time, auth_key []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["id"] = id
	token.Claims["exp"] = expires.Unix()
	// Sign and get the complete encoded token as a string
	return token.SignedString(auth_key)
}
func mockLoginTokenWithoutID(expires time.Time, auth_key []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["exp"] = expires.Unix()
	// Sign and get the complete encoded token as a string
	return token.SignedString(auth_key)
}
func mockLoginTokenRS(id string, expires time.Time, auth_key []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	// Set some claims
	token.Claims["id"] = id
	token.Claims["exp"] = expires.Unix()
	// Sign and get the complete encoded token as a string
	return token.SignedString(auth_key)
}

func createGinContext(tokenText string) *gin.Context {
	c := &gin.Context{}
	c.Request = &http.Request{}
	c.Request.Header = make(http.Header)
	c.Request.Header.Set(AuthTokenHeaderKey, tokenText)

	return c
}

func TestAuthenticationSetUp(t *testing.T) {
	assert.Panics(t, func() {
		Authentication([]byte(""))
	})
	assert.NotPanics(t, func() {
		Authentication(TEST_AUTH_KEY)
	})

	validTokenText, err := mockLoginToken("10", time.Now().Add(10*time.Second), TEST_AUTH_KEY)
	assert.NoError(t, err)

	handler := Authentication(TEST_AUTH_KEY)

	// Run the handler with a valid token
	c := createGinContext(validTokenText)
	assert.NotPanics(t, func() {
		handler(c)
	})

	assert.Empty(t, c.Errors)
	assert.Equal(t, "10", c.MustGet(AuthUserIDKey))
	assert.IsType(t, &jwt.Token{}, c.MustGet(AuthTokenKey))

	// Run the handler with invalid token
	// There is no easy way to create a gin writer to properly detect response
	// Go will panic when Authentication tries to add error to an empty Writer object
	// So we can assume that the validation failed and user is not allowed through :)
	c = createGinContext("RANDOM TOKEN TEXT")
	assert.Panics(t, func() {
		handler(c)
	})
}

func TestAuthentication(t *testing.T) {
	validTokenText, err := mockLoginToken("10", time.Now().Add(10*time.Second), TEST_AUTH_KEY)
	assert.NoError(t, err)

	// This is the only one where the token is valid
	token, err := ValidateToken(validTokenText, TEST_AUTH_KEY)
	assert.NoError(t, err)
	assert.IsType(t, &jwt.Token{}, token)

	// Empty token
	token, err = ValidateToken("", TEST_AUTH_KEY)
	assert.Error(t, err)
	assert.Nil(t, token)

	// Random token
	token, err = ValidateToken("afawfaw fawfaw faw fawfaw fwaf.wfwfwf.wfwfawf", TEST_AUTH_KEY)
	assert.Error(t, err)
	assert.Nil(t, token)

	// Expired token
	invalidTokenText, err := mockLoginToken("10", time.Now().Add(-10*time.Minute), TEST_AUTH_KEY)
	assert.NoError(t, err)

	token, err = ValidateToken(invalidTokenText, TEST_AUTH_KEY)
	assert.Error(t, err)
	assert.Nil(t, token)

	// Wrong encryption key
	invalidTokenText, err = mockLoginToken("10", time.Now().Add(10*time.Second), []byte("WRONG KEY"))
	assert.NoError(t, err)

	token, err = ValidateToken(invalidTokenText, TEST_AUTH_KEY)
	assert.Error(t, err)
	assert.Nil(t, token)

	// Token doesn't have id
	invalidTokenText, err = mockLoginTokenWithoutID(time.Now().Add(10*time.Second), TEST_AUTH_KEY)
	assert.NoError(t, err)

	token, err = ValidateToken(invalidTokenText, TEST_AUTH_KEY)
	assert.Error(t, err)
	assert.Nil(t, token)

	// Different signing algorithm
	invalidTokenText, err = mockLoginTokenRS("10", time.Now().Add(10*time.Second), TEST_PRIVATE_KEY)
	assert.NoError(t, err)

	token, err = ValidateToken(invalidTokenText, TEST_AUTH_KEY)
	assert.Error(t, err)
	assert.Nil(t, token)

}
