package keepass_parser

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tobischo/gokeepasslib"
	"golang.org/x/crypto/ssh"
	"net/http"
	"os"
	"scritps_runner/internal/config"
	"scritps_runner/internal/handlers"
	"strings"
	"time"
)

type handler struct {
}

func NewHandler() handlers.Handlers {
	return &handler{}
}

func (h *handler) Register(router *gin.Engine) {
	router.GET("/get_pass", h.parsePasswordByTags)
}

func (h *handler) parsePasswordByTags(context *gin.Context) {
	var keepassDTO KeePassParserDTO
	if err := context.ShouldBindJSON(&keepassDTO); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "can't parse structure",
		})
		return
	}
	cfg := config.GetConfig()
	keepassFile, _ := os.Open(cfg.KeepassFile)

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(cfg.KeepassPassword)
	_ = gokeepasslib.NewDecoder(keepassFile).Decode(db)

	err := db.UnlockProtectedEntries()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "can't unlock kp file",
		})
		return
	}
	allGroupsInDb := db.Content.Root.Groups
	passwords := findPasswordsByTags(keepassDTO.Tags, allGroupsInDb)
	res := checkPasswordsUseSSH(keepassDTO.Ip, keepassDTO.Login, passwords)
	context.JSON(http.StatusOK, res)
	return

}

func findPasswordsByTags(tags []string, groups []gokeepasslib.Group) map[string]string {
	res := make(map[string]string)

	for _, group := range groups {

		findName := strings.ToLower(group.Name)
		for _, tag := range tags {
			if strings.Contains(findName, strings.ToLower(tag)) {
				temp := getPasswordsFromGroup(&group)
				for k, v := range temp {
					res[k] = v
				}
				if len(group.Groups) != 0 {
					for _, v := range group.Groups {
						temp = getPasswordsFromGroup(&v)
						for k, v := range temp {
							res[k] = v
						}
					}
				}
			}
		}
		if len(group.Groups) != 0 {
			temp := findPasswordsByTags(tags, group.Groups)
			for k, v := range temp {
				res[k] = v
			}
		}
	}
	return res
}

func getPasswordsFromGroup(group *gokeepasslib.Group) map[string]string {
	res := make(map[string]string)
	for _, item := range group.Entries {
		res[item.GetTitle()] = item.GetPassword()
	}
	return res
}

func checkPasswordsUseSSH(ip string, login string, passwords map[string]string) ResultKeePassParserDTO {
	res := ResultKeePassParserDTO{
		Title:    "not found",
		Password: "not found",
	}

	for title, password := range passwords {
		sshConfig := &ssh.ClientConfig{
			User: login,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         3 * time.Second,
		}

		client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip, 22), sshConfig)
		if err != nil {
			continue
		} else {
			res = ResultKeePassParserDTO{
				Title:    title,
				Password: password,
			}
			err := client.Close()
			if err != nil {
				return res
			}
			break
		}
	}
	return res
}
