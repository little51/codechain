package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AssetBody struct
type AssetBody struct {
	PublicKey string `json:"publickey" binding:required`
	Sign      string `json:"sign" binding:required`
	Msg       string `json:"msg" binding:required`
}

// Asset struct.
type Asset struct {
	Key   string `json:"key" binding:required`
	Value string `json:"value" binding:required`
}

// Query struct
type Query struct {
	Key string `json:"key" binding:required`
}

// NewAsset create new asset by publickey,sign and asset
func NewAsset(c *gin.Context) {
	//parse post body
	var jsonAssetBody AssetBody
	if err := c.ShouldBindJSON(&jsonAssetBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//send message to chain core
	// url := "http://localhost:26657/broadcast_tx_commit?tx=\"" + _asset.Key + "=" + _asset.Value + "\""
	url := "http://localhost:26657"

	var baseInitData = "{" +
		"\"publickey\":\"" + jsonAssetBody.PublicKey + "\"," +
		"\"sign\":\"" + jsonAssetBody.Sign + "\"," +
		"\"msg\":\"" + jsonAssetBody.Msg + "\"" +
		"}"
	fmt.Println(baseInitData)
	var baseInput = []byte(baseInitData)
	var encodingString = base64.StdEncoding.EncodeToString(baseInput)
	var post = "{\"method\":\"broadcast_tx_commit\",\"jsonrpc\":\"2.0\",\"params\":{\"tx\":\"" + encodingString + "\"},\"id\":\"\"}"
	var jsonStr = []byte(post)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		c.JSON(200, gin.H{
			"result": false,
			"info":   "",
			"error":  err.Error(),
		})
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		c.JSON(200, gin.H{
			"result": false,
			"info":   "",
			"error":  err.Error(),
		})
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll((resp.Body))
	if resp.StatusCode == 200 {
		c.JSON(200, gin.H{
			"result": true,
			"info":   string(body),
			"error":  "",
		})
	} else {
		c.JSON(200, gin.H{
			"result": false,
			"info":   "",
			"error":  err.Error(),
		})
	}
}

// Query value by key
func QueryAsset(c *gin.Context) {
	var json Query
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_key := json.Key
	url := "http://localhost:26657/abci_query?data=\"key=" + _key + "\""
	// url := "http://172.16.62.48:26659/abci_query?data=\"key=" + _key + "\""
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(200, gin.H{
			"result": false,
			"info":   "",
			"error":  err.Error(),
		})
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll((resp.Body))
	if resp.StatusCode == 200 {
		c.JSON(200, gin.H{
			"result": true,
			"info":   string(body),
			"error":  "",
		})
	} else {
		c.JSON(200, gin.H{
			"result": false,
			"info":   "",
			"error":  err.Error(),
		})
	}
}
