package main

import (
	"fmt"
)

type TokenTx struct {
	Token  string `json:"token"`
	From   string `json:"from"`
	To     string `json:"to"`
	Amount uint32 `json:"amount"`
}

type TokenTxString struct {
	Token  string `json:"token"`
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
}

func DoTokenTx(app *CoreApplication, tokenTx TokenTx) uint32 {
	fmt.Println(tokenTx)
	return 0
}
