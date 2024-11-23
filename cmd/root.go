package cmd

import (
	"fmt"
	"github.com/shirokurostone/curl-to/lib"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func buildCurlParam(
	url string,
	request string,
	headers []string,
	data []string,
	form []string,
) (lib.CurlParam, error) {
	var hs []lib.KV
	for _, h := range headers {
		parts := strings.SplitN(h, ":", 2)
		if parts == nil || len(parts) != 2 {
			return lib.CurlParam{}, fmt.Errorf("invalid header value: %s", h)
		}
		hs = append(hs, lib.KV{strings.TrimSuffix(parts[0], " "), strings.TrimPrefix(parts[1], " ")})
	}

	var ds []lib.KV
	for _, d := range data {
		parts := strings.SplitN(d, "=", 2)
		if parts == nil || len(parts) != 2 {
			return lib.CurlParam{}, fmt.Errorf("invalid data value: %s", d)
		}
		ds = append(ds, lib.KV{parts[0], parts[1]})
	}

	var fs []lib.Form
	for _, f := range form {
		parts := strings.Split(f, ";")
		form := lib.Form{}
		for i, p := range parts {
			kv := strings.SplitN(p, "=", 2)
			if kv == nil || len(kv) != 2 {
				return lib.CurlParam{}, fmt.Errorf("invalid form value: %s", f)
			}
			if i == 0 {
				form.Name = kv[0]
				form.Value = kv[1]
			} else {
				switch kv[0] {
				case "type":
					form.TypeValue = kv[1]
				case "filename":
					form.Filename = kv[1]
				case "headers":
					headerParts := strings.SplitN(kv[1], ":", 2)
					if headerParts == nil || len(headerParts) != 2 {
						return lib.CurlParam{}, fmt.Errorf("invalid form value: %s", f)
					}
					form.Headers = append(form.Headers, lib.KV{strings.TrimSuffix(headerParts[0], " "), strings.TrimPrefix(headerParts[1], " ")})
				}
			}
		}
		fs = append(fs, form)
	}

	param := lib.CurlParam{
		URL:     url,
		Method:  request,
		Headers: hs,
		Data:    ds,
		Form:    fs,
	}

	return param, nil
}

var rootCmd = &cobra.Command{
	Use:  "curl-to",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		lang := args[0]
		url := args[1]

		param, err := buildCurlParam(url, request, headers, data, form)
		if err != nil {
			return err
		}

		var code string
		switch lang {
		case "ruby":
			code, err = lib.GenerateRubyCode(param)
		default:
			err = fmt.Errorf("unsupported language: %s", lang)
		}
		if err != nil {
			return err
		}

		fmt.Println(code)
		return nil
	},
}

var request string
var headers []string
var data []string
var form []string

func init() {
	rootCmd.PersistentFlags().StringVarP(&request, "request", "X", "GET", "")
	rootCmd.PersistentFlags().StringArrayVarP(&headers, "header", "H", nil, "")
	rootCmd.PersistentFlags().StringArrayVarP(&data, "data", "d", nil, "")
	rootCmd.PersistentFlags().StringArrayVarP(&form, "form", "F", nil, "")
	rootCmd.MarkFlagsMutuallyExclusive("data", "form")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
