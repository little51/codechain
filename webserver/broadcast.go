package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

type MsgTx struct {
	PrivateKey string `json:"privatekey" binding:required`
	PublicKey  string `json:"publickey" binding:required`
	Msg        string `json:"msg" binding:required`
}

type MsgDetail struct {
	Msg string `json:"msg" binding:required`
}

func BroadCastGitClone(c *gin.Context) {
	var postMsg MsgDetail
	if err := c.ShouldBindJSON(&postMsg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("BroadCastGitClone")
	fmt.Println(postMsg)
	// git cache according to msg ...
}

// post {method jsonrpc params id} to 26657/broadcast_tx_commit
func BroadCastMsg(c *gin.Context) {
	var json MsgTx
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var base64msg = base64.StdEncoding.EncodeToString([]byte(json.Msg))

	_privatekey, _ := hex.DecodeString(json.PrivateKey)
	var privateKey ed25519.PrivKeyEd25519
	copy(privateKey[:], _privatekey)
	signStr, err := privateKey.Sign([]byte(base64msg))
	if err == nil {
		sign := hex.EncodeToString(signStr)
		// after getting sign,then post { public sign msg}
		url := "http://localhost:26657"
		var baseInitData = "{" +
			"\"publickey\":\"" + json.PublicKey + "\"," +
			"\"sign\":\"" + sign + "\"," +
			"\"msg\":\"" + base64msg + "\"" +
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
	} else {
		error := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{
			"body":  signStr,
			"error": error,
		})
	}
}
