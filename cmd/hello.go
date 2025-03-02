package cmd

import (
	"context"
	"fmt"
	"github.com/acidsailor/ncli/internal/utils"
	"github.com/scrapli/scrapligo/driver/netconf"
	"github.com/scrapli/scrapligo/driver/options"
	"github.com/spf13/cobra"
)

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "Send hello request",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		d, err := netconf.NewDriver(
			driverOpts.Host,
			options.WithAuthNoStrictKey(),
			options.WithAuthUsername(driverOpts.Username),
			options.WithAuthPassword(driverOpts.Password),
			options.WithPort(driverOpts.Port),
			options.WithLogger(logger),
		)
		if err != nil {
			return err
		}

		err = d.Channel.Open()
		if err != nil {
			return err
		}
		defer d.Channel.Close()

		r, err := d.Channel.ReadUntilPrompt(ctx)
		if err != nil {
			return err
		}

		fmt.Println(utils.NetconfStrip(string(r)))

		return nil
	},
}
