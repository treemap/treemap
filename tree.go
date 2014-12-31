package main

import (
	"log"
)

type Tree struct {
	Id         int64    `json:"id"`
	CommonName string   `json:"common_name"`
	LatinName  string   `json:"latin_name"`
	GeomData   []string `json:"geom",sql:"-"`
	Area       float64  `json:"area",sql:"-"`
	Center     string   `json:"center",sql:"-"`
}

func AllTrees() (trees []Tree) {
	err := db.Model(Tree{}).Select("id, latin_name, common_name").Scan(&trees)
	if err != nil {
		log.Println(err)
	}

	return
}

func (t *Tree) GetGeodata() {
	rows, err := db.Table("tree_geoms").Select("ST_AsGeoJSON(ST_CollectionExtract(geom, 3)) as geom2").Where("latin_name = ?", t.LatinName).Rows()
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		var geodata string
		rows.Scan(&geodata)
		t.GeomData = append(t.GeomData, geodata)
	}
}

func (t *Tree) GetArea() {
	var a struct {
		Area float64
	}
	db.Table("tree_geoms").Select("SUM(ST_Area(ST_Transform(geom, 900913))) as area").Where("latin_name = ?", t.LatinName).Scan(&a)

	t.Area = a.Area * 0.000189394 * 0.000189394 // Get the miles
	log.Println("Area:", t.Area)
}

func (t *Tree) GetCenter() {
	var a struct {
		Center string
	}
	db.Table("tree_geoms").Select("ST_AsGeoJSON(ST_Centroid(geom)) as center").Where("latin_name = ?", t.LatinName).Scan(&a)

	t.Center = a.Center
	log.Println("Center:", t.Center)
}
