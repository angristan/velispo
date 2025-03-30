package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Stations []Station `yaml:"stations"`
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search and save Velib stations",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := http.Get(stationInfoURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var info StationInformation
		if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
			return err
		}

		idx, err := fuzzyfinder.Find(
			info.Data.Stations,
			func(i int) string {
				return info.Data.Stations[i].Name
			},
		)
		if err != nil {
			return err
		}

		selected := info.Data.Stations[idx]

		// Load or create config
		configPath := filepath.Join(os.Getenv("HOME"), ".velispo.yaml")
		var config Config
		if data, err := os.ReadFile(configPath); err == nil {
			yaml.Unmarshal(data, &config)
		}

		// Add station if not exists
		exists := false
		for _, s := range config.Stations {
			if s.StationID == selected.StationID {
				exists = true
				break
			}
		}
		if !exists {
			config.Stations = append(config.Stations, selected)
		}

		// Save config
		data, err := yaml.Marshal(&config)
		if err != nil {
			return err
		}
		return os.WriteFile(configPath, data, 0o644)
	},
}
