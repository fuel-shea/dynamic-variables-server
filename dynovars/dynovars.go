package dynovars

import (
	"bitbucket.org/kardianos/osext"
	"encoding/csv"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"math"
	"os"
	"path"
)

type DynoVarSource struct {
	RuleSet RuleSet
	MgoSess *mgo.Session
	MgoDB   *mgo.Database
	MgoColl *mgo.Collection
}

type RuleSet struct {
	GameID        string
	FeatureNames  []string
	VariableNames []string
	Rules         []Rule
}

type Rule struct {
	Features  map[string]Feature
	Variables map[string]interface{}
}

func NewRule() Rule {
	return Rule{
		Features:  make(map[string]Feature),
		Variables: make(map[string]interface{}),
	}
}

type Feature struct {
	Value string
	Exp   string
}

func NewDynoVarSource() (DynoVarSource, error) {
	dvs := DynoVarSource{}
	err := dvs.Init()
	return dvs, err
}

func (dvs *DynoVarSource) VarsFromFeatures(featureMatch map[string]string) (map[string]interface{}, error) {

	// go and find the first row that matches this criteria
	result := map[string]interface{}{}

	ruleSet, err := dvs.retrieveRules()
	if err != nil {
		return nil, err
	}
	dvs.RuleSet = ruleSet

	// run through the map for each critera and match against the incoming information
	for _, rule := range dvs.RuleSet.Rules {
		matched := true
		for featureName, feature := range rule.Features {
			matchValue := featureMatch[featureName]

			if feature.Value != "any" {
				if matchValue == "" {
					matched = false
					break
				} else {
					// use the expression to perform the test
					switch feature.Exp {
					case "<":
						if matchValue >= feature.Value {
							matched = false
							break
						}
					case ">":
						if matchValue <= feature.Value {
							matched = false
							break
						}
					case "=":
						if matchValue != feature.Value {
							matched = false
							break
						}
					}
				}
			}
		}
		if matched {
			result = rule.Variables
			break
		}
	}

	return result, nil
}

func (dvs *DynoVarSource) Init() error {
	gameID := "53e256d96170706e28063201" // TODO move into proper place

	ruleSet := RuleSet{GameID: gameID}

	varStartIdx := math.MaxUint32
	rawCSVData, err := readCSVByGameId(gameID)
	if err != nil {
		return err
	}
	for rowIdx, row := range rawCSVData {
		if rowIdx == 0 {
			for colIdx, value := range row {
				if value == "|" {
					varStartIdx = colIdx
				}
				if colIdx > varStartIdx {
					// this is a variable
					ruleSet.VariableNames = append(ruleSet.VariableNames, value)
				} else if colIdx < varStartIdx {
					if colIdx%2 == 0 {
						ruleSet.FeatureNames = append(ruleSet.FeatureNames, value)
					}
				}
			}

		} else {
			rule := NewRule()
			for colIdx, value := range row {
				if colIdx > varStartIdx {
					// this is a variable
					variableName := ruleSet.VariableNames[colIdx-varStartIdx-1]
					rule.Variables[variableName] = value

				} else if colIdx < varStartIdx {
					featureName := ruleSet.FeatureNames[colIdx/2]
					var feature Feature
					if _, wasFound := rule.Features[featureName]; wasFound {
						feature = rule.Features[featureName]
					} else {
						feature = Feature{}
					}

					if colIdx%2 == 0 {
						// this is a feature value
						feature.Value = value
					} else {
						// this is a feature modifier
						feature.Exp = value
					}
					rule.Features[featureName] = feature
				}
			}

			ruleSet.Rules = append(ruleSet.Rules, rule)
		}
	}

	if err := dvs.initDB(); err != nil {
		fmt.Println(err) // TMP
		return err
	}

	dvs.RuleSet = ruleSet
	if err := dvs.persistRules(); err != nil {
		fmt.Println(err) // TMP
		dvs.RuleSet = RuleSet{}
		return err
	}

	return nil
}

func readCSVByGameId(gameId string) ([][]string, error) {
	execDir, err := osext.ExecutableFolder()
	if err != nil {
		return nil, err
	}
	csvPath := path.Join(execDir, gameId+".csv")

	csvfile, err := os.Open(csvPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1 // see the Reader struct information below

	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return rawCSVdata, nil
}

func (dvs *DynoVarSource) initDB() error {
	if dvs.MgoSess == nil {
		sess, err := mgo.Dial("localhost")
		if err != nil {
			return err
		}
		dvs.MgoSess = sess
	}
	if dvs.MgoDB == nil {
		dvs.MgoDB = dvs.MgoSess.DB("dynovarsDB")
	}
	if dvs.MgoColl == nil {
		dvs.MgoColl = dvs.MgoDB.C("dynovarsColl")
	}
	return nil
}

func (dvs DynoVarSource) persistRules() error {
	return dvs.MgoColl.Insert(dvs.RuleSet)
}

func (dvs DynoVarSource) retrieveRules() (RuleSet, error) {
	result := RuleSet{}
	err := dvs.MgoColl.Find(bson.M{}).One(&result)
	return result, err
}
