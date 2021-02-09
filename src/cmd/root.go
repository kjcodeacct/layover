/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"layover/tcp"
	"layover/udp"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var log *zap.SugaredLogger

var (
	cfgFile string

	ProxyHost string
	ProxyPort int
	Protocol  string
	ServePort int
	DebugMode int
	LogDir    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "layover",
	Short: "Proxy layover process",
	Long: `# Layover
single port proxy process
<https://github.com/kjcodeacct/layover>

# Env Variables
* LAYOVER_PROXYHOST - default:"0.0.0.0"
	* the host layover is proxying from, unless specifying to a different host machine uses the default

* LAYOVER_PROXYPORT - required:true
	* the port layover is proxying *FROM*
	* this is *typically* the port not in the container system

* LAYOVER_PROTOCOL - default:"tcp"
	* IP protocol used by the specified port
	* options available
		* "tcp"
		* "udp"

* LAYOVER_SERVEPORT default - default:"8080"
	* the port layover is proxying *TO* and is serving
	* if running in a container typically does *not* need to be specified

* LAYOVER_DEBUGMODE default - "0"
	* options available
		* 0 - off
		* 1 - basic logging of IP connecting and warnings
		* 2 - full logging including data (please don't use in production)

* LAYOVER_LOGDIR
	* directory to place logfiles created by enabling the LAYOVER_DEBUGMOD
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.layover.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	/* for testing*/

	ProxyHost = "127.0.0.1"
	ProxyPort = 8090
	ServePort = 8081

	initLog()
}

func initLog() {
	zapLogger, _ := zap.NewProduction()

	defer zapLogger.Sync()
	log = zapLogger.Sugar()

	tcp.SetLog(log)
	udp.SetLog(log)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".layover" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".layover")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
