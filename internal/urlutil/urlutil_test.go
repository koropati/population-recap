package urlutil_test

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/koropati/population-recap/internal/urlutil"
	"github.com/stretchr/testify/assert"
)

const (
	errMsgUnexpectedResult = "Unexpected result"
	urlUser                = "/users"
	exampleHost            = "example.com"
)

func TestGetBaseURL(t *testing.T) {
	// Kasus uji untuk request HTTP
	reqHTTP := &http.Request{
		Host: exampleHost,
	}

	// Memeriksa hasil GetBaseURL untuk request HTTP
	resultHTTP := urlutil.GetBaseURL(reqHTTP)
	assert.Equal(t, "http://example.com", resultHTTP, "Unexpected result for HTTP request")

	// Kasus uji untuk request HTTPS
	reqHTTPS := &http.Request{
		Host: "secure.example.com",
		TLS:  &tls.ConnectionState{},
	}

	// Memeriksa hasil GetBaseURL untuk request HTTPS
	resultHTTPS := urlutil.GetBaseURL(reqHTTPS)
	assert.Equal(t, "https://secure.example.com", resultHTTPS, "Unexpected result for HTTPS request")
}

func TestCreateUrlVerificationEmail(t *testing.T) {
	request := &http.Request{
		Host: exampleHost,
	}
	verificationToken := "abc123"

	expectedResult := "http://example.com/verify-email?token=abc123"

	result := urlutil.CreateUrlVerificationEmail(request, verificationToken)
	assert.Equal(t, expectedResult, result, errMsgUnexpectedResult)
}

func TestRemoveAPIVersionMiddleware(t *testing.T) {
	// Kasus uji untuk pathURL dengan API version
	pathURLWithAPIVersion := "/v1/users"
	resultWithAPIVersion := urlutil.RemoveAPIVersionMiddleware(pathURLWithAPIVersion)
	assert.Equal(t, urlUser, resultWithAPIVersion, errMsgUnexpectedResult)

	// Kasus uji untuk pathURL tanpa API version
	pathURLWithoutAPIVersion := urlUser
	resultWithoutAPIVersion := urlutil.RemoveAPIVersionMiddleware(pathURLWithoutAPIVersion)
	assert.Equal(t, urlUser, resultWithoutAPIVersion, errMsgUnexpectedResult)
}
func TestGetFirstPathName(t *testing.T) {
	tests := []struct {
		name     string
		pathURL  string
		expected string
	}{
		{
			name:     "Dengan API version dan prefix 'my/'",
			pathURL:  "/v1/my/users",
			expected: "users",
		},
		{
			name:     "Dengan API version tanpa prefix 'my/'",
			pathURL:  "/v1/users",
			expected: "users",
		},
		{
			name:     "Tanpa API version dan prefix 'my/'",
			pathURL:  "/my/users",
			expected: "my",
		},
		{
			name:     "Tanpa API version dan tanpa prefix 'my/'",
			pathURL:  urlUser,
			expected: "users",
		},
		{
			name:     "Path hanya berisi satu segment",
			pathURL:  "/v1/",
			expected: "",
		},
		{
			name:     "Path dengan multiple segments",
			pathURL:  "/v1/my/users/details",
			expected: "users",
		},
		{
			name:     "Path kosong",
			pathURL:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := urlutil.GetFirstPathName(tt.pathURL)
			assert.Equal(t, tt.expected, result, errMsgUnexpectedResult)
		})
	}
}

func TestGenerateRedirectForgotPassword(t *testing.T) {
	baseUrlFrontEnd := "http://example.com"
	pathUrl := "/reset-password"
	isSuccess := true
	msg := "Berhasil"
	forgotPasswordToken := "abc123"

	expectedResult := "http://example.com/reset-password?success=true&msg=Berhasil&t=abc123&name="

	result := urlutil.GenerateRedirectForgotPassword(baseUrlFrontEnd, pathUrl, isSuccess, msg, forgotPasswordToken, "")
	assert.Equal(t, expectedResult, result, errMsgUnexpectedResult)
}

func TestCreateUrlForgotPassword(t *testing.T) {
	request := &http.Request{
		Host: exampleHost,
	}
	forgotPasswordToken := "abc123"

	expectedResult := "http://example.com/verify-forgot-password?token=abc123"

	result := urlutil.CreateUrlForgotPassword(request, forgotPasswordToken)
	assert.Equal(t, expectedResult, result, errMsgUnexpectedResult)
}
