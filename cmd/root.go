package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:  "curl-to",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		lang := args[0]
		url := args[1]
		fmt.Println(lang, url, request)
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
