package main

import (
	"bytes"
	"context"
	"fmt"

	abcitypes "github.com/tendermint/tendermint/abci/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// AssetsApplication 数据库变量.
type AssetsApplication struct {
	db *mongo.Client
}

var _ abcitypes.Application = (*AssetsApplication)(nil)

// NewAssetsApplication AssetsApplication构造函数，db变量由main.go传入 .
func NewAssetsApplication(db *mongo.Client) *AssetsApplication {
	return &AssetsApplication{db: db}
}

// Info .
func (AssetsApplication) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{}
}

// SetOption .
func (AssetsApplication) SetOption(req abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
	return abcitypes.ResponseSetOption{}
}

// isValid 校验交易串是否合法
func (app *AssetsApplication) isValid(tx []byte) (code uint32) {
	parts := bytes.Split(tx, []byte("="))
	if len(parts) != 2 {
		return 1
	}
	return 0
}

// DeliverTx 交易送达响应 .
func (app *AssetsApplication) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	fmt.Println(string(req.Tx))
	code := app.isValid(req.Tx)
	if code != 0 {
		return abcitypes.ResponseDeliverTx{Code: code}
	}
	parts := bytes.Split(req.Tx, []byte("="))
	key, value := string(parts[0]), string(parts[1])
	collection := app.db.Database("chain").Collection("assets")
	assets := bson.M{"key": string(key), "value": string(value)}
	fmt.Println(assets)
	insertResult, err := collection.InsertOne(context.TODO(), assets)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted a single document:", insertResult.InsertedID)
	return abcitypes.ResponseDeliverTx{Code: 0}
}

// CheckTx 交易校验响应 .
func (app *AssetsApplication) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	code := app.isValid(req.Tx)
	return abcitypes.ResponseCheckTx{Code: code, GasWanted: 1}
}

// Commit 交易提交响应 .
func (app *AssetsApplication) Commit() abcitypes.ResponseCommit {
	return abcitypes.ResponseCommit{Data: []byte{}}
}

// Query 查询交易 .
func (app *AssetsApplication) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	parts := bytes.Split(reqQuery.Data, []byte("="))
	value := string(parts[1])
	filter := bson.M{"key": string(value)}
	collection := app.db.Database("chain").Collection("assets")
	assets := bson.M{}
	err := collection.FindOne(context.TODO(), filter).Decode(&assets)
	if err != nil {
		error := fmt.Sprintf("%s", err)
		resQuery.Code = 1
		resQuery.Log = error
		resQuery.Value = nil
	} else {
		if value, ok := assets["value"].(string); ok {
			resQuery.Value = []byte(value)
			resQuery.Info = value
			resQuery.Code = 0
			resQuery.Log = ""
		} else {
			resQuery.Value = nil
			resQuery.Code = 1
			resQuery.Log = "error type"
		}
	}
	return
}

// InitChain 初始化链 .
func (app *AssetsApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	collection := app.db.Database("chain").Collection("assets")
	collection.Drop(context.TODO())
	return abcitypes.ResponseInitChain{}
}

// BeginBlock .
func (app *AssetsApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	return abcitypes.ResponseBeginBlock{}
}

// EndBlock .
func (AssetsApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
}
