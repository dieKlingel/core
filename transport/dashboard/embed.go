package dashboard

import "embed"

//go:embed html/*
var files embed.FS

func Files() embed.FS {
	return files
}
