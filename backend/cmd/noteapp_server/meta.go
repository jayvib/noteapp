package main

import "time"

var (
	// Version is the version of the current server
	Version = "development"
	// BuildCommit is the git build recent commit during server build.
	BuildCommit = "development"
	// BuildDate is the timestamp of when the server last build.
	BuildDate = time.Now().Truncate(time.Second).UTC()
)
