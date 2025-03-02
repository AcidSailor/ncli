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
	getSchemaIdentifier string
	getSchemaVersion    string
	getSchemaFormat     string

	getSchemaCmd = &cobra.Command{
		Use:   "get-schema",
		Short: "Get schema with specified identifier",
		Long:  "Get schema with specified identifier\nYou can retrieve all schema identifiers with \"ncli <required args> get --filter /netconf-state/schemas\"",
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

			var st string

			tagToValue := map[string]string{
				"identifier": getSchemaIdentifier,
				"version":    getSchemaVersion,
				"format":     getSchemaFormat,
			}

			for k, v := range tagToValue {
				if v != "" {
					st += utils.WrapWithTags(v, k)
				}
			}

			schemaPayload := fmt.Sprintf(
				"<get-schema xmlns=\"urn:ietf:params:xml:ns:yang:ietf-netconf-monitoring\">%s</get-schema>",
				st,
			)

			if withLock != "" {
				if _, err = d.Lock(withLock); err != nil {
					return err
				}
				defer d.Unlock(withLock)
			}

			r, err := d.RPC(
				opoptions.WithFilter(schemaPayload),
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
	getSchemaCmd.Flags().StringVar(&getSchemaIdentifier, "identifier", "", "schema identifier")
	getSchemaCmd.Flags().StringVar(&getSchemaVersion, "version", "", "schema version")
	getSchemaCmd.Flags().StringVar(&getSchemaFormat, "format", "yang", "schema format")
	getSchemaCmd.MarkFlagRequired("identifier")
}
