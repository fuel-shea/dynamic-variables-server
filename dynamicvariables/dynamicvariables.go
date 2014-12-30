package dynamicvariables

import (
	"errors"
	"github.com/fuel-shea/fuel-go-utils/fuelutils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DynoVarFactory struct {
	MgoSess *mgo.Session
	DBName  string
}

func NewDynoVarFactory(DBHost, DBName string) (DynoVarFactory, error) {
	factory := DynoVarFactory{}
	err := factory.Init(DBHost, DBName)
	return factory, err
}

func (factory *DynoVarFactory) Init(DBHost, DBName string) error {
	if factory.MgoSess == nil {
		sess, err := mgo.Dial(DBHost)
		if err != nil {
			return err
		}
		factory.MgoSess = sess
	}
	factory.DBName = DBName
	return nil
}

func (factory DynoVarFactory) NewDynoVarSource() (DynoVarSource, error) {
	dvSource := DynoVarSource{}
	sourceDB := factory.MgoSess.Copy().DB(factory.DBName)
	err := dvSource.Init(*sourceDB)
	return dvSource, err
}

type DynoVarSource struct {
	MgoDB        *mgo.Database
	GameDataColl *mgo.Collection
	FeatsColl    *mgo.Collection
	VarsColl     *mgo.Collection
}

func (dvSource *DynoVarSource) VarsFromFeatures(featureMatches map[string]interface{}, gameID string) (map[string]interface{}, error) {
	blankReturnVal := make(map[string]interface{})

	var gameDataRes bson.M
	if err := dvSource.GameDataColl.
		Find(bson.M{"game_id": gameID}).
		One(&gameDataRes); err != nil {
		return blankReturnVal, err
	}
	nRules := int(gameDataRes["num_rules"].(float64))
	featTypes := fuelutils.InterfaceArr2StringArr(gameDataRes["feature_types"].([]interface{}))
	varTypes := fuelutils.InterfaceArr2StringArr(gameDataRes["variable_types"].([]interface{}))

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
	if err := dvSource.VarsColl.Find(bson.M{"rule_idx": winningRuleIdx}).One(&winningRuleVars); err != nil {
		return blankReturnVal, err
	}

	result := make(map[string]interface{})
	for _, varType := range varTypes {
		result[varType] = winningRuleVars[varType]
	}

	return result, nil
}

func (dvSource *DynoVarSource) Init(mgoDB mgo.Database) error {
	dvSource.VarsColl = mgoDB.C("variables")
	dvSource.FeatsColl = mgoDB.C("features")
	dvSource.GameDataColl = mgoDB.C("game_rule_data")
	return nil
}
