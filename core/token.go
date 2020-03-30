package main

type TokenTx struct {
	token  string `json:"token"`
	from   string `json:"from"`
	to     string `json:to`
	amount uint32 `json:amount`
}

func Tx(app *CoreApplication, token TokenTx) uint32 {
	return 0
}
