package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/avast/retry-go"
)

const clientID = "2a2a596e779f06ec61b6"
const scope = "repo,read:user,user:email"
const grantType = "urn:ietf:params:oauth:grant-type:device_code"

type githubDeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// Request a code for a device flow authentication
func Github_AskDeviceCode() (githubDeviceCodeResponse, error) {
	httpURL := fmt.Sprintf("https://github.com/login/device/code?client_id=%s&scope=%s", url.QueryEscape(clientID), url.QueryEscape(scope))
	req, err := http.NewRequest("POST", httpURL, nil)
	if err != nil {
		return githubDeviceCodeResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	var githubDeviceCodeResponse githubDeviceCodeResponse

	err = retry.Do(func() error {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return fmt.Errorf("response status code is not 200")
		}
		// Read body
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		// Unmarshal
		err = json.Unmarshal(content, &githubDeviceCodeResponse)
		if err != nil {
			return err
		}
		// Check if the response is valid
		if githubDeviceCodeResponse.DeviceCode == "" {
			return fmt.Errorf("invalid response")
		}
		if githubDeviceCodeResponse.UserCode == "" {
			return fmt.Errorf("invalid response")
		}

		return nil

	}, retry.Attempts(12), retry.Delay(1*time.Second))
	if err != nil {
		return githubDeviceCodeResponse, err
	}
	return githubDeviceCodeResponse, nil

}

// Define errors

var (
	ErrExpiredToken = errors.New("expired token")
)

type githubPollResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// Poll GitHub Api to check if the code has been entered by the user
func GitHub_PollToken(deviceCode string) (string, error) {
	httpURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&device_code=%s&grant_type=%s", url.QueryEscape(clientID), url.QueryEscape(deviceCode), url.QueryEscape(grantType))
	req, err := http.NewRequest("POST", httpURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// We use recursion in case of an error

	res, err := http.DefaultClient.Do(req)
	// Case the request failed
	if err != nil {
		time.Sleep(6 * time.Second)
		return GitHub_PollToken(deviceCode)
	}
	defer res.Body.Close()
	// Case the response is not valid
	if res.StatusCode != 200 {
		time.Sleep(6 * time.Second)
		return GitHub_PollToken(deviceCode)
	}
	// Read body
	content, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	// If the token is not ready yet, the json response will contains authorization_pending or expired_token
	if strings.Contains(string(content), "authorization_pending") {
		time.Sleep(6 * time.Second)
		return GitHub_PollToken(deviceCode)
	}
	// After 900 seconds, the token will expire. We need to handle this case
	if strings.Contains(string(content), "expired_token") {
		return "", ErrExpiredToken
	}

	// Unmarshal
	var decodedBody githubPollResponse
	err = json.Unmarshal(content, &decodedBody)
	if err != nil {
		return "", fmt.Errorf("couldn't decode response: %s", err)
	}
	// Check if the response is valid (in case we didn't catch an error earlier)
	if decodedBody.AccessToken == "" {
		return "", fmt.Errorf("invalid response from github. access_token is empty")
	}
	return decodedBody.AccessToken, nil

}

type githubEmail struct {
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
	Primary  bool   `json:"primary"`
}

// Get emails of a user from an oAuth token
func Github_GetEmails(token string) ([]string, error) {
	// Create request
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return nil, err
	}

	// Headers
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("X-Github-Api-Version", "2022-11-28")

	// Retry
	var emails []string
	err = retry.Do(func() error {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return fmt.Errorf("response status code is not 200")
		}
		// Read body
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		// Unmarshal
		var decodedBody []githubEmail
		err = json.Unmarshal(content, &decodedBody)
		if err != nil {
			return err
		}
		for _, email := range decodedBody {
			if email.Verified {
				emails = append(emails, email.Email)
			}
		}
		return nil
	}, retry.Attempts(12), retry.DelayType(retry.BackOffDelay))
	if err != nil {
		return nil, err
	}
	return emails, nil
}

type githubUser struct {
	Login string `json:"login"`
}

// Get the username of the user from an oAuth token
func GitHub_GetUserName(token string) (string, error) {
	// Create request
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return "", err
	}

	// Headers
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("X-Github-Api-Version", "2022-11-28")

	// Retry
	var user githubUser
	retry.Do(func() error {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return fmt.Errorf("response status code is not 200")
		}

		// Read body
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		// Unmarshal
		err = json.Unmarshal(content, &user)
		if err != nil {
			return err
		}
		return nil
	}, retry.Attempts(12), retry.DelayType(retry.BackOffDelay))
	return user.Login, nil

}
