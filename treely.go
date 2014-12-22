package main

import (
	gocache "github.com/abhiyerra/gowebcommons/cache"
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
