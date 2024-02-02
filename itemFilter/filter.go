package itemFilter

import (
	"gw2-tradingpost-tracker/internal/gw2Client/gw2models"
	"sync"
)

const (
	Exotic   = "Exotic"
	Rare     = "Rare"
	MinLevel = 76
	MaxLevel = 80
	Weapon   = "Weapon"
	Armor    = "Armor"
)

type syncMap struct {
	m    map[int]string
	lock sync.Locker
}

func (s syncMap) Write(id int, name string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.m[id] = name
}

func (s syncMap) Remove(id int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.m, id)
}

func (s syncMap) Get(id int) string {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.m[id]
}

type Filter struct {
	exotics syncMap
	rares   syncMap
}

func (f *Filter) FilterRaresAndExoticsSalvagables(item []gw2models.BaseItemData) {
	for _, v := range item {
		if !v.ItemLevel(MinLevel, MaxLevel) {
			continue
		}
		if !v.ItemType(Weapon) && !v.ItemType(Armor) {
			continue
		}

		if v.ItemRarity(Rare) {
			f.rares.Write(v.GetId(), v.GetName())
		}

		if v.ItemRarity(Exotic) {
			f.exotics.Write(v.GetId(), v.GetName())
			continue
		}
	}
}
