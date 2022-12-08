package storage

type SnapshotWriterInterface interface {
	Load()
	Receive()
	Close()
}
