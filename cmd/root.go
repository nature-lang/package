package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"package/src"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolVarP(&src.Verbose, "verbose", "v", false, "verbose output")
}

const (
	Version = "0.1.0-beta"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hugo",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("nature package %v\n", Version)
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync by package.toml",
	Run: func(cmd *cobra.Command, args []string) {
		src.Sync()
	},
}

var rootCmd = &cobra.Command{}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
