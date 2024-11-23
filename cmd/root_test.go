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
			lib.CurlParam{},
			fmt.Errorf("invalid form value: key=value;headers=invalid"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actual, actualErr := buildCurlParam(tt.url, tt.request, tt.headers, tt.data, tt.form)
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.expectedErr, actualErr)
		})
	}
}
