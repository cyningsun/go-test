/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cyningsun/go-test/20220402-coverage-tool/coverage"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:   "xcover",
		Short: "Xcover is a strengthen coverage profile analysis tool",
		Run: func(cmd *cobra.Command, args []string) {
			parser := coverage.NewParser()
			if err := parser.Parse(cover); err != nil {
				log.Fatalf("parse coverage profile failed, err:%v", err)
			}
			switch prefix {
			case "*":
				fmt.Printf("%v:%v\n", filepath.Dir(cover), parser.TotalCov)
			}
		},
	}

	cover  string
	prefix string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.20220402-coverage-tool.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVarP(&cover, "coverprofile", "c", "", "coverage profile path")
	rootCmd.MarkFlagRequired("coverprofile")

	rootCmd.Flags().StringVarP(&prefix, "group", "g", "*", "group by code path")
}
