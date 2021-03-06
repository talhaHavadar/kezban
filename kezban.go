package kezban

import (
	"time"
	"reflect"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"errors"
)

type KezQu struct {
	Query interface{}
	Limit int
}

type Model struct {
	model          interface{}
	Id             bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at,omitempty"`
	collectionName string

}

var Database *mgo.Session = nil
var APPNAME string
func Initialize(uri string, appname string) (error) {
	database, err := mgo.Dial(uri)
	if err != nil {
		fmt.Errorf("Database initialization failed!!\n")
		return err
	}
	database.SetMode(mgo.Monotonic, true)
	Database = database
	APPNAME = appname
	fmt.Println("Database initialization is completed.")
	return nil
}

func (self *Model) SetItself(model interface{}) {
	self.model = model
}

func (self *Model) uniqueFieldCheck() error {
	uniqueMap := GetFields(self.model, "unique")
	/**
	* TODO: Unique check will be implemented via func (*Collection) EnsureIndex
	*/
	if len(uniqueMap) > 0 {
		newModel := createEmptyStruct(self.model)
		FillStruct(newModel, uniqueMap)
		err := self.FindOne(newModel, newModel)
		if err != nil {
			if err.Error() == "not found" {
				return nil
			} else {
				return err
			}
		}
		return errors.New("unique field duplicate")
	} else {
		return nil
	}

}

func (self *Model) Save() (*Model, error) {
	self.UpdatedAt = time.Now()
	if !self.checkAndSetCollectionName() {
		fmt.Errorf("Something went wrong while trying to fetch collection name.\n")
		return nil, errors.New("Something went wrong while trying to fetch collection name.")
	}
	if err := self.uniqueFieldCheck(); err != nil {
		fmt.Errorf(err.Error() + "\n")
		return nil, err
	}
	if !self.Id.Valid() { // first time creation
		self.Id = bson.NewObjectId()
		self.CreatedAt = time.Now()
		err := Database.DB(APPNAME).C(self.collectionName).Insert(&self.model)
		if err != nil {
			fmt.Errorf(err.Error() + "\n")
			return nil, err
		}
	} else {
		err := Database.DB(APPNAME).C(self.collectionName).Update(
			bson.M{"_id" : self.Id},
			self.model,
		)
		if err != nil {
			fmt.Errorf(err.Error() + "\n")
			return nil, err
		}
	}
	return self, nil
}

func (self *Model) FindOne(query interface{}, model interface{}) (err error) {
	if !self.checkAndSetCollectionName() {
		panic("Collection name was not set.")
	}
	var q *mgo.Query
	if q, err = self.constructQuery(query); err != nil {
		return err
	}
	fmt.Println("FindOne: q=", q, "self=", self)
	return q.One(model)
}


/*
 * @param query for specific filters
 * @param models needs to be pointer of model array
 * @return err
 */
func (self *Model) FindAll(query KezQu, models interface{}) (error) {
	mQuery, err := self.constructQuery(query.Query);
	fmt.Println("FindAll: mQuery=", mQuery)
	if err != nil {
		return err
	}
	if query.Limit > 0 {
		mQuery.Limit(query.Limit)
	}
	return mQuery.All(models)
}

func (self *Model) Search(query KezQu, indexes []string, models interface{}) (error) {
	mQuery := self.constructSearchQuery(query.Query.(bson.M), indexes)
	fmt.Println("Search: mQuery=", mQuery)
	if query.Limit > 0 {
		mQuery.Limit(query.Limit)
	}
	return mQuery.All(models)
}

func (self *Model) Remove(query interface{}) (err error) {
	if !self.checkAndSetCollectionName() {
		panic("Collection name was not set.")
	}

	var bsonQ bson.M
	if bsonQ, err = docToBson(query); err != nil {
		return err
	}

	return Database.DB(APPNAME).C(self.collectionName).Remove(bsonQ)
}

func (self *Model) Delete() (err error) {
	if !self.checkAndSetCollectionName() {
		panic("Collection name was not set.")
	}

	var bsonQ bson.M = bson.M{
		"_id": self.Id,
	}

	return Database.DB(APPNAME).C(self.collectionName).Remove(bsonQ)
}

func (self *Model) Find(q bson.M) *mgo.Query {
	if !self.checkAndSetCollectionName() {
		panic("Collection name was not set.")
	}
	return Database.DB(APPNAME).C(self.collectionName).Find(q)
}

func (self *Model) getMethodViaReflection(methodName string) (reflect.Value, string) {
	modelVal := reflect.ValueOf(self.model)
	function := modelVal.Elem().Addr().MethodByName(methodName)
	if function.IsValid() {
		return function, ""
	}
	return reflect.Zero(reflect.TypeOf(function)), methodName + " is invalid"
}

func (self *Model) constructSearchQuery(query bson.M, indexes []string) (*mgo.Query) {
	if !self.checkAndSetCollectionName() {
		panic("Collection name was not set.")
	}
	c := Database.DB(APPNAME).C(self.collectionName)
	index := mgo.Index{
		Key: indexes,
	}
	c.EnsureIndex(index)
	return c.Find(query)
}

func (self *Model) constructQuery(queryDoc interface{}) (*mgo.Query, error) {
	if !self.checkAndSetCollectionName() {
		panic("Collection name was not set.")
	}
	var query bson.M
	var err error

	if queryDoc == nil {
		queryDoc = self.model
	}
	if query, err = docToBson(queryDoc); err != nil {
		return nil, err
	}
	fmt.Println("constructQuery:", query, "AppName:", APPNAME, "CollectionName:", self.collectionName)
	return Database.DB(APPNAME).C(self.collectionName).Find(query), nil
}

func (self *Model) checkAndSetCollectionName() bool {
	if self.collectionName == "" {
		fn, err := self.getMethodViaReflection("GetCollectionName")
		if err != "" {
			panic(err)
		}
		result := fn.Call(nil)

		if len(result) > 0 && result[0].String() != "" {
			self.collectionName = result[0].String()
			return true
		}
		return false
	} else {
		return true
	}
}

func docToBson(doc interface{}) (bsonData bson.M, err error) {
	if bsonData, ok := doc.(bson.M); ok {
		return bsonData, nil
	}
	var tmpBlob []byte
	if tmpBlob, err = bson.Marshal(doc); err != nil {
		return
	}
	if err = bson.Unmarshal(tmpBlob, &bsonData); err != nil {
		return
	}
	return
}

func ToBSON(doc interface{}) (bsonData bson.M, err error) {
	if bsonData, ok := doc.(bson.M); ok {
		return bsonData, nil
	}
	var tmpBlob []byte
	if tmpBlob, err = bson.Marshal(doc); err != nil {
		return
	}
	if err = bson.Unmarshal(tmpBlob, &bsonData); err != nil {
		return
	}
	return
}