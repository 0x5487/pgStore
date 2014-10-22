package main

import (
	//"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

func EnableApi(router *gin.Engine) {

	router.Use(AUTH())

	v1 := router.Group("/api/v1")
	{
		v1.GET("/products/:productId", getProductEndpointV1)
		v1.GET("/products", getProductsEndpointV1)
		v1.POST("/products", createProductEndpointV1)
		v1.PUT("/products/:productId", updateProductEndpointV1)
		v1.DELETE("/products/:productId", deleteProductEndpointV1)

		v1.GET("/collections/:collectionId", getCollectionEndpointV1)
		v1.GET("/collections", getCollectionsEndpointV1)
		v1.POST("/collections", createCollectionEndpointV1)
		v1.PUT("/collections/:collectionId", updateCollectionEndpointV1)
		v1.DELETE("/collections/:collectionId", deleteCollectionEndpointV1)

		v1.GET("/orders/:orderId", getOrderEndpointV1)
		v1.POST("/orders", createOrderEndpointV1)
		v1.PUT("/orders/:orderId/lineitems", updateOrderLineItemsEndpointV1)
		v1.PUT("/orders/:orderId/coupons", updateOrderCouponsEndpointV1)
		v1.PUT("/orders/:orderId", updateOrderEndpointV1)
		v1.DELETE("/orders/:orderId", deleteOrderEndpointV1)
	}
}

func AUTH() gin.HandlerFunc {

	return func(c *gin.Context) {
		logInfo("Auth middleware")
		var store = Store{Name: "jason"}
		c.Set("_store", store)

		db, err := GetDB()
		PanicIf(err)
		defer db.Close()

		var dbLayer = new(DbLayer)
		dbLayer.Conn = db
		c.Set("_dbLayer", dbLayer)

		var catalogService = NewCatalogService(dbLayer, store)
		c.Set("_catalogService", catalogService)

		c.Next()
	}
}

func getProductsEndpointV1(c *gin.Context) {

	/*catalogService := c.MustGet("_catalogService").(*CatalogService)

	result, err := catalogService.GetProducts()
	if err != nil {
		c.JSON(500, gin.H{"message": "Internal Server Error", "status": 500})
		return
	}

	if result == nil {
		c.JSON(200, result)
		return
	}

	c.JSON(200, result)*/
}

func getProductEndpointV1(c *gin.Context) {
	param_productId := c.Params.ByName("productId")
	logDebug(fmt.Sprintf("productId: %s", param_productId))

	//validation
	productId, err := ToInt64(param_productId)
	if err != nil {
		c.JSON(404, gin.H{"message": "Not Found", "status": 404})
		return
	}

	catalogService := c.MustGet("_catalogService").(*CatalogService)

	result, err := catalogService.GetProduct(productId)
	if err != nil {
		c.JSON(500, gin.H{"message": "Internal Server Error", "status": 500})
		return
	}

	c.JSON(200, result)
}

func createProductEndpointV1(c *gin.Context) {
	var productJSON Product

	c.Bind(&productJSON)

	c.JSON(200, gin.H{"state": "ok"})
}

func updateProductEndpointV1(c *gin.Context) {

}

func deleteProductEndpointV1(c *gin.Context) {

}

func getCollectionsEndpointV1(c *gin.Context) {
	catalogService := c.MustGet("_catalogService").(*CatalogService)

	result, err := catalogService.GetCollections()
	if err != nil {
		c.JSON(500, gin.H{"message": "Internal Server Error", "status": 500})
		return
	}

	c.JSON(200, result)
}

func getCollectionEndpointV1(c *gin.Context) {
	param_collectionId := c.Params.ByName("collectionId")
	logDebug(fmt.Sprintf("collectionId: %s", param_collectionId))

	//validation
	collectionId, err := ToInt64(param_collectionId)
	if err != nil {
		c.JSON(404, gin.H{"message": "Not Found", "status": 404})
		return
	}

	catalogService := c.MustGet("_catalogService").(*CatalogService)

	result, err := catalogService.GetCollection(collectionId)
	if err != nil {
		c.JSON(500, gin.H{"message": "Internal Server Error", "status": 500})
		return
	}

	c.JSON(200, result)
}

func createCollectionEndpointV1(c *gin.Context) {

}

func updateCollectionEndpointV1(c *gin.Context) {

}

func deleteCollectionEndpointV1(c *gin.Context) {

}

func getOrderEndpointV1(c *gin.Context) {

}

func createOrderEndpointV1(c *gin.Context) {

}

func updateOrderLineItemsEndpointV1(c *gin.Context) {

}

func updateOrderCouponsEndpointV1(c *gin.Context) {

}

func updateOrderEndpointV1(c *gin.Context) {

}

func deleteOrderEndpointV1(c *gin.Context) {

}
