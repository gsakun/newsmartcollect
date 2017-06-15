package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//+ `,"smartkey": "` + smartkey
func pushIt(value, timestamp, metric, tags, containerId, counterType, endpoint string) error {

	postThing := `[{"metric": "` + metric + `", "endpoint": "` + endpoint + `", "timestamp": ` + timestamp + `,"step": ` + "60" + `,"smartkey": "` + containerId + `","value": "` + value + `","counterType": "` + counterType + `","tags": "` + tags + `"}]`
	//LogRun(plu_name + "*****" + postThing)

	//push data to falcon-agent
	url := "http://127.0.0.1:1988/v1/push"
	resp, err := http.Post(url,
		"application/x-www-form-urlencoded",
		strings.NewReader(postThing))
	if err != nil {
		//LogErr("Post err in pushIt", err)
		return err
	}
	defer resp.Body.Close()
	body, err1 := ioutil.ReadAll(resp.Body)
	//LogRun(string(body))
	fmt.Println(body)
	if err1 != nil {
		LogErr("ReadAll err in pushIt", err1)
		return err1
	}

	return nil
}
