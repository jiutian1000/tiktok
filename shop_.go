package tiktok

import (
	"encoding/json"
)

type Shop struct {
	ShopID     string      `json:"shop_id"`
	ShopName   string      `json:"shop_name"`
	Region     string      `json:"region"`
	ShopCode   string      `json:"shop_code"`
	ShopCipher string      `json:"shop_cipher"`
	Type       json.Number `json:"type"`
}

type ShopList struct {
	Shops []Shop `json:"shop_list"`
}
