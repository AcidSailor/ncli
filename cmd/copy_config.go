package cmd

import (
	"fmt"
	"github.com/scrapli/scrapligo/driver/netconf"
	"github.com/spf13/cobra"
	"os"
)

var (
	copyConfigSource string
	copyConfigTarget string

	copyConfigCmd = &cobra.Command{
		Use:   "copy-config",
		Short: "Copy configuration from source to target datastore",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			d, err := netconf.NewDriver(
				driverOpts.Host,
				DriverCommonOptions()...,
			)
			if err != nil {
				return err
			}

			err = d.Open()
			if err != nil {
				return err
			}
			defer d.Close()

			if withLock != "" {
				if _, err = d.Lock(withLock); err != nil {
					return err
				}
				defer d.Unlock(withLock)
			}

			r, err := d.CopyConfig(
				copyConfigSource,
				copyConfigTarget,
			)
			if err != nil {
				return err
			}

			if r.Failed != nil {
				fmt.Fprintln(os.Stderr, r.Result)
				os.Exit(1)
			}
			fmt.Println(r.Result)

			return nil
		},
	}
)

func init() {
	copyConfigCmd.Flags().StringVar(&copyConfigSource, "source", "", "config source")
	copyConfigCmd.Flags().StringVar(&copyConfigTarget, "target", "", "config target")
	copyConfigCmd.MarkFlagRequired("source")
	copyConfigCmd.MarkFlagRequired("target")
	copyConfigCmd.MarkFlagsRequiredTogether("source", "target")
}
