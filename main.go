package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	stationInfoURL   = "https://velib-metropole-opendata.smovengo.cloud/opendata/Velib_Metropole/station_information.json"
	stationStatusURL = "https://velib-metropole-opendata.smovengo.cloud/opendata/Velib_Metropole/station_status.json"
)

var rootCmd = &cobra.Command{
	Use:   "velispo",
	Short: "Check Velib availability near you",
}

func init() {
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(checkCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
