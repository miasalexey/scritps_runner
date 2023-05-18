package mikrotik

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
	"log"
	"net/http"
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
	router.GET("/enable_winbox", h.enableWinbox)
}

func (h *handler) enableWinbox(context *gin.Context) {
	var mikrotik EnableWinboxDTO
	if err := context.ShouldBindJSON(&mikrotik); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "can't parse structure",
		})
	}
	host := mikrotik.Ip
	port := 22
	config := &ssh.ClientConfig{
		User: mikrotik.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(mikrotik.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         3 * time.Second,
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		res := fmt.Sprintf("can't connect to %s", host)
		context.JSON(http.StatusBadRequest, gin.H{

			"error": res,
		})
		return
	}

	session, err := client.NewSession()
	if err != nil {
		res := fmt.Sprintf("can't connect to %s", host)
		context.JSON(http.StatusBadRequest, gin.H{

			"error": res,
		})
		return
	}
	defer func(session *ssh.Session) {
		err := session.Close()
		if err != nil {
			log.Println("can't close session on host " + host)
		}
	}(session)

	var buildCommand strings.Builder
	buildCommand.WriteString(`/ip firewall filter set dst-port=22,8291 comment="Allow SSH,Winbox FROM trusted hosts" [find comment~"SSH"]`)
	buildCommand.WriteString("\n")
	buildCommand.WriteString(`/ip ser enable [find name=winbox]`)
	buildCommand.WriteString("\n")
	command := buildCommand.String()
	_, err = session.CombinedOutput(command)
	if err != nil {
		log.Printf("Failed to execute command '%s': %v", command, err)
		context.JSON(http.StatusBadRequest, gin.H{

			"error": "res",
		})
		return
	}

	res := fmt.Sprintf("winbox enabled on host %s enabled successfully", host)
	context.JSON(http.StatusOK, gin.H{
		"response": res,
	})
}
