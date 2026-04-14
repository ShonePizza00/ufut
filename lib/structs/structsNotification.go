package structsUFUT

type InventoryOrderNotification struct {
	Action            string   `json:"action"`
	ItemsAvailability []bool   `json:"itemsAvailability"`
	ItemsIDs          []string `json:"itemsIDs"`
}
