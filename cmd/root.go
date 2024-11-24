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
	dataAscii []string,
	dataBinary []string,
	form []string,
	user string,
	basic bool,
	digest bool,
	userAgent string,
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
	if data != nil {
		for _, d := range data {
			parts := strings.SplitN(d, "=", 2)
			if parts == nil || len(parts) != 2 {
				return lib.CurlParam{}, fmt.Errorf("invalid data value: %s", d)
			}
			ds = append(ds, lib.KV{parts[0], parts[1]})
		}
	} else if dataAscii != nil {
		for _, d := range dataAscii {
			parts := strings.SplitN(d, "=", 2)
			if parts == nil || len(parts) != 2 {
				return lib.CurlParam{}, fmt.Errorf("invalid data-ascii value: %s", d)
			}
			ds = append(ds, lib.KV{parts[0], parts[1]})
		}
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

	var dbs string
	if dataBinary != nil {
		dbs = strings.Join(dataBinary, "&")
	}

	if userAgent != "" {
		hs = append(hs, lib.KV{"User-Agent", userAgent})
	}

	param := lib.CurlParam{
		URL:        url,
		Method:     request,
		Headers:    hs,
		Data:       ds,
		DataBinary: dbs,
		Form:       fs,
		AuthType:   lib.AuthNone,
	}

	if user != "" {
		parts := strings.SplitN(user, ":", 2)
		if parts == nil || len(parts) != 2 {
			return lib.CurlParam{}, fmt.Errorf("invalid user value: %s", user)
		}
		param.User = parts[0]
		param.Password = parts[1]
		if basic && !digest {
			param.AuthType = lib.AuthBasic
		} else if !basic && digest {
			param.AuthType = lib.AuthDigest
		} else if !basic && !digest {
			param.AuthType = lib.AuthBasic
		} else {
			return lib.CurlParam{}, fmt.Errorf("no auth type")
		}
	}

	return param, nil
}

var rootCmd = &cobra.Command{
	Use:  "curl-to",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		lang := args[0]
		url := args[1]

		param, err := buildCurlParam(url, request, headers, data, dataAscii, dataBinary, form, user, basic, digest, userAgent)
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
var dataAscii []string
var dataBinary []string
var form []string
var user string
var basic bool
var digest bool
var userAgent string

func init() {
	rootCmd.PersistentFlags().StringVarP(&request, "request", "X", "GET", "")
	rootCmd.PersistentFlags().StringArrayVarP(&headers, "header", "H", nil, "")
	rootCmd.PersistentFlags().StringArrayVarP(&data, "data", "d", nil, "")
	rootCmd.PersistentFlags().StringArrayVar(&dataAscii, "data-ascii", nil, "")
	rootCmd.PersistentFlags().StringArrayVar(&dataBinary, "data-binary", nil, "")
	rootCmd.PersistentFlags().StringArrayVarP(&form, "form", "F", nil, "")
	rootCmd.MarkFlagsMutuallyExclusive("data", "data-ascii", "data-binary", "form")
	rootCmd.PersistentFlags().StringVarP(&user, "user", "u", "", "")
	rootCmd.PersistentFlags().BoolVar(&basic, "basic", false, "")
	rootCmd.PersistentFlags().BoolVar(&digest, "digest", false, "")
	rootCmd.MarkFlagsMutuallyExclusive("basic", "digest")
	rootCmd.PersistentFlags().StringVarP(&userAgent, "user-agent", "A", "", "")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
