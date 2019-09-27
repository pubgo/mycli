package cmds

import (
	"fmt"
	"github.com/pubgo/mycli/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "db2rest version info",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("version", version.Version)
			fmt.Println("commitV", version.CommitV)
			fmt.Println("buildV", version.BuildV)
			return nil
		},
	})
}
