package controllers

import (
	"net/http"

	"github.com/jag6/fotogen/models"
)

type Galleries struct {
	Template struct {
		New Template
	}
	GalleryService *models.GalleryService
}

func (g Galleries) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")
	g.Template.New.Execute(w, r, data)
}
