package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go-expt/gdv"
	"net/http"
)

func main() {
	gameID := "53e256d96170706e28063201"
	gdvs, err := gdv.ReadDynamicsMap(gameID)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.Path("/features").
		HandlerFunc(FeaturesHandler(gdvs)).
		Methods("POST").
		Headers("Content-Type", "application/json").
		Name("setFeatures")

	http.Handle("/", r)
	http.ListenAndServe(":3031", nil)
}

func FeaturesHandler(gdvs gdv.FeatureMap) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var jsonData map[string]interface{}
		decodeErr := json.NewDecoder(r.Body).Decode(&jsonData)
		if decodeErr != nil {
			sendError(w, "Invalid JSON")
			return
		}

		params := buildParams(jsonData, gdvs)
		gVars := gdv.ResultFromFeatures(params, gdvs)
		sendSuccess(w, gVars)
		return
	}
}

func buildParams(jsonData map[string]interface{}, gdvs gdv.FeatureMap) map[string]string {
	firstRow := gdvs[0]["Criteria"]
	features := make([]string, len(firstRow))
	featIdx := 0
	for critKey, _ := range firstRow {
		features[featIdx] = critKey
		featIdx++
	}

	params := map[string]string{}
	for _, feat := range features {
		if val, ok := jsonData[feat]; ok {
			params[feat] = val.(string)
		}
	}

	return params
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
