[![Build Status](https://api.cirrus-ci.com/github/imdhruva/splunk.svg)](https://cirrus-ci.com/github/imdhruva/splunk)

# splunk-rest
Rest Client implementation for Splunk

## Example

The invocation can be illustrated as under:
```golang
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/imdhruva/splunk"
)

const (
	Username := "admin"
	SearchString := "index=_internal | stats count by splunk_server"
)

func main() {
	// read user credentials
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter pwd: ")
	text, _ := reader.ReadString('\n')
	pwd = strings.TrimSpace(text)

	// initialize user credentials
	var user splunk.User
	var splunkURL splunk.URL
	user.SetUsername(Username)
	user.SetPassword(pwd)

	// perform search operation and return the body this encapsulates the operations
	// 1. for performing basic authentication
	// 2. trigerring search and returning a search-id
	// 3. returing the body of the search
	body, err := Search(url, user, SearchString)
	if err != nil {
		t.Error("FAIL : ", err)
	} else if body == nil {
		t.Error("FAIL : Empty body returned")
	}
	fmt.Printf("Search OP :" + string(body))
}
```
