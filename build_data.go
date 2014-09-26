package main

import (
//"fmt"
)

type aaa struct {
	Sku string `json:"sku"`
}

func build_data() {

	db, err := GetDB()
	PanicIf(err)
	defer db.Close()

	var dbLayer = new(DbLayer)
	dbLayer.Conn = db

	//delete admin and jason schema
	jason_schema, err := dbLayer.GetSchema("jason")
	if jason_schema != nil {
		dbLayer.DeleteSchema("jason", "force")
	}

	admin_schema, err := dbLayer.GetSchema("admin")
	if admin_schema != nil {
		err = dbLayer.DeleteSchema("admin", "force")
		PanicIf(err)
	}

	admin_schema, err = dbLayer.CreateSchema("admin")
	PanicIf(err)

	stores, err := admin_schema.CreateCollection("stores")
	PanicIf(err)

	err = stores.CreateIndex("unique_storename_idx", true, "(data->>'name')")
	PanicIf(err)

	//create jason store
	adminService := NewAdminService(dbLayer)

	var jasonStore = Store{Name: "jason"}
	_, err = adminService.CreateStore(jasonStore)
	PanicIf(err)

	//create products
	catalogService, err := NewCatalogService(dbLayer, jasonStore)
	PanicIf(err)

	var sku1 = Sku{Sku: "abc001", ListPrice: 100, Price: 90}
	var product1 = Product{ResourceId: "men-shoe", Name: "men shoe", Skus: []Sku{sku1}}
	_, err = catalogService.InsertProduct(product1)
	PanicIf(err)

	var sku2 = Sku{Sku: "abc002", ListPrice: 120, Price: 90}
	var sku3 = Sku{Sku: "abc003", ListPrice: 110, Price: 10050000}
	var product2 = Product{ResourceId: "men-shirt", Name: "men shirt", Skus: []Sku{sku2, sku3}}
	_, err = catalogService.InsertProduct(product2)
	PanicIf(err)

	//create collection
	var collection1 = Collection{ResourceId: "men", Name: "men collection"}
	_, err = catalogService.InsertCollection(collection1)

}
