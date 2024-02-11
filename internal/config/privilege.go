package config

import (
	"github.com/nekogravitycat/DiscordAI/internal/jsondata"
)

type privilegeConfig struct {
	Models []string `json:"models"`
}

func newPrivilegeConfig() privilegeConfig {
	pc := privilegeConfig{
		Models: []string{},
	}
	return pc
}

var privilegeData = map[string]privilegeConfig{"0": newPrivilegeConfig()}

func GetPrivilegeConfig(level string) (c privilegeConfig, ok bool) {
	c, ok = privilegeData[level]
	return c, ok
}

const PRIVILEGEFILE = "./configs/privilege.json"

func loadPrivilegeConfig() {
	jsondata.Check(PRIVILEGEFILE, privilegeData)
	jsondata.Load(PRIVILEGEFILE, &privilegeData)
}
