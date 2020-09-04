package cmd

import (
	"log"

	"github.com/j4ng5y/dhwh/pkg/server"
	"github.com/spf13/cobra"
)

func Run() {
	var (
		rootCMD = &cobra.Command{
			Use:     "dhwh",
			Version: "0.0.1",
			Run: func(ccmd *cobra.Command, args []string) {
				S := server.NewWithOptions(server.WithDefaults())
				S.Run()
			},
		}
	)

	if err := rootCMD.Execute(); err != nil {
		log.Fatal(err)
	}
}
