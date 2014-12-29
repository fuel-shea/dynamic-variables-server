package main

import (
	"dynamic-variables-server/dynamicvariables"
	"encoding/json"
	"fmt"
	"github.com/fuel-shea/fuel-go-utils/fuelresponder"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.Path("/features").
		HandlerFunc(FeaturesHandler).
		Methods("POST").
		Name("setFeatures")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":3030", nil))
}

func FeaturesHandler(w http.ResponseWriter, r *http.Request) {
	dvSrc, err := dynamicvariables.NewDynoVarSource()
	if err != nil {
		fuelresponder.SendError(w, fuelresponder.ErrTypes["empty_result"])
		return
	}

	var jsonData map[string]interface{}
	decodeErr := json.NewDecoder(r.Body).Decode(&jsonData)
	if decodeErr != nil {
		fuelresponder.SendError(w, fuelresponder.ErrTypes["invalid_request"])
		return
	}

	gameIDJSON, found := jsonData["game_id"]
	if !found {
		fuelresponder.SendError(w, fuelresponder.ErrTypes["invalid_request"])
		return
	}
	gameID := gameIDJSON.(string)

	dVars, err := dvSrc.VarsFromFeatures(jsonData, gameID)
	if err != nil {
		fmt.Println(err)
		fuelresponder.SendError(w, fuelresponder.ErrTypes["general_error"])
		return
	}

	fuelresponder.SendSuccess(w, dVars)
	return
}
