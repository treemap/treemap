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
	default:
		return "160934" // Default to 100 miles

	}
}

func (zc *Zipcode) GetInfo() {
	db.Select("geoid10 as number, ST_AsGeoJSON(ST_Centroid(geom)) as center, ST_AsGeoJSON(ST_CollectionExtract(geom, 3)) as geom_data").Where("geoid10 = ?", zc.Number).First(zc)
}

func (zc *Zipcode) Parks(distance uint) (parks []NationalPark) {
	err := db.Model(NationalPark{}).
		Select("ST_AsGeoJSON(ST_CollectionExtract(national_parks.geom, 3)) as geom_data, national_parks.unit_name, national_parks.unit_code").
		Joins(fmt.Sprintf("INNER JOIN zipcodes ON zipcodes.geoid10 = '%s' AND ST_DWithin(zipcodes.geom, national_parks.geom, %s, true)", zc.Number, milesToMeters(distance))).
		Scan(&parks)

	if err != nil {
		log.Println(err)
	}

	return
}

func (zc *Zipcode) Trees(distance uint) (trees []Tree) {

	err := db.Model(Tree{}).Select("distinct trees.id, trees.latin_name, trees.common_name").
		// TODO: Sql injection here. Need to sanatize this.
		Joins(fmt.Sprintf("INNER JOIN tree_geoms ON tree_geoms.latin_name = trees.latin_name INNER JOIN zipcodes ON zipcodes.geoid10 = '%s' AND ST_DWithin(zipcodes.geom, tree_geoms.geom, %s, true)", zc.Number, milesToMeters(distance))).
		Order("trees.latin_name asc").Scan(&trees)

	if err != nil {
		log.Println(err)
	}

	return
}

func (zc *Zipcode) Hydrology(hydroType string, distance uint) (hydrology []Hydrology) {
	err := db.Table(hydroType).
		Select(fmt.Sprintf("ST_AsGeoJSON(ST_CollectionExtract(%s.geom, 3)) as geom_data, %s.name", hydroType, hydroType)).
		Joins(fmt.Sprintf("INNER JOIN zipcodes ON zipcodes.geoid10 = '%s' AND ST_DWithin(zipcodes.geom, %s.geom, %s, true)", zc.Number, milesToMeters(distance), hydroType)).
		Scan(&hydrology)

	if err != nil {
		log.Println(err)
	}

	return
}

func showZipCodeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	zipcode := vars["zipcode"]

	z := cache.Get("zipcode/"+zipcode, func() interface{} {
		z := Zipcode{Number: zipcode}
		z.GetInfo()

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
		case "trees":
			return zc.Trees(50)
		case "lakes":
			return zc.Hydrology("lakes", 50)
		case "rivers":
			return zc.Hydrology("rivers", 50)
		}

		return nil
	})

	render.RenderJson(w, parks)
}
