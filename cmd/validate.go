package cmd

import (
	"fmt"
	"github.com/scrapli/scrapligo/driver/netconf"
	"github.com/spf13/cobra"
)

var (
	validateSource string
	validateCmd    = &cobra.Command{
		Use:   "validate",
		Short: "Validate changes in specified datastore",
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

			if withLock != "" {
				if _, err = d.Lock(withLock); err != nil {
					return err
				}
				defer func() { _, _ = d.Unlock(withLock) }()
			}

			r, err := d.Validate(validateSource)
			if err != nil {
				return err
			}

			if r.Failed != nil {
				return r.Failed
			}
			fmt.Println(r.Result)

			return nil
		},
	}
)

func init() {
	validateCmd.Flags().StringVar(&validateSource, "source", "candidate", "source to validate")
}
