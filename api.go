package main

import (
	"github.com/gin-gonic/gin"
)

func EnableApi(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.GET("/products/:productId", getProductsEndpointV1)
		v1.GET("/products", getProductEndpointV1)
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

func Api() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func getProductsEndpointV1(c *gin.Context) {

}

func getProductEndpointV1(c *gin.Context) {
	//productId := c.Params.ByName("productId")
}

func createProductEndpointV1(c *gin.Context) {

}

func updateProductEndpointV1(c *gin.Context) {

}

func deleteProductEndpointV1(c *gin.Context) {

}

func getCollectionsEndpointV1(c *gin.Context) {

}

func getCollectionEndpointV1(c *gin.Context) {
	//productId := c.Params.ByName("productId")
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
