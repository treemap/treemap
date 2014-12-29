package main

import (
	render "github.com/abhiyerra/gowebcommons/render"
	"net/http"
)

type Hydrology struct {
	Name     string `json:"name"`
	GeomData string `json:"geom"`
}

func hydrologyHandler(hydroType string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		hydrology := cache.Get(hydroType, func() interface{} {
			var hydrology []Hydrology
			db.Table(hydroType).Select("ST_AsGeoJSON(ST_CollectionExtract(geom, 3)) as geom_data, name").Scan(&hydrology)

			return hydrology
		})

		render.RenderJson(w, hydrology)
	}
}
