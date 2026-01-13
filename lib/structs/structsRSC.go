package structsUFUT

type ItemsRequestRSC struct {
	Category   string `json:"category"`
	Price      int    `json:"price"`
	StartIndex int    `json:"startIndex"`
	Count      int    `json:"count"`
	OrderBy    string `json:"orderBy"`
}
type ItemsResponseRSC struct {
	ItemsIDs []string `json:"itemsID"`
}

type ItemDataRSC struct {
	ItemID      string `json:"itemID"`
	SellerID    string `json:"sellerID"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Category    string `json:"category"`
	Status      string `json:"status"`
	Quantity    int    `json:"quantity"`
}
