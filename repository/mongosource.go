// Copyright 2010 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package repo

import (
	"context"
	"github.com/wonderivan/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

//
type MongoClient struct {
	//
	client *mongo.Client
	//
	cancel context.CancelFunc
}

//
func NewMongoClient(url string, timeout time.Duration) (c *MongoClient, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err == nil {
		c = &MongoClient{
			client: client,
			cancel: cancel,
		}
	}
	return
}

//
func (c *MongoClient) Ping() error {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err := c.client.Ping(ctx, readpref.Primary())
	return err
}

//
func (c *MongoClient) Close() {
	//c.client.Disconnect()
}

//
func (c *MongoClient) Db(name string) *mongo.Database {
	db := c.client.Database(name)
	return db
}

//
func (c *MongoClient) Collection(name string, dname string) *mongo.Collection {

	return c.client.Database(dname).Collection(name)
}

//
func (c *MongoClient) Find(coll *mongo.Collection, query interface{}, opts ...*options.FindOptions) (cursor *mongo.Cursor, err error) {
	cursor, err = coll.Find(
		context.Background(),
		query,
		opts...,
	)

	return
}

//
func (c *MongoClient) Range(cursor *mongo.Cursor, iter func(doc interface{})) {
	for cursor.Next(context.Background()) {
		//iter(cursor.Decode(&st))
	}
}

//
func (c *MongoClient) InsertMany(coll *mongo.Collection, docs []interface{}, opts ...*options.InsertManyOptions) (result *mongo.InsertManyResult, err error) {
	result, err = coll.InsertMany(context.Background(), docs)

	return
}

//
func (c *MongoClient) InsertOne(coll *mongo.Collection, doc interface{}, opts ...*options.InsertOneOptions) (result *mongo.InsertOneResult, err error) {
	result, err = coll.InsertOne(context.Background(), doc)

	return
}

//
func init() {
	RegisterDataSource("mongo", &MongoSource{})
}

//
var db *mongo.Database

//
var collMap map[string]*mongo.Collection

//
func getColl(name string) *mongo.Collection {

	coll, ok := collMap[name]
	if !ok {
		coll = db.Collection(name)
	}
	return coll
}

//
type MongoSource struct {
}

//
func (e *MongoSource) Init(query map[string]interface{}) error {
	//e.CollName =query["CollName"].(string)

	return nil
}

//
func (e *MongoSource) Load(collName string, id string, v interface{}) bool {

	//return e.Query(collName, map[string]interface{}{
	//	"id":id,
	//},v)

	result := getColl(collName).FindOne(context.Background(), map[string]interface{}{
		"id": id,
	})
	if result.Err() != nil {

		return false
	}
	return true
}

//
func (e *MongoSource) Query(collName string, query map[string]interface{}) (result interface{}, err error) {
	result, err = getColl(collName).Find(context.Background(), query, nil)
	if err != nil {
		return
	}

	return
}

//
func (e *MongoSource) Update(collName string, id string, v interface{}) bool {
	//Id :=v["Id"].(primitive.ObjectID)
	oid, _ := primitive.ObjectIDFromHex(id)
	err := e.updateOne(collName, oid, v)
	return err == nil
}

//
func (e *MongoSource) Insert(collName string, id string, v interface{}) bool {
	err := e.insertOne(collName, v)

	return err == nil
}

//
func (e *MongoSource) Remove(collName string, id string) bool {
	oid, _ := primitive.ObjectIDFromHex(id)
	err := e.remove(collName, oid)

	return err == nil
}

//
func (e *MongoSource) insertOne(collName string, doc interface{}, opts ...*options.InsertOneOptions) error {

	_, err := getColl(collName).InsertOne(context.Background(), doc, opts...)

	return err
}

//
func (e *MongoSource) updateOne(collName string, id primitive.ObjectID, doc interface{}, opts ...*options.UpdateOptions) error {

	b := bson.D{{
		"_id", id,
	},
	}

	v, _ := bson.Marshal(doc)
	b2 := bson.D{
		{
			"$set", v,
		},
	}
	_, err := getColl(collName).UpdateOne(context.Background(), b2, b, opts...)

	return err
}

//
func (e *MongoSource) findOne(collName string, filter interface{}, opts ...*options.FindOneOptions) (result *mongo.SingleResult) {

	result = getColl(collName).FindOne(context.Background(), filter, opts...)

	return
}

//
func (e *MongoSource) FindOneByParams(collName string, params map[string]interface{}) (result *mongo.SingleResult) {
	b := bson.D{}
	for key, value := range params {
		b = append(b, bson.E{
			key,
			value,
		})
	}
	result = e.findOne(collName, b)

	return
}

//
func (e *MongoSource) remove(collName string, id primitive.ObjectID, opts ...*options.DeleteOptions) error {
	b := bson.D{{
		"_id", id,
	},
	}
	_, err := getColl(collName).DeleteOne(context.Background(), b, opts...)
	return err
}

//
func InitMongoSource(client *MongoClient, dbname string) {
	db = client.Db(dbname)
	if db == nil {
		logger.Debug("use db fail", dbname)
	} else {
		logger.Debug("init mongo client", dbName, client)
	}
}
