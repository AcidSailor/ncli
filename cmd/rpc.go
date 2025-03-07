package cmd

import (
	"fmt"
	"github.com/acidsailor/ncli/internal/utils"
	"github.com/scrapli/scrapligo/driver/netconf"
	"github.com/scrapli/scrapligo/driver/opoptions"
	"github.com/spf13/cobra"
	"os"
)

var (
	rpcPath string
	rpcFile string

	rpcCmd = &cobra.Command{
		Use:   "rpc",
		Short: "Send rpc request",
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

			var f string

			if rpcPath != "" {
				f = utils.FlatPathToSubtreeWithValue(rpcPath, "")
			}

			if rpcFile != "" {
				fb, err := os.ReadFile(rpcFile)
				if err != nil {
					return err
				}
				f = string(fb)
			}

			if withLock != "" {
				if _, err = d.Lock(withLock); err != nil {
					return err
				}
				defer d.Unlock(withLock)
			}

			r, err := d.RPC(
				opoptions.WithFilter(f),
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
	rpcCmd.Flags().StringVar(&rpcPath, "rpc", "", "rpc path")
	rpcCmd.Flags().StringVar(&rpcFile, "rpc-file", "", "rpc file")
	rpcCmd.MarkFlagsOneRequired("rpc", "rpc-file")
	rpcCmd.MarkFlagsMutuallyExclusive("rpc", "rpc-file")
}
