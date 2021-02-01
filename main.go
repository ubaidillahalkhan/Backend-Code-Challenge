package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

type Village struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	District string `json:"district"`
	Regency  string `json:"regency"`
	Province string `json:"province"`
}

func prime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bilangan := mux.Vars(r)["id"]

	n := new(big.Int)
	n, ok := n.SetString(bilangan, 10)
	if !ok {
		fmt.Println("SetString: error")
		return
	}
	if n.ProbablyPrime(0) {
		fmt.Fprintln(w, "is prime")
	} else {
		fmt.Fprintln(w, "is not prime")
	}
}

func getVillages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	// Tolong diganti user, password dan database sesuai dengan user dan password di lokal
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/database")
	if err != nil {
		panic(err.Error())
	}

	result, err := db.Query("select v.id as id, v.name as name, d.name as district, r.name as regency, p.name as province from villages v join districts d on v.district_id = d.id join regencies r on d.regency_id = r.id join provinces p on p.id = r.province_id where v.name = ?", params["name"])
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	var village Village
	for result.Next() {
		err := result.Scan(&village.ID, &village.Name, &village.District, &village.Regency, &village.Province)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(village)
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/prime/{id}", prime).Methods("GET")
	router.HandleFunc("/villages/{name}", getVillages).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
