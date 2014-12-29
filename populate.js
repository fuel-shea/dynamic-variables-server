(function addAllData() {
    db.dropDatabase();

    var gameId = "gid1";
    var rules = [
        {
            rule_idx: 0,
            game_id: gameId,
            features: [
                { "type" : "UserID", "value" : "any", "mod" : "=" },
                { "type" : "Country", "value" : "any", "mod" : "=" },
                { "type" : "Device", "value" : "iOS", "mod" : "=" },
                { "type" : "DeviceHardware", "value" : "any", "mod" : "=" },
                { "type" : "OSVersion", "value" : "any", "mod" : "=" },
                { "type" : "UserAge", "value" : "any", "mod" : "=" },
                { "type" : "Region", "value" : "any", "mod" : "=" },
                { "type" : "LTV", "value" : "any", "mod" : "=" },
                { "type" : "GameVersion", "value" : "any", "mod" : "=" },
            ],
            variables: {
                "randomMax" : "950",
                "whammyChance" : "6",
                "floaterDuration" : "0.9"
            }
        },
        {
            rule_idx: 1,
            game_id: gameId,
            features: [
                { "type" : "Device", "value" : "any", "mod" : "=" },
                { "type" : "LTV", "value" : "any", "mod" : "=" },
                { "type" : "UserAge", "value" : "any", "mod" : "=" },
                { "type" : "UserID", "value" : "any", "mod" : "=" },
                { "type" : "Country", "value" : "CA", "mod" : "=" },
                { "type" : "Region", "value" : "any", "mod" : "=" },
                { "type" : "DeviceHardware", "value" : "any", "mod" : "=" },
                { "type" : "OSVersion", "value" : "any", "mod" : "=" },
                { "type" : "GameVersion", "value" : "any", "mod" : "=" },
            ],
            variables: {
                "randomMax" : "1400",
                "whammyChance" : "6",
                "floaterDuration" : "0.9"
            }
        },
        {
            rule_idx: 2,
            game_id: gameId,
            features: [
                { "type" : "UserID", "value" : "any", "mod" : "=" },
                { "type" : "Country", "value" : "any", "mod" : "=" },
                { "type" : "Region", "value" : "any", "mod" : "=" },
                { "type" : "Device", "value" : "any", "mod" : "=" },
                { "type" : "LTV", "value" : "any", "mod" : "=" },
                { "type" : "UserAge", "value" : "5", "mod" : ">" },
                { "type" : "DeviceHardware", "value" : "any", "mod" : "=" },
                { "type" : "OSVersion", "value" : "any", "mod" : "=" },
                { "type" : "GameVersion", "value" : "any", "mod" : "=" },
            ],
            variables: {
                "randomMax" : "1500",
                "whammyChance" : "6",
                "floaterDuration" : "0.9"
            }
        },
        {
            rule_idx: 3,
            game_id: gameId,
            features: [
                { "type" : "UserID", "value" : "any", "mod" : "=" },
                { "type" : "Country", "value" : "any", "mod" : "=" },
                { "type" : "DeviceHardware", "value" : "any", "mod" : "=" },
                { "type" : "GameVersion", "value" : "any", "mod" : "=" },
                { "type" : "Region", "value" : "any", "mod" : "=" },
                { "type" : "Device", "value" : "any", "mod" : "=" },
                { "type" : "OSVersion", "value" : "any", "mod" : "=" },
                { "type" : "LTV", "value" : "any", "mod" : "=" },
                { "type" : "UserAge", "value" : "any", "mod" : "=" },
            ],
            variables: {
                "randomMax" : "950",
                "whammyChance" : "6",
                "floaterDuration" : "1"
            }
        },
    ];

    var featureTypes = [];
    rules.forEach(function(r) {
        r.features.forEach(function(f) {
            featureTypes.push(f.type);
            var fObj = {
                game_id: r.game_id,
                rule_idx: r.rule_idx,
                type: f.type,
                value: f.value,
                mod: f.mod,
            }
            db.features.insert(fObj);
        });
    });

    var varTypes = [];
    rules.forEach(function(r) {
        var vObj = {
            game_id: r.game_id,
            rule_idx: r.rule_idx
        };
        for (var vKey in r.variables) {
            varTypes.push(vKey);
            vObj[vKey] = r.variables[vKey]
        }
        db.variables.insert(vObj);
    });

    db.game_rule_data.insert({
        game_id: gameId,
        variable_types: varTypes,
        feature_types: featureTypes,
        num_rules: rules.length,
    });
})();

