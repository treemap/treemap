package main

import (
	"fmt"
	render "github.com/abhiyerra/gowebcommons/render"
	"log"
	"net/http"
)

type Hydrology struct {
	Name     string `json:"name"`
	GeomData string `json:"geom"`
}

func nearbyHydrologyHandler(hydroType string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var hydrology []Hydrology

		longitude := r.URL.Query().Get("long")
		latitude := r.URL.Query().Get("lat")
		log.Println("Hydro Type:", hydroType, "Long:", longitude, "Lat:", latitude)

		err := db.Table(hydroType).
			Select("ST_AsGeoJSON(ST_CollectionExtract(geom, 3)) as geom_data, name").
			Where(fmt.Sprintf("ST_DWithin(ST_GeomFromText('POINT(%s %s)' , 4326)::geography, geom, 160934, true)", longitude, latitude)). // Within 100 miles -> 160934 meters
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
