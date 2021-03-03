package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const processName = "layover"
const localhost = "127.0.0.1"

var logger *zap.SugaredLogger

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   processName,
	Short: "Proxy layover process",
	Long: `Layover
Single port proxy process
https://github.com/kjcodeacct/layover
`,
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

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.layover.yaml)")
	rootCmd.PersistentFlags().StringVar(&LogDir, "logdir", cwd, "directory to place log files in")

	initLog()
}

func initLog() {
	zapLogger, _ := zap.NewProduction()

	defer zapLogger.Sync()
	logger = zapLogger.Sugar()
}
