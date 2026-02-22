package cmd

import (
	"fmt"
	"os"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tagsCmd = &cobra.Command{
	Args:  cobra.ExactArgs(1),
	Use:   "tags <image_name>",
	Short: "List tags for a specific image",
	Run: func(cmd *cobra.Command, args []string) {
		repoName := args[0]
		
		regURL := viper.GetString("registry")
		if regURL == "" {
			fmt.Fprintf(os.Stderr, "Error: registry URL is required\n")
			os.Exit(1)
		}

		opts := []name.Option{}
		if viper.GetBool("insecure") {
			opts = append(opts, name.Insecure)
		}

		repo, err := name.NewRepository(regURL + "/" + repoName, opts...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing repository: %v\n", err)
			os.Exit(1)
		}

		tags, err := remote.List(repo, getOptions()...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing tags: %v\n", err)
			os.Exit(1)
		}

		for _, tag := range tags {
			fmt.Println(tag)
		}
	},
}

func init() {
	rootCmd.AddCommand(tagsCmd)
}
