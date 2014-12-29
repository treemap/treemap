package main

import (
	"fmt"
	render "github.com/abhiyerra/gowebcommons/render"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Zipcode struct {
	Number   string `json:"number"`
	GeomData string `json:"geom"`
	Center   string `json:"center"`
}

// Yeah, yeah this is a hack...
func milesToMeters(distance uint) string {
	switch distance {
	case 50:
		return "80467"
	case 100:
		return "160934"
	case 200:
		return "321869"
	}

	return "0"
}

func (zc *Zipcode) Parks(distance uint) (parks []NationalPark) {
	log.Println("Parks")

	err := db.Model(NationalPark{}).
		Select("ST_AsGeoJSON(ST_CollectionExtract(national_parks.geom, 3)) as geom_data, national_parks.unit_name, national_parks.unit_code").
		Joins(fmt.Sprintf("INNER JOIN zipcodes ON zipcodes.geoid10 = '%s' AND ST_DWithin(zipcodes.geom, national_parks.geom, %s, true)", zc.Number, milesToMeters(distance))).
		Scan(&parks)

	if err != nil {
		log.Println(err)
	}

	return
}

func zipcodeHydrologyHandler(hydroType string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		zipcode := vars["zipcode"]

		hydrology := cache.Get("hydrology/"+hydroType+"/"+zipcode, func() interface{} {
			var hydrology []Hydrology

			log.Println("Hydro Type:", hydroType, zipcode)

			err := db.Table(hydroType).
				Select(fmt.Sprintf("ST_AsGeoJSON(ST_CollectionExtract(%s.geom, 3)) as geom_data, %s.name", hydroType, hydroType)).
				Joins(fmt.Sprintf("INNER JOIN zipcodes ON zipcodes.geoid10 = '%s' AND ST_DWithin(zipcodes.geom, %s.geom, 80934 , true)", zipcode, hydroType)).
				Scan(&hydrology)

			if err != nil {
				log.Println(err)
			}

			return hydrology
		})

		render.RenderJson(w, hydrology)
	}
}

func showZipCodeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	zipcode := vars["zipcode"]

	z := cache.Get("zipcode/"+zipcode, func() interface{} {
		z := Zipcode{Number: zipcode}
		db.Select("geoid10 as number, ST_AsGeoJSON(ST_Centroid(geom)) as center, ST_AsGeoJSON(ST_CollectionExtract(geom, 3)) as geom_data").Where("geoid10 = ?", zipcode).First(&z)

		return z
	})

	render.RenderJson(w, z)
}

func zipcodeTableHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	zipcode := vars["zipcode"]
	table := vars["table"]

	log.Println("Zipcode", zipcode)
	parks := cache.Get(table+"/"+zipcode, func() interface{} {
		zc := Zipcode{Number: zipcode}
		switch table {
		case "parks":
			return zc.Parks(50)
		}

		return nil
	})

	render.RenderJson(w, parks)
}
