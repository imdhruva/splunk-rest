// Package splunk enabled to perform restful calls to splunk instances
package splunk

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type User struct {
	username, password string
	sessionKey         SessionKey
}

// sets the username required for Splunk authentication
func (user *User) SetUsername(username string) {
	username = strings.TrimSpace(username)
	user.username = username
}

// sets the password for splunk authentication
func (user *User) SetPassword(password string) {
	user.password = password
}

// SessionKey represents json struct for the login response from Splunk
type SessionKey struct {
	Value string `json:"sessionKey"`
}

// BasicAuth performs http form if basic authentication and returns token
// for subsequent http calls within the session
func (url URL) BasicAuth(user User) error {
	payload := url.GetAuthPayload(user)

	authEndPoint, err := url.AuthEndPoint()
	if err != nil {
		return err
	}
	resp, err := user.HttpCall("POST", authEndPoint, payload)

	if err != nil {
		return err
	}
	bytes := []byte(resp)
	user.SetSessionToken(bytes)
	return nil
}

// HttpCall returns the byte slices of http calls
func (user *User) HttpCall(method string, url string, payload io.Reader) ([]byte, error) {
	client := HttpClient()
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	user.AddAuthHeader(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	bytes := []byte(body)

	return bytes, nil
}

// SetSessionToken returns the session token if none is set right now
func (user *User) SetSessionToken(bytes []byte) error {
	var key SessionKey
	unmarshalError := json.Unmarshal(bytes, &key)
	user.sessionKey = key
	return unmarshalError
}

// AddAuthHeader returns http-request after adding header for sessionKey if it exists
func (user *User) AddAuthHeader(request *http.Request) {
	if user.sessionKey.Value != "" {
		request.Header.Add("Authorization", fmt.Sprintf("Splunk %s", user.sessionKey.Value))
	} else {
		request.SetBasicAuth(user.Username(), user.Password())
	}
}

// HttpClient returns the http.client required for performing the login
func HttpClient() http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	return *client
}

// GetAuthPayload returns the payload for performing the basic authentication
func (u *URL) GetAuthPayload(user User) io.Reader {
	var payload io.Reader
	val := url.Values{}
	val.Add("username", user.Username())
	val.Add("password", user.Password())
	val.Add("output_mode", u.OutputFormat())

	if val != nil {
		payload = bytes.NewBufferString(val.Encode())
	}
	return payload
}
