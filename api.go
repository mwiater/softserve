// api.go
package softserve

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type MockResponse struct {
	Status  int               `yaml:"status"`
	Headers map[string]string `yaml:"headers"`
	Body    string            `yaml:"body"`
}

var apiRoutes map[string]MockResponse

func LoadAPIResponses() error {
	apiRoutes = make(map[string]MockResponse)

	data, err := os.ReadFile("api.yaml")
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("api.yaml not found")
		}
		return err
	}

	return yaml.Unmarshal(data, &apiRoutes)
}

func HandleAPIRequest(w http.ResponseWriter, r *http.Request) bool {
	if !GetConfig().API {
		return false
	}

	prefix := GetConfig().APIPrefix
	if prefix == "" {
		prefix = "/api/"
	}
	if !strings.HasPrefix(r.URL.Path, prefix) {
		return false
	}

	key := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
	mock, ok := apiRoutes[key]
	if !ok {
		http.NotFound(w, r)
		return true
	}

	for k, v := range mock.Headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(mock.Status)
	w.Write([]byte(mock.Body))
	return true
}
