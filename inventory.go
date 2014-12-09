package main

type Warehouse struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type StockItem struct {
	Id          int64 `json:"id"`
	WarehouseId int64 `json:"warehouse_id"`
	VariationId int64 `json:"variation_id"`
	OrderId     int64 `json:"order_id"`
}
