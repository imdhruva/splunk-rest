package splunk

import (
	"testing"
)

func TestUser(t *testing.T) {
	var user User
	for _, tc := range userCases {
		user.SetUsername(tc.username)
		user.SetPassword(tc.password)
		if tc.username == user.Username() {
			t.Log("PASS : ", tc.description, tc.username, user.Username())
		}
		if tc.password == user.Password() {
			t.Log("PASS : ", tc.description)
		}
	}
	t.Log("Tested ", len(userCases), " test cases")
}

// TestBasicAuth tests basic authentication types for splunk
func TestBasicAuth(t *testing.T) {
	var user User
	var url URL
	for _, tc := range basicAuthCases {
		want := tc.want
		user.SetUsername(tc.user)
		user.SetPassword(tc.password)
		url.SetHost(tc.host)
		url.SetPort(tc.port)
		err := url.BasicAuth(user)
		if err != nil {
			if want != err.Error() {
				t.Error("FAIL: Error mismatch; got error", err, "but want error", want)
			}
		}
		t.Log(`PASS: `, tc.description)
	}
}

// Test the splunk search
func TestSearch(t *testing.T) {
	var user User
	var url URL
	for _, tc := range searchCases {
		user.SetUsername(tc.user)
		user.SetPassword(tc.password)
		err := url.BasicAuth(user)
		if err != nil {
			t.Error("FAIL : ", err)
		}
		body, err := Search(url, user, tc.search)
		if err != nil {
			t.Error("FAIL : ", err)
		} else if body == nil {
			t.Error("FAIL : Empty body returned")
		}
		t.Log("PASS : ", tc.description)
	}
}
