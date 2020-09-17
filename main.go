package main

import (
	"log"
	"strconv"
	"image"
	"image/png"
	"github.com/cufee/am-stats/stats"
	"github.com/cufee/am-stats/render"

	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
)

type request struct {
	PlayerID	int		`json:"player_id"`
	Realm		string	`json:"realm"`
	Days		int		`json:"days"`
}

func handler() {
	log.Println("Starting webserver on", 6969)
	hostPORT := ":" + strconv.Itoa(6969)

	myRouter := mux.NewRouter().StrictSlash(true)
	// myRouter.HandleFunc("/clans", updateClanActivity)
	myRouter.HandleFunc("/player", handlePlayerRequest).Methods("GET")

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
	img, err := render.ImageFromStats(export)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	repondWithImage(w, http.StatusOK, img)
}

func main() {
	handler()
}