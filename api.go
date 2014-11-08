package main

import (
	//"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func EnableApi(router *gin.Engine) {

	router.Use(AUTH())

	v1 := router.Group("/api/v1")
	{
		v1.GET("/products/:productId", getProductEndpointV1, getProductCountV1, notFound)
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

func notFound(c *gin.Context) {
	logInfo("[Api] Not Found")
	c.JSON(404, gin.H{"message": "Not Found"})
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

func getProductCountV1(c *gin.Context) {
	logInfo("[Api][start] getProductCountV1")

	param_count := c.Params.ByName("productId")
	if len(param_count) > 0 && param_count != "count" {
		logInfo("[Api][Pass] getProductCountV1")
		c.Next()
		return
	}

	var collectionId_str string

	//get collectionId from querystring
	if len(c.Request.URL.Query()["collection"]) > 0 {
		collectionId_str = c.Request.URL.Query()["collection"][0]
	}

	catalogService := c.MustGet("_catalogService").(*CatalogService)

	var result int
	var err error
	if len(collectionId_str) > 0 {
		logInfo(fmt.Sprintf("collectionid: %s", collectionId_str))

		collectionId, err := ToInt64(collectionId_str)
		if err != nil {
			logError(err.Error())
			c.JSON(404, gin.H{"message": "Not Found"})
			c.Fail(404, err)
		}

		result, err = catalogService.GetProductCount(collectionId)
	} else {
		result, err = catalogService.GetProductCount()
	}

	if err != nil {
		logError(err.Error())
		c.JSON(500, gin.H{"message": "Internal Server Error"})
		c.Fail(500, err)
		return
	}

	c.JSON(200, result)
	c.Abort(-1)
	logInfo("[Api][end] getProductCountV1")
}

func getProductsEndpointV1(c *gin.Context) {
	logInfo("[API]getProductsEndpointV1")

	catalogService := c.MustGet("_catalogService").(*CatalogService)

	result, err := catalogService.GetProducts()
	if err != nil {
		logError(err.Error())
		c.JSON(500, gin.H{"message": "Internal Server Error"})
		return
	}

	c.JSON(200, result)
}

func getProductEndpointV1(c *gin.Context) {
	logInfo("[Api] getProductEndpointV1")

	param_productId := c.Params.ByName("productId")
	logInfo(fmt.Sprintf(" productId param: %s", param_productId))

	//validation
	productId, err := ToInt64(param_productId)
	if err != nil {
		logInfo("[Api][Pass] getProductEndpointV1")
		c.Next()
		return
	}

	catalogService := c.MustGet("_catalogService").(*CatalogService)

	result, err := catalogService.GetProduct(productId)
	if err != nil {
		logError(err.Error())
		c.JSON(500, gin.H{"message": "Internal Server Error"})
		c.Fail(500, err)
		return
	}

	c.JSON(200, result)
	c.Abort(-1)
}

func createProductEndpointV1(c *gin.Context) {
	logInfo("[Api]createProductEndpointV1")

	//bind JSON
	var product = Product{}
	if !c.BindWith(&product, binding.JSON) {
		c.JSON(400, gin.H{"message": "json format is invalid."})
		return
	}

	//validation
	if len(product.Name) == 0 {
		c.JSON(400, gin.H{"message": "name is required"})
		return
	}

	//act
	catalogService := c.MustGet("_catalogService").(*CatalogService)
	productId, err := catalogService.CreateProduct(product)
	if err != nil {
		logError(err.Error())

		if v, ok := err.(appError); ok {
			c.JSON(500, gin.H{"message": v.Error()})
			return
		}

		c.JSON(500, gin.H{"message": "Internal Server Error"})
		return
	}

	c.JSON(201, gin.H{"Location": fmt.Sprintf("/products/%d", productId)})
}

func updateProductEndpointV1(c *gin.Context) {
	logInfo("[Api]updateProductEndpointV1")

	//bind JSON
	var product = Product{}
	if !c.BindWith(&product, binding.JSON) {
		c.JSON(400, gin.H{"message": "json format is invalid."})
		return
	}

	//validation
	if len(product.Name) == 0 {
		c.JSON(400, gin.H{"message": "name is required"})
		return
	}

	//act
	catalogService := c.MustGet("_catalogService").(*CatalogService)
	err := catalogService.UpdateProduct(product)
	if err != nil {
		logError(err.Error())

		if v, ok := err.(appError); ok {
			c.JSON(500, gin.H{"message": v.Error()})
			return
		}

		c.JSON(500, gin.H{"message": "Internal Server Error"})
		return
	}

	c.JSON(200, "")
}

func deleteProductEndpointV1(c *gin.Context) {

}

func getCollectionsEndpointV1(c *gin.Context) {
	catalogService := c.MustGet("_catalogService").(*CatalogService)

	result, err := catalogService.GetCollections()
	if err != nil {
		logError(err.Error())
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
		logError(err.Error())
		c.JSON(404, gin.H{"message": "Not Found", "status": 404})
		return
	}

	catalogService := c.MustGet("_catalogService").(*CatalogService)

	result, err := catalogService.GetCollection(collectionId)
	if err != nil {
		logError(err.Error())
		c.JSON(500, gin.H{"message": "Internal Server Error", "status": 500})
		return
	}

	c.JSON(200, result)
}

func createCollectionEndpointV1(c *gin.Context) {
	logInfo("creating collection")

	//bind JSON
	var collection = Collection{}
	if !c.BindWith(&collection, binding.JSON) {
		c.Abort(400)
		return
	}

	//validation
	if len(collection.Name) == 0 {
		c.JSON(400, gin.H{"message": "name is required"})
		return
	}

	catalogService := c.MustGet("_catalogService").(*CatalogService)
	collectionId, err := catalogService.CreateCollection(collection)
	if err != nil {
		logError(err.Error())
		c.JSON(500, gin.H{"message": "Internal Server Error"})
		return
	}

	c.JSON(201, gin.H{"Location": fmt.Sprintf("/collections/%d", collectionId)})
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
