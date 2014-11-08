package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type ValidationError struct {
	Message    string                 `json:"message"`
	ModelState map[string]interface{} `json:"ModelState"`
}

type Key struct {
	Sku string `json:"sku"`
}

type Image struct {
	Url string `json:"url"`
}

type IFilter interface {
	getSQL() string
}

type Collection struct {
	Id              int64         `json:"id"`
	Name            string        `json:"name"`
	IsVisible       bool          `json:"is_visible"`
	Content         string        `json:"content"`
	Image           interface{}   `json:"image"`
	Tags            []string      `json:"tags"`
	SortOrder       int           `json:"sort_order"`
	UrlName         string        `json:"url_name"`
	PageTitle       string        `json:"page_title"`
	MetaDescription string        `json:"meta_description"`
	CustomFields    []CustomField `json:"custom_fields"`
	Products        []int64       `json:"products"`
}

type Product struct {
	Id int64 `json:"id"`

	//details
	Name                      string   `json:"name"`
	Content                   string   `json:"content"`
	Tags                      []string `json:"tags"`
	Vendor                    string   `json:"vendor"`
	ListPrice                 Money    `json:"list_price"`
	Price                     Money    `json:"price"`
	Weight                    int      `json:"weight"`
	SortOrder                 int      `json:"sort_order"`
	IsPurchasable             bool     `json:"is_purchaseable"`
	IsVisible                 bool     `json:"is_visible"`
	IsBackOrderEnabled        bool     `json:"is_backorder_enabled"`
	IsPreOrderEnabled         bool     `json:"is_preorder_enabled"`
	IsShippingAddressRequired bool     `json:"is_shipping_address_required"`

	//sku
	Skus    []Sku    `json:"skus"`
	Options []Option `json:"options"`

	//seo
	UrlName         string `json:"url_name"`
	PageTitle       string `json:"page_title"`
	MetaDescription string `json:"meta_description"`

	Images       []Image       `json:"images"`
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
	Label  string   `json:"label"`
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

type CollectionFilter struct {
	Ids []int64
}

func (source *CollectionFilter) getSQL() string {
	var sql string

	if source.Ids != nil {
		var all_ids []string
		for _, value := range source.Ids {
			id_str := strconv.FormatInt(value, 10)
			all_ids = append(all_ids, id_str)
		}
		ids_str := strings.Join(all_ids, ",")
		sql += fmt.Sprintf("( id in %s)", ids_str)
	}

	return sql
}

func NewCatalogService(dbLayer *DbLayer, store Store) *CatalogService {
	schema, err := dbLayer.GetSchema(store.Name)
	if err != nil {
		logError(err.Error())
	}

	service := new(CatalogService)
	service.DB = dbLayer
	service.Schema = schema
	service.Store = store

	return service
}

func (source *CatalogService) GetProductCount(collectionIds ...int64) (int, error) {
	schema := source.Schema

	productsCollection, err := schema.GetCollection("products")
	if err != nil {
		return 0, err
	}

	if len(collectionIds) > 0 {
		collectionId := collectionIds[0]
		countSQL := fmt.Sprintf("(data @> '{\"collections\":[%d]}')", collectionId)
		return productsCollection.Count(countSQL)
	}

	return productsCollection.Count()
}

// Error Code
// 1: collection is not existing
func (source *CatalogService) CreateProduct(product Product) (int64, error) {
	//validate product
	if len(product.Name) <= 0 {
		return 0, errors.New("name can't be empty")
	}

	if len(product.Skus) > 0 {
		for _, sku := range product.Skus {
			if len(sku.Sku) <= 0 {
				return 0, errors.New("The sku field can't be empty")
			}
		}
	}

	for _, value := range product.Collections {
		collection, err := source.GetCollection(value)
		if err != nil {
			return 0, err
		}

		if collection == nil {
			vErr := appError{}
			vErr.Code = 1
			vErr.Message = fmt.Sprintf("collection:%d is not existing", value)
			return 0, vErr
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

	productJSON, err := toJSON(&product)
	if err != nil {
		return 0, err
	}

	msg_log := fmt.Sprintf("[Insert Product] %s", productJSON)
	logDebug(msg_log)

	doc := new(JDocument)
	doc.data = productJSON

	err = products.Insert(doc)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return doc.id, nil
}

// Error Code
// 1: product doesn't exist
// 2: collection doesn't exist
// 3: variation has already existed
func (source *CatalogService) UpdateProduct(product Product) error {
	//validate product
	if len(product.Name) <= 0 {
		return errors.New("name can't be empty")
	}

	vErr := appError{}
	oldProduct, err := source.GetProduct(product.Id)
	if err != nil {
		return err
	}

	if oldProduct == nil {
		vErr.Code = 1
		vErr.Message = fmt.Sprintf("product:%d doesn't exist", product.Id)
		return vErr
	}

	for _, value := range product.Collections {
		collection, err := source.GetCollection(value)
		if err != nil {
			return err
		}

		if collection == nil {
			vErr.Code = 2
			vErr.Message = fmt.Sprintf("collection:%d doesn't exist", value)
			return vErr
		}
	}

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

	schema := &JSchema{DB: dbLayer, Name: source.Store.Name}
	unique_skus, err := schema.GetCollection("unique_skus")
	if err != nil {
		return err
	}

	for _, sku := range product.Skus {
		if len(sku.Sku) <= 0 {
			tx.Rollback()
			return errors.New("The sku field can't be empty")
		}

		var isFound bool
		for _, old_variation := range oldProduct.Skus {
			if sku.Sku == old_variation.Sku {
				isFound = true
			}
		}

		if isFound == false {
			findSQL := fmt.Sprintf("{\"sku\":\"%s\"}", sku.Sku)
			target, err := unique_skus.FindOne(findSQL)
			if err != nil {
				tx.Rollback()
				return err
			}

			if target != nil {
				tx.Rollback()
				vErr.Code = 3
				vErr.Message = fmt.Sprintf("variation:%s has already existed", sku.Sku)
				return vErr
			}

			key := Key{Sku: sku.Sku}
			json, err := toJSON(key)
			if err != nil {
				tx.Rollback()
				return err
			}
			doc := new(JDocument)
			doc.data = json

			err = unique_skus.Insert(doc)
			if err != nil {
				tx.Rollback()
				return err
			}

		}
	}

	//update the product
	products, err := schema.GetCollection("products")
	if err != nil {
		tx.Rollback()
		return err
	}

	productJSON, err := toJSON(&product)
	if err != nil {
		tx.Rollback()
		return err
	}

	msg_log := fmt.Sprintf("[Update Product] %s", productJSON)
	logDebug(msg_log)

	doc := new(JDocument)
	doc.id = product.Id
	doc.data = productJSON

	err = products.Update(doc)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (source *CatalogService) GetProducts() (*[]Product, error) {
	schema := source.Schema

	productsCollection, err := schema.GetCollection("products")
	if err != nil {
		return nil, err
	}

	docs, err := productsCollection.Find()
	if err != nil {
		return nil, err
	}

	if len(*docs) <= 0 {
		return nil, nil
	}

	result := []Product{}

	for _, doc := range *docs {
		product := new(Product)
		err = fromJSON(product, doc.data)
		if err != nil {
			return nil, err
		}
		product.Id = doc.id
		result = append(result, *product)
	}

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

	result := new(Product)
	err = fromJSON(result, doc.data)
	if err != nil {
		return nil, err
	}
	result.Id = doc.id
	return result, nil
}

func (source *CatalogService) CreateCollection(collection Collection) (int64, error) {
	//validate collection
	if len(collection.Name) <= 0 {
		return 0, errors.New("name can't be empty")
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

func (source *CatalogService) GetCollections(filters ...IFilter) (*[]Collection, error) {
	schema := source.Schema

	collections, err := schema.GetCollection("collections")
	if err != nil {
		return nil, err
	}

	docs := new([]JDocument)

	if len(filters) > 0 {
		filter := filters[0]
		docs, err = collections.Find(filter.getSQL())
	} else {
		docs, err = collections.Find()
	}

	if err != nil {
		return nil, err
	}

	if len(*docs) <= 0 {
		return nil, nil
	}

	result := []Collection{}

	for _, doc := range *docs {
		collection := new(Collection)
		err = fromJSON(collection, doc.data)
		if err != nil {
			return nil, err
		}
		collection.Id = doc.id
		result = append(result, *collection)
	}

	return &result, nil
}

func (source *CatalogService) GetCollection(collectionId int64) (*Collection, error) {
	schema := source.Schema

	collections, err := schema.GetCollection("collections")
	if err != nil {
		return nil, err
	}

	doc, err := collections.FindById(collectionId)
	if err != nil {
		return nil, err
	}

	if doc == nil {
		return nil, nil
	}

	result := new(Collection)
	err = fromJSON(result, doc.data)
	if err != nil {
		return nil, err
	}
	result.Id = doc.id
	return result, nil
}
