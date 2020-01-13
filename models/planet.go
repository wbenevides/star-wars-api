package models

import "gopkg.in/mgo.v2/bson"

type Planet struct {
	ID      bson.ObjectId `bson:"_id" json:"id"`
	Name    string        `bson:"name" json:"name"`
	Climate string        `bson:"climate" json:"climate"`
	Terrain string        `bson:"terrain" json:"terrain"`
	Films   int           `bson:"films" json:"films"`
}
