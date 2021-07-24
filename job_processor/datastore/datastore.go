package datastore

type DataStore interface {
	UpdateStatus(objectID, status string) error
}
