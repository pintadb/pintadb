package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/columbusearch/pintadb/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile     string
	cfg         server.Config
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "pintadb",
		Short: "PintaDB is a lightweight, fast text vector database",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pinta/config.yaml)")
	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "Author name for copyright attribution")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "Name of license for the project (can provide `licensetext` in config)")
	rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
	rootCmd.PersistentFlags().StringVarP(&cfg.FullPath, "fullPath", "", "", "Path to the database file (default is $HOME/.pinta/data/pinta.db)")
	rootCmd.PersistentFlags().Uint64VarP(&cfg.Dimension, "dimension", "", 300, "Dimension of the text vectors")
	rootCmd.PersistentFlags().Uint64VarP(&cfg.HTTPPort, "httpPort", "", 4880, "Port for the HTTP server")
	rootCmd.PersistentFlags().Uint64VarP(&cfg.GRPCPort, "grpcPort", "", 4882, "Port for the GRPC server")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	viper.SetDefault("license", "apache")
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".pinta  " (without extension).
		viper.AddConfigPath(path.Join(home, ".pinta"))
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
