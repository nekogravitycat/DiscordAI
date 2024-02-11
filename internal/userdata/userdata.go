package userdata

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/nekogravitycat/DiscordAI/internal/config"
	"github.com/nekogravitycat/DiscordAI/internal/jsondata"
	openai "github.com/sashabaranov/go-openai"
)

type UserInfo struct {
	Model          string  `json:"model"`
	Credit         float32 `json:"credit"`
	PrivilegeLevel int     `json:"privilege-level"`
}

func NewUserInfo() UserInfo {
	u := UserInfo{
		Model:          openai.GPT3Dot5Turbo,
		Credit:         config.InitCredits,
		PrivilegeLevel: config.InitPrivilege,
	}
	return u
}

func (u UserInfo) HasPrivilege(model string) bool {
	c, ok := config.GetPrivilegeConfig(u.PrivilegeLevel)
	if !ok {
		fmt.Println("Unrecognized privilege level: " + strconv.Itoa(u.PrivilegeLevel))
		return false
	}

	return slices.Contains(c.Models, model)
}

var users = map[string]UserInfo{"0": NewUserInfo()}

func GetUser(discordID string) (user UserInfo, ok bool) {
	user, ok = users[discordID]
	return user, ok
}

func SetUser(discordID string, user UserInfo) UserInfo {
	users[discordID] = user
	return users[discordID]
}

const USERFILE = "./data/users.json"

func LoadUserData() {
	jsondata.Check(USERFILE, users)
	jsondata.Load(USERFILE, &users)
}

func SaveUserData() {
	jsondata.Save(USERFILE, users)
}
