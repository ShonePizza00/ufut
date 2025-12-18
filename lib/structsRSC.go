package structsUFUT

type ItemsRequestRSC struct {
	Category   string  `json:"category"`
	Price      float32 `json:"price"`
	StartIndex int     `json:"startIndex"`
	Count      int     `json:"count"`
	OrderBy    string  `json:"orderBy"`
}
type ItemsResponseRSC struct {
	ItemsIDs []string `json:"itemsIDs"`
}

type ItemDataRSC struct {
	ItemID      int     `json:"itemID"`
	SellerID    string  `json:"sellerID"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Category    string  `json:"category"`
}
