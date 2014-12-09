package main

//"fmt"

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

	catalogService := NewCatalogService(dbLayer, jasonStore)

	//create collection
	var collection1 = Collection{Permalink: "men", Name: "men collection"}
	_, err = catalogService.CreateCollection(collection1)

	//create products
	var sku1 = Variation{Sku: "abc001", Price: Money{90, 99}}
	var product1 = Product{Permalink: "men-shoe", Price: Money{91, 89}, Name: "men shoe", Variations: []Variation{sku1}}
	product1.Collections = []int64{1}
	_, err = catalogService.CreateProduct(product1)
	PanicIf(err)

	var sku2 = Variation{Sku: "abc002", ListPrice: Money{92, 0}, Price: Money{93, 0}}
	var sku3 = Variation{Sku: "abc003", ListPrice: Money{94, 0}, Price: Money{95, 0}}
	var product2 = Product{Permalink: "men-shirt", Name: "men shirt", Variations: []Variation{sku2, sku3}}
	product2.Collections = []int64{1}
	_, err = catalogService.CreateProduct(product2)
	PanicIf(err)

	//create order

}
