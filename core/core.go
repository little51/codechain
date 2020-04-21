package main

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

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

// DeliverTx check it and save to mongodb.
func (app *CoreApplication) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	_msgString, ok := app.signVerify(req.Tx)
	if !ok {
		return abcitypes.ResponseDeliverTx{Code: 1, Info: "DeliverTx verification failed"}
	}
	tokenObj, err := app.DecodeMsg(_msgString)
	if err != nil {
		return abcitypes.ResponseDeliverTx{Code: 1, Info: "DeliverTx DecodeMsg failed"}
	}
	_to := tokenObj.To
	_from := tokenObj.From
	_token := tokenObj.Token
	_amount := tokenObj.Amount
	if _to == "" || _to == _from {
		// create new token
		_, err := app.MongoDB_Query_CodeName(string(_token))
		if err == nil {
			return abcitypes.ResponseDeliverTx{Code: 1, Info: "DeliverTx CodeName has existed"}
		} else {
			if _, err := app.MongoDB_Add_CodeName(string(_token)); err != nil {
				return abcitypes.ResponseDeliverTx{Code: 1, Info: "DeliverTx MongoDB_Add_CodeName failed"}
			}
		}
		// add asset in assets
		assetNew := Asset{Publickey: _from, Token: _token, Amount: _amount}
		if _, err := app.MongoDB_Update_Assets(_from, _token, assetNew); err != nil {
			return abcitypes.ResponseDeliverTx{Code: 1, Info: "DeliverTx MongoDB_Update_Assets failed"}
		}
		return abcitypes.ResponseDeliverTx{Code: 0}
	}
	if _to != _from {
		fromPublic, err := app.MongoDB_Query_Assets(_from, _token)
		if err != nil {
			info := "you have any code of " + _token
			return abcitypes.ResponseDeliverTx{Code: 1, Info: info}
		}
		if fromPublic.Amount < _amount {
			return abcitypes.ResponseDeliverTx{Code: 1, Info: "your amount is not enough"}
		}
		fromAssets := Asset{Publickey: _from, Token: _token, Amount: fromPublic.Amount - _amount}

		toPublic, err := app.MongoDB_Query_Assets(_to, _token)
		var toAssets Asset
		if err != nil {
			toAssets = Asset{Publickey: _to, Token: _token, Amount: _amount}
		} else {
			toAssets = Asset{Publickey: _to, Token: _token, Amount: toPublic.Amount + _amount}
		}

		app.MongoDB_Update_Assets(_from, _token, fromAssets)
		app.MongoDB_Update_Assets(_to, _token, toAssets)
	}
	return abcitypes.ResponseDeliverTx{Code: 0}
}

// CheckTx check tx format .
func (app *CoreApplication) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	_, ok := app.signVerify(req.Tx)
	if !ok {
		return abcitypes.ResponseCheckTx{Code: 1, GasWanted: 1, Info: "CheckTx failedly carried out the signVerify"}
	}
	return abcitypes.ResponseCheckTx{Code: 0, GasWanted: 1, Info: "CheckTx successfully carried out the signVerify"}
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
	assetsArray, err := app.MongoDB_QueryAllKindAssetsFromPublicKey(value)
	if err != nil {
		resQuery.Code = 1
		resQuery.Log = "MongoDB_QueryAllKindAssetsFromPublicKey has something wrong"
		resQuery.Value = nil
	} else {
		if len(assetsArray) > 0 {
			var arrayCore = ""
			for i := 0; i < len(assetsArray); i++ {
				obj := "{\"publickey\":\"" + string(assetsArray[i].Publickey) + "\",\"token\":\"" + string(assetsArray[i].Token) + "\",\"amount\":" + strconv.Itoa(int(assetsArray[i].Amount)) + "}"

				if i != len(assetsArray)-1 {
					arrayCore += obj + ","
				} else {
					arrayCore += obj
				}
			}
			arrayAll := "{\"array\": [" + arrayCore + "]}"
			resQuery.Value = []byte(arrayAll)
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
