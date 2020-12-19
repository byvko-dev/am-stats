package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"strconv"

	"github.com/cufee/am-stats/config"
	"github.com/cufee/am-stats/render"
	"github.com/cufee/am-stats/stats"
	externalapis "github.com/cufee/am-stats/wargamingapi"
	"github.com/fogleman/gg"

	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type request struct {
	PlayerID  int    `json:"player_id"`
	Premium   bool   `json:"premium"`
	Verified  bool   `json:"verified"`
	Realm     string `json:"realm"`
	Days      int    `json:"days"`
	Sort      string `json:"sort_key"`
	TankLimit int    `json:"detailed_limit"`
	BgURL     string `json:"bg_url"`
}

const currentBG string = "bg_event.jpg"

func handler() {
	log.Println("Starting webserver on", config.APIport)
	hostPORT := ":" + strconv.Itoa(config.APIport)

	myRouter := mux.NewRouter().StrictSlash(true)
	// myRouter.HandleFunc("/clans", updateClanActivity)
	myRouter.HandleFunc("/player", handlePlayerRequest).Methods("GET")
	myRouter.HandleFunc("/stats", handleStatsRequest).Methods("GET")

	log.Fatal(http.ListenAndServe(hostPORT, myRouter))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	log.Println("Request - ", code)
}

func repondWithImage(w http.ResponseWriter, code int, image image.Image) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(code)
	png.Encode(w, image)
	log.Println("Request - ", code)
}

func handlePlayerRequest(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in handlePlayerRequest", r)
		}
	}()

	var request request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	export, err := stats.ExportSessionAsStruct(request.PlayerID, request.Realm, request.Days)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if export.PlayerDetails == (externalapis.PlayerProfile{}) || export.PlayerDetails.Name == "" {
		log.Printf("%+v", request)
		respondWithError(w, http.StatusNotFound, "bad player data")
		return
	}
	if request.TankLimit == 0 {
		request.TankLimit = 10
	}

	// Get bg Image
	var bgImage image.Image
	if request.BgURL != "" {
		response, _ := http.Get(request.BgURL)
		if response != nil {
			bgImage, _, err = image.Decode(response.Body)
			defer response.Body.Close()
		} else {
			log.Printf("bad bg image for %v", request.PlayerID)
			err = fmt.Errorf("bad bg image")
		}
	}
	if err != nil || request.BgURL == "" {
		bgImage, err = gg.LoadImage(config.AssetsPath + currentBG)
		if err != nil {
			log.Println(err)
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
	}

	img, err := render.ImageFromStats(export, request.Sort, request.TankLimit, request.Premium, request.Verified, bgImage)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	repondWithImage(w, http.StatusOK, img)
}

func handleStatsRequest(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r)
		}
	}()
	var request request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	export, err := stats.ExportSessionAsStruct(request.PlayerID, request.Realm, request.Days)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if export.PlayerDetails.Name == "" {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, export)
}

func main() {
	handler()
}
