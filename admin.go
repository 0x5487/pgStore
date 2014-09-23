package main

import (
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

	var target Store
	id, err := stores.FindOne(queryName, &target)

	if err != nil {
		return 0, nil, err
	}

	return id, &target, nil
}

func (source *AdminService) CreateStore(store Store) (int, error) {
	return 0, nil
}

func (source *AdminService) DeleteStore(name string) error {
	//schema := JSchema{DB: source.DB, Name: "admin"}

	//id, store, err := source.GetStore(name)
	//if err != nil {
	//	return err
	//}

	//delete the store from store's collection

	//delete the store's schema

	//db.DeleteSchema("jason", "force")
	return nil
}
