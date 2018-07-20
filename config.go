// Package splunk enabled to perform restful calls to splunk instances
package splunk

import (
	"errors"
	"regexp"
	"strings"
)

const (
	// HostIPRegex describes regular expression for a valid IP based hostname
	HostIPRegex = `^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`
	// PortRegex describes regular expression for a valid Port
	PortRegex = `^([0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$`
	// HostNameRegex describes regular expression for a valid name based hostname
	HostNameRegex = `^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`

	defaultHost         = "localhost"
	defaultPort         = "8089"
	defaultBaseURL      = "https://" + defaultHost + ":" + defaultPort
	authBaseEP          = "/services/auth/login"
	defaultOutputFormat = "json"
	searchEndPoint      = "/services/search/jobs"
)

type URL struct {
	host, port, baseURL, authEndPoint, outputFormat string
}

func (url *URL) SetHost(host string) {
	url.host = host
}

func (url *URL) SetPort(port string) {
	url.port = port
}

func (user *User) Username() string {
	return user.username
}

func (pwd *User) Password() string {
	return pwd.password
}

// Returns host name for the instance
func (url *URL) Host() (string, error) {
	if url.host == "" {
		url.host = defaultHost
	}
	url.host = strings.TrimSpace(url.host)
	matchIPHost, err := regexp.MatchString(HostIPRegex, url.host)
	if err != nil {
		return "", err
	}
	matchNameHost, err := regexp.MatchString(HostNameRegex, url.host)
	if err != nil {
		return "", err
	}
	if matchIPHost != true && matchNameHost != true {
		err = errors.New("Invalid hostname :" + url.host + " provided")
	}
	return url.host, err
}

// Port returns port name for the instance
func (url *URL) Port() (string, error) {
	if url.port == "" {
		url.port = defaultPort
	}
	url.port = strings.TrimSpace(url.port)
	match, err := regexp.MatchString(PortRegex, url.port)
	if match != true {
		err = errors.New("Invalid port : " + url.port + " provided")
	} else if err != nil {
		return "", err
	}
	return url.port, err
}

// BaseUrl returns the https://host-port values for the url
func (url *URL) BaseUrl() (string, error) {
	host, err := url.Host()
	if err != nil {
		return "", err
	}
	port, err := url.Port()
	if err != nil {
		return "", err
	}
	if url.baseURL == "" {
		url.baseURL = CreateBaseUrl(host, port)
	}
	return url.baseURL, err
}

// AuthEndPoint reuturns baseURL + /auth/login
func (url URL) AuthEndPoint() (string, error) {
	baseURL, err := url.BaseUrl()
	if err != nil {
		return "", err
	}
	if url.authEndPoint == "" {
		url.authEndPoint = baseURL + authBaseEP
	}
	return url.authEndPoint, err
}

// SearchEndPoint returns protocol-host-port concat'ied search-end-point
func (url *URL) SearchEndPoint() (string, error) {
	baseURL, err := url.BaseUrl()
	if err != nil {
		return "", err
	}
	return baseURL + searchEndPoint, err
}

// OutputFormat sets the return format of the rest calls
func (url *URL) OutputFormat() string {
	if url.outputFormat == "" {
		url.outputFormat = defaultOutputFormat
	}
	return url.outputFormat
}

// CreateBaseUrl returns the protocol-host-port string
func CreateBaseUrl(host, port string) string {
	return "https://" + host + ":" + port
}
