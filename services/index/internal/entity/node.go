package entity

type Node struct {
	ID             string
	ProfileURL     string
	ProfileHash    *string
	Status         string
	LastValidated  *int64
	FailureReasons *[]string
	Version        *int32
	CreatedAt      int64
	ProfileStr     string
}

type Nodes []Node
