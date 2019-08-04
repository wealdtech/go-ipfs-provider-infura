package infura

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	ma "github.com/multiformats/go-multiaddr"
	mh "github.com/multiformats/go-multihash"
	provider "github.com/wealdtech/go-ipfs-provider"
)

const (
	baseURL    = "https://ipfs.infura.io:5001/api/v0"
	gatewayURL = "https://ipfs.infura.io"
)

// Ping is an internal method to ensure the endpoint is accessible.
// This returns true if the endpoint is accessible, otherwise false.
// This returns an error for a network or authentication problem.
func (p *Provider) Ping() (bool, error) {
	res, err := p.get(fmt.Sprintf("%s/version", baseURL), "")
	if err != nil {
		return false, err
	}

	if msg, exists := res["Version"]; exists {
		if len(msg.(string)) > 0 {
			return true, nil
		}
		return false, nil
	}

	return false, errors.New("unexpected failure")
}

// List lists all content pinned to this provider.
func (p *Provider) List() ([]*provider.ItemStatistics, error) {
	return nil, errors.New("not supported by this provider")
}

// ItemStats returns information on an IPFS hash pinned to this provider.
// Note that this does not return the file name.
func (p *Provider) ItemStats(hash string) (*provider.ItemStatistics, error) {
	res, err := p.get(fmt.Sprintf("%s/object/stat?arg=%s", baseURL, hash), "")
	if err != nil {
		return nil, err
	}

	item := &provider.ItemStatistics{
		Hash: hash,
	}

	if c, exists := res["CumulativeSize"]; exists {
		item.Size = uint64(c.(float64))
	}

	return item, nil
}

// ServiceStats provides statistics for this provider.
func (p *Provider) ServiceStats() (*provider.SiteStatistics, error) {
	return nil, errors.New("not supported by this provider")
}

// PinContent pins content to this provider.
func (p *Provider) PinContent(name string, content io.Reader, opts *provider.ContentOpts) (string, error) {
	var b bytes.Buffer
	var contentType string

	if name != "" && content != nil {
		// Add content

		// Defer closing the content if it is closeable
		if x, ok := content.(io.Closer); ok {
			defer x.Close()
		}

		// Set up the form field
		w := multipart.NewWriter(&b)
		var fw io.Writer
		fw, err := w.CreateFormFile("file", name)
		if err != nil {
			return "", err
		}
		io.Copy(fw, content)
		w.Close()
		contentType = w.FormDataContentType()
	}

	url := fmt.Sprintf("%s/add?pin=true", baseURL)
	if opts != nil && opts.StoreInDirectory {
		url = url + "&wrap-with-directory=true"
	}
	res, err := p.post(url, contentType, &b)
	if err != nil {
		return "", err
	}

	// If we pinned multiple files return the one without a name
	results, exists := res["results"]
	if exists {
		for _, result := range results.([]*map[string]interface{}) {
			name, exists := (*result)["Name"]
			if exists && name == "" {
				return (*result)["Hash"].(string), nil
			}
		}
		return "", errors.New("no hash returned")
	}

	msg, exists := res["Hash"]
	if exists {
		return msg.(string), nil
	}
	return "", errors.New("no hash returned")
}

// Pin pins existing IPFS content to this provider.
func (p *Provider) Pin(hash string) error {
	_, err := p.get(fmt.Sprintf("%s/pin/add?arg=%s", baseURL, hash), "")
	return err
}

// Unpin removes content from this provider.
func (p *Provider) Unpin(hash string) error {
	return errors.New("not supported by this provider")
}

// GatewayURL provides a gateway URL for the given input.
// The input may be an existing gateway URL, a multiaddr, or
// a plain hash
func (p *Provider) GatewayURL(input string) (string, error) {
	// Multiaddr
	_, err := ma.NewMultiaddr(input)
	if err == nil {
		return fmt.Sprintf("%s%s", gatewayURL, input), nil
	}

	// URI
	if strings.HasPrefix(input, "ipfs://") {
		return fmt.Sprintf("%s/ipfs/%s", gatewayURL, input[7:]), nil
	}
	if strings.HasPrefix(input, "ipns://") {
		return fmt.Sprintf("%s/ipns/%s", gatewayURL, input[7:]), nil
	}

	// Existing gateway URL
	index := strings.Index(input, "/ipfs/")
	if index == -1 {
		index = strings.Index(input, "/ipns/")
	}
	if index != -1 {
		return fmt.Sprintf("%s%s", gatewayURL, input[index:]), nil
	}

	// Plain hash
	_, err = mh.FromB58String(input)
	if err == nil {
		return fmt.Sprintf("%s/ipfs/%s", gatewayURL, input), nil
	}

	return "", errors.New("unrecognised format")
}
