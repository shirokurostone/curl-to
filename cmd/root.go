package cmd

import (
	"fmt"
	"github.com/shirokurostone/curl-to/lib"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:  "curl-to",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		lang := args[0]
		url := args[1]

		var hs [][2]string
		for _, h := range headers {
			parts := strings.SplitN(h, ":", 2)
			if parts == nil || len(parts) != 2 {
				return fmt.Errorf("invalid header value: %s", h)
			}
			hs = append(hs, [2]string{strings.TrimSuffix(parts[0], " "), strings.TrimPrefix(parts[1], " ")})
		}

		var ds [][2]string
		for _, d := range data {
			parts := strings.SplitN(d, "=", 2)
			if parts == nil || len(parts) != 2 {
				return fmt.Errorf("invalid data value: %s", d)
			}
			ds = append(ds, [2]string{parts[0], parts[1]})
		}

		param := lib.CurlParam{
			URL:     url,
			Method:  request,
			Headers: hs,
			Data:    ds,
		}

		var code string
		var err error
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

func init() {
	rootCmd.PersistentFlags().StringVarP(&request, "request", "X", "GET", "")
	rootCmd.PersistentFlags().StringArrayVarP(&headers, "header", "H", nil, "")
	rootCmd.PersistentFlags().StringArrayVarP(&data, "data", "d", nil, "")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
