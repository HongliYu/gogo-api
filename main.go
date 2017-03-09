package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gogotattoo/common/models"
	"github.com/gogotattoo/gogo-api/artwork"
	"github.com/gorilla/mux"
)

var tattoos []models.Tattoo
var hennas []models.Henna
var piercing []models.Piercing
var designs []models.Design

var artistWorks = make(map[string]models.Artworks)

// ArtistArtworkRefresh returns the list of all tattoos
// TODO: add a timer, allow only every 5-10 mins
func ArtistArtworkRefresh(artType string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		artistName := mux.Vars(r)["name"]
		artistWorks[artistName+"/"+artType] = artwork.Refresh(artistName, artType)
		json.NewEncoder(w).Encode(artistWorks[artistName+"/"+artType])
	}
}

// ArtistArtwork returns the list of all artists' tattoos actually published to git repos
func ArtistArtwork(artType string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(artistWorks[mux.Vars(r)["name"]+"/"+artType])
	}
}

var refreshLastCalled = time.Now()

// ArtistArtworkRefreshAll refreshes all known art types for given artistName
// returns {"tattoo: {...}, design: {}, etc..."} json
func ArtistArtworkRefreshAll(w http.ResponseWriter, r *http.Request) {
	artistName := mux.Vars(r)["name"]
	result := make(map[string]models.Artworks)
	for _, artType := range []string{"tattoo", "henna", "piercing", "design"} {
		if time.Now().After(refreshLastCalled.Add(time.Minute * 3)) {
			artistWorks[artistName+"/"+artType] = artwork.Refresh(artistName, artType)
		}
		result[artType] = artistWorks[artistName+"/"+artType]
	}
	refreshLastCalled = time.Now()
	json.NewEncoder(w).Encode(result)
}

// Tattoo shows info on a single tattoo work by id
func Tattoo(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, item := range tattoos {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(models.NewTattoo("", "brr", "", ""))
}

// TattooToml shows info of a single tattoo work by id in toml format
func TattooToml(w http.ResponseWriter, req *http.Request) {
	toml.NewEncoder(w).Encode(tattoos)
}

// Tattoos returns the list of all tattoos
func Tattoos(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(tattoos)
}

var myClient = &http.Client{Timeout: 100 * time.Second}

func getJSON(url string, target interface{}) (io.ReadCloser, error) {
	r, err := myClient.Get(url)
	if err != nil {
		return nil, err
	}
	//defer r.Body.Close()
	t := r.Body
	return t, json.NewDecoder(r.Body).Decode(target)
}

// CreateTattoo adds a new tattoo to the memory
func CreateTattoo(w http.ResponseWriter, req *http.Request) {
	log.Println("POST /tattoo")
	params := mux.Vars(req)
	defer req.Body.Close()
	var tat models.Tattoo
	err := json.NewDecoder(req.Body).Decode(&tat)
	log.Println("TITLE\n", tat.Title)
	if err != nil {
		log.Println("ERROR\n", err)
		json.NewEncoder(w).Encode(err)
		return
	}
	tat.ID = params["id"]
	tattoos = append(tattoos, tat)
	m, _ := json.Marshal(tat)
	log.Println("TATTOO\n", string(m)+"\n")
	json.NewEncoder(w).Encode(tat)
}

// DeleteTattoo deletes a tattoo by id from the memory
func DeleteTattoo(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range tattoos {
		if item.ID == params["id"] {
			tattoos = append(tattoos[:index], tattoos[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(tattoos)
}

// Hennas returns the list of all hennas
func Hennas(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(hennas)
}

// CreateHenna adds a new henna object
func CreateHenna(w http.ResponseWriter, req *http.Request) {
	log.Println("POST /henna")
	params := mux.Vars(req)
	defer req.Body.Close()
	var hen models.Henna
	err := json.NewDecoder(req.Body).Decode(&hen)
	log.Println("TITLE\n", hen.Title)
	if err != nil {
		log.Println("ERROR\n", err)
		json.NewEncoder(w).Encode(err)
		return
	}
	hen.ID = params["id"]
	hennas = append(hennas, hen)
	m, _ := json.Marshal(hen)
	log.Println("HENNA\n", string(m)+"\n")
	json.NewEncoder(w).Encode(hen)
}

// CreateDesign adds a new design object
func CreateDesign(w http.ResponseWriter, req *http.Request) {
	log.Println("POST /design")
	params := mux.Vars(req)
	defer req.Body.Close()
	var des models.Design
	err := json.NewDecoder(req.Body).Decode(&des)
	log.Println("TITLE\n", des.Title)
	if err != nil {
		log.Println("ERROR\n", err)
		json.NewEncoder(w).Encode(err)
		return
	}
	des.ID = params["id"]
	designs = append(designs, des)
	m, _ := json.Marshal(des)
	log.Println("DESIGN\n", string(m)+"\n")
	json.NewEncoder(w).Encode(des)
}

// Designs returns the list of all designs
func Designs(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(designs)
}

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

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/tattoo", Tattoos).Methods("GET")
	//router.HandleFunc("/tattoo/{id}", Tattoo).Methods("GET")
	router.HandleFunc("/tattoo.toml", TattooToml).Methods("GET")
	router.HandleFunc("/tattoo/{id}", CreateTattoo).Methods("POST")
	router.HandleFunc("/tattoo/{id}", DeleteTattoo).Methods("DELETE")

	router.HandleFunc("/henna", Hennas).Methods("GET")
	router.HandleFunc("/henna/{id}", CreateHenna).Methods("POST")

	router.HandleFunc("/design", Designs).Methods("GET")
	router.HandleFunc("/design/{id}", CreateDesign).Methods("POST")

	router.HandleFunc("/piercing", Piercing).Methods("GET")
	router.HandleFunc("/piercing/{id}", CreatePiercing).Methods("POST")

	for _, t := range []string{"tattoo", "henna", "piercing", "design"} {
		router.HandleFunc("/"+t+"/{name}", ArtistArtwork(t)).Methods("GET")
		router.HandleFunc("/"+t+"/{name}/refresh", ArtistArtworkRefresh(t)).Methods("GET")
	}

	router.HandleFunc("/all/{name}/refresh", ArtistArtworkRefreshAll).Methods("GET")

	for _, artistName := range []string{"gogo", "aid", "xizi"} {
		for _, artType := range []string{"tattoo", "henna", "piercing", "design"} {
			artistWorks[artistName+"/"+artType] = artwork.Refresh(artistName, artType)
		}
	}

	log.Fatal(http.ListenAndServe(":12345", Log(router)))
}

// Log prints basic http request info
func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
