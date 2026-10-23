package main

import "snippetbox.peter/internal/models"

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}
