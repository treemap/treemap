package main

import (
	gocache "github.com/abhiyerra/gowebcommons/cache"
	render "github.com/abhiyerra/gowebcommons/render"
	"github.com/coreos/go-etcd/etcd"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
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
		return zc.TableData(table, 50)
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

func seaRiseTableHandler(w http.ResponseWriter, r *http.Request) {

	z := cache.Get("sea_rise", func() interface{} {
		return AllSeaRise()
	})

	render.RenderJson(w, z)
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
	if len(os.Args) < 2 {
		log.Println("Need serve or generate")
		return
	}

	switch os.Args[1] {
	case "serve":
		r := mux.NewRouter()
		r.HandleFunc("/v1/sea_rise", seaRiseTableHandler).Methods("GET")
		r.HandleFunc("/v1/zipcodes/{zipcode}/{table}", zipcodeTableHandler).Methods("GET")
		r.HandleFunc("/v1/{table}/{resourceId}", showHandler).Methods("GET")
		r.HandleFunc("/v1/{table}", indexHandler).Methods("GET")
		r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

		http.Handle("/", r)
		http.ListenAndServe(":3001", nil)
	case "generate":
		if len(os.Args) < 3 {
			log.Println("What to generate? zipcodes")
			return
		}

		switch os.Args[2] {
		case "zipcodes":
			GenerateZipcodes()
		case "trees":
			GenerateTrees()
		case "parks":
			GenerateParks()
		default:
			GenerateZipcodeTable(os.Args[2], 50)
		}
	}
}
