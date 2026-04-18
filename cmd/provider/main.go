package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/lmittmann/tint"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background() // Create base context

	if err := run(ctx); err != nil {
		msg := fmt.Sprintf("Application exited with error: %+v", err)
		slog.ErrorContext(ctx, msg)
		os.Exit(1)
	}

	slog.InfoContext(ctx, "Application exited")
}

// run manages app lifecycle, signal handling, and runner.
func run(ctx context.Context) error {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen,
	})))

	pvr, cleanup, err := NewProvider(ctx)
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		defer stop()

		return providerserver.Serve(ctx, func() provider.Provider {
			return pvr
		}, providerserver.ServeOpts{
			// NOTE: This is not a typical Terraform Registry provider address,
			// such as registry.terraform.io/hashicorp/hashicups. This specific
			// provider address is used in these tutorials in conjunction with a
			// specific Terraform CLI configuration for manual development testing
			// of this provider.
			Address: "hashicorp.com/edu/hashicups",
			Debug:   debug,
		})
	})

	g.Go(func() error {
		defer cleanup()

		<-ctx.Done()

		return nil
	})

	return g.Wait()
}
