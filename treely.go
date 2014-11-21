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
}

func (t *Tree) Geodata() (coordinateJson []string) {
	rows, err := db.Table("tree_geoms").Select("ST_AsGeoJSON(ST_CollectionExtract(geom, 3)) as geom2").Where("latin_name = ?", t.LatinName).Rows()
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		var geodata string
		rows.Scan(&geodata)
		coordinateJson = append(coordinateJson, geodata)
	}

	t.GeomData = coordinateJson

	return coordinateJson
}

func FindTree(latinName string) (t *Tree) {

	return nil
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

func init() {
	treelyDatabaseUrl := os.Getenv("TREELY_DATABASE_URL")

	log.Println("Database:", treelyDatabaseUrl)

	var err error
	db, err = gorm.Open("postgres", treelyDatabaseUrl)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	r := mux.NewRouter()
	// r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/trees/{treeId}", showTreesHandler)
	r.HandleFunc("/trees", treesHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.Handle("/", r)
	http.ListenAndServe(":3001", nil)
}
