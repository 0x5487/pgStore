package main

type Collection struct {
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

	//details
	Name                      string `json:"name"`
	Content                   string `json:"content"`
	Tags                      string `json:"tags"`
	Vendor                    string `json:"vendor"`
	ListPrice                 int64  `json:"list_price"`
	Price                     int64  `json:"price"`
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
	ListPrice             int64  `json:"list_price"`
	Price                 int64  `json:"price"`
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
