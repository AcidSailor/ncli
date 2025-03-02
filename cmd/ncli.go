package cmd

import (
	"fmt"
	"github.com/scrapli/scrapligo/logging"
	"github.com/spf13/cobra"
	"log"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type DriverOpts struct {
	Host     string
	Port     int
	Username string
	Password string
}

var (
	driverOpts       = &DriverOpts{}
	withLock         string
	withLoggingLevel string
	logger           *logging.Instance

	ncli = &cobra.Command{
		Use:     "ncli",
		Short:   "Simple netconf command line client",
		Version: fmt.Sprintf("%s commit %s date %s ", version, commit, date),
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
