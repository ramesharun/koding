package models

import (
	"labix.org/v2/mgo/bson"
)

type KodingKeys struct {
	Id       bson.ObjectId `bson:"_id" json:"-"`
	Key      string        `bson:"key"`
	Hostname string        `bson:"hostname"`
	Kitename string        `bson:"kitename"`
	Owner    string        `bson:"owner"`
}
