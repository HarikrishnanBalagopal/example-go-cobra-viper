package main

import (
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Replicas    int
	VersionLong bool `mapstructure:"long"`
}

var config = Config{}

func rootHandler() {
	fmt.Println("starting", config.Replicas, "replicas")
}

func getRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:  "myapp",
		Args: cobra.NoArgs,
		Run: func(*cobra.Command, []string) {
			rootHandler()
		},
	}
	rootCmd.PersistentFlags().StringP("config", "c", "", "path to config file")
	rootCmd.Flags().IntP("replicas", "r", 2, "no. of replicas")
	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.BindPFlags(rootCmd.Flags())
	return rootCmd
}

func versionHandler() {
	fmt.Println("version v1.0.0")
	if config.VersionLong {
		fmt.Println("gitCommitHash: 21321321")
	}
}

func getVersionCommand() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:  "version",
		Args: cobra.NoArgs,
		Run: func(*cobra.Command, []string) {
			versionHandler()
		},
	}
	versionCmd.Flags().BoolP("long", "l", false, "verbose version info")
	viper.BindPFlags(versionCmd.Flags())
	return versionCmd
}

func readConfigFile() {
	viper.SetConfigType("yaml")
	if !viper.IsSet("config") {
		return
	}
	configFilePath := viper.GetString("config")
	fmt.Println("reading config from file at path", configFilePath)
	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("failed to read the config file at path %s . Error: %q", configFilePath, err)
	}
}

func setupViper() {
	viper.SetEnvPrefix("M2K")
	viper.AutomaticEnv()
	readConfigFile()
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("failed to unmarshal the config. Error: %q", err)
	}
}

func setupCobraAndRun() error {
	rootCmd := getRootCommand()
	rootCmd.AddCommand(getVersionCommand())
	cobra.OnInitialize(setupViper)
	return rootCmd.Execute()
}

func main() {
	if err := setupCobraAndRun(); err != nil {
		logrus.Fatal(err)
	}
}
