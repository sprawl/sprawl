package api

type order struct {
	orderID      uint
	timestamp    uint
	asset        uint
	counterAsset uint
	amount       uint
	price        string
}

