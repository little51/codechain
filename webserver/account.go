package main

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

// SignBody struct .
type SignBody struct {
	PrivateKey string `json:"privatekey" binding:required`
	Msg        string `json:"msg" binding:required`
}

// NewAccount create key pair .
func NewAccount(c *gin.Context) {
	privateKey := ed25519.GenPrivKey()
	publicKey := privateKey.PubKey()
	_address := publicKey.Address()
	//copy from 5th
	_privatekey := hex.EncodeToString(privateKey.Bytes()[5:])
	var _publickey = fmt.Sprintf("%s", publicKey)
	_publickey = strings.Replace(_publickey, "PubKeyEd25519{", "", -1)
	_publickey = strings.Replace(_publickey, "}", "", -1)
	c.JSON(200, gin.H{
		"privateKey": _privatekey,
		"publicKey":  _publickey,
		"address":    _address,
		"error":      "",
	})
}

// Sign string by privateKey
func Sign(c *gin.Context) {
	var json SignBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_privatekey, _ := hex.DecodeString(json.PrivateKey)
	var privateKey ed25519.PrivKeyEd25519
	copy(privateKey[:], _privatekey)
	signStr, err := privateKey.Sign([]byte(json.Msg))
	if err == nil {
		c.JSON(200, gin.H{
			"sign":  hex.EncodeToString(signStr),
			"error": "",
		})
	} else {
		error := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{
			"body":  signStr,
			"error": error,
		})
	}
}
