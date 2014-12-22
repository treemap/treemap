package main

import (
	"fmt"
	render "github.com/abhiyerra/gowebcommons/render"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type NationalPark struct {
	UnitName string `json:"name"`
	UnitCode string `json:"code"`
	GeomData string `json:"geom"`
}

func nearbyParksHandler(w http.ResponseWriter, r *http.Request) {
	var parks []NationalPark

	longitude := r.URL.Query().Get("long")
	latitude := r.URL.Query().Get("lat")
	log.Println("Long:", longitude, "Lat:", latitude)

	err := db.Model(NationalPark{}).
		Select("ST_AsGeoJSON(ST_CollectionExtract(geom, 3)) as geom_data, unit_name, unit_code").
		Where(fmt.Sprintf("ST_DWithin(ST_GeomFromText('POINT(%s %s)' , 4326)::geography, geom, 160934, true)", longitude, latitude)). // Within 100 miles -> 160934 meters
		Scan(&parks)
	if err != nil {
		log.Println(err)
	}

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
