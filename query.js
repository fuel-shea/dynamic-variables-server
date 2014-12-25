var gameId = "gid1";
var matchFeatures = { Country: "CA" };
var featTypes = db.feat_types.findOne({game_id: gameId})["types"];

var nRules = db.var_types.count();
var eligibleRules = [];
for (var i=0; i<nRules; i++) {
    (function(n) { eligibleRules.push(n) })(i);
}

featTypes.forEach(function(ft) {
    var matchVal;
    if (matchFeatures.hasOwnProperty(ft)) {
        matchVal = matchFeatures[ft];
    } else {
        matchVal = "any";
    }

    var p = [
        {
            $match: {
                "game_id": gameId,
                "rule_idx": {"$in": eligibleRules},
                "type": ft,
                "$or": [
                        { "value": "any" },
                        { "$and": [ { "modifier": "="}, {"value": matchVal} ] },
                        { "$and": [ { "modifier": ">"}, {"value": {"$gt": matchVal}} ] },
                        { "$and": [ { "modifier": "<"}, {"value": {"$lt": matchVal}} ] },
                    ]
            }
        },
        {$group: { _id: { rule_idx: "$rule_idx" }}},
        {$sort: {"_id.rule_idx": 1}},
        {$project: {rule_idx: "$_id.rule_idx", _id: 0}},
    ];

    var ruleIdxs = db.features.aggregate(p);
    eligibleRules = db.features.aggregate(p).toArray().map(function(ruleIdxRes){
        return String.ParseInt(ruleIdxRes.rule_idx)
    });
});

if (eligibleRules.length !== 1) {
    throw new Error("eligibleRules has " + eligibleRules.length.toString() + " elements!");
}

var variables = db.variables.findOne({game_id: gameId, rule_idx: eligibleRules[0]});

var varObj = {};
var varTypes = db.var_types.findOne({game_id: gameId});
varTypes.types.forEach(function(vt) {
    varObj[vt] = variables[vt];
});

varObj;



/*
var featureTypeHandlers = [];
var featureMatchObj = {$match: {$or: featureTypeHandlers}};
featTypes.forEach(function(ft) {
    var matchVal;
    if (matchFeatures.hasOwnProperty(ft)) {
        matchVal = matchFeatures[ft]
    } else {
        matchVal = "any"
    }
    featureTypeHandlers.push({
        $and: [
            { type: ft },
            { $or: [
                    { value: "any" },
                    { $and: [ { modifier: "="}, {value: matchVal} ] },
                    { $and: [ { modifier: ">"}, {value: {$gt: matchVal}} ] },
                    { $and: [ { modifier: "<"}, {value: {$lt: matchVal}} ] },
                ]
            }
        ]
    });
});

var q = [
    featureMatchObj,
    {$group: { _id: { rule_idx: "$rule_idx" }}},
    {$sort: {"_id.rule_idx": 1}},
    {$limit: 1},
    {$project: {rule_idx: "$_id.rule_idx", _id: 0}},
]

var ruleIdxRes = db.features.aggregate(q);
var variables = db.variables.findOne({game_id: gameId, rule_idx: ruleIdxRes.next().rule_idx});
*/
