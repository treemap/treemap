package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	db gorm.DB
)

type Tree struct {
	Id         int64    `json:"id"`
	CommonName string   `json:"common_name"`
	LatinName  string   `json:"latin_name"`
	GeomData   []string `json:"geom",sql:"-"`
	Area       float64  `json:"area",sql:"-"`
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

type NationalPark struct {
	UnitName string `json:"name"`
	GeomData string `json:"geom"`
}

func renderJson(w http.ResponseWriter, page interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	b, err := json.Marshal(page)
	if err != nil {
		log.Println("error:", err)
		fmt.Fprintf(w, "")
	}

	w.Write(b)
}

func showTreesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	treeId, _ := strconv.ParseInt(vars["treeId"], 10, 64)

	tree := Tree{Id: int64(treeId)}
	db.First(&tree)
	tree.GetGeodata()
	tree.GetArea()

	renderJson(w, tree)
}

func treesHandler(w http.ResponseWriter, r *http.Request) {
	var trees []Tree

	err := db.Model(Tree{}).Select("id, latin_name, common_name").Scan(&trees)
	if err != nil {
		log.Println(err)
	}

	renderJson(w, trees)
}

func parksHandler(w http.ResponseWriter, r *http.Request) {
	var parks []NationalPark
	db.Model(NationalPark{}).Select("ST_AsGeoJSON(ST_CollectionExtract(geom, 3)) as geom_data, unit_name").Scan(&parks)

	renderJson(w, parks)
}

func init() {
	databaseUrl := os.Getenv("TREELY_DATABASE_URL")
	if databaseUrl == "" {
		databaseUrl = "user=ayerra dbname=treely_development sslmode=disable"
	}

	log.Println("Database:", databaseUrl)

	var err error
	db, err = gorm.Open("postgres", databaseUrl)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	r := mux.NewRouter()
	// r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/trees/{treeId}", showTreesHandler)
	r.HandleFunc("/trees", treesHandler)
	r.HandleFunc("/parks", parksHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.Handle("/", r)
	http.ListenAndServe(":3001", nil)
}
