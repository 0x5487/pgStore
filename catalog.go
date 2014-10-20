package main

import (
	"errors"
)

type Key struct {
	Sku string `json:"sku"`
}

type Collection struct {
	Id              int64         `json:"id"`
	ResourceId      string        `json:"resource_id"`
	Name            string        `json:"name"`
	IsVisible       bool          `json:"is_visible"`
	Content         string        `json:"content"`
	Image           interface{}   `json:"image"`
	Tags            []string      `json:"tags"`
	SortOrder       int           `json:"sort_order"`
	PageTitle       string        `json:"page_title"`
	MetaDescription string        `json:"meta_description"`
	CustomFields    []CustomField `json:"custom_fields"`
	Products        []int         `json:"products"`
}

type Product struct {
	Id int64 `json:"id"`

	//details
	Name                      string `json:"name"`
	Content                   string `json:"content"`
	Tags                      string `json:"tags"`
	Vendor                    string `json:"vendor"`
	ListPrice                 Money  `json:"list_price"`
	Price                     Money  `json:"price"`
	Weight                    int    `json:"weight"`
	SortOrder                 int    `json:"sort_order"`
	IsPurchasable             bool   `json:"is_purchaseable"`
	IsVisible                 bool   `json:"is_visible"`
	IsBackOrderEnabled        bool   `json:"is_backorder_enabled"`
	IsPreOrderEnabled         bool   `json:"is_preorder_enabled"`
	IsShippingAddressRequired bool   `json:"is_shipping_address_required"`

	//sku
	Skus    []Sku    `json:"skus"`
	Options []Option `json:"options"`

	//seo
	ResourceId      string `json:"resource_id"`
	PageTitle       string `json:"page_title"`
	MetaDescription string `json:"meta_description"`

	Images       interface{}   `json:"images"`
	CustomFields []CustomField `json:"custom_fields"`
	Collections  []int64       `json:"collections"`
}

type Sku struct {
	Sku                   string `json:"sku"`
	ListPrice             Money  `json:"list_price"`
	Price                 Money  `json:"price"`
	SortOrder             int    `json:"sort_order"`
	InventoryQuantity     int    `json:"inventory_quantity"`
	ManageInventoryMethod int    `json:"manage_inventory_method"`
	Weight                int    `json:"weight"`
}

type Option struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type CustomField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CatalogService struct {
	DB     *DbLayer
	Schema *JSchema
	Store  Store
}

func NewCatalogService(dbLayer *DbLayer, store Store) (*CatalogService, error) {
	schema, err := dbLayer.GetSchema(store.Name)
	if err != nil {
		return nil, err
	}

	service := &CatalogService{dbLayer, schema, store}
	return service, nil
}

func (source *CatalogService) InsertProduct(product Product) (int64, error) {
	//validate product
	if len(product.Name) <= 0 {
		return 0, errors.New("name can't be empty")
	}

	if len(product.Skus) <= 0 {
		return 0, errors.New("no skus with the product")
	} else {
		for _, sku := range product.Skus {
			if len(sku.Sku) <= 0 {
				return 0, errors.New("The sku field can't be empty")
			}
		}
	}

	//insert the product
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

	schema := &JSchema{DB: dbLayer, Name: source.Store.Name}
	unique_skus, err := schema.GetCollection("unique_skus")
	if err != nil {
		return 0, err
	}

	for _, sku := range product.Skus {
		key := Key{Sku: sku.Sku}

		json, err := toJSON(key)
		if err != nil {
			return 0, err
		}
		doc := new(JDocument)
		doc.data = json

		err = unique_skus.Insert(doc)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	products, err := schema.GetCollection("products")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	json, err := toJSON(product)
	if err != nil {
		return 0, err
	}
	doc := new(JDocument)
	doc.data = json

	err = products.Insert(doc)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return doc.id, nil
}

func (source *CatalogService) InsertCollection(collection Collection) (int64, error) {
	//validate collection
	if len(collection.Name) <= 0 {
		return 0, errors.New("name can't be empty")
	}
	if len(collection.ResourceId) <= 0 {
		return 0, errors.New("resource_id can't be empty")
	}

	//insert the collection
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

	schema := &JSchema{DB: dbLayer, Name: source.Store.Name}

	collections, err := schema.GetCollection("collections")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	json, err := toJSON(collection)
	if err != nil {
		return 0, err
	}
	doc := new(JDocument)
	doc.data = json

	err = collections.Insert(doc)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return doc.id, nil
}

func (source *CatalogService) GetCollections() (*[]Collection, error) {
	schema := source.Schema

	collections, err := schema.GetCollection("collections")
	if err != nil {
		return nil, err
	}

	docs, err := collections.Find()
	if err != nil {
		return nil, err
	}

	if len(*docs) <= 0 {
		return nil, nil
	}

	result := []Collection{}

	return &result, nil
}

func (source *CatalogService) GetProduct(productId int64) (*Product, error) {
	schema := source.Schema

	products, err := schema.GetCollection("products")
	if err != nil {
		return nil, err
	}

	doc, err := products.FindById(productId)
	if err != nil {
		return nil, err
	}

	if doc == nil {
		return nil, nil
	}

	result := Product{}
	err = fromJSON(&result, doc.data)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
