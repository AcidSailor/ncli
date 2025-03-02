package cmd

import (
	"fmt"
	"github.com/acidsailor/ncli/internal/utils"
	"github.com/scrapli/scrapligo/driver/netconf"
	"github.com/scrapli/scrapligo/driver/options"
	"github.com/spf13/cobra"
	"os"
)

var (
	editConfigTarget string
	editConfigPath   string
	editConfigValue  string
	editConfigFile   string

	editConfigCmd = &cobra.Command{
		Use:   "edit-config",
		Short: "Send edit-config rpc to specified target datastore",
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

			var cfg string

			if editConfigTarget != "" {
				cfg = utils.WrapWithTags(
					utils.FlatPathToSubtreeWithValue(editConfigPath, editConfigValue),
					"config",
				)
			}

			if editConfigFile != "" {
				cb, err := os.ReadFile(editConfigFile)
				if err != nil {
					return err
				}
				cfg = string(cb)
			}

			if withLock != "" {
				if _, err = d.Lock(withLock); err != nil {
					return err
				}
				defer d.Unlock(withLock)
			}

			r, err := d.EditConfig(
				editConfigTarget,
				cfg,
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
	editConfigCmd.Flags().StringVar(&editConfigPath, "path", "", "config path")
	editConfigCmd.Flags().StringVar(&editConfigValue, "value", "", "config value for specified path")
	editConfigCmd.Flags().StringVar(&editConfigFile, "config-file", "", "config file")
	editConfigCmd.Flags().StringVar(&editConfigTarget, "target", "", "config target")
	editConfigCmd.MarkFlagRequired("target")
	editConfigCmd.MarkFlagsOneRequired("path", "config-file")
	editConfigCmd.MarkFlagsRequiredTogether("path", "value")
	editConfigCmd.MarkFlagsMutuallyExclusive("path", "config-file")
}
