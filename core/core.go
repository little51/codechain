package main

import (
	"bytes"
	"context"
	"fmt"

	"encoding/json"

	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/version"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type State struct {
	Height  int64  `json:"height"`
	AppHash []byte `json:"app_hash"`
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
	collection := app.db.Database("chain").Collection("state")
	_state := bson.M{}
	filter := bson.M{"key": "laststate"}
	err := collection.FindOne(context.TODO(), filter).Decode(&_state)
	var state State
	if err == nil {
		if value, ok := _state["value"].(bson.M); ok {
			if _appHash, ok := value["app_hash"].(primitive.Binary); ok {
				state.AppHash = _appHash.Data
			}
			if _height, ok := value["height"].(int64); ok {
				state.Height = _height
			}
		}
	}
	fmt.Println("load state from mongodb : ", state)
	return state
}

// SaveState .
func SaveState(app *CoreApplication) {
	_stateNew := bson.M{"key": "laststate", "value": bson.M{"app_hash": app.state.AppHash, "height": app.state.Height}}
	_stateOld := bson.M{}
	collection := app.db.Database("chain").Collection("state")
	filter := bson.M{"key": "laststate"}
	err := collection.FindOne(context.TODO(), filter).Decode(&_stateOld)
	if err == nil {
		collection.UpdateOne(context.TODO(), _stateOld, bson.M{"$set": _stateNew})
	} else {
		collection.InsertOne(context.TODO(), _stateNew)
	}
	fmt.Println("save state to mongodb : ", _stateNew)
	return
}

// NewCoreApplication mongodb connection come from main.go .
func NewCoreApplication(db *mongo.Client) *CoreApplication {
	return &CoreApplication{db: db}
}

// Info interface.
func (app *CoreApplication) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	app.state = LoadState(app)
	if app.state.Height == 0 {
		return abcitypes.ResponseInfo{Data: "codechain v0.0.1",
			Version:    version.ABCIVersion,
			AppVersion: 20200330,
		}
	}
	return abcitypes.ResponseInfo{Data: "codechain v0.0.1",
		Version:          version.ABCIVersion,
		AppVersion:       20200330,
		LastBlockAppHash: app.state.AppHash,
		LastBlockHeight:  app.state.Height,
	}
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
	//find tx,if exists then update else insert
	filter := bson.M{"key": string(key)}
	assetOld := bson.M{}
	assetNew := bson.M{"key": string(key), "value": string(value)}
	collection := app.db.Database("chain").Collection("assets")
	err := collection.FindOne(context.TODO(), filter).Decode(&assetOld)
	if err == nil {
		collection.UpdateOne(context.TODO(), assetOld, bson.M{"$set": assetNew})
	} else {
		collection.InsertOne(context.TODO(), assetNew)
	}
	var _code uint32
	//tokens parse
	var tokenTx TokenTx
	errParse := json.Unmarshal([]byte(string(value)), &tokenTx)
	if errParse == nil {
		_code = DoTokenTx(app, tokenTx)
	}
	return abcitypes.ResponseDeliverTx{Code: _code}
}

// CheckTx check tx format .
func (app *CoreApplication) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	code := app.isValid(req.Tx)
	return abcitypes.ResponseCheckTx{Code: code, GasWanted: 1}
}

// Commit interface .
func (app *CoreApplication) Commit() abcitypes.ResponseCommit {
	//appHash := make([]byte, 32)
	//binary.PutVarint(appHash, app.state.Height)
	//make empty apphash avoid create empty block
	app.state.AppHash = []byte{} //appHash
	app.state.Height++
	SaveState(app)
	return abcitypes.ResponseCommit{Data: app.state.AppHash}
}

// Query  query document from mongledb.
func (app *CoreApplication) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	parts := bytes.Split(reqQuery.Data, []byte("="))
	value := string(parts[1])
	filter := bson.M{"key": string(value)}
	collection := app.db.Database("chain").Collection("assets")
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
	collection := app.db.Database("chain").Collection("assets")
	collection.Drop(context.TODO())
	return abcitypes.ResponseInitChain{}
}

// BeginBlock interface.
func (app *CoreApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	return abcitypes.ResponseBeginBlock{}
}

// EndBlock interface.
func (app *CoreApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
	/*
		dynamic add validator
		if len(app.validators) == 0 || req.Height <= 21 {
			return abcitypes.ResponseEndBlock{}
		}
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
		return abcitypes.ResponseEndBlock{ValidatorUpdates: app.validators}*/
}
