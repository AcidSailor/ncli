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
	getFilterPath  string
	getFilterValue string
	getFilterFile  string
	getFilterType  string

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Send get rpc with specified filter or filter file",
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

			switch getFilterType {
			case "subtree":
				f = utils.FlatPathToSubtreeWithValue(getFilterPath, getFilterValue)
			case "xpath":
				f = getFilterPath
			}

			if getFilterFile != "" {
				fb, err := os.ReadFile(getFilterFile)
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

			r, err := d.Get(
				f,
				opoptions.WithFilterType(getFilterType),
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
	getCmd.Flags().StringVar(&getFilterType, "filter-type", "subtree", "filter type - subtree or xpath")
	getCmd.Flags().StringVar(&getFilterPath, "path", "", "filter path")
	getCmd.Flags().StringVar(&getFilterValue, "value", "", "filter value")
	getCmd.Flags().StringVar(&getFilterFile, "filter-file", "", "filter file")
	getCmd.MarkFlagsOneRequired("path", "filter-file")
	getCmd.MarkFlagsMutuallyExclusive("path", "filter-file")
}
