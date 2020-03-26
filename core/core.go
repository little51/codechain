package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	abcitypes "github.com/tendermint/tendermint/abci/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type State struct {
	Height  int32  `json:"height"`
	AppHash string `json:"app_hash"`
}

// CoreApplication mongodb connection.
type CoreApplication struct {
	db         *mongo.Client
	validators []abcitypes.ValidatorUpdate
	state      State
}

var _ abcitypes.Application = (*CoreApplication)(nil)

// LoadState .
func LoadState(app *CoreApplication) State {
	collection := app.db.Database("chain").Collection("State")
	//test
	//_state_w := bson.M{"key": "laststate", "value": bson.M{"app_hash": "132", "height": 3}}
	//	collection.InsertOne(context.TODO(), _state_w)
	_state := bson.M{}
	filter := bson.M{"key": "laststate"}
	err := collection.FindOne(context.TODO(), filter).Decode(&_state)
	var state State
	if err == nil {
		if value, ok := _state["value"].(bson.M); ok {
			if _appHash, ok := value["app_hash"].(string); ok {
				state.AppHash = _appHash
			}
			if _height, ok := value["height"].(int32); ok {
				state.Height = _height
			}
		}
	}
	return state
}

// NewCoreApplication mongodb connection come from main.go .
func NewCoreApplication(db *mongo.Client) *CoreApplication {
	return &CoreApplication{db: db}
}

// Info interface.
func (app *CoreApplication) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	app.state = LoadState(app)
	return abcitypes.ResponseInfo{}
	/*return abcitypes.ResponseInfo{Data: "codechain v0.0.1",
		Version:          version.ABCIVersion,
		AppVersion:       20200326,
		LastBlockAppHash: app.state.AppHash,
		LastBlockHeight:  app.state.Height,
	}*/
}

// SetOption interface.
func (CoreApplication) SetOption(req abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
	return abcitypes.ResponseSetOption{}
}

// isValid
func (app *CoreApplication) isValid(tx []byte) (code uint32) {
	parts := bytes.Split(tx, []byte("="))
	if len(parts) != 2 {
		return 1
	}
	return 0
}

// DeliverTx check it and save to mongodb.
func (app *CoreApplication) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	code := app.isValid(req.Tx)
	if code != 0 {
		return abcitypes.ResponseDeliverTx{Code: code}
	}
	parts := bytes.Split(req.Tx, []byte("="))
	key, value := string(parts[0]), string(parts[1])
	collection := app.db.Database("chain").Collection("Core")
	Core := bson.M{"key": string(key), "value": string(value)}
	fmt.Println(Core)
	_, err := collection.InsertOne(context.TODO(), Core)
	if err != nil {
		panic(err)
	}
	return abcitypes.ResponseDeliverTx{Code: 0}
}

// CheckTx check tx format .
func (app *CoreApplication) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	code := app.isValid(req.Tx)
	return abcitypes.ResponseCheckTx{Code: code, GasWanted: 1}
}

// Commit interface .
func (app *CoreApplication) Commit() abcitypes.ResponseCommit {
	return abcitypes.ResponseCommit{Data: []byte{}}
}

// Query  query document from mongledb.
func (app *CoreApplication) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	parts := bytes.Split(reqQuery.Data, []byte("="))
	value := string(parts[1])
	filter := bson.M{"key": string(value)}
	collection := app.db.Database("chain").Collection("Core")
	Core := bson.M{}
	err := collection.FindOne(context.TODO(), filter).Decode(&Core)
	if err != nil {
		error := fmt.Sprintf("%s", err)
		resQuery.Code = 1
		resQuery.Log = error
		resQuery.Value = nil
	} else {
		if value, ok := Core["value"].(string); ok {
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

// InitChain drop collection .
func (app *CoreApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	app.validators = req.Validators
	collection := app.db.Database("chain").Collection("Core")
	collection.Drop(context.TODO())
	return abcitypes.ResponseInitChain{}
}

// BeginBlock interface.
func (app *CoreApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	return abcitypes.ResponseBeginBlock{}
}

// EndBlock interface.
func (app *CoreApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	//fmt.Printf("%+v", req) test dynamic add validator
	if len(app.validators) == 0 || req.Height <= 21 {
		return abcitypes.ResponseEndBlock{}
	}
	fmt.Println(len(app.validators))
	var v abcitypes.ValidatorUpdate
	// test new validator's public key
	v.Power = 10
	v.PubKey.Type = "ed25519"
	v.PubKey.Data, _ = base64.StdEncoding.DecodeString("BsY96CRY2RK+vcVbMFpOiGQSLJARQTlDB00BbyZuwM0=")
	//
	keyExists := false
	for i := 0; i < len(app.validators); i++ {
		if bytes.Compare(app.validators[i].PubKey.Data, v.PubKey.Data) == 0 {
			keyExists = true
			break
		}
	}
	if keyExists {
		return abcitypes.ResponseEndBlock{}
	}
	app.validators = append(app.validators, v)
	return abcitypes.ResponseEndBlock{ValidatorUpdates: app.validators}
}
