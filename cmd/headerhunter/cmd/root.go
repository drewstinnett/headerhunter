/*
Package cmd holds the commands for the executable
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "headerhunter",
		Short:         "Server to inspect http headers",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose logging")
	cmd.AddCommand(newServeCmd())
	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd := newRootCmd()
	cobra.OnInitialize(func() {
		initConfig(rootCmd)
	})
	if err := rootCmd.Execute(); err != nil {
		slog.Warn("fatal error occurred", "error", err)
		os.Exit(1)
	}
}

func initConfig(cmd *cobra.Command) {
	opts := log.Options{
		ReportTimestamp: true,
		Prefix:          "headerhunter 🫨 ",
	}
	if mustGetCmd[bool](*cmd, "verbose") {
		opts.Level = log.DebugLevel
	}
	logger := slog.New(log.NewWithOptions(os.Stderr, opts))
	slog.SetDefault(logger)
}

/*
func init() { //nolint: gochecknoinits // We have to initialize the cobra config on load
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// Slog
	opts := log.Options{
		ReportTimestamp: true,
		Prefix:          "headerhunter 🫨 ",
	}
	if verbose {
		opts.Level = log.DebugLevel
	}
	logger := slog.New(log.NewWithOptions(os.Stderr, opts))
	slog.SetDefault(logger)
}
*/
