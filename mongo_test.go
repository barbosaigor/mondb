package mondb

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestConn(t *testing.T) {
	d := New("numbers", "testing")
	err := d.Conn(DefaultMongoURL)
	defer d.Discn()
	if err != nil {
		t.Errorf("Conn: Error to connect to database: %s", err)
		return
	}
}

func TestNew(t *testing.T) {
	d := New("numbers", "testing")
	if d == nil {
		t.Error("New: Expected Mongo instance but got nil")
	}
}

func TestNewWithQueryTimeout(t *testing.T) {
	d := NewWithQueryTimeout("numbers", "testing", 2)
	if d == nil {
		t.Error("NewWithQueryTimeout: Expected Mongo instance but got nil")
		return
	}
	if d.QueryTimeout != 2 {
		t.Errorf("NewWithQueryTimeout: Expected %d but got %d", 2, d.QueryTimeout)
	}
}

func TestInsertOne(t *testing.T) {
	d := New("numbers", "testing")
	defer d.Discn()
	if err := d.Conn(DefaultMongoURL); err != nil {
		t.Errorf("InsertOne: Error to connect to database: %s", err)
		return
	}
	if err := d.InsertOne(bson.M{"name": "pi", "value": 3.14159}); err != nil {
		t.Errorf("InsertOne: Error to insert a document to database: %s", err)
	}
}

func TestFindOne(t *testing.T) {
	d := New("numbers", "testing")
	defer d.Discn()
	if err := d.Conn(DefaultMongoURL); err != nil {
		t.Errorf("FindOne: Error to connect to database: %s", err)
		return
	}
	filter := bson.M{"name": "pi"}
	if inst, err := d.FindOne(filter); err == ErrDocumentNotFound {
		t.Errorf("FindOne: Error no document was found")
	} else if err != nil {
		t.Errorf("FindOne: Error to find a document: %s", err)
	} else {
		t.Log(inst)
	}

	filter = bson.M{"value": 3.14159}
	if inst, err := d.FindOne(filter); err == ErrDocumentNotFound {
		t.Errorf("FindOne: Error no document was found")
	} else if err != nil {
		t.Errorf("FindOne: Error to find a document: %s", err)
	} else {
		t.Log(inst)
	}

	filter = bson.M{"value": 3.148}
	if _, err := d.FindOne(filter); err == nil {
		t.Error("FindOne: should return an error but got: nil")
	}
}

func TestFindMany(t *testing.T) {
	d := New("numbers", "testing")
	defer d.Discn()
	if err := d.Conn(DefaultMongoURL); err != nil {
		t.Errorf("Find: Error to connect to database: %s", err)
		return
	}
	filter := bson.M{"name": "pi"}
	if inst, err := d.FindMany(filter); err != nil {
		t.Errorf("Find: Error to find documents: %s", err)
	} else {
		t.Log(inst)
	}
	filter = bson.M{"value": 3.14159}
	if inst, err := d.FindMany(filter); err != nil {
		t.Errorf("Find: Error no document was found")
	} else {
		t.Log(inst)
	}

	filter = bson.M{"value": 3.148}
	if _, err := d.FindMany(filter); err == nil {
		t.Error("Find: should return an error but got: nil")
	}
}

func TestUpdateOne(t *testing.T) {
	d := New("numbers", "testing")
	defer d.Discn()
	if err := d.Conn(DefaultMongoURL); err != nil {
		t.Errorf("UpdateOne: Error to connect to database: %s", err)
		return
	}
	filter := bson.M{"value": 3.14159}
	updated, err := d.UpdateOne(filter, map[string]interface{}{"value": 500})
	if err != nil {
		t.Errorf("UpdateOne: fail to update a document: %s", err)
	}
	if !updated {
		t.Error("UpdateOne: Expected updated true but got false")
	}
}

func TestDeleteOne(t *testing.T) {
	d := New("numbers", "testing")
	defer d.Discn()
	if err := d.Conn(DefaultMongoURL); err != nil {
		t.Errorf("DeleteOne: Error to connect to database: %s", err)
		return
	}
	if err := d.InsertOne(bson.M{"name": "pi", "value": 3.14159}); err != nil {
		t.Errorf("DeleteOne: Error to insert a document to database: %s", err)
		return
	}
	deleted, err := d.DeleteOne(map[string]interface{}{"value": 3.14159})
	if err != nil {
		t.Errorf("DeleteOne: fail to delete a document: %s", err)
	}
	if !deleted {
		t.Errorf("DeleteOne: Expected deleted true but got false")
	}
}
