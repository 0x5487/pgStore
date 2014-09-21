package main

import (
	"fmt"
)

func build_data() {

	db, err := GetDB()
	PanicIf(err)
	defer db.Close()

	//delete jason schema
	err = db.DeleteSchema("jason", "force")
	PanicIf(err)

	schema, err := db.CreateSchema("jason")
	PanicIf(err)

	//create product's collection
	products, err := schema.CreateCollection("products")
	PanicIf(err)

	err = products.CreateIndex("product_sku_index", true, "(data->>'sku')")
	PanicIf(err)

	//insert products
	skus := []Sku{{Sku: "abc123-123"}}
	product1 := Product{Name: "Jason", Skus: skus}

	oid1, err := products.Insert(product1)
	fmt.Println(oid1)
	PanicIf(err)

}
