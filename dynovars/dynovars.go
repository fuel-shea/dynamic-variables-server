package dynovars

import (
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

	for _, featureType := range featureTypes {
		matchVal, found := featureMatches[featureType]
		if !found {
			matchVal = "any"
		} else {
			matchVal = matchVal.(string)
		}

		pipe := []bson.M{
			// match the features
			bson.M{
				"$match": bson.M{
					"game_id":  gameID,
					"rule_idx": bson.M{"$in": eligibleRules},
					"type":     featureType,
					"$or": []bson.M{
						// the value can be "any"
						bson.M{"value": "any"},
						// if the modifier is "=", the value can exactly match the feature
						bson.M{"$and": []bson.M{
							bson.M{"mod": "="}, bson.M{"value": matchVal}},
						},
						// if the modifier is ">", the value can be greater than the feature
						bson.M{"$and": []bson.M{
							bson.M{"mod": ">"}, bson.M{"value": bson.M{"$gt": matchVal}}},
						},
						// if the modifier is "<", the value can be less than the feature
						bson.M{"$and": []bson.M{
							bson.M{"mod": "<"}, bson.M{"value": bson.M{"$lt": matchVal}}},
						},
					},
				},
			},

			// make the entire unique and only include "rule_idx" field
			bson.M{
				"$group": bson.M{
					"_id": bson.M{"rule_idx": "$rule_idx"},
				},
			},

			// sort them by index/priority
			bson.M{
				"$sort": bson.M{"_id.rule_idx": 1},
			},

			// flatten it
			bson.M{
				"$project": bson.M{
					"rule_idx": "$_id.rule_idx",
					"_id":      0,
				},
			},
		}

		var ruleIdxRes []bson.M
		if err := dvs.FeatsColl.Pipe(pipe).All(&ruleIdxRes); err != nil {
			return blankReturnVal, err
		}
		newEligibleRules := make([]int, len(ruleIdxRes))
		for i, ruleIdx := range ruleIdxRes {
			newEligibleRules[i] = int(ruleIdx["rule_idx"].(float64))
		}
		eligibleRules = newEligibleRules
	}

	// TODO ensure rule exists
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
	if dvs.MgoDB == nil {
		dvs.MgoDB = dvs.MgoSess.DB("dynovars2")
	}
	if dvs.VarsColl == nil {
		dvs.VarsColl = dvs.MgoDB.C("variables")
	}
	if dvs.FeatsColl == nil {
		dvs.FeatsColl = dvs.MgoDB.C("features")
	}
	if dvs.FeatTypesColl == nil {
		dvs.FeatTypesColl = dvs.MgoDB.C("feat_types")
	}
	if dvs.VarTypesColl == nil {
		dvs.VarTypesColl = dvs.MgoDB.C("var_types")
	}
	return nil
}
