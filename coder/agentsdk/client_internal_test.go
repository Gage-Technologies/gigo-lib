package agentsdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
)

const jsonCT = "application/json"

func Test_Client(t *testing.T) {
	t.Parallel()

	const method = http.MethodPost
	const path = "/ok"
	const token = "token"
	const wsId = 420
	const reqBody = `{"msg": "request body"}`
	const resBody = `{"status": "ok"}`

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, method, r.Method)
		assert.Equal(t, path, r.URL.Path)
		assert.Equal(t, token, r.Header.Get(AgentTokenHeader))
		assert.Equal(t, fmt.Sprintf("%d", wsId), r.Header.Get(WorkspaceIDHeader))
		for k, v := range r.Header {
			t.Logf("header %q: %q", k, strings.Join(v, ", "))
		}

		w.Header().Set("Content-Type", jsonCT)
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, resBody)
	}))

	u, err := url.Parse(s.URL)
	require.NoError(t, err)
	client := New(u)
	client.SetSessionAuth(wsId, token)

	logBuf := bytes.NewBuffer(nil)
	client.Logger = slog.Make(sloghuman.Sink(logBuf)).Leveled(slog.LevelDebug)
	client.LogBodies = true

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resp, err := client.Request(ctx, method, path, []byte(reqBody))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, jsonCT, resp.Header.Get("Content-Type"))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, resBody, string(body))

	// logStr := logBuf.String()
	// require.Contains(t, logStr, "agentsdk request")
	// require.Contains(t, logStr, method)
	// require.Contains(t, logStr, path)
	// require.Contains(t, logStr, strings.ReplaceAll(reqBody, `"`, `\"`))
	// require.Contains(t, logStr, "agentsdk response")
	// require.Contains(t, logStr, "200")
	// require.Contains(t, logStr, strings.ReplaceAll(resBody, `"`, `\"`))
}

func Test_readBodyAsError(t *testing.T) {
	t.Parallel()

	exampleURL := "http://example.com"
	simpleResponse := Response{
		Message: "test",
		Detail:  "hi",
	}

	longResponse := ""
	for i := 0; i < 2000; i++ {
		longResponse += "a"
	}

	unexpectedJSON := marshal(map[string]any{
		"hello": "world",
		"foo":   "bar",
	})

	//nolint:bodyclose
	tests := []struct {
		name   string
		req    *http.Request
		res    *http.Response
		assert func(t *testing.T, err error)
	}{
		{
			name: "JSONWithRequest",
			req:  httptest.NewRequest(http.MethodGet, exampleURL, nil),
			res:  newResponse(http.StatusNotFound, jsonCT, marshal(simpleResponse)),
			assert: func(t *testing.T, err error) {
				sdkErr := assertSDKError(t, err)

				assert.Equal(t, simpleResponse, sdkErr.Response)
				assert.ErrorContains(t, err, sdkErr.Response.Message)
				assert.ErrorContains(t, err, sdkErr.Response.Detail)

				assert.Equal(t, http.StatusNotFound, sdkErr.StatusCode())
				assert.ErrorContains(t, err, strconv.Itoa(sdkErr.StatusCode()))

				assert.Equal(t, http.MethodGet, sdkErr.method)
				assert.ErrorContains(t, err, sdkErr.method)

				assert.Equal(t, exampleURL, sdkErr.url)
				assert.ErrorContains(t, err, sdkErr.url)

				assert.Empty(t, sdkErr.Helper)
			},
		},
		{
			name: "JSONWithoutRequest",
			req:  nil,
			res:  newResponse(http.StatusNotFound, jsonCT, marshal(simpleResponse)),
			assert: func(t *testing.T, err error) {
				sdkErr := assertSDKError(t, err)

				assert.Equal(t, simpleResponse, sdkErr.Response)
				assert.Equal(t, http.StatusNotFound, sdkErr.StatusCode())
				assert.Empty(t, sdkErr.method)
				assert.Empty(t, sdkErr.url)
				assert.Empty(t, sdkErr.Helper)
			},
		},
		{
			name: "UnauthorizedHelper",
			req:  nil,
			res:  newResponse(http.StatusUnauthorized, jsonCT, marshal(simpleResponse)),
			assert: func(t *testing.T, err error) {
				sdkErr := assertSDKError(t, err)

				assert.Contains(t, sdkErr.Helper, "Try logging in")
				assert.ErrorContains(t, err, sdkErr.Helper)
			},
		},
		{
			name: "NonJSON",
			req:  nil,
			res:  newResponse(http.StatusNotFound, "text/plain; charset=utf-8", "hello world"),
			assert: func(t *testing.T, err error) {
				sdkErr := assertSDKError(t, err)

				assert.Contains(t, sdkErr.Response.Message, "unexpected non-JSON response")
				assert.Equal(t, "hello world", sdkErr.Response.Detail)
			},
		},
		{
			name: "NonJSONLong",
			req:  nil,
			res:  newResponse(http.StatusNotFound, "text/plain; charset=utf-8", longResponse),
			assert: func(t *testing.T, err error) {
				sdkErr := assertSDKError(t, err)

				assert.Contains(t, sdkErr.Response.Message, "unexpected non-JSON response")

				expected := longResponse[0:1024] + "..."
				assert.Equal(t, expected, sdkErr.Response.Detail)
			},
		},
		{
			name: "JSONNoBody",
			req:  nil,
			res:  newResponse(http.StatusNotFound, jsonCT, ""),
			assert: func(t *testing.T, err error) {
				sdkErr := assertSDKError(t, err)

				assert.Contains(t, sdkErr.Response.Message, "empty response body")
			},
		},
		{
			name: "JSONNoMessage",
			req:  nil,
			res:  newResponse(http.StatusNotFound, jsonCT, unexpectedJSON),
			assert: func(t *testing.T, err error) {
				sdkErr := assertSDKError(t, err)

				assert.Contains(t, sdkErr.Response.Message, "unexpected status code")
				assert.Contains(t, sdkErr.Response.Message, "has no message")
				assert.Equal(t, unexpectedJSON, sdkErr.Response.Detail)
			},
		},
	}

	for _, c := range tests {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			c.res.Request = c.req

			err := readBodyAsError(c.res)
			c.assert(t, err)
		})
	}
}

func assertSDKError(t *testing.T, err error) *Error {
	t.Helper()

	var sdkErr *Error
	require.Error(t, err)
	require.True(t, xerrors.As(err, &sdkErr))

	return sdkErr
}

func newResponse(status int, contentType string, body interface{}) *http.Response {
	var r io.ReadCloser
	switch v := body.(type) {
	case string:
		r = io.NopCloser(strings.NewReader(v))
	case []byte:
		r = io.NopCloser(bytes.NewReader(v))
	case io.ReadCloser:
		r = v
	case io.Reader:
		r = io.NopCloser(v)
	default:
		panic(fmt.Sprintf("unknown body type: %T", body))
	}

	return &http.Response{
		Status:     http.StatusText(status),
		StatusCode: status,
		Header: http.Header{
			"Content-Type": []string{contentType},
		},
		Body: r,
	}
}

func marshal(res any) string {
	b, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	return string(b)
}
