package cmd

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/drewstinnett/headerhunter"
	"github.com/spf13/cobra"
)

const (
	defaultTimeOutSeconds      = 5
	defaultReadTimeoutMinutes  = 10
	defaultWriteTimeoutMinutes = 10
)

func hunterWithCmd(cmd *cobra.Command, args []string) (*headerhunter.Hunter, error) {
	var opt headerhunter.Option
	if strings.HasPrefix(args[0], "http") {
		opt = headerhunter.WithProxyURL(args[0])
	} else {
		opt = headerhunter.WithStaticDir(args[0])
	}

	return headerhunter.New(opt, headerhunter.WithPrefix(mustGetCmd[string](*cmd, "prefix")))
}

func newServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "serve DIR|URL",
		Args: cobra.ExactArgs(1),
		Example: `$ headerhunter serve /srv/public
$ headehunter serve https://www.example.com
`,
		Short: "Inspect headers and serve up a static directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := hunterWithCmd(cmd, args)
			if err != nil {
				return err
			}

			listener, err := net.Listen("tcp", mustGetCmd[string](*cmd, "addr"))
			if err != nil {
				return err
			}

			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			s := &http.Server{
				Handler:        h,
				ReadTimeout:    mustGetCmd[time.Duration](*cmd, "read-timeout"),
				WriteTimeout:   mustGetCmd[time.Duration](*cmd, "write-timeout"),
				MaxHeaderBytes: 0,
			}

			go func() {
				<-ctx.Done()
				shutdownCtx, cancel := context.WithTimeout(context.Background(), defaultTimeOutSeconds*time.Second)
				defer cancel()
				if err = s.Shutdown(shutdownCtx); err != nil {
					slog.Error("error shutting down server", "error", err)
				}
			}()

			cert := mustGetCmd[string](*cmd, "cert")
			if cert != "" {
				key := mustGetCmd[string](*cmd, "key")
				slog.Info("launching https server", "addr", listener.Addr().(*net.TCPAddr))
				return s.ServeTLS(listener, cert, key)
			}
			slog.Info("launching http server", "addr", listener.Addr().(*net.TCPAddr))
			return s.Serve(listener)
		},
	}
	cmd.PersistentFlags().StringP("addr", "a", ":3000", "address to listen on")
	cmd.PersistentFlags().StringP("prefix", "p", "/", "prefix to route requests to")
	cmd.PersistentFlags().Duration("read-timeout", defaultReadTimeoutMinutes*time.Minute, "read timeout for the server")
	cmd.PersistentFlags().
		Duration("write-timeout", defaultWriteTimeoutMinutes*time.Minute, "write timeout for the server")
	// TLS Options
	cmd.PersistentFlags().StringP("cert", "c", "", "TLS Cert")
	cmd.PersistentFlags().StringP("key", "k", "", "TLS Key")
	cmd.MarkFlagsRequiredTogether("cert", "key")
	return cmd
}
