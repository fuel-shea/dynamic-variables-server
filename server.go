package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	r.Path("/features").
		HandlerFunc(FeaturesHandler).
		Methods("POST").
		Headers("Content-Type", "application/json").
		Name("setFeatures")

	http.Handle("/", r)
	http.ListenAndServe(":3030", nil)
}

func FeaturesHandler(w http.ResponseWriter, r *http.Request) {
	var jsonData interface{}
	decodeErr := json.NewDecoder(r.Body).Decode(&jsonData)
	if decodeErr != nil {
		sendError(w, "Invalid JSON")
		return
	}

	gameVars, gvErr := getGameVars(jsonData)
	if gvErr != nil {
		sendError(w, "Error parsing CSV")
		return
	}

	sendSuccess(w, gameVars)
	return
}

func getGameVars(data interface{}) (map[string]interface{}, error) {
	res := map[string]interface{}{
		"max_score":     "1000",
		"whammy_chance": "5",
	}
	return res, nil
}

func sendSuccess(w http.ResponseWriter, data map[string]interface{}) {
	succObj := SuccRespObj{
		Result: data,
	}
	succObj.Init()

	json.NewEncoder(w).Encode(succObj)
}

func sendError(w http.ResponseWriter, msg string) {
	errObj := ErrRespObj{
		ErrorObj: ErrDetailsObj{
			Message:   msg,
			ErrorCode: "ERROR",
		},
	}
	errObj.Init()

	json.NewEncoder(w).Encode(errObj)
}

type SuccRespObj struct {
	Success bool                   `json:"success"`
	Result  map[string]interface{} `json:"result"`
}

func (sro *SuccRespObj) Init() {
	sro.Success = true
}

type ErrRespObj struct {
	Success  bool          `json:"success"`
	ErrorObj ErrDetailsObj `json:"error"`
}

func (ero *ErrRespObj) Init() {
	ero.Success = false
}

type ErrDetailsObj struct {
	Message   string `json:"message"`
	ErrorCode string `json:"errorcode"`
}
