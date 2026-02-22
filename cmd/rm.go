package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rmCmd = &cobra.Command{
	Args:  cobra.ExactArgs(1),
	Use:   "rm <image_name>[:tag]",
	Short: "Delete a specific tag or image by digest",
	Run: func(cmd *cobra.Command, args []string) {
		imageName := args[0]
		
		regURL := viper.GetString("registry")
		if regURL == "" {
			fmt.Fprintf(os.Stderr, "Error: registry URL is required\n")
			os.Exit(1)
		}

		// If tag is missing, assume latest
		fullImageName := regURL + "/" + imageName
		if !strings.Contains(imageName, ":") && !strings.Contains(imageName, "@") {
			fullImageName += ":latest"
		}

		opts := []name.Option{}
		if viper.GetBool("insecure") {
			opts = append(opts, name.Insecure)
		}

		ref, err := name.ParseReference(fullImageName, opts...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing reference: %v\n", err)
			os.Exit(1)
		}

		// Get Digest
		img, err := remote.Head(ref, getOptions()...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting digest: %v\n", err)
			os.Exit(1)
		}
		
		digest := img.Digest.String()
		fmt.Printf("Found digest: %s for %s\n", digest, ref.String())

		// To Delete, we MUST use digest reference
		repo := ref.Context()
		deleteRef, err := name.NewDigest(repo.String() + "@" + digest)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating digest reference: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Deleting %s...\n", deleteRef.String())
		if err := remote.Delete(deleteRef, getOptions()...); err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting image: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Successfully deleted.")
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
