package cmd

import (
	"fmt"
	"github.com/acidsailor/ncli/internal/utils"
	"github.com/scrapli/scrapligo/driver/netconf"
	"github.com/scrapli/scrapligo/driver/opoptions"
	"github.com/scrapli/scrapligo/driver/options"
	"github.com/spf13/cobra"
	"os"
)

var (
	getConfigSource      string
	getConfigFilterPath  string
	getConfigFilterValue string
	getConfigFilterFile  string
	getConfigFilterType  string

	getConfigCmd = &cobra.Command{
		Use:   "get-config",
		Short: "Send get-config rpc with specified filter and source datasource",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

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

			err = d.Open()
			if err != nil {
				return err
			}
			defer d.Close()

			var f string

			switch getConfigFilterType {
			case "subtree":
				f = utils.FlatPathToSubtreeWithValue(getConfigFilterPath, getConfigFilterValue)
			case "xpath":
				f = getConfigFilterPath
			}

			if getConfigFilterFile != "" {
				fb, err := os.ReadFile(getConfigFilterFile)
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

			r, err := d.GetConfig(
				getConfigSource,
				opoptions.WithFilter(f),
				opoptions.WithFilterType(getConfigFilterType),
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
	getConfigCmd.Flags().StringVar(&getConfigFilterType, "filter-type", "subtree", "filter type - subtree or xpath")
	getConfigCmd.Flags().StringVar(&getConfigFilterPath, "path", "", "filter path")
	getConfigCmd.Flags().StringVar(&getConfigFilterValue, "value", "", "filter value")
	getConfigCmd.Flags().StringVar(&getConfigFilterFile, "filter-file", "", "filter file")
	getConfigCmd.Flags().StringVar(&getConfigSource, "source", "", "config source")
	getConfigCmd.MarkFlagRequired("source")
	getConfigCmd.MarkFlagsOneRequired("path", "filter-file")
	getConfigCmd.MarkFlagsMutuallyExclusive("path", "filter-file")
}
