package handler

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func Proxy(proxyPort string) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Panic("configuration error: " + err.Error())
	}

	// Fetch content of private key from AWS Secrets Manager

	svc := secretsmanager.NewFromConfig(cfg)

	secretName := "bastion-ssh-key"

	res, err := svc.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	})

	if err != nil {
		log.Panic("failed to get secret: " + err.Error())
	}

	// Write private key to temp file

	file, err := os.CreateTemp("", "")

	if err != nil {
		log.Panic("failed to create temp file: " + err.Error())
	}

	log.Println("Downloaded private key file to:", file.Name())

	// Remove the file when the program exits
	defer os.Remove(file.Name())

	file.WriteString(*res.SecretString)

	// Find bastion instance

	bastion := getBastion()

	args := []string{
		"-D",
		proxyPort,
		"-C",
		"-N",
		"ec2-user@" + bastion,
		"-o",
		fmt.Sprintf(`ProxyCommand sh -c "aws ssm start-session --target %v --document-name AWS-StartSSHSession --parameters 'portNumber=22'"`, bastion),
		"-i",
		file.Name(),
	}

	command := exec.Command("ssh", args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	log.Printf("Opening Proxy at localhost:%s. Use ALL_PROXY=socks5://localhost:%s", proxyPort, proxyPort)

	err = command.Run()

	if err != nil {
		log.Panic("failed to open proxy: " + err.Error())
	}
}
