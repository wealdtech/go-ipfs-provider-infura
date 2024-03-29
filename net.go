package infura

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

// get is a helper that runs an HTTP GET with appropriate headers and returns
// the body as a generic JSON interface.
func (p *Provider) get(url string, contentType string) (map[string]interface{}, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := make(map[string]interface{})
	return res, handleBody(data, &res)
}

// post is a helper that runs an HTTP POST with appropriate headers and returns
// the resultant IPFS hash
func (p *Provider) post(url string, contentType string, content *bytes.Buffer) (map[string]interface{}, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, content)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		return nil, errors.New(resp.Status)
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := make(map[string]interface{})
	return res, handleBody(data, &res)
}

// Handle the body returned by a GET or POST, including various error conditions.
func handleBody(data []byte, res *map[string]interface{}) error {
	if len(data) == 0 {
		// Nothing back is valid
		return nil
	}
	data = bytes.TrimSpace(data)

	// If there is more than 1 line then this is newline-delimited
	if bytes.Contains(data, []byte{'\n'}) {
		results := make([]*map[string]interface{}, 0)
		for _, line := range bytes.Split(data, []byte{'\n'}) {
			lineRes := make(map[string]interface{})
			err := json.Unmarshal(line, &lineRes)
			if err != nil {
				if line[0] != '{' {
					// Not JSON, return the plain text as the message
					(*res)["message"] = strings.TrimSpace(string(data))
					return nil
				}
				// Failed to unmarshal; return the raw value as an error
				return errors.New(string(data))
			}
			results = append(results, &lineRes)
		}
		(*res)["results"] = results
		return nil
	}

	// Single result
	err := json.Unmarshal(data, &res)
	if err != nil {
		if data[0] != '{' {
			// Not JSON, return the plain text as the message
			(*res)["message"] = strings.TrimSpace(string(data))
			return nil
		}
		// Failed to unmarshal; return the raw value as an error
		return errors.New(string(data))
	}

	if msg, exists := (*res)["Type"]; exists {
		if msg.(string) == "error" {
			return errors.New((*res)["Message"].(string))
		}
	}
	return nil
}
