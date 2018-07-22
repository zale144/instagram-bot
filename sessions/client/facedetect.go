package client

import (
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
	"log"
	"github.com/zale144/instagram-bot/sessions/model"
)

type Response struct {
	Result int `json:"result"`
}

// GetNumberOfFaces will make a rpc call to a facedetect service
func GetNumberOfFaces(url string) (int, error) {
	request := map[string]interface{}{
		"id":"0",
		"method":  "Num.Faces",
		"params": []map[string]string{{"url":  url }},
	}
	rsp, err := rpcCall(request)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	tResp := Response{}
	json.Unmarshal([]byte(rsp), &tResp)
	return tResp.Result, nil
}

// rpcCall will execute the rpc call with provided parameters
func rpcCall(req map[string]interface{}) (string, error) {
	b, _ := json.Marshal(req)
	rsp, err := http.Post(model.RpcURI, "application/json", bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	b, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
