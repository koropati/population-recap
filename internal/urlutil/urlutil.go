package urlutil

import (
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func GetBaseURL(request *http.Request) string {
	// Dapatkan skema (http atau https)
	scheme := "http"
	if request.TLS != nil {
		scheme = "https"
	}

	// Dapatkan host (termasuk port jika ada)
	host := request.Host

	// Gabungkan skema dan host untuk membentuk base URL
	baseURL := scheme + "://" + host

	return baseURL
}

func RemoveAPIVersionMiddleware(pathURL string) string {
	regex := regexp.MustCompile(`^/v[0-9]+/`)
	newPath := regex.ReplaceAllString(pathURL, "/")
	return newPath
}

func GetFirstPathName(pathURL string) string {
	var newPath string

	if strings.Contains(pathURL, "my/") {
		regex := regexp.MustCompile(`^/v[0-9]+/my/`)
		newPath = regex.ReplaceAllString(pathURL, "/")
	} else {
		regex := regexp.MustCompile(`^/v[0-9]+/`)
		newPath = regex.ReplaceAllString(pathURL, "/")
	}
	if strings.Contains(newPath, "/") {
		arrPath := strings.Split(newPath, "/")
		if len(arrPath) > 0 {
			if len(arrPath) == 1 && arrPath[0] != "" {
				return arrPath[0]
			} else {
				return arrPath[1]
			}
		} else {
			return ""
		}
	} else {
		return newPath
	}

}

func CreateUrlVerificationEmail(request *http.Request, verificationToken string) (result string) {
	baseUrl := GetBaseURL(request)
	result = baseUrl + "/verify-email?token=" + url.QueryEscape(verificationToken)
	return result
}

func GenerateRedirectForgotPassword(baseUrlFrontEnd string, pathUrl string, isSuccess bool, msg string, forgotPasswordToken string, name string) (urlRedirectToFE string) {
	urlRedirectToFE = baseUrlFrontEnd + pathUrl + "?success=" + strconv.FormatBool(isSuccess) + "&msg=" + msg + "&t=" + url.QueryEscape(forgotPasswordToken) + "&name=" + name
	return
}

func CreateUrlForgotPassword(request *http.Request, forgotPasswordToken string) (result string) {
	baseUrl := GetBaseURL(request)
	result = baseUrl + "/forgot-password/verify?token=" + url.QueryEscape(forgotPasswordToken)
	return result
}
