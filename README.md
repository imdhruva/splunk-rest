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

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter pwd: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	var user splunk.User
	var splunkURL splunk.URL
	user.SetUsername("admin")
	user.SetPassword(text)
	err := splunkURL.BasicAuth(user)
	if err != nil {
		fmt.Printf("ERROR : %s", err)
	}
	sid, err := splunkURL.Search("index=_internal | stats count by splunk_server", user)
	if err != nil {
		fmt.Printf("FAIL : %s", err)
	} else if sid == "" {
		fmt.Printf("FAIL : Empty body returned")
	} else {
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
