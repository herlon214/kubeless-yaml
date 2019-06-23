package kubeless

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kubeless/kubeless/pkg/functions"
)

// Handler handles the request
func Handler(event functions.Event, context functions.Context) (string, error) {
	resp, _ := http.Get("https://github.com/")
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return string(body), nil
}
