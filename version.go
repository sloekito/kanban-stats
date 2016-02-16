package main

const ApplicationName = "kanban-stats"
const Version string = "1.0"

var GitCommit string = ""

type VersionInfo struct{}

func GetVersion() string {
	return Version + "." + GitCommit
}
