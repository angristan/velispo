package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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
			return fmt.Errorf("no stations saved, use 'velispo stations add' command first")
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
		t := table.NewWriter()
		t.SetStyle(table.StyleRounded)
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Station", "Mechanical", "E-bike", "Docks", "Last Update"})

		for _, station := range config.Stations {
			if status, ok := statusMap[station.StationCode]; ok {
				mechanical := status.NumBikesAvailableTypes[0].Mechanical
				ebike := status.NumBikesAvailableTypes[1].Ebike
				lastReported := time.Since(time.Unix(status.LastReported, 0)).Round(time.Minute)

				// Apply red color to zero values
				mechanicalStr := fmt.Sprintf("%d", mechanical)
				if mechanical == 0 {
					mechanicalStr = text.FgRed.Sprint(mechanicalStr)
				}

				ebikeStr := fmt.Sprintf("%d", ebike)
				if ebike == 0 {
					ebikeStr = text.FgRed.Sprint(ebikeStr)
				}

				docksStr := fmt.Sprintf("%d", status.NumDocksAvailable)
				if status.NumDocksAvailable == 0 {
					docksStr = text.FgRed.Sprint(docksStr)
				}

				t.AppendRow(table.Row{
					station.Name,
					mechanicalStr,
					ebikeStr,
					docksStr,
					lastReported.String(),
				})
			}
		}

		t.Render()
		return nil
	},
}
