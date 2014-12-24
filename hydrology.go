package main

import (
	"fmt"
	render "github.com/abhiyerra/gowebcommons/render"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Hydrology struct {
	Name     string `json:"name"`
	GeomData string `json:"geom"`
}

func zipcodeHydrologyHandler(hydroType string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		zipcode := vars["zipcode"]

		var hydrology []Hydrology

		longitude := r.URL.Query().Get("long")
		latitude := r.URL.Query().Get("lat")
		log.Println("Hydro Type:", hydroType, "Long:", longitude, "Lat:", latitude)

		err := db.Table(hydroType).
			Select(fmt.Sprintf("ST_AsGeoJSON(ST_CollectionExtract(%s.geom, 3)) as geom_data, %s.name", hydroType, hydroType)).
			Joins(fmt.Sprintf("INNER JOIN zipcodes ON zipcodes.geoid10 = '%s' AND ST_DWithin(zipcodes.geom, %s.geom, 160934 , true)", zipcode, hydroType)).
			Scan(&hydrology)
		if err != nil {
			log.Println(err)
		}

		render.RenderJson(w, hydrology)
	}
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
