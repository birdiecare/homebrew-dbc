package handler

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/config"
)

// Create SSM Session Output: returns a StreamURL and Token to open a WebSocket connection (SSM)
func createSession(t string, h string, p string, lp string) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Panic("configuration error: " + err.Error())
	}

	args := []string{
		"ssm",
		"start-session",
		"--region",
		cfg.Region,
		"--target",
		t,
		"--document-name",
		"AWS-StartPortForwardingSessionToRemoteHost",
		"--parameters",
		fmt.Sprintf("host=%s,portNumber=%s,localPortNumber=%s", h, p, lp),
	}

	command := exec.Command("aws", args...)

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	log.Printf("Opening connection on localhost:%s", lp)

	err = command.Run()

	if err != nil {
		log.Panic("failed to open port-forwarding connection: " + err.Error())
	}
}
