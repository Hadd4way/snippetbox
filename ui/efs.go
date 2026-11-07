package ui

import (
	"embed"
)

//this comment means instructs go to store the files
// from ui/static folder in a embedded filesystem
// referenced by the global variable Files

//go:embed "static"
var Files embed.FS
