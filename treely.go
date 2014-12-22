package main

import (
	"fmt"
	gocache "github.com/abhiyerra/gowebcommons/cache"
	render "github.com/abhiyerra/gowebcommons/render"
	"github.com/coreos/go-etcd/etcd"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

const (
	DatabaseUrlKey = "/treemap/database_url"
)

var (
	db    gorm.DB
	cache gocache.Cache
)

type Hydrology struct {
	Name     string `json:"name"`
	GeomData string `json:"geom"`
}

func nearbyHydrologyHandler(hydroType string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var hydrology []Hydrology

		longitude := r.URL.Query().Get("long")
		latitude := r.URL.Query().Get("lat")
		log.Println("Hydro Type:", hydroType, "Long:", longitude, "Lat:", latitude)

		err := db.Table(hydroType).
			Select("ST_AsGeoJSON(ST_CollectionExtract(geom, 3)) as geom_data, name").
			Where(fmt.Sprintf("ST_DWithin(ST_GeomFromText('POINT(%s %s)' , 4326)::geography, geom, 160934, true)", longitude, latitude)). // Within 100 miles -> 160934 meters
			Scan(&hydrology)
		if err != nil {
			log.Println(err)
		}

		render.RenderJson(w, hydrology)
	}
}

func hydrologyHandler(hydroType string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		hydrology := cache.Get(hydroType, func() interface{} {
			var hydrology []Hydrology
			db.Table(hydroType).Select("ST_AsGeoJSON(ST_CollectionExtract(geom, 3)) as geom_data, name").Scan(&hydrology)

			return hydrology
		})

		render.RenderJson(w, hydrology)
	}
}

type Zipcode struct {
	Number   string `json:"number"`
	GeomData string `json:"geom"`
	Center   string `json:"center"`
}

func showZipCodeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	zipcode := vars["zipcode"]

	z := cache.Get("zipcode/"+zipcode, func() interface{} {
		z := Zipcode{Number: zipcode}
		db.Select("geoid10 as number, ST_AsGeoJSON(ST_Centroid(geom)) as center, ST_AsGeoJSON(ST_CollectionExtract(geom, 3)) as geom_data").Where("geoid10 = ?", zipcode).First(&z)

		return z
	})

	render.RenderJson(w, z)
}

func dbConnect(databaseUrl string) {
	log.Println("Connecting to database:", databaseUrl)
	var err error
	db, err = gorm.Open("postgres", databaseUrl)
	if err != nil {
		log.Println(err)
	}
	db.LogMode(true)
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
	r.HandleFunc("/trees/nearby", nearbyTreesHandler).Methods("GET")
	r.HandleFunc("/trees/{treeId}", showTreesHandler).Methods("GET")
	r.HandleFunc("/trees", treesHandler).Methods("GET")

	r.HandleFunc("/parks/nearby", nearbyParksHandler).Methods("GET")
	r.HandleFunc("/parks", parksHandler).Methods("GET")

	r.HandleFunc("/lakes/nearby", nearbyHydrologyHandler("lakes")).Methods("GET")
	r.HandleFunc("/lakes", hydrologyHandler("lakes")).Methods("GET")

	r.HandleFunc("/rivers/nearby", nearbyHydrologyHandler("rivers")).Methods("GET")
	r.HandleFunc("/rivers", hydrologyHandler("rivers")).Methods("GET")

	r.HandleFunc("/zipcode/{zipcode}", showZipCodeHandler).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":3001", nil)
}
