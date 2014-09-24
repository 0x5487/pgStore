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

	var myStore = Store{Name: "jason"}
	_, err = adminService.CreateStore(myStore)
	PanicIf(err)
}
