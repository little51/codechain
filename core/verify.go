package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/tendermint/tendermint/crypto/ed25519"
)

// tx struct
type TxStruct struct {
	PublicKey string `json:"publickey"`
	Sign      string `json:"sign"`
	Msg       string `json:"msg"`
}

// isTxJson
func (app *CoreApplication) isTxJson(jsonByte []byte) (jsonObj TxStruct, err error) {
	var _txData TxStruct
	if err := json.Unmarshal(jsonByte, &_txData); err != nil {
		return _txData, err
	} else {
		return _txData, nil
	}
}

// Verify sign
func (app *CoreApplication) signVerify(req []byte) (string, bool) {
	var _msg string
	// restore public key
	reqData, err := app.isTxJson(req)
	if err != nil {
		return _msg, false
	}
	_publickey := reqData.PublicKey
	_sign, _ := hex.DecodeString(reqData.Sign)
	_msg = reqData.Msg
	var publicKey ed25519.PubKeyEd25519
	temp, _ := hex.DecodeString(_publickey)
	copy(publicKey[:], temp)
	//verify sign
	b := publicKey.VerifyBytes([]byte(_msg), []byte(_sign))
	if !b {
		return _msg, false
	}
	return _msg, true
}

// judgement the Type of Value in Msg
func (app *CoreApplication) DecodeMsg(value string) (TokenTx, error) {
	var tokenString TokenTxString
	var tokenObj TokenTx
	decodeBytes, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		fmt.Println("decodeString msg wrong")
		return tokenObj, err
	}
	if err := json.Unmarshal([]byte(decodeBytes), &tokenString); err != nil {
		fmt.Println("value string is not a right JSON")
		return tokenObj, err
	}
	tokenObj.Token = tokenString.Token
	tokenObj.From = tokenString.From
	tokenObj.To = tokenString.To
	tempIntNum, _ := strconv.Atoi(tokenString.Amount)
	tokenObj.Amount = uint32(tempIntNum)
	return tokenObj, nil
}
