package cmd

import (
	"fmt"
	"github.com/scrapli/scrapligo/driver/netconf"
	"github.com/scrapli/scrapligo/driver/opoptions"
	"github.com/spf13/cobra"
)

var (
	sessionId int

	killSessionCmd = &cobra.Command{
		Use:   "kill-session",
		Short: "Kills session with the specified id",
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

			r, err := d.RPC(
				opoptions.WithFilter(
					fmt.Sprintf(
						"<kill-session><session-id>%d</session-id></kill-session>", sessionId),
				),
			)
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
	killSessionCmd.Flags().IntVar(&sessionId, "session-id", 0, "session id")
	_ = killSessionCmd.MarkFlagRequired("session-id")
}
