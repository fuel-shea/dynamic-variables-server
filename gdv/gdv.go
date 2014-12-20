package gdv

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
)

func RunCSVReader() {

	fmt.Printf("Generating FeatureMap\n\n")

	gameFeatureMap := readDynamicsMap("53e256d96170706e28063201")

	fmt.Printf("Testing FeatureMap\n\n")

	// run a test
	// match against the following values and return the correct lines

	testValues1 := map[string]string{
		"Country": "CA",
	}
	fmt.Printf("testFeatures %s\n", testValues1)
	testResult1 := resultFromFeatures(testValues1, gameFeatureMap)
	fmt.Printf("result1 %s\n", testResult1)

	testValues2 := map[string]string{
		"Country": "CA",
		"Device":  "iOS",
	}
	fmt.Printf("testFeatures %s\n", testValues2)
	testResult2 := resultFromFeatures(testValues2, gameFeatureMap)
	fmt.Printf("result2 %s\n", testResult2)

	testValues3 := map[string]string{
		"Country": "CA",
		"Device":  "Android",
	}
	fmt.Printf("testFeatures %s\n", testValues3)
	testResult3 := resultFromFeatures(testValues3, gameFeatureMap)
	fmt.Printf("result3 %s\n", testResult3)

	testValues4 := map[string]string{
		"Device":  "Android",
		"UserAge": "5",
	}
	fmt.Printf("testFeatures %s\n", testValues4)
	testResult4 := resultFromFeatures(testValues4, gameFeatureMap)
	fmt.Printf("result4 %s\n", testResult4)

	testValues5 := map[string]string{
		"Device":  "Android",
		"UserAge": "6",
	}
	fmt.Printf("testFeatures %s\n", testValues5)
	testResult5 := resultFromFeatures(testValues5, gameFeatureMap)
	fmt.Printf("result5 %s\n", testResult5)

	testValues6 := map[string]string{
		"Device":  "iOS",
		"UserAge": "6",
	}
	fmt.Printf("testFeatures %s\n", testValues6)
	testResult6 := resultFromFeatures(testValues6, gameFeatureMap)
	fmt.Printf("result6 %s\n", testResult6)

}

func resultFromFeatures(featureMatch map[string]string, featureMap []map[string]map[string]interface{}) map[string]interface{} {

	// go and find the first row that matches this criteria
	result := map[string]interface{}{}

	// run through the map for each critera and match against the incoming information
	for _, featureLine := range featureMap {
		matched := true
		for criteriaName, criteria := range featureLine["Criteria"] {
			// compare the criteria against the incoming
			matchValue := featureMatch[criteriaName]
			criteriaStringMap := criteria.(map[string]string)
			criteriaValue := criteriaStringMap["Value"]
			criteriaExp := criteriaStringMap["Exp"]

			if criteriaValue != "any" {
				if matchValue == "" {
					matched = false
					break
				} else {
					// use the expression to perform the test
					switch criteriaExp {
					case "<":
						if matchValue >= criteriaValue {
							matched = false
							break
						}
					case ">":
						if matchValue <= criteriaValue {
							matched = false
							break
						}
					case "=":
						if matchValue != criteriaValue {
							matched = false
							break
						}
					}

				}
			}
		}
		if matched {
			matchResult := featureLine["Result"]
			result = matchResult
			break
		}

	}

	return result
}

func readDynamicsMap(gameId string) []map[string]map[string]interface{} {

	csvfile, err := os.Open("/home/shea/go/src/go-expt/gdv/" + gameId + ".csv")

	featureMap := []map[string]map[string]interface{}{}

	if err != nil {
		fmt.Println(err)
		return featureMap
	}

	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1 // see the Reader struct information below

	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// sanity check, display to standard output

	valueStartIndex := math.MaxUint32
	featureNames := []string{}
	variableNames := []string{}
	for rowIndex, row := range rawCSVdata {
		if rowIndex == 0 {
			for colIndex, value := range row {
				if value == "|" {
					valueStartIndex = colIndex
				}
				if colIndex > valueStartIndex {
					// this is a variable
					variableNames = append(variableNames, value)
				} else if colIndex < valueStartIndex {
					if colIndex%2 == 0 {
						featureNames = append(featureNames, value)
					}
				}
			}
		} else {
			featureLine := map[string]map[string]interface{}{}
			features := map[string]interface{}{}
			variables := map[string]interface{}{}

			for colIndex, value := range row {
				if colIndex > valueStartIndex {
					// this is a variable
					variableName := variableNames[colIndex-valueStartIndex-1]
					variables[variableName] = value
				} else if colIndex < valueStartIndex {
					featureName := featureNames[colIndex/2]
					if colIndex%2 == 0 {
						// this is a feature value
						if features[featureName] == nil {
							feature := map[string]string{}
							features[featureName] = feature
						}
						existingFeature := features[featureName].(map[string]string)
						existingFeature["Value"] = value
					} else {
						// this is a feature modifier
						if features[featureName] == nil {
							feature := map[string]string{}
							features[featureName] = feature
						}
						existingFeature := features[featureName].(map[string]string)
						existingFeature["Exp"] = value
					}
				}
			}

			featureLine["Criteria"] = features
			featureLine["Result"] = variables

			featureMap = append(featureMap, featureLine)

		}
	}

	fmt.Printf("featureMap %s\n\n", featureMap)

	return featureMap
}
