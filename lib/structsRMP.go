package structsUFUT

type OrderRequestRMP struct {
	OrderID int    `json:"orderID"`
	UserID  string `json:"userID"`
	Status  string `json:"status"`
}

type OrdersResponseRMP struct {
	OrderID []int    `json:"ordersID"`
	Status  []string `json:"statuses"`
}

type ItemRequestRMP struct {
	UserID   string `json:"userID"`
	ItemID   string `json:"itemID"`
	Quantity int    `json:"quantity"`
}

type ShoppingCartRMP struct {
	UserID     string   `json:"userID"`
	ItemsID    []string `json:"itemsID"`
	Quantities []int    `json:"quantities"`
}
