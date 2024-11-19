/*
 * @Author: Liu Sainan
 * @Date: 2023-12-10 16:46:36
 */

package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Collection struct {
	Id           string             `json:"id" bson:"_id"`
	LastmodEpoch primitive.ObjectID `json:"lastmod_epoch" bson:"lastmodEpoch"`
	Lastmod      primitive.DateTime `json:"lastmod" bson:"lastmod"`
	Dropped      bool               `json:"dropped" bson:"dropped"`
	Key          bson.D             `json:"key" bson:"key"`
	Unique       bool               `json:"unique" bson:"unique"`
	//Uuid         string              `json:"uuid" bson:"uuid"`
}
