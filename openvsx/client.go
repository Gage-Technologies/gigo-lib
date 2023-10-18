package openvsx

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var (
	ErrExtensionNotFound = fmt.Errorf("extension not found")
)

// Client
//
//	Implements the Client interface for the OpenVSX API
type Client struct {
	// the base URL for the OpenVSX API
	baseUrl string
	// the HTTP client to use for requests
	client *http.Client
}

// NewClient
//
//	Create a new OpenVSX API client
func NewClient(baseUrl string, client *http.Client) *Client {
	if baseUrl == "" {
		baseUrl = "https://open-vsx.org"
	}
	if client == nil {
		client = http.DefaultClient
	}

	return &Client{
		baseUrl: baseUrl,
		client:  client,
	}
}

// prepRequest
//
//	Prepares a request for the OpenVSX API
func (c *Client) newRequest(method string, url string, body io.Reader) (*http.Request, error) {
	// create request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "gigo-openvsx-client; contact@gigo.dev")

	return req, nil
}

// readBodyAsError
//
//	Reads the body of a response as an error
func readBodyAsError(res *http.Response) error {
	// read body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	// return error
	return fmt.Errorf("request failed: %s", string(body))
}

// decodeBody
//
//	Decodes the body of a response into the given interface
func decodeBody(res *http.Response, v interface{}) error {
	// decode body
	err := json.NewDecoder(res.Body).Decode(v)
	if err != nil {
		return fmt.Errorf("failed to decode body: %w", err)
	}

	return nil
}

// GetMetadata
//
//	Get the metadata for an extension
func (c *Client) GetMetadata(extensionId string, version string) (*Extension, error) {
	// split the extension id by . to get the publisher and name
	parts := strings.Split(extensionId, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid extension id: %s", extensionId)
	}

	// default to no version
	url := fmt.Sprintf("%s/api/%s/%s", c.baseUrl, parts[0], parts[1])
	if version != "" {
		url = fmt.Sprintf("%s/api/%s/%s/%s", c.baseUrl, parts[0], parts[1], version)
	}

	// create request
	req, err := c.newRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// send request
	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// check status code
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return nil, ErrExtensionNotFound
		}
		return nil, readBodyAsError(res)
	}

	// decode response
	var metadata Extension
	err = decodeBody(res, &metadata)
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}

// DownloadExtension
//
//	Downloads an extension
func (c *Client) DownloadExtension(extensionId string, version string) (io.ReadCloser, int64, error) {
	// set version to latest if not specified
	if version == "" {
		version = "latest"
	}

	// split the extension id by . to get the publisher and name
	parts := strings.Split(extensionId, ".")
	if len(parts) != 2 {
		return nil, -1, fmt.Errorf("invalid extension id: %s", extensionId)
	}

	// we have to retrieve the metadata if we are working with the latest version
	var url string
	if version == "latest" {
		// retrieve the metadata for the extension
		metadata, err := c.GetMetadata(extensionId, "")
		if err != nil {
			return nil, -1, fmt.Errorf("failed to get extension metadata: %w", err)
		}

		// default to the latest version - create the url if the version is not latest
		var ok bool
		url, ok = metadata.Files["download"]
		if !ok {
			return nil, -1, fmt.Errorf("extension does not have a download url")
		}
	} else {
		url = fmt.Sprintf(
			"%s/api/%s/%s/%s/file/%s.%s-%s.vsix",
			c.baseUrl, parts[0], parts[1], version, parts[0], parts[1], version,
		)
	}

	// create request
	req, err := c.newRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, -1, fmt.Errorf("failed to create request: %w", err)
	}

	// send request
	res, err := c.client.Do(req)
	if err != nil {
		return nil, -1, fmt.Errorf("failed to send request: %w", err)
	}

	// check status code
	if res.StatusCode != http.StatusOK {
		if res.Body != nil {
			_ = res.Body.Close()
		}
		if res.StatusCode == http.StatusNotFound {
			return nil, -1, ErrExtensionNotFound
		}
		return nil, -1, readBodyAsError(res)
	}

	// get the content length header and convert it to an int64
	contentLength := res.Header.Get("Content-Length")
	if contentLength == "" {
		if res.Body != nil {
			_ = res.Body.Close()
		}
		return nil, -1, fmt.Errorf("failed to get content length header")
	}
	contentLengthInt, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		if res.Body != nil {
			_ = res.Body.Close()
		}
		return nil, -1, fmt.Errorf("failed to parse content length header: %w", err)
	}

	return res.Body, contentLengthInt, nil
}
