package main

import (
	"dynamic-variables-server/dynamicvariables"
	"encoding/json"
	"fmt"
	"github.com/fuel-shea/fuel-go-utils/fuelconfig"
	"github.com/fuel-shea/fuel-go-utils/fuelresponder"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	config, err := fuelconfig.CreateConfig("dynamic-variables-server")
	if err != nil {
		panic(err)
	}
	dvFact, err := dynamicvariables.NewDynoVarFactory(config.DBHost, config.DBName)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.Path("/features").
		HandlerFunc(FeaturesHandler(dvFact)).
		Methods("POST").
		Name("setFeatures")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":3030", nil))
}

func FeaturesHandler(dvFact dynamicvariables.DynoVarFactory) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		dvSrc, err := dvFact.NewDynoVarSource()
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
}
