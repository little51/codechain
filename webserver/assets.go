package main

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tendermint/tendermint/crypto/ed25519"
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

// NewAsset create new asset by publickey,sign and asset
func NewAsset(c *gin.Context) {
	//parse post body
	var jsonAssetBody AssetBody
	if err := c.ShouldBindJSON(&jsonAssetBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//parse asset to : key = value
	_msg := jsonAssetBody.Msg
	var _asset Asset
	err := json.Unmarshal([]byte(_msg), &_asset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//restore public key
	_publickey := jsonAssetBody.PublicKey
	_sign, _ := hex.DecodeString(jsonAssetBody.Sign)
	var publicKey ed25519.PubKeyEd25519
	temp, _ := hex.DecodeString(_publickey)
	copy(publicKey[:], temp)
	//verify sign
	b := publicKey.VerifyBytes([]byte(_msg), []byte(_sign))
	if !b {
		c.JSON(200, gin.H{
			"result": false,
			"info":   "",
			"error":  "sign error",
		})
		return
	}
	//send message to chain core
	url := "http://localhost:26657/broadcast_tx_commit?tx=\"" + _asset.Key + "=" + _asset.Value + "\""
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
