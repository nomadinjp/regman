package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List repositories",
	Run: func(cmd *cobra.Command, args []string) {
		reg, err := getRegistry()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing registry: %v\n", err)
			os.Exit(1)
		}

		ctx := context.Background()
		repos, err := remote.Catalog(ctx, reg, getOptions()...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching catalog: %v\n", err)
			os.Exit(1)
		}

		for _, repo := range repos {
			fmt.Println(repo)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
