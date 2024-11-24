package cmd

import (
	"fmt"
	"github.com/shirokurostone/curl-to/lib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildCurlParam(t *testing.T) {

	testCases := []struct {
		name        string
		url         string
		request     string
		headers     []string
		data        []string
		dataAscii   []string
		dataBinary  []string
		form        []string
		user        string
		basic       bool
		digest      bool
		userAgent   string
		expected    lib.CurlParam
		expectedErr error
	}{
		{
			name:    "minimum",
			url:     "http://localhost/",
			request: "GET",
			expected: lib.CurlParam{
				URL:    "http://localhost/",
				Method: "GET",
			},
		},
		{
			name:    "headers: valid",
			url:     "http://localhost/",
			request: "GET",
			headers: []string{"Content-Type: application/json"},
			expected: lib.CurlParam{
				URL:     "http://localhost/",
				Method:  "GET",
				Headers: []lib.KV{{"Content-Type", "application/json"}},
			},
		},
		{
			name:        "headers: invalid value",
			url:         "http://localhost/",
			request:     "GET",
			headers:     []string{"invalid"},
			expectedErr: fmt.Errorf("invalid header value: invalid"),
		},
		{
			name:    "data: valid",
			url:     "http://localhost/",
			request: "GET",
			data:    []string{"key=value"},
			expected: lib.CurlParam{
				URL:    "http://localhost/",
				Method: "GET",
				Data:   []lib.KV{{"key", "value"}},
			},
		},
		{
			name:        "data: invalid value",
			url:         "http://localhost/",
			request:     "GET",
			data:        []string{"invalid"},
			expectedErr: fmt.Errorf("invalid data value: invalid"),
		},
		{
			name:      "dataAscii: valid",
			url:       "http://localhost/",
			request:   "GET",
			dataAscii: []string{"key=value"},
			expected: lib.CurlParam{
				URL:    "http://localhost/",
				Method: "GET",
				Data:   []lib.KV{{"key", "value"}},
			},
		},
		{
			name:        "dataAscii: invalid value",
			url:         "http://localhost/",
			request:     "GET",
			dataAscii:   []string{"invalid"},
			expectedErr: fmt.Errorf("invalid data-ascii value: invalid"),
		},
		{
			name:       "dataBinary: valid",
			url:        "http://localhost/",
			request:    "GET",
			dataBinary: []string{"value1", "value2"},
			expected: lib.CurlParam{
				URL:        "http://localhost/",
				Method:     "GET",
				DataBinary: "value1&value2",
			},
		},
		{
			name:    "form: valid",
			url:     "http://localhost/",
			request: "GET",
			form:    []string{"key=value;type=application/json;filename=sample.json;headers=X-Header: Value"},
			expected: lib.CurlParam{
				URL:    "http://localhost/",
				Method: "GET",
				Form: []lib.Form{{
					"key",
					"value",
					"application/json",
					"sample.json",
					[]lib.KV{{"X-Header", "Value"}},
				}},
			},
		},
		{
			name:        "form: invalid value",
			url:         "http://localhost/",
			request:     "GET",
			form:        []string{"invalid"},
			expectedErr: fmt.Errorf("invalid form value: invalid"),
		},
		{
			name:        "form: invalid headers",
			url:         "http://localhost/",
			request:     "GET",
			form:        []string{"key=value;headers=invalid"},
			expectedErr: fmt.Errorf("invalid form value: key=value;headers=invalid"),
		},
		{
			name:    "user: valid: basic",
			url:     "http://localhost/",
			request: "GET",
			user:    "user:password",
			basic:   false,
			digest:  false,
			expected: lib.CurlParam{
				URL:      "http://localhost/",
				Method:   "GET",
				AuthType: lib.AuthBasic,
				User:     "user",
				Password: "password",
			},
		},
		{
			name:    "user: valid: digest",
			url:     "http://localhost/",
			request: "GET",
			user:    "user:password",
			basic:   false,
			digest:  true,
			expected: lib.CurlParam{
				URL:      "http://localhost/",
				Method:   "GET",
				AuthType: lib.AuthDigest,
				User:     "user",
				Password: "password",
			},
		},
		{
			name:        "user: invalid: user",
			url:         "http://localhost/",
			request:     "GET",
			user:        "user",
			expectedErr: fmt.Errorf("invalid user value: user"),
		},
		{
			name:        "user: invalid: no auth type",
			url:         "http://localhost/",
			request:     "GET",
			user:        "user:password",
			basic:       true,
			digest:      true,
			expectedErr: fmt.Errorf("no auth type"),
		},
		{
			name:      "userAgent: valid",
			url:       "http://localhost/",
			request:   "GET",
			userAgent: "curl-to",
			expected: lib.CurlParam{
				URL:     "http://localhost/",
				Method:  "GET",
				Headers: []lib.KV{{"User-Agent", "curl-to"}},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actual, actualErr := buildCurlParam(tt.url, tt.request, tt.headers, tt.data, tt.dataAscii, tt.dataBinary, tt.form, tt.user, tt.basic, tt.digest, tt.userAgent)
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.expectedErr, actualErr)
		})
	}
}
