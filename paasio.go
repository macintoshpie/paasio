package paasio

import (
	"io"
	"sync"
)

type readCounter struct {
	nRead    *int64
	nopsRead *int
	muRead   *sync.RWMutex
	io.Reader
}

type writeCounter struct {
	nWrite    *int64
	nopsWrite *int
	muWrite   *sync.RWMutex
	io.Writer
}

type readWriteCounter struct {
	ReadCounter
	WriteCounter
}

func (rc readCounter) ReadCount() (int64, int) {
	rc.muRead.RLock()
	defer rc.muRead.RUnlock()
	return *rc.nRead, int(*rc.nopsRead)
}

func (rc readCounter) Read(p []byte) (n int, err error) {
	n, err = rc.Reader.Read(p)
	rc.muRead.Lock()
	defer rc.muRead.Unlock()
	*rc.nRead += int64(n)
	*rc.nopsRead += 1
	return n, err
}

func (wc writeCounter) WriteCount() (int64, int) {
	wc.muWrite.RLock()
	defer wc.muWrite.RUnlock()
	return *wc.nWrite, int(*wc.nopsWrite)
}

func (wc writeCounter) Write(p []byte) (n int, err error) {
	n, err = wc.Writer.Write(p)
	wc.muWrite.Lock()
	defer wc.muWrite.Unlock()
	*wc.nWrite += int64(n)
	*wc.nopsWrite += 1
	return n, err
}

func (rwc readWriteCounter) ReadCount() (int64, int) {
	return rwc.ReadCounter.ReadCount()
}

func (rwc readWriteCounter) WriteCount() (int64, int) {
	return rwc.WriteCounter.WriteCount()
}

func NewReadCounter(reader io.Reader) ReadCounter {
	n := new(int64)
	nops := new(int)
	return readCounter{nRead: n, nopsRead: nops, muRead: new(sync.RWMutex), Reader: reader}
}

func NewWriteCounter(writer io.Writer) WriteCounter {
	n := new(int64)
	nops := new(int)
	return writeCounter{nWrite: n, nopsWrite: nops, muWrite: new(sync.RWMutex), Writer: writer}
}

func NewReadWriteCounter(readWriter io.ReadWriter) ReadWriteCounter {
	return readWriteCounter{NewReadCounter(readWriter), NewWriteCounter(readWriter)}
}
