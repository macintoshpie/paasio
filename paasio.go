package paasio

import "io"

type readCounter struct {
	io.Reader
}

type writeCounter struct {
	io.Writer
}

type readWriteCounter struct {
	io.ReadWriter
}

func (rc readCounter) ReadCount() (int64, int) {
	return 0, 0
}

func (wc writeCounter) WriteCount() (int64, int) {
	return 0, 0
}

func (rwc readWriteCounter) ReadCount() (int64, int) {
	return 0, 0
}

func (rwc readWriteCounter) WriteCount() (int64, int) {
	return 0, 0
}

func NewReadCounter(reader io.Reader) ReadCounter {
	return readCounter{reader}
}

func NewWriteCounter(writer io.Writer) WriteCounter {
	return writeCounter{writer}
}

func NewReadWriteCounter(readWriter io.ReadWriter) ReadWriteCounter {
	return readWriteCounter{readWriter}
}
