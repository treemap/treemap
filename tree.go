package main

import (
	"fmt"
	render "github.com/abhiyerra/gowebcommons/render"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Tree struct {
	Id         int64    `json:"id"`
	CommonName string   `json:"common_name"`
	LatinName  string   `json:"latin_name"`
	GeomData   []string `json:"geom",sql:"-"`
	Area       float64  `json:"area",sql:"-"`
	Center     string   `json:"center",sql:"-"`
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

func showTreesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	treeId, _ := strconv.ParseInt(vars["treeId"], 10, 64)

	tree := cache.Get("tree/"+vars["treeId"], func() interface{} {
		tree := Tree{Id: int64(treeId)}
		db.First(&tree)
		tree.GetGeodata()
		tree.GetArea()
		tree.GetCenter()

		return tree
	})

	render.RenderJson(w, tree)
}

func nearbyTreesHandler(w http.ResponseWriter, r *http.Request) {
	var trees []Tree

	longitude := r.URL.Query().Get("long")
	latitude := r.URL.Query().Get("lat")
	log.Println("Long:", longitude, "Lat:", latitude)

	err := db.Model(Tree{}).Select("distinct trees.id, trees.latin_name, trees.common_name").
		Joins(fmt.Sprintf("INNER JOIN tree_geoms ON tree_geoms.latin_name = trees.latin_name AND ST_DWithin(ST_GeomFromText('POINT(%s %s)' , 4326)::geography, tree_geoms.geom, 160934 , true)", longitude, latitude)).
		Order("trees.latin_name asc").Scan(&trees)

	if err != nil {
		log.Println(err)
	}

	render.RenderJson(w, trees)
}

func treesHandler(w http.ResponseWriter, r *http.Request) {
	trees := cache.Get("trees", func() interface{} {
		var trees []Tree

		err := db.Model(Tree{}).Select("id, latin_name, common_name").Scan(&trees)
		if err != nil {
			log.Println(err)
		}

		return trees
	})

	render.RenderJson(w, trees)
}
