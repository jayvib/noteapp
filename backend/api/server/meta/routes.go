package meta

import "time"

// Metadata contains the information for the API server.
type Metadata struct {
	Version     string    `json:"version,omitempty"`
	BuildCommit string    `json:"build_commit,omitempty"`
	BuildDate   time.Time `json:"build_date,omitempty"`
}
