package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check availability at saved stations",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load config
		configPath := filepath.Join(os.Getenv("HOME"), ".velispo.yaml")
		var config Config
		data, err := os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("no stations saved, use search command first")
		}
		if err := yaml.Unmarshal(data, &config); err != nil {
			return err
		}

		// Get status
		resp, err := http.Get(stationStatusURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var status StationStatus
		if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
			return err
		}

		// Create status map
		statusMap := make(map[string]Status)
		for _, s := range status.Data.Stations {
			statusMap[s.StationCode] = s
		}

		// Create table
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Station", "Mechanical", "E-bike", "Docks", "Last Update"})

		for _, station := range config.Stations {
			if status, ok := statusMap[station.StationCode]; ok {
				mechanical := status.NumBikesAvailableTypes[0].Mechanical
				ebike := status.NumBikesAvailableTypes[1].Ebike
				lastReported := time.Since(time.Unix(status.LastReported, 0)).Round(time.Minute)

				table.Append([]string{
					station.Name,
					fmt.Sprintf("%d", mechanical),
					fmt.Sprintf("%d", ebike),
					fmt.Sprintf("%d", status.NumDocksAvailable),
					lastReported.String(),
				})
			}
		}

		table.Render()
		return nil
	},
}
