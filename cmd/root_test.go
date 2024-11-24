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
		form        []string
		user        string
		basic       bool
		digest      bool
		userAgent   string
		expected    lib.CurlParam
		expectedErr error
	}{
		{
			"minimum",
			"http://localhost/",
			"GET",
			nil,
			nil,
			nil,
			"",
			true,
			false,
			"",
			lib.CurlParam{
				URL:     "http://localhost/",
				Method:  "GET",
				Headers: nil,
				Data:    nil,
				Form:    nil,
			},
			nil,
		},
		{
			"headers: valid",
			"http://localhost/",
			"GET",
			[]string{"Content-Type: application/json"},
			nil,
			nil,
			"",
			true,
			false,
			"",
			lib.CurlParam{
				URL:     "http://localhost/",
				Method:  "GET",
				Headers: []lib.KV{{"Content-Type", "application/json"}},
				Data:    nil,
				Form:    nil,
			},
			nil,
		},
		{
			"headers: invalid value",
			"http://localhost/",
			"GET",
			[]string{"invalid"},
			nil,
			nil,
			"",
			true,
			false,
			"",
			lib.CurlParam{},
			fmt.Errorf("invalid header value: invalid"),
		},
		{
			"data: valid",
			"http://localhost/",
			"GET",
			nil,
			[]string{"key=value"},
			nil,
			"",
			true,
			false,
			"",
			lib.CurlParam{
				URL:     "http://localhost/",
				Method:  "GET",
				Headers: nil,
				Data:    []lib.KV{{"key", "value"}},
				Form:    nil,
			},
			nil,
		},
		{
			"data: invalid value",
			"http://localhost/",
			"GET",
			nil,
			[]string{"invalid"},
			nil,
			"",
			true,
			false,
			"",
			lib.CurlParam{},
			fmt.Errorf("invalid data value: invalid"),
		},
		{
			"form: valid",
			"http://localhost/",
			"GET",
			nil,
			nil,
			[]string{"key=value;type=application/json;filename=sample.json;headers=X-Header: Value"},
			"",
			true,
			false,
			"",
			lib.CurlParam{
				URL:     "http://localhost/",
				Method:  "GET",
				Headers: nil,
				Data:    nil,
				Form: []lib.Form{{
					"key",
					"value",
					"application/json",
					"sample.json",
					[]lib.KV{{"X-Header", "Value"}},
				}},
			},
			nil,
		},
		{
			"form: invalid value",
			"http://localhost/",
			"GET",
			nil,
			nil,
			[]string{"invalid"},
			"",
			true,
			false,
			"",
			lib.CurlParam{},
			fmt.Errorf("invalid form value: invalid"),
		},
		{
			"form: invalid headers",
			"http://localhost/",
			"GET",
			nil,
			nil,
			[]string{"key=value;headers=invalid"},
			"",
			true,
			false,
			"",
			lib.CurlParam{},
			fmt.Errorf("invalid form value: key=value;headers=invalid"),
		},
		{
			"user: valid: basic",
			"http://localhost/",
			"GET",
			nil,
			nil,
			nil,
			"user:password",
			false,
			false,
			"",
			lib.CurlParam{
				URL:      "http://localhost/",
				Method:   "GET",
				Headers:  nil,
				Data:     nil,
				Form:     nil,
				AuthType: lib.AuthBasic,
				User:     "user",
				Password: "password",
			},
			nil,
		},
		{
			"user: valid: digest",
			"http://localhost/",
			"GET",
			nil,
			nil,
			nil,
			"user:password",
			false,
			true,
			"",
			lib.CurlParam{
				URL:      "http://localhost/",
				Method:   "GET",
				Headers:  nil,
				Data:     nil,
				Form:     nil,
				AuthType: lib.AuthDigest,
				User:     "user",
				Password: "password",
			},
			nil,
		},
		{
			"user: invalid: user",
			"http://localhost/",
			"GET",
			nil,
			nil,
			nil,
			"user",
			true,
			false,
			"",
			lib.CurlParam{},
			fmt.Errorf("invalid user value: user"),
		},
		{
			"user: invalid: no auth type",
			"http://localhost/",
			"GET",
			nil,
			nil,
			nil,
			"user:password",
			true,
			true,
			"",
			lib.CurlParam{},
			fmt.Errorf("no auth type"),
		},
		{
			"userAgent: valid",
			"http://localhost/",
			"GET",
			nil,
			nil,
			nil,
			"",
			true,
			true,
			"curl-to",
			lib.CurlParam{
				URL:     "http://localhost/",
				Method:  "GET",
				Headers: []lib.KV{{"User-Agent", "curl-to"}},
				Data:    nil,
				Form:    nil,
			},
			nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actual, actualErr := buildCurlParam(tt.url, tt.request, tt.headers, tt.data, tt.form, tt.user, tt.basic, tt.digest, tt.userAgent)
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.expectedErr, actualErr)
		})
	}
}
