package main

import (
	"errors"
	"fmt"
)

type Store struct {
	Name string `json:"name"`
}

type AdminService struct {
	DB *DbLayer
}

func NewAdminService(dbLayer *DbLayer) *AdminService {
	service := AdminService{DB: dbLayer}
	return &service
}

func (source *AdminService) GetStore(name string) (int, *Store, error) {
	schema := JSchema{DB: source.DB, Name: "admin"}
	stores := JCollection{Schema: &schema, Name: "stores"}

	var queryName = fmt.Sprintf("{\"name\":\"%s\"}", name)

	var target *Store
	id, err := stores.FindOne(queryName, target)

	if err != nil {
		return 0, nil, err
	}

	return id, target, nil
}

func (source *AdminService) CreateStore(store Store) (int, error) {
	//validate store object
	if len(store.Name) == 0 {
		return 0, errors.New("store name can't be empty")
	}

	//insert the store into admin
	db, err := GetDB()
	if err != nil {
		return 0, err
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	var dbLayer = new(DbLayer)
	dbLayer.Conn = tx

	schema := &JSchema{DB: dbLayer, Name: "admin"}
	stores := &JCollection{Schema: schema, Name: "stores"}

	id, err := stores.Insert(store)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	//create store schema
	schema, err = dbLayer.CreateSchema(store.Name)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	products, err := schema.CreateCollection("products")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = products.CreateIndex("unique_resourceId_idx", true, "(data->>'resource_id')")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	return id, nil
}

func (source *AdminService) DeleteStore(id int) error {
	//validate store object
	if id <= 0 {
		return errors.New("id can't be less than 0")
	}

	//create transaction
	db, err := GetDB()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	var dbLayer = new(DbLayer)
	dbLayer.Conn = tx

	//delete the store from store's collection
	schema := &JSchema{DB: dbLayer, Name: "admin"}
	stores := &JCollection{Schema: schema, Name: "stores"}
	err = stores.Remove(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	//delete the store's schema
	err = dbLayer.DeleteSchema("jason", "force")
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}