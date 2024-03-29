/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/

package cmd

import (
	"github.com/spf13/cobra"
	"go.olapie.com/tools/gocode/domain"
	"go.olapie.com/utils"
	"strings"
)

// domainCmd represents the domain command
var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Generate domain model code",
	Long:  `Generate domain model code with XML model. E.g.\n` + domain.ExampleXML,
	Run: func(cmd *cobra.Command, args []string) {
		inputFilename := utils.MustGet(cmd.Flags().GetString("filename"))
		outputFilename := utils.MustGet(cmd.Flags().GetString("output"))
		if strings.HasSuffix(inputFilename, ".xml") {
			domain.Generate(inputFilename, outputFilename)
		} else {
			domain.GenerateBatch(inputFilename, outputFilename)
		}
	},
}

func init() {
	rootCmd.AddCommand(domainCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// domainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	domainCmd.Flags().StringP("filename", "f", "domain.xml", "Domain model XML filename")
	domainCmd.Flags().StringP("output", "o", "domain.gen.go", "Generated filename")
}
