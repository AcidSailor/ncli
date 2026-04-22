package cmd

import (
	"fmt"
	"os"

	"github.com/acidsailor/ncli/internal/utils"
	"github.com/scrapli/scrapligo/driver/netconf"
	"github.com/spf13/cobra"
)

var (
	editConfigTarget   string
	editConfigPath     string
	editConfigValue    string
	editConfigFile     string
	editConfigValidate bool
	editConfigDiscard  bool
	editConfigCommit   bool

	editConfigCmd = &cobra.Command{
		Use:   "edit-config",
		Short: "Send edit-config rpc to specified target datastore",
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
			defer func() { _ = d.Close() }()

			var cfg string

			if editConfigTarget != "" {
				cfg = utils.WrapWithTags(
					utils.FlatPathToSubtreeWithValue(
						editConfigPath,
						editConfigValue,
					),
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
				defer func() { _, _ = d.Unlock(withLock) }()
			}

			r, err := d.EditConfig(
				editConfigTarget,
				cfg,
			)
			if err != nil {
				return err
			}
			if r.Failed != nil {
				return r.Failed
			}
			fmt.Println(r.Result)

			if editConfigValidate {
				r, err = d.Validate(editConfigTarget)
				if err != nil {
					return err
				}
				if r.Failed != nil {
					return r.Failed
				}
				fmt.Println(r.Result)
			}

			if editConfigDiscard {
				r, err = d.Discard()
				if err != nil {
					return err
				}
				if r.Failed != nil {
					return r.Failed
				}
				fmt.Println(r.Result)
			}

			if editConfigCommit {
				r, err = d.Commit()
				if err != nil {
					return err
				}
				if r.Failed != nil {
					return r.Failed
				}
				fmt.Println(r.Result)
			}

			return nil
		},
	}
)

func init() {
	editConfigCmd.Flags().StringVar(&editConfigPath, "path", "", "config path")
	editConfigCmd.Flags().
		StringVar(&editConfigValue, "value", "", "config value for specified path")
	editConfigCmd.Flags().
		StringVar(&editConfigFile, "config-file", "", "config file")
	editConfigCmd.Flags().
		StringVar(&editConfigTarget, "target", "", "config target")
	editConfigCmd.Flags().BoolVar(&editConfigValidate, "validate", false,
		"execute validate operation after edit-config")
	editConfigCmd.Flags().BoolVar(&editConfigCommit, "commit", false,
		"execute commit operation after edit-config")
	editConfigCmd.Flags().BoolVar(&editConfigDiscard, "discard", false,
		"execute discard operation after edit-config")
	_ = editConfigCmd.MarkFlagRequired("target")
	editConfigCmd.MarkFlagsOneRequired("path", "config-file")
	editConfigCmd.MarkFlagsRequiredTogether("path", "value")
	editConfigCmd.MarkFlagsMutuallyExclusive("path", "config-file")
	editConfigCmd.MarkFlagsMutuallyExclusive("commit", "discard")
}
