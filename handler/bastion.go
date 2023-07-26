package handler

import (
	"context"
	"fmt"

	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// Return Bastion Instance by finding by `tag:service = bastion`
func getBastion() string {

	var bastion string

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Panic("configuration error: " + err.Error())
	}

	svc := ec2.NewFromConfig(cfg)

	params := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:service"),
				Values: []string{"bastion"},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{"running"},
			},
		},
	}

	instance_list, err := svc.DescribeInstances(context.TODO(), params)
	if err != nil {
		log.Fatal(err)
	}

	bastion = *instance_list.Reservations[0].Instances[0].InstanceId

	log.Println(fmt.Sprintf("Using bastion: %s", bastion))
	return bastion
}
