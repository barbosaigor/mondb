package mondb

// Interface defines an interface for nosql document databases.
type Interface interface {
	// FindMany using a filter find all documents that is matching
	FindMany(filter map[string]interface{}) ([]map[string]interface{}, error)
	// FindOne using a filter finds a document that matches the filter
	FindOne(filter map[string]interface{}) (map[string]interface{}, error)
	// InsertOne inserts a new document to the database
	InsertOne(obj map[string]interface{}) error
	// UpdateOne updates an existing document in the database which matches filter.
	// If no document was found to update, then it returns false and nil.
	// If there is an error, then it returns false and the error.
	UpdateOne(filter map[string]interface{}, obj map[string]interface{}) (bool, error)
	// DeleteOne deletes a document which matches filter
	// If no document was found to delete, then it returns false and nil.
	// If there is an error, then it returns false and the error.
	DeleteOne(filter map[string]interface{}) (bool, error)
	// Conn creates a connection to database
	Conn(url string) error
	// Discn closes the connection to database
	Discn()
}
