package entity

import "github.com/MurmurationsNetwork/MurmurationsServices/common/jsonapi"

type Node struct {
	ID             string
	ProfileURL     string
	ProfileHash    *string
	Status         string
	LastUpdated    *int64
	FailureReasons *[]jsonapi.Error
	Version        *int32
	CreatedAt      int64
	ProfileStr     string
}

type Nodes []Node
