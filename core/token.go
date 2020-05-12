package main

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

type MsgTx struct {
	Msg string `json:"msg"`
}
