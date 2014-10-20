package main

import (
	"time"
)

type Order struct {
	Id                        int64         `json:"id"`
	ShopperId                 int           `json:"shopper_id"`
	LineItems                 []LineItem    `json:"line_items"`
	SubTotalExcTax            Money         `json:"subtotal_exc_tax"`
	SubTotalIncTax            Money         `json:"subtotal_inc_tax"`
	SubTotalTax               Money         `json:"subtotal_tax"`
	ShippingOptionId          int           `json:"shipping_optionid"`
	ShippingFeeExcTax         Money         `json:"shipping_fee_exc_tax"`
	ShippingFeeIncTax         Money         `json:"shipping_fee_inc_tax"`
	ShippingFeeTax            Money         `json:"shipping_fee_tax"`
	HandlingFeeExcTax         Money         `json:"handling_fee_exc_tax"`
	HandlingFeeIncTax         Money         `json:"handling_fee_inc_tax"`
	HandlingFeeTax            Money         `json:"handling_fee_tax"`
	Coupons                   []Coupon      `json:"coupons"`
	Discount                  Money         `json:"discount"`
	TotalExculsiveTax         Money         `json:"total_exc_tax"`
	TotalInclusiveTax         Money         `json:"total_inc_tax"`
	TotalTax                  Money         `json:"total_tax"`
	ItemsCount                int           `json:"items_count"`
	ItemsShipped              int           `json:"items_shipped"`
	PaymentMethod             string        `json:"payment_method"`
	PaymentState              int           `json:"payment_state"`
	BillingAddress            ContactInfo   `json:"billing_address"`
	ShippingAddress           ContactInfo   `json:"shipping_address"`
	IsShippingAddressRequired bool          `json:"is_shipping_address_required"`
	StateId                   int           `json:"state_id"`
	State                     string        `json:"status"`
	CusomterMessage           string        `json:"customer_message"`
	IPAddress                 string        `json:"ip_address"`
	CurrencyCode              string        `json:"currency_code"`
	IsTestOrder               bool          `json:"is_test_order"`
	PaidDateUtc               time.Time     `json:"paid_date_utc"`
	CustomFields              []CustomField `json:"custom_fields"`
	Notes                     string        `json:"notes"`
	CreatedAt                 time.Time     `json:"created_at"`
	UpdatedAt                 time.Time     `json:"updated_at"`
}

type LineItem struct {
	ProductId            int           `json:"product_id"`
	Sku                  string        `json:"sku"`
	Name                 string        `json:"name"`
	Quantity             int           `json:"quantity"`
	UnitPrice            Money         `json:"unit_price"`
	UnitPriceWithTax     Money         `json:"unit_price_with_tax"`
	ExtendedPrice        Money         `json:"extended_price"`
	ExtendedPriceWithTax Money         `json:"extended_price_with_tax"`
	UnitWeight           int           `json:"unit_weight"`
	ExtendedWeight       int           `json:"extended_weight"`
	CustomFields         []CustomField `json:"custom_fields"`
}

type ContactInfo struct {
	FirstName    string        `json:"first_name"`
	LastName     string        `json:"last_name"`
	Company      string        `json:"company"`
	Address1     string        `json:"address1"`
	Address2     string        `json:"address2"`
	City         string        `json:"city"`
	State        string        `json:"state"`
	PostCode     string        `json:"postcode"`
	Country      string        `json:"country"`
	CountryCode  string        `json:"country_code"`
	Email        string        `json:"email"`
	Phone        string        `json:"phone"`
	CustomFields []CustomField `json:"custom_fields"`
}

type Coupon struct {
	Type     int    `json:"type"`
	Code     string `json:"code"`
	Discount Money  `json:"discount"`
}

type OrderService struct {
	DB     *DbLayer
	Schema *JSchema
	Store  Store
}

type IOrderService interface {
	CreateOrder(order Order) (int64, error)
}

func NewOrderService(dbLayer *DbLayer, store Store) (*OrderService, error) {
	schema, err := dbLayer.GetSchema(store.Name)
	if err != nil {
		return nil, err
	}

	service := &OrderService{dbLayer, schema, store}
	return service, nil
}

func (source *OrderService) CreateOrder(lineItems []LineItem) (int64, error) {

	orderCollection, err := source.Schema.GetCollection("orders")
	if err != nil {
		return 0, err
	}

	doc := &JDocument{}
	err = orderCollection.Insert(doc)
	if err != nil {
		return 0, err
	}

	return doc.id, nil
}

func (source *OrderService) GetOrder(orderId int64) (*Order, error) {
	orderCollection, err := source.Schema.GetCollection("orders")
	if err != nil {
		return nil, err
	}

	_, err = orderCollection.FindById(orderId)
	if err != nil {
		return nil, err
	}

	return &Order{}, nil
}

func (source *OrderService) UpdateOrder(order Order) error {
	return nil
}

func (source *OrderService) Getlineitems(orderId int64) (*[]LineItem, error) {
	return nil, nil
}

func (source *OrderService) Updatelineitems(orderId int64, lineItems []LineItem) error {
	return nil
}
