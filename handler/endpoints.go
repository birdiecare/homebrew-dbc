package handler

import (
	"context"
	"fmt"

	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/ktr0731/go-fuzzyfinder"
)

type db struct {
	DBId      string
	Type      string
	Endpoints []string
	IAM       bool
}

// Return DB Type list of Database ClusterId's, InstanceId's and Endpoints.
func getEndpoints() []db {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Panic("configuration error: " + err.Error())
	}

	svc := rds.NewFromConfig(cfg)

	log.Println("Grabbing RDS Endpoints...")

	i_params := &rds.DescribeDBInstancesInput{}
	instance_list, err := svc.DescribeDBInstances(context.TODO(), i_params)
	if err != nil {
		log.Fatal(err)
	}

	c_params := &rds.DescribeDBClustersInput{}
	cluster_list, err := svc.DescribeDBClusters(context.TODO(), c_params)
	if err != nil {
		log.Fatal(err)
	}

	var endpoints []db

	// Get clusters
	for _, i := range cluster_list.DBClusters {
		endpoints = append(endpoints, db{
			DBId:      *i.DBClusterIdentifier,
			Type:      "Cluster",
			Endpoints: []string{*i.ReaderEndpoint, *i.Endpoint},
			IAM:       *i.IAMDatabaseAuthenticationEnabled,
		})
	}

	// Get Instances not in Clusters
	for _, i := range instance_list.DBInstances {
		if i.DBClusterIdentifier == nil {
			endpoints = append(endpoints, db{
				DBId:      *i.DBInstanceIdentifier,
				Type:      "Instance",
				Endpoints: []string{*i.Endpoint.Address},
				IAM:       *i.IAMDatabaseAuthenticationEnabled,
			})
		}
	}

	return endpoints
}

func FuzzEndpoints() (string, string) {

	endpoints := getEndpoints()

	idx, err := fuzzyfinder.Find(
		endpoints,
		func(i int) string {
			return endpoints[i].DBId
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("DB: %s\nType: %s\nEndpoints: %+q\nIAM Auth: %t",
				endpoints[i].DBId,
				endpoints[i].Type,
				endpoints[i].Endpoints,
				endpoints[i].IAM)
		}))
	if err != nil {
		log.Fatal(err)
	}

	if len(endpoints[idx].Endpoints) > 1 {
		return endpoints[idx].DBId, fuzzCluster(endpoints[idx].Endpoints)
	}

	return endpoints[idx].DBId, endpoints[idx].Endpoints[0]
}

func fuzzCluster(e []string) string {
	idx, err := fuzzyfinder.Find(e, func(i int) string {
		return e[i]
	},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("Role: %s", getRole(e, i))
		}))
	if err != nil {
		log.Fatal(err)
	}

	return e[idx]
}

func getRole(e []string, i int) string {
	if e[i] == e[0] {
		return "Reader"
	} else {
		return "Writer"
	}
}
