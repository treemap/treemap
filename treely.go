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
	r.HandleFunc("/{table}/{resourceId}", showHandler).Methods("GET")
	r.HandleFunc("/trees", treesHandler).Methods("GET")

	r.HandleFunc("/parks", parksHandler).Methods("GET")
	r.HandleFunc("/lakes", hydrologyHandler("lakes")).Methods("GET")
	r.HandleFunc("/rivers", hydrologyHandler("rivers")).Methods("GET")

	r.HandleFunc("/zipcodes/{zipcode}/{table}", zipcodeTableHandler).Methods("GET")
	r.HandleFunc("/{table}/{resourceId}", showHandler).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":3001", nil)
}
