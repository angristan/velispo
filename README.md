# Velispo 🚴‍♂️

A command-line tool to check Velib (Paris bike-sharing) station availability near you.

## Features

- 🔍 **Fuzzy search** for Velib stations across Paris
- ⭐ **Save favorite stations** for quick monitoring
- 📊 **Real-time availability** display with color-coded alerts
- 🚨 **Zero availability highlighting** - stations with no bikes/docks are shown in red
- 💻 **Clean terminal interface** with beautiful tables

## Installation

### Prerequisites

- Go 1.24 or later

### Build from source

```bash
git clone <repository-url>
cd velispo
go build
```

This will create a `velispo` executable in the current directory.

## Usage

### Managing Stations

#### Add a Station to Favorites

```bash
velispo stations add
```

This will:

1. Open a fuzzy finder with all available Velib stations
2. Let you search and select a station
3. Ask for confirmation (press Enter to confirm, or type 'n' to cancel)
4. Save the station to your favorites

#### Remove a Station from Favorites

```bash
velispo stations remove
```

This will:

1. Show a fuzzy finder with your saved stations
2. Let you select which station to remove
3. Ask for confirmation (press Enter to confirm, or type 'n' to cancel)
4. Remove the station from your favorites

### Checking Station Availability

```bash
velispo check
```

Displays a table showing real-time availability for all your saved stations:

- **Station**: Station name
- **Mechanical**: Number of regular bikes available
- **E-bike**: Number of electric bikes available
- **Docks**: Number of available docking spots
- **Last Update**: How long ago the data was updated

**Color coding**: Any cell showing `0` will be displayed in **red** to quickly identify unavailable resources.

### Help

```bash
velispo --help              # General help
velispo stations --help     # Help for stations command
velispo stations add --help # Help for add subcommand
```

## Example Output

```
╭────────────────────────────────┬────────────┬────────┬───────┬─────────────╮
│ STATION                        │ MECHANICAL │ E-BIKE │ DOCKS │ LAST UPDATE │
├────────────────────────────────┼────────────┼────────┼───────┼─────────────┤
│ Place des Ternes.              │ 3          │ 1      │ 41    │ 36m0s       │
│ Parc Monceau                   │ 2          │ 4      │ 12    │ 42m0s       │
│ Place de la Madeleine - Royale │ 13         │ 10     │ 16    │ 36m0s       │
╰────────────────────────────────┴────────────┴────────┴───────┴─────────────╯
```

In the example above, the `0` in the Mechanical column for "Place Saint-Ferdinand" would appear in red.

## Configuration

Station favorites are stored in `~/.velispo.yaml`. You can manually edit this file if needed, but it's recommended to use the `velispo stations` commands instead.

## Data Source

This tool uses the official Velib Métropole Open Data API:

- Station information: Real-time station locations and details
- Station status: Live availability data for bikes and docks

## Dependencies

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [go-fuzzyfinder](https://github.com/ktr0731/go-fuzzyfinder) - Interactive fuzzy finder
- [go-pretty](https://github.com/jedib0t/go-pretty) - Table formatting and colors
- [yaml.v2](https://gopkg.in/yaml.v2) - YAML configuration handling

---

**Note**: This is an unofficial tool and is not affiliated with Velib or Smovengo.
