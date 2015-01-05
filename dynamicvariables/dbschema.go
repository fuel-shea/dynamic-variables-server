package dynamicvariables

type GameRuleData struct {
	GameID        string   `bson:"game_id"`
	VariableTypes []string `bson:"variable_types"`
	FeatureTypes  []string `bson:"feature_types"`
	NumRules      int      `bson:"num_rules"`
}

type Feature struct {
	GameID  string `bson:"game_id"`
	RuleIdx int    `bson:"rule_idx"`
	Type    string `bson:"type"`
	Val     string `bson:"value"`
	Mod     string `bson:"mod"`
}

type Variable struct {
	GameID       string `bson:"game_id"`
	RuleIdx      int    `bson:"rule_idx"`
	RandomMax    string `bson:"randomMax"`
	WhammyChance string `bson:"whammyChance"`
}
