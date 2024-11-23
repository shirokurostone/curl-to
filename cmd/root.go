package cmd

import (
	"fmt"
	"github.com/shirokurostone/curl-to/lib"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:  "curl-to",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		lang := args[0]
		url := args[1]
		param := lib.CurlParam{
			URL:    url,
			Method: request,
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

func init() {
	rootCmd.PersistentFlags().StringVarP(&request, "request", "X", "GET", "")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
