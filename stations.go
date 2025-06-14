package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Stations []Station `yaml:"stations"`
}

var stationsCmd = &cobra.Command{
	Use:   "stations",
	Short: "Manage your saved Velib stations",
}

var stationsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new station to your favorites",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Fetch station information
		resp, err := http.Get(stationInfoURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var info StationInformation
		if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
			return err
		}

		// Use fuzzy finder to select station
		idx, err := fuzzyfinder.Find(
			info.Data.Stations,
			func(i int) string {
				return info.Data.Stations[i].Name
			},
		)
		if err != nil {
			return err
		}

		selected := info.Data.Stations[idx] // Confirm selection
		fmt.Printf("Add station '%s' to your favorites? (Y/n): ", selected.Name)
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response == "n" || response == "no" {
			fmt.Println("Station not added.")
			return nil
		}

		// Load or create config
		configPath := filepath.Join(os.Getenv("HOME"), ".velispo.yaml")
		var config Config
		if data, err := os.ReadFile(configPath); err == nil {
			yaml.Unmarshal(data, &config)
		}

		// Check if station already exists
		for _, s := range config.Stations {
			if s.StationID == selected.StationID {
				fmt.Printf("Station '%s' is already in your favorites.\n", selected.Name)
				return nil
			}
		}

		// Add station
		config.Stations = append(config.Stations, selected)

		// Save config
		data, err := yaml.Marshal(&config)
		if err != nil {
			return err
		}

		if err := os.WriteFile(configPath, data, 0o644); err != nil {
			return err
		}

		fmt.Printf("Station '%s' added to your favorites!\n", selected.Name)
		return nil
	},
}

var stationsRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a station from your favorites",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load config
		configPath := filepath.Join(os.Getenv("HOME"), ".velispo.yaml")
		var config Config
		data, err := os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("no stations saved yet, use 'velispo stations add' first")
		}
		if err := yaml.Unmarshal(data, &config); err != nil {
			return err
		}

		if len(config.Stations) == 0 {
			fmt.Println("No stations saved yet.")
			return nil
		}

		// Use fuzzy finder to select station to remove
		idx, err := fuzzyfinder.Find(
			config.Stations,
			func(i int) string {
				return config.Stations[i].Name
			},
		)
		if err != nil {
			return err
		}

		selected := config.Stations[idx] // Confirm removal
		fmt.Printf("Remove station '%s' from your favorites? (Y/n): ", selected.Name)
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response == "n" || response == "no" {
			fmt.Println("Station not removed.")
			return nil
		}

		// Remove station
		config.Stations = append(config.Stations[:idx], config.Stations[idx+1:]...)

		// Save config
		data, err = yaml.Marshal(&config)
		if err != nil {
			return err
		}

		if err := os.WriteFile(configPath, data, 0o644); err != nil {
			return err
		}

		fmt.Printf("Station '%s' removed from your favorites!\n", selected.Name)
		return nil
	},
}

func init() {
	stationsCmd.AddCommand(stationsAddCmd)
	stationsCmd.AddCommand(stationsRemoveCmd)
}
