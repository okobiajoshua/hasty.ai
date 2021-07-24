package datastore

type DataStore interface {
	Save(job Job) error
	GetByObjectID(objectID string) (*Job, error)
	GetByJobID(jobID string) (*Job, error)
}
