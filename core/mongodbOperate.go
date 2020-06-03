package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CodeName struct {
	Name string `json:"name"`
}

type Asset struct {
	Publickey string `json:"publickey"`
	Token     string `json:"token""`
	Amount    uint32 `json:"amount"`
	Repostory string `json:"repostory"`
}

// query data from Assets according to publickey and token
func (app *CoreApplication) MongoDB_Query_Assets(publickey string, token string, repostory string) (Asset, error) {
	collection := app.db.Database("chain").Collection("assets")
	var assetsResult Asset
	filter := bson.M{"publickey": string(publickey), "token": string(token), "repostory": string(repostory)}
	err := collection.FindOne(context.TODO(), filter).Decode(&assetsResult)
	if err != nil {
		return assetsResult, err
	} else {
		return assetsResult, nil
	}
}

// update data from Assets according to assetsNew
func (app *CoreApplication) MongoDB_Update_Assets(publickey string, token string, repostory string, assetNew Asset) (Asset, error) {
	collection := app.db.Database("chain").Collection("assets")
	assetOld, err := app.MongoDB_Query_Assets(publickey, token, repostory)
	if err != nil {
		if _, err := collection.InsertOne(context.TODO(), assetNew); err != nil {
			return assetNew, err
		}
		return assetNew, nil
	} else {
		_, err := collection.UpdateOne(context.TODO(), assetOld, bson.M{"$set": assetNew})
		if err != nil {
			return assetNew, err
		}
		return assetNew, nil
	}
}

// query all data from Assets according to public
func (app *CoreApplication) MongoDB_QueryAllKindAssetsFromPublicKey(publickey string) ([]Asset, error) {
	collection := app.db.Database("chain").Collection("assets")
	findOptions := options.Find()
	var results []Asset
	filter := bson.M{"publickey": string(publickey)}
	cur, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return results, err
	}
	for cur.Next(context.TODO()) {
		var elem Asset
		err := cur.Decode(&elem)
		if err != nil {
			return results, err
		}
		results = append(results, elem)
	}
	if err := cur.Err(); err != nil {
		return results, err
	}

	cur.Close(context.TODO())
	// fmt.Println("MongoDB_QueryAllDataFromPublicKey")
	// fmt.Println(results)
	return results, nil
}

// query data from CodeName according to codeName
func (app *CoreApplication) MongoDB_Query_CodeName(codeName string) (CodeName, error) {
	collection := app.db.Database("chain").Collection("codeName")
	filter := bson.M{"name": string(codeName)}
	var tokenNameResult CodeName
	err := collection.FindOne(context.TODO(), filter).Decode(&tokenNameResult)
	if err != nil {
		return tokenNameResult, err
	} else {
		return tokenNameResult, nil
	}
}

// add data ine CodeName according codeName
func (app *CoreApplication) MongoDB_Add_CodeName(codeName string) (CodeName, error) {
	collection := app.db.Database("chain").Collection("codeName")
	codeOne := CodeName{codeName}
	_, err := collection.InsertOne(context.TODO(), codeOne)
	if err != nil {
		return codeOne, err
	}
	return codeOne, nil
}
