package caching

import (
	"encoding/json"
	"github.com/YarikRevich/NewPhoto/log"
)

type Definer struct {
	Model interface{}
	Data  string
}

func (d Definer) Define() interface{} {
	if err := json.Unmarshal([]byte(d.Data), &d.Model); err != nil {
		log.Logger.CFatalln("CacheDefiner", err)
	}

	return d.Model
}
