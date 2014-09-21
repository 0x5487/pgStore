package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const ()

type JDB struct {
	*sql.DB
}

type JSchema struct {
	db   *JDB
	Name string
}

type JCollection struct {
	Schema *JSchema
	Name   string
}

func GetDB() (*JDB, error) {
	conn, err := sql.Open("postgres", "host=localhost dbname=ecStore user=postgres sslmode=disable")
	if err != nil {
		return nil, err
	}
	db := JDB{conn}
	return &db, nil
}

func (source *JDB) GetSchema(name string) (*JSchema, error) {
	var get_schema_sql string = "select count(schema_name) from information_schema.schemata where schema_name = $1"

	var count int
	err := source.QueryRow(get_schema_sql, name).Scan(&count)
	if err != nil {
		return nil, err
	}

	if count < 1 {
		return nil, nil
	}

	schema := JSchema{db: source, Name: name}
	return &schema, nil
}

func (source *JDB) CreateSchema(name string) (*JSchema, error) {
	var create_schema_sql string = fmt.Sprintf("CREATE SCHEMA %s", name)

	stmt, err := source.Prepare(create_schema_sql)
	PanicIf(err)
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}

	var schema JSchema = JSchema{db: source, Name: name}
	return &schema, nil
}

func (source *JDB) DeleteSchema(name string, params ...string) error {
	var delete_schema_sql string = fmt.Sprintf("DROP SCHEMA %s ", name)

	if len(params) > 0 {
		switch params[0] {
		case "force":
			delete_schema_sql = fmt.Sprintf("DROP SCHEMA %s CASCADE", name)
		}
	}

	logInfo(delete_schema_sql)
	_, err := source.Exec(delete_schema_sql)
	if err != nil {
		return err
	}

	return nil
}

func (source *JSchema) CreateCollection(name string) (*JCollection, error) {
	var createTableSQL string = fmt.Sprintf("CREATE TABLE %s.%s (oid serial primary key, data jsonb)", source.Name, name)

	logInfo(createTableSQL)
	_, err := source.db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	collection := JCollection{Schema: source, Name: name}
	return &collection, nil
}

func (source *JSchema) DeleteCollection(name string) error {
	var deleteTableSQL string = fmt.Sprintf("DROP TABLE %s", name)

	logInfo(deleteTableSQL)
	_, err := source.db.Exec(deleteTableSQL)
	if err != nil {
		return err
	}

	return nil
}

func (source *JCollection) CreateIndex(indexName string, isUnique bool, argu string) error {
	var createIndexSQL string

	if isUnique {
		createIndexSQL = fmt.Sprintf("CREATE UNIQUE INDEX %s ON %s.%s( %s )", indexName, source.Schema.Name, source.Name, argu)
	} else {
		createIndexSQL = fmt.Sprintf("CREATE INDEX %s ON %s.%s( %s )", indexName, source.Schema.Name, source.Name, argu)
	}

	logInfo(createIndexSQL)
	_, err := source.Schema.db.Exec(createIndexSQL)
	if err != nil {
		return err
	}

	return nil
}

func (source *JCollection) Insert(target interface{}) (int64, error) {
	db := source.Schema.db
	targetJSON, err := toJSON(target)
	//var insertSQL = fmt.Sprintf("INSERT INTO %s.%s (data) VALUES ('%s')", source.Schema.Name, source.Name, targetJSON)

	var insertSQL = fmt.Sprintf("INSERT INTO %s.%s (data) VALUES ($1) RETURNING oid", source.Schema.Name, source.Name)
	logInfo(insertSQL)

	var oid int64
	err = db.QueryRow(insertSQL, targetJSON).Scan(&oid)
	if err != nil {
		PanicIf(err)
	}

	return oid, nil
}
