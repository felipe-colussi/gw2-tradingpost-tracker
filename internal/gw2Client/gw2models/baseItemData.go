package gw2models

type BaseItemData struct {
	Id     int    `json:"id"`
	Level  int    `json:"level"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Rarity string `json:"rarity"`
}

func (b BaseItemData) GetId() int {
	return b.Id
}

func (b BaseItemData) GetName() string {
	return b.Name
}

func (b BaseItemData) ItemRarity(s string) bool {
	return b.Rarity == s
}

func (b BaseItemData) ItemType(s string) bool {
	return b.Type == s
}

func (b BaseItemData) ItemLevel(min, max int) bool {
	return b.Level >= min && b.Level <= max
}
