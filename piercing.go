package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/gogotattoo/common/models"
	"github.com/gorilla/mux"
)

// CreatePiercing adds a new design artwork
func CreatePiercing(w http.ResponseWriter, req *http.Request) {
	log.Println("POST /piercing")
	params := mux.Vars(req)
	defer req.Body.Close()
	var per models.Piercing
	err := json.NewDecoder(req.Body).Decode(&per)
	log.Println("TITLE\n", per.Title)
	if err != nil {
		log.Println("ERROR\n", err)
		json.NewEncoder(w).Encode(err)
		return
	}
	per.ID = params["id"]
	piercing = append(piercing, per)
	m, _ := json.Marshal(per)
	log.Println("PIERCING\n", string(m)+"\n")
	json.NewEncoder(w).Encode(per)
}

// Piercing returns the list of all piercing works
func Piercing(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(piercing)
}

// PiercingToml shows info of a single piercing work by id in toml format
func PiercingToml(w http.ResponseWriter, req *http.Request) {
	toml.NewEncoder(w).Encode(piercing[len(piercing)-1])
}

// DeletePiercing deletes a piercing by id from the posted works in memory
func DeletePiercing(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range piercing {
		if item.ID == params["id"] {
			piercing = append(piercing[:index], piercing[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(piercing)
}
