# MonDB
MonDB library defines an interface for mongodb driver.  
It makes it easy to operate the database only with Go native types.  

# Installation
```bash 
go get github.com/barbosaigor/mondb
```   

## Usage
The user must make a connection to database and after all close it.  
```golang
db := New("dealership", "cars")
if err := db.Conn(DefaultMongoURL); err != nil {
    // ...
}
defer db.Discn()
// Now you can have fun
```  

e.g Find documents  
```golang
db := New("dealership", "cars")
// ...
filter := map[string]interface{}{"year": 1964}
cars, err := db.FindMany(filter)
```  

e.g Find a document  
```golang
filter := map[string]interface{}{"name": "A car", "price": 1500.7}
car, err := db.FindOne(filter)
```  

e.g Insert a document  
```golang
car := map[string]interface{}{
    "name": "A car", 
    "price": 1000.49, 
    "year": 1943
}
err := db.InsertOne(car)
```  

e.g update params of a document  
```golang
filter := map[string]interface{}{"name": "A car"}
car := map[string]interface{}{
    "name": "The car", 
    "year": 2001
}
wasUpdated, err := db.UpdateOne(filter, car)
```  

e.g delete a document  
```golang
filter := map[string]interface{}{"name": "The car"}
wasDeleted, err := db.DeleteOne(filter)
```  

For more information check out [documentation](https://pkg.go.dev/github.com/barbosaigor).  