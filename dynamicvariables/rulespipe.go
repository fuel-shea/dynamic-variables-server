package dynamicvariables

import (
	"gopkg.in/mgo.v2/bson"
)

type LoopPipe struct {
	Pipe   []bson.M
	GameID string
}

func PipeSkeleton(gameID string) LoopPipe {
	return LoopPipe{
		Pipe: []bson.M{
			bson.M{
			// placeholder for game/rule matching
			},

			bson.M{
			// placeholder for feature matching
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
		},
		GameID: gameID,
	}
}

func (lp LoopPipe) UpdateForLoop(eligibleRules []int, featureType string, matchVal interface{}) {
	// ensure the feature is from the right game and from an eligible rule
	lp.Pipe[0] = bson.M{
		"$match": bson.M{
			"game_id":  lp.GameID,
			"rule_idx": bson.M{"$in": eligibleRules},
		},
	}

	// ensure the feature actually matches the matchVal
	lp.Pipe[1] = bson.M{
		"$match": bson.M{
			"type": featureType,
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
	}
}
