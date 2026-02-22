package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "regman",
	Short: "A simple Docker Registry manager",
	Long:  "A tool for managing private Docker registries. Configurable via flags, ENV, or config file (~/.regman.yaml).",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.regman.yaml)")
	rootCmd.PersistentFlags().String("registry", "", "Registry URL (e.g., https://my-registry.com)")
	rootCmd.PersistentFlags().String("user", "", "Registry username")
	rootCmd.PersistentFlags().String("pass", "", "Registry password")
	rootCmd.PersistentFlags().Bool("insecure", false, "Allow HTTP or skip TLS verification")

	// Bind flags to viper
	viper.BindPFlag("registry", rootCmd.PersistentFlags().Lookup("registry"))
	viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("pass", rootCmd.PersistentFlags().Lookup("pass"))
	viper.BindPFlag("insecure", rootCmd.PersistentFlags().Lookup("insecure"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".regman")
	}

	viper.SetEnvPrefix("REGMAN")
	viper.AutomaticEnv() // Read REGMAN_REGISTRY, etc.

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// Silently ignore if config file not found, but report other errors
	}
}

func getOptions() []remote.Option {
	var options []remote.Option

	// Authentication priority: 
	// 1. Explicit Flags/Viper
	// 2. Default Keychain (Docker login)
	user := viper.GetString("user")
	pass := viper.GetString("pass")

	if user != "" || pass != "" {
		auth := &authn.Basic{
			Username: user,
			Password: pass,
		}
		options = append(options, remote.WithAuth(auth))
	} else {
		// Fallback to Docker's credentials (~/.docker/config.json)
		options = append(options, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	}

	// Insecure Transport
	if viper.GetBool("insecure") {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		options = append(options, remote.WithTransport(tr))
	}

	return options
}

func getRegistry() (name.Registry, error) {
	regURL := viper.GetString("registry")
	if regURL == "" {
		return name.Registry{}, fmt.Errorf("registry URL is required (via --registry, REGMAN_REGISTRY env, or config file)")
	}

	// Clean protocol if present (name.NewRegistry expects domain)
	regURL = strings.TrimPrefix(regURL, "https://")
	regURL = strings.TrimPrefix(regURL, "http://")

	opts := []name.Option{}
	if viper.GetBool("insecure") {
		opts = append(opts, name.Insecure)
	}
	return name.NewRegistry(regURL, opts...)
}
