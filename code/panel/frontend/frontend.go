package frontend

import "embed"

//go:embed static/*
var FileSystem embed.FS
