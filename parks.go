package main

import (
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

func parksHandler(w http.ResponseWriter, r *http.Request) {
	parks := cache.Get("parks", func() interface{} {
		var parks []NationalPark
		db.Model(NationalPark{}).Select("ST_AsGeoJSON(ST_CollectionExtract(geom, 3)) as geom_data, unit_name, unit_code").Scan(&parks)
		return parks
	})

	render.RenderJson(w, parks)
}
