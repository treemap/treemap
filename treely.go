package main

import (
	"encoding/json"
	gocache "github.com/abhiyerra/gowebcommons/cache"
	render "github.com/abhiyerra/gowebcommons/render"
	"github.com/coreos/go-etcd/etcd"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	DatabaseUrlKey = "/treemap/database_url"
)

var (
	db    gorm.DB
	cache gocache.Cache
)

func dbConnect(databaseUrl string) {
	log.Println("Connecting to database:", databaseUrl)
	var err error
	db, err = gorm.Open("postgres", databaseUrl)
	if err != nil {
		log.Println(err)
	}
	db.LogMode(true)
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

func showHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]
	resource := vars["resourceId"]

	z := cache.Get(table+"/"+resource, func() interface{} {
		switch table {
		case "zipcodes":
			z := Zipcode{Number: resource}
			z.GetInfo()

			return z
		case "trees":
			treeId, _ := strconv.ParseInt(resource, 10, 64)

			tree := Tree{Id: int64(treeId)}
			db.First(&tree)
			tree.GetGeodata()
			tree.GetArea()
			tree.GetCenter()

			return tree
		}

		return nil
	})

	render.RenderJson(w, z)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]

	z := cache.Get(table, func() interface{} {
		switch table {
		case "trees":
			var trees []Tree

			err := db.Model(Tree{}).Select("id, latin_name, common_name").Scan(&trees)
			if err != nil {
				log.Println(err)
			}

			return trees
		case "parks":
			var parks []NationalPark
			db.Model(NationalPark{}).Select("ST_AsGeoJSON(ST_CollectionExtract(geom, 3)) as geom_data, unit_name, unit_code").Scan(&parks)
			return parks
		case "lakes":
			var hydrology []Hydrology
			db.Table("lakes").Select("ST_AsGeoJSON(ST_CollectionExtract(geom, 3)) as geom_data, name").Scan(&hydrology)

			return hydrology

		case "rivers":
			var hydrology []Hydrology
			db.Table("rivers").Select("ST_AsGeoJSON(ST_CollectionExtract(geom, 3)) as geom_data, name").Scan(&hydrology)

			return hydrology

		}

		return nil
	})

	render.RenderJson(w, z)
}

func GenerateZipcodesOnS3() {
	log.Println("Gathering Zipcodes")
	zipcodes := AllZipcodes()

	log.Println("Got Zipcodes")
	for i := range zipcodes {
		log.Println("Writing", zipcodes[i].Number)
		zipcodes[i].GetInfo()

		os.MkdirAll("static/data/zipcodes", os.ModeDir|os.ModePerm)

		b, err := json.Marshal(zipcodes[i])
		if err != nil {
			log.Println("error:", err)
		}

		// Write the file
		err = ioutil.WriteFile("static/data/zipcodes/"+zipcodes[i].Number+".json", b, 0644)

		// for _, k := range []string{"trees", "parks", "lakes", "rivers"} {
		// 	for _, j := range []string{"50", "100", "200"} {

		// 		// Mkdir if it doesn't exist
		// 		os.MkdirAll("data/"+k+"/"+j, os.ModeDir)

		// 		// Get the Tree information
		// 		trees := zipcodes[i].Trees(j)
		// 		b, err := json.Marshal(trees)
		// 		if err != nil {
		// 			log.Println("error:", err)
		// 		}

		// 		// Generate Parks

		// 		// Write the file to disk.
		// 	}
		// }
	}
}

func init() {
	etcdHosts := os.Getenv("ETCD_HOSTS")
	if etcdHosts == "" {
		etcdHosts = "http://127.0.0.1:4001"
	}

	etcdClient := etcd.NewClient([]string{etcdHosts})

	resp, err := etcdClient.Get(DatabaseUrlKey, false, false)
	if err != nil {
		panic(err)
	}

	databaseUrl := resp.Node.Value
	dbConnect(databaseUrl)

	cache = gocache.New()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/zipcodes/{zipcode}/{table}", zipcodeTableHandler).Methods("GET")
	r.HandleFunc("/{table}/{resourceId}", showHandler).Methods("GET")
	r.HandleFunc("/{table}", indexHandler).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":3001", nil)
}
