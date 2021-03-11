package app

import "embed"

//go:embed embed/*
var EmbedFS embed.FS
