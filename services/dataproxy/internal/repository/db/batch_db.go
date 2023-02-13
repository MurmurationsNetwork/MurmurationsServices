package db

type BatchRepository interface {
}

type batchRepository struct{}

func NewBatchRepository() BatchRepository {
	return &batchRepository{}
}
