package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var pool *redis.Pool

func main() {
	pool = newPool()
	if err := populateRedis(); err != nil {
		panic(err)
	}

	http.HandleFunc("/cohortvalues", cohortValsForFeats)
	http.ListenAndServe(":8080", nil)
}

func cohortValsForFeats(w http.ResponseWriter, r *http.Request) {
	conn := pool.Get()
	defer conn.Close()

	params, err := buildParams(r.URL.Query(), conn)
	if err != nil {
		fmt.Fprintf(w, "%#v\n", err)
		return
	}
	row := getMatchingRow(params)
	vals := getRowVals(row)
	fmt.Fprintf(w, "%#v\n", vals)
}

func buildParams(givenParams url.Values, conn redis.Conn) (map[string]string, error) {
	features, err := redis.Strings(conn.Do("SMEMBERS", "features"))
	if err != nil {
		return nil, err
	}
	builtParams := make(map[string]string)
	for _, f := range features {
		f = strings.ToLower(f)
		if givenParamVals, found := givenParams[f]; found {
			givenParamVal := strings.ToLower(givenParamVals[0])
			builtParams[f] = givenParamVal
		} else {
			builtParams[f] = "any"
		}
	}
	return builtParams, nil
}

func getMatchingRow(params map[string]string) int {
	// STUB
	return 0
}

func getRowVals(index int) map[string]string {
	// STUB
	return map[string]string{"val1": "1000", "val2": "5"}
}

func jsonifyMap(m interface{}) string {
	b, _ := json.Marshal(m)
	return string(b[:])
}

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func populateRedis() error {
	conn := pool.Get()
	defer conn.Close()

	if _, err := conn.Do("FLUSHALL"); err != nil {
		return err
	}

	feats := []string{
		"country",
		"device",
		"gender",
	}
	for _, feat := range feats {
		if _, err := conn.Do("SADD", "features", feat); err != nil {
			return err
		}
	}

	vals := []string{
		"maxscore",
		"random",
	}
	for _, val := range vals {
		if _, err := conn.Do("SADD", "values", val); err != nil {
			return err
		}
	}

	cells := map[string][]string{
		"country":  []string{"any", "CAN"},
		"device":   []string{"any", "iOS"},
		"gender":   []string{"any", "any"},
		"maxscore": []string{"1000", "1100"},
		"random":   []string{"5", "6"},
	}
	for colheader, colvals := range cells {
		for _, colval := range colvals {
			colval = strings.ToLower(colval)
			if _, err := conn.Do("RPUSH", colheader, colval); err != nil {
				return err
			}
		}
	}

	return nil
}
