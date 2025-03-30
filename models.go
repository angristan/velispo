package main

type StationInformation struct {
	LastUpdatedOther int `json:"lastUpdatedOther"`
	TTL              int `json:"ttl"`
	Data             struct {
		Stations []Station `json:"stations"`
	} `json:"data"`
}

type Station struct {
	StationID   int64   `json:"station_id"`
	StationCode string  `json:"stationCode"`
	Name        string  `json:"name"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Capacity    int     `json:"capacity"`
}

type StationStatus struct {
	LastUpdatedOther int `json:"lastUpdatedOther"`
	TTL              int `json:"ttl"`
	Data             struct {
		Stations []Status `json:"stations"`
	} `json:"data"`
}

type Status struct {
	StationID              int64 `json:"station_id"`
	NumBikesAvailable      int   `json:"num_bikes_available"`
	NumBikesAvailableTypes []struct {
		Mechanical int `json:"mechanical,omitempty"`
		Ebike      int `json:"ebike,omitempty"`
	} `json:"num_bikes_available_types"`
	NumDocksAvailable int    `json:"num_docks_available"`
	LastReported      int64  `json:"last_reported"`
	StationCode       string `json:"stationCode"`
}
