# splunk-rest
Rest Client implementation for Splunk

## Example

The consumption can be illustrated as under:
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

	// perform authentication splunk
	err := splunkURL.BasicAuth(user)
	if err != nil {
		fmt.Printf("ERROR : %s", err)
	}

	// trigger the search for the search string; this inherently returns the sid
	sid, err := splunkURL.Search(SearchString, user)
	if err != nil {
		fmt.Printf("FAIL : %s", err)
	} else if sid == "" {
		fmt.Printf("FAIL : Empty body returned")
	} else {
		// fetch the result for the sid; this inherently performs the status check for the search
		// search will therefore only be returned if the dispatchStatus of the search is 'DONE'
		// this returns the json formatted/marshalled search output
		body, err := splunkURL.GetSearchResult(sid, user)
		if err != nil {
			fmt.Printf("FAIL : %s", err)
		} else if body == nil {
			fmt.Printf("FAIL : Empty response")
		}
		fmt.Printf("Search OP :" + string(body))
	}
}
```
