/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"

	"github.com/cyningsun/go-test/20220402-coverage-tool/coverage"
	"github.com/spf13/cobra"
)

// checkIgnoreCmd represents the checkIgnore command
var (
	checkIgnoreCmd = &cobra.Command{
		Use:   "check-ignore <pathname>",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			files, err := coverage.IgnoreFiles(args)
			if err != nil {
				log.Fatalf("read path failed, err:%v", err)
			}
			for _, each := range files {
				if verbose {
					log.Printf("%v:%v", each.Pattern, each.Path)
				} else {
					log.Printf("%v", each.Path)
				}
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(checkIgnoreCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkIgnoreCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkIgnoreCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	checkIgnoreCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, `Instead of printing the paths that are excluded, for each path that matches an exclude pattern, print the exclude pattern together with the path. (Matching an exclude pattern usually means the path is excluded, but if the pattern begins with ! then it is a negated pattern and matching it means the path is NOT excluded.)`)
}
