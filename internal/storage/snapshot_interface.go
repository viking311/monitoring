package storage

type SnapshotWriterInterface interface {
	Load()
	Close()
	Receive()
}
