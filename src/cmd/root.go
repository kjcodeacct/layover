package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const processName = "layover"
const localhost = "127.0.0.1"

var log *zap.SugaredLogger

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   processName,
	Short: "Proxy layover process",
	Long: `# Layover
single port proxy process
<https://github.com/kjcodeacct/layover>
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.layover.yaml)")

	initLog()
}

func initLog() {
	zapLogger, _ := zap.NewProduction()

	defer zapLogger.Sync()
	log = zapLogger.Sugar()
}
