package main

import (
	"time"
)

type Order struct {
	Id               int64         `json:"id"`
	ShopperId        int           `json:"shopper_id"`
	LineItems        []LineItem    `json:"line_items"`
	ItemsCount       int           `json:"items_count"`
	ShippingOptionId int           `json:"shipping_optionid"`
	SubTotal         Money         `json:"subtotal"`
	SubTotalWithTax  Money         `json:"subtotal_with_tax"`
	ShippingFee      Money         `json:"shipping_fee"`
	Coupons          []Coupon      `json:"coupons"`
	Discount         Money         `json:"discount"`
	Total            Money         `json:"total"`
	TotalWithTax     Money         `json:"total_with_tax"`
	PaymentMethod    string        `json:"payment_method"`
	PaymentState     int           `json:"payment_state"`
	BillingAddress   ContactInfo   `json:"billing_address"`
	ShippingAddress  ContactInfo   `json:"shipping_address"`
	StateId          int           `json:"state_id"`
	State            string        `json:"status"`
	IPAddress        string        `json:"ip_address"`
	CurrencyCode     string        `json:"currency_code"`
	IsTestOrder      bool          `json:"is_test_order"`
	PaidDateUtc      time.Time     `json:"paid_date_utc"`
	CustomFields     []CustomField `json:"custom_fields"`
	Memo             string        `json:"memo"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
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
}

type OrderService struct {
}
