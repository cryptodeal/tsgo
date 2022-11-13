package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/cryptodeal/tsgo/config"
	"github.com/cryptodeal/tsgo/tsgo"
	"github.com/spf13/cobra"
)

func Execute() {
	var rootCmd = &cobra.Command{
		Use:   "tsgo",
		Short: "Tool for generating Typescript from Go types",
		Long:  `TSGo generates Typescript interfaces and constants from Go files by parsing them.`,
	}

	rootCmd.PersistentFlags().String("config", "tsgo.yaml", "config file to load (default is tsgo.yaml in the current folder)")
	rootCmd.Version = Version() + " " + Target() + " (" + CommitDate() + ") " + Commit()
	rootCmd.PersistentFlags().BoolP("debug", "D", false, "Debug mode (prints debug messages)")

	rootCmd.AddCommand(&cobra.Command{
		Use:   "generate",
		Short: "Generate and write to disk",
		Run:   generate,
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func generate(cmd *cobra.Command, args []string) {
	cfgFilepath, err := cmd.Flags().GetString("config")
	if err != nil {
		log.Fatal(err)
	}
	tsgoConfig := config.ReadFromFilepath(cfgFilepath)
	t := tsgo.New(&tsgoConfig)

	err = t.Generate()
	if err != nil {
		log.Fatalf("TSGo failed: %v", err)
	}
}
