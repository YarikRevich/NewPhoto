package caching

import (
	"encoding/json"
	"github.com/YarikRevich/NewPhoto/log"
)

type DataConfigurator struct {
	Model interface{}
}

func (dc DataConfigurator) Configure() []byte {
	by, err := json.Marshal(dc.Model)
	if err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("CacheConfigure", err)
	}
	return by
}
