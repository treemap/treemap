package main

import (
	"fmt"
	render "github.com/abhiyerra/gowebcommons/render"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type NationalPark struct {
	UnitName string `json:"name"`
	UnitCode string `json:"code"`
	GeomData string `json:"geom"`
}

func zipcodeParksHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	zipcode := vars["zipcode"]

	log.Println("Zipcode", zipcode)

	parks := cache.Get("parks/"+vars["zipcode"], func() interface{} {
		var parks []NationalPark

		err := db.Model(NationalPark{}).
			Select("ST_AsGeoJSON(ST_CollectionExtract(national_parks.geom, 3)) as geom_data, national_parks.unit_name, national_parks.unit_code").
			Joins(fmt.Sprintf("INNER JOIN zipcodes ON zipcodes.geoid10 = '%s' AND ST_DWithin(zipcodes.geom, national_parks.geom, 80934 , true)", zipcode)).
			Scan(&parks)

		if err != nil {
			log.Println(err)
		}

		return parks
	})

	render.RenderJson(w, parks)
}

func parksHandler(w http.ResponseWriter, r *http.Request) {
	parks := cache.Get("parks", func() interface{} {
		var parks []NationalPark
		db.Model(NationalPark{}).Select("ST_AsGeoJSON(ST_CollectionExtract(geom, 3)) as geom_data, unit_name, unit_code").Scan(&parks)
		return parks
	})

	render.RenderJson(w, parks)
}
