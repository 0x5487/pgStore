package main

type Order struct {
	Total           int64
	ShippingCost    int64
	BillingAddress  ContactInfo `json:"billing_address"`
	ShippingAddress ContactInfo `json:"shipping_address"`
	Status          string
	IPAddress       string
	CurrencyCode    string
}

type LineItem struct {
	ProductId     int64
	VariationSku  string
	Quantity      int
	UnitPrice     int64
	ExtendedPrice int64
}

type ContactInfo struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Company     string `json:"company"`
	Address1    string `json:"address1"`
	Address2    string `json:"address2"`
	City        string `json:"city"`
	State       string `json:"state"`
	PostCode    string `json:"postcode"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
}
