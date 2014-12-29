package dynamicvariables

import (
	"errors"
	"github.com/fuel-shea/fuel-go-utils/fuelutils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DynoVarSource struct {
	MgoSess       *mgo.Session
	MgoDB         *mgo.Database
	FeatsColl     *mgo.Collection
	FeatTypesColl *mgo.Collection
	VarsColl      *mgo.Collection
	VarTypesColl  *mgo.Collection
}

func NewDynoVarSource() (DynoVarSource, error) {
	dvs := DynoVarSource{}
	err := dvs.Init()
	return dvs, err
}

func (dvs *DynoVarSource) VarsFromFeatures(featureMatches map[string]interface{}, gameID string) (map[string]interface{}, error) {
	blankReturnVal := make(map[string]interface{})

	var featTypeRes bson.M
	if err := dvs.FeatTypesColl.
		Find(bson.M{"game_id": gameID}).
		One(&featTypeRes); err != nil {
		return blankReturnVal, err
	}
	featureTypes := fuelutils.InterfaceArr2StringArr(featTypeRes["types"].([]interface{}))

	nRules, err := dvs.VarsColl.Count()
	if err != nil {
		return blankReturnVal, err
	}
	eligibleRules := make([]int, nRules)
	for i := range eligibleRules {
		eligibleRules[i] = i
	}

	pipe := PipeSkeleton(gameID)

	for _, featureType := range featureTypes {
		matchVal, found := featureMatches[featureType]
		if !found {
			matchVal = "any"
		} else {
			matchVal = matchVal.(string)
		}

		var ruleIdxRes []bson.M
		pipe.UpdateForLoop(eligibleRules, featureType, matchVal)
		if err := dvs.FeatsColl.Pipe(pipe.Pipe).All(&ruleIdxRes); err != nil {
			return blankReturnVal, err
		}
		newEligibleRules := make([]int, len(ruleIdxRes))
		for i, ruleIdx := range ruleIdxRes {
			newEligibleRules[i] = int(ruleIdx["rule_idx"].(float64))
		}
		eligibleRules = newEligibleRules

		if len(eligibleRules) == 1 {
			break
		}

		if len(eligibleRules) == 0 {
			return blankReturnVal, errors.New("No rules matched query")
		}
	}
	winningRuleIdx := eligibleRules[0]

	var winningRuleVars bson.M
	if err := dvs.VarsColl.Find(bson.M{"rule_idx": winningRuleIdx}).One(&winningRuleVars); err != nil {
		return blankReturnVal, err
	}

	var varTypesRes bson.M
	if err := dvs.VarTypesColl.Find(bson.M{"game_id": gameID}).One(&varTypesRes); err != nil {
		return blankReturnVal, err
	}
	varTypes := fuelutils.InterfaceArr2StringArr(varTypesRes["types"].([]interface{}))

	result := make(map[string]interface{})
	for _, varType := range varTypes {
		result[varType] = winningRuleVars[varType]
	}

	return result, nil
}

func (dvs *DynoVarSource) Init() error {
	if dvs.MgoSess == nil {
		sess, err := mgo.Dial("localhost")
		if err != nil {
			return err
		}
		dvs.MgoSess = sess
	}
	dvs.MgoDB = dvs.MgoSess.DB("dynamicvariables")
	dvs.VarsColl = dvs.MgoDB.C("variables")
	dvs.FeatsColl = dvs.MgoDB.C("features")
	dvs.FeatTypesColl = dvs.MgoDB.C("feat_types")
	dvs.VarTypesColl = dvs.MgoDB.C("var_types")
	return nil
}

func (dvs *DynoVarSource) CleanUp() {
	dvs.MgoSess.Close()
}
