package dynamicvariables

import (
	"errors"
	"github.com/fuel-shea/fuel-go-utils/fuelutils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type DynoVarFactory struct {
	MgoSess *mgo.Session
	DBName  string
}

func NewDynoVarFactory(DBHost, DBName string) (DynoVarFactory, error) {
	factory := DynoVarFactory{}
	if factory.MgoSess == nil {
		sess, err := mgo.DialWithTimeout(DBHost, time.Duration(1)*time.Second)
		if err != nil {
			return factory, err
		}
		factory.MgoSess = sess
	}
	factory.DBName = DBName
	return factory, nil
}

type DynoVarSource struct {
	MgoDB        *mgo.Database
	GameDataColl *mgo.Collection
	FeatsColl    *mgo.Collection
	VarsColl     *mgo.Collection
}

func (factory DynoVarFactory) NewDynoVarSource() DynoVarSource {
	dvSource := DynoVarSource{}
	sourceDB := factory.MgoSess.Copy().DB(factory.DBName)
	dvSource.VarsColl = sourceDB.C("variables")
	dvSource.FeatsColl = sourceDB.C("features")
	dvSource.GameDataColl = sourceDB.C("game_rule_data")
	return dvSource
}

func (dvSource *DynoVarSource) VarsFromFeatures(featureMatches map[string]interface{}, gameID string) (map[string]interface{}, error) {
	blankReturnVal := make(map[string]interface{})

	var gameDataRes bson.M
	if err := dvSource.GameDataColl.
		Find(bson.M{"game_id": gameID}).
		One(&gameDataRes); err != nil {
		return blankReturnVal, err
	}
	nRules := gameDataRes["num_rules"].(int)
	featTypes, err := fuelutils.InterfaceArr2StringArr(gameDataRes["feature_types"].([]interface{}))
	if err != nil {
		return blankReturnVal, err
	}
	varTypes, err := fuelutils.InterfaceArr2StringArr(gameDataRes["variable_types"].([]interface{}))
	if err != nil {
		return blankReturnVal, err
	}

	eligibleRules := make([]int, nRules)
	for i := range eligibleRules {
		eligibleRules[i] = i
	}

	pipe := PipeSkeleton(gameID)

	for _, featureType := range featTypes {
		matchVal, found := featureMatches[featureType]
		if !found {
			matchVal = "any"
		} else {
			matchVal = matchVal.(string)
		}

		var ruleIdxRes []bson.M
		pipe.UpdateForLoop(eligibleRules, featureType, matchVal)
		if err := dvSource.FeatsColl.Pipe(pipe.Pipe).All(&ruleIdxRes); err != nil {
			return blankReturnVal, err
		}
		newEligibleRules := make([]int, len(ruleIdxRes))
		for i, ruleIdx := range ruleIdxRes {
			newEligibleRules[i] = ruleIdx["rule_idx"].(int)
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
	if err := dvSource.VarsColl.Find(bson.M{"rule_idx": winningRuleIdx}).One(&winningRuleVars); err != nil {
		return blankReturnVal, err
	}

	result := make(map[string]interface{})
	for _, varType := range varTypes {
		result[varType] = winningRuleVars[varType]
	}

	return result, nil
}
