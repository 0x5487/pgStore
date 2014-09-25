package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type DbWrapper interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
}

type DbLayer struct {
	Conn interface{}
}

type JSchema struct {
	DB   *DbLayer
	Name string
}

type JCollection struct {
	Schema *JSchema
	Name   string
}

type JDocument struct {
	id   int
	data string
}

func GetDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "host=localhost dbname=ecStore user=postgres sslmode=disable")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (source *DbLayer) GetSchema(name string) (*JSchema, error) {
	var conn = source.Conn.(DbWrapper)
	var get_schema_sql string = "select count(schema_name) from information_schema.schemata where schema_name = $1"

	var count int
	err := conn.QueryRow(get_schema_sql, name).Scan(&count)
	if err != nil {
		return nil, err
	}

	if count < 1 {
		return nil, nil
	}

	schema := JSchema{DB: source, Name: name}
	return &schema, nil
}

func (source *DbLayer) CreateSchema(name string) (*JSchema, error) {
	var db = source.Conn.(DbWrapper)
	var create_schema_sql string = fmt.Sprintf("CREATE SCHEMA %s", name)

	stmt, err := db.Prepare(create_schema_sql)
	PanicIf(err)
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}

	var schema JSchema = JSchema{DB: source, Name: name}
	return &schema, nil
}

func (source *DbLayer) DeleteSchema(name string, params ...string) error {
	var db = source.Conn.(DbWrapper)
	var delete_schema_sql string = fmt.Sprintf("DROP SCHEMA %s ", name)

	if len(params) > 0 {
		switch params[0] {
		case "force":
			delete_schema_sql = fmt.Sprintf("DROP SCHEMA %s CASCADE", name)
		}
	}

	logInfo(delete_schema_sql)
	_, err := db.Exec(delete_schema_sql)
	if err != nil {
		return err
	}

	return nil
}

func (source *JSchema) GetCollection(name string) (*JCollection, error) {
	collection := &JCollection{Schema: source, Name: name}
	return collection, nil
}

func (source *JSchema) CreateCollection(name string) (*JCollection, error) {
	var db = source.DB.Conn.(DbWrapper)
	var createTableSQL string = fmt.Sprintf("CREATE TABLE %s.%s (id serial primary key, data jsonb)", source.Name, name)

	logInfo(createTableSQL)
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	collection := JCollection{Schema: source, Name: name}
	return &collection, nil
}

func (source *JSchema) DeleteCollection(name string) error {
	var db = source.DB.Conn.(DbWrapper)
	var deleteTableSQL string = fmt.Sprintf("DROP TABLE %s", name)

	logInfo(deleteTableSQL)
	_, err := db.Exec(deleteTableSQL)
	if err != nil {
		return err
	}

	return nil
}

func (source *JCollection) CreateIndex(indexName string, isUnique bool, argu string) error {
	var db = source.Schema.DB.Conn.(DbWrapper)
	var createIndexSQL string

	if isUnique {
		createIndexSQL = fmt.Sprintf("CREATE UNIQUE INDEX %s ON %s.%s( %s )", indexName, source.Schema.Name, source.Name, argu)
	} else {
		createIndexSQL = fmt.Sprintf("CREATE INDEX %s ON %s.%s( %s )", indexName, source.Schema.Name, source.Name, argu)
	}

	logInfo(createIndexSQL)
	_, err := db.Exec(createIndexSQL)
	if err != nil {
		return err
	}

	return nil
}

func (source *JCollection) Insert(target interface{}) (int, error) {
	var db = source.Schema.DB.Conn.(DbWrapper)
	targetJSON, err := toJSON(target)
	//var insertSQL = fmt.Sprintf("INSERT INTO %s.%s (data) VALUES ('%s')", source.Schema.Name, source.Name, targetJSON)

	var insertSQL = fmt.Sprintf("INSERT INTO %s.%s (data) VALUES ($1) RETURNING id", source.Schema.Name, source.Name)
	logInfo(insertSQL)
	logInfo(targetJSON)
	var id int
	err = db.QueryRow(insertSQL, targetJSON).Scan(&id)

	PanicIf(err)

	return id, nil
}

func (source *JCollection) FindOne(query string, target interface{}) (int, error) {
	var db = source.Schema.DB.Conn.(DbWrapper)

	var querySQL string = fmt.Sprintf("SELECT * FROM %s.%s WHERE (data @> '%s') limit 1;", source.Schema.Name, source.Name, query)

	var (
		id   int
		json string
	)

	logInfo(querySQL)
	err := db.QueryRow(querySQL).Scan(&id, &json)
	logInfo(fmt.Sprintf("id: %d, json: %s", id, json))

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	return id, nil
}

func (source *JCollection) Remove(id int) error {
	var db = source.Schema.DB.Conn.(DbWrapper)

	var deleteSQL string = fmt.Sprintf("DELETE FORM %s.%s WHERE id=%d;", source.Schema.Name, source.Name, id)

	logInfo(deleteSQL)
	_, err := db.Exec(deleteSQL)
	if err != nil {
		return err
	}

	return nil
}
