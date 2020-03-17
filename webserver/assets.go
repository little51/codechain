package main

import (
	"encoding/hex"
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

// NewAsset create new asset by publickey,sign and msg
func NewAsset(c *gin.Context) {
	var json AssetBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_publickey := json.PublicKey
	_sign, _ := hex.DecodeString(json.Sign)
	_msg := json.Msg
	var publicKey ed25519.PubKeyEd25519
	temp, _ := hex.DecodeString(_publickey)
	copy(publicKey[:], temp)
	b := publicKey.VerifyBytes([]byte(_msg), []byte(_sign))
	c.JSON(200, gin.H{
		"result": b,
		"error":  "",
	})
}
