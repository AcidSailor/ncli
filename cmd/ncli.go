package cmd

import (
	"fmt"
	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/logging"
	"github.com/scrapli/scrapligo/util"
	"github.com/spf13/cobra"
	"log"
)

type DriverOpts struct {
	Host      string
	Port      int
	Username  string
	Password  string
	NcVersion string
}

func DriverCommonOptions() []util.Option {
	return []util.Option{
		options.WithAuthNoStrictKey(),
		options.WithNetconfPreferredVersion(driverOpts.NcVersion),
		options.WithTransportType("standard"),
		options.WithAuthUsername(driverOpts.Username),
		options.WithAuthPassword(driverOpts.Password),
		options.WithPort(driverOpts.Port),
		options.WithLogger(logger),
	}
}

var (
	driverOpts       = &DriverOpts{}
	withLock         string
	withLoggingLevel string
	logger           *logging.Instance

	ncli = &cobra.Command{
		Use:   "ncli",
		Short: "Simple netconf command line client",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if withLoggingLevel != "" {
				logger, err = logging.NewInstance(
					logging.WithLevel(withLoggingLevel),
					logging.WithLogger(log.Print),
				)
			}
			return err
		},
	}
)

func SetVersionInfo(version, commit, date string) {
	ncli.Version = fmt.Sprintf("%s commit %s date %s ", version, commit, date)
}

func Execute() error {
	return ncli.Execute()
}

func init() {
	ncli.PersistentFlags().StringVar(&driverOpts.Host, "host", "", "hostname or address of the device")
	ncli.PersistentFlags().IntVar(&driverOpts.Port, "port", 830, "port of the device")
	ncli.MarkFlagRequired("host")

	ncli.PersistentFlags().StringVar(&driverOpts.Username, "username", "", "username for authentication")
	ncli.PersistentFlags().StringVar(&driverOpts.Password, "password", "", "password for authentication")
	ncli.MarkFlagsRequiredTogether("username", "password")

	ncli.PersistentFlags().StringVar(&withLock, "lock", "", "wrap calls with lock/unlock - if applicable")
	ncli.PersistentFlags().StringVar(&driverOpts.NcVersion, "with-nc-version", "1.0", "netconf version (1.0 or 1.1)")
	ncli.PersistentFlags().StringVar(&withLoggingLevel, "logging-level", "", "set logging level - info,debug,critical")

	commands := []*cobra.Command{
		helloCmd,
		getConfigCmd,
		getCmd,
		rpcCmd,
		killSessionCmd,
		copyConfigCmd,
		deleteConfigCmd,
		discardChangesCmd,
		commitCmd,
		validateCmd,
		editConfigCmd,
		getSchemaCmd,
	}

	for _, c := range commands {
		ncli.AddCommand(c)
	}
}
