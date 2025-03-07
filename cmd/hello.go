package cmd

import (
	"context"
	"fmt"
	"github.com/acidsailor/ncli/internal/utils"
	"github.com/scrapli/scrapligo/driver/netconf"
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
			DriverCommonOptions()...,
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
