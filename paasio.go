package paasio

import (
	"io"
	"sync"
)

// readCounter implements ReadCounter
type readCounter struct {
	nRead    *int64        // nRead tracks the number of bytes read
	nopsRead *int          // nopsRead tracks the number of times Read is called
	muRead   *sync.RWMutex // muRead is used to protect the critical sections referencing the values above
	io.Reader
}

// writeCounter implements WriteCounter
type writeCounter struct {
	nWrite    *int64        // nWrite tracks the number of bytes written
	nopsWrite *int          // nopsWrite tracks the number of times Write is called
	muWrite   *sync.RWMutex // muWrite is used to protect the critical sections referencing the values above
	io.Writer
}

// readWriteCounter implements ReadWriteCounter
type readWriteCounter struct {
	ReadCounter
	WriteCounter
}

// ReadCount returns the number of bytes read and number of times Read has been called
func (rc readCounter) ReadCount() (n int64, nops int) {
	rc.muRead.RLock()
	defer rc.muRead.RUnlock()
	return *rc.nRead, int(*rc.nopsRead)
}

// Read reads data into p
func (rc readCounter) Read(p []byte) (n int, err error) {
	// including the Read in the critical section is slower, but it ensures a more accurate
	// reporting of bytes read
	rc.muRead.Lock()
	defer rc.muRead.Unlock()
	n, err = rc.Reader.Read(p)
	*rc.nRead += int64(n)
	*rc.nopsRead += 1
	return n, err
}

// WriteCount returns the number of bytes written and number of times Write has been called
func (wc writeCounter) WriteCount() (int64, int) {
	wc.muWrite.RLock()
	defer wc.muWrite.RUnlock()
	return *wc.nWrite, int(*wc.nopsWrite)
}

// Write writes data into p
func (wc writeCounter) Write(p []byte) (n int, err error) {
	// including the Write in the critical section is slower, but it ensures a more accurate
	// reporting of bytes written
	wc.muWrite.Lock()
	defer wc.muWrite.Unlock()
	n, err = wc.Writer.Write(p)
	*wc.nWrite += int64(n)
	*wc.nopsWrite += 1
	return n, err
}

// ReadCount returns the number of bytes read and number of times Read has been called
func (rwc readWriteCounter) ReadCount() (int64, int) {
	return rwc.ReadCounter.ReadCount()
}

// WriteCount returns the number of bytes written and number of times Write has been called
func (rwc readWriteCounter) WriteCount() (int64, int) {
	return rwc.WriteCounter.WriteCount()
}

// NewReadCounter creates a new ReadCounter
func NewReadCounter(reader io.Reader) ReadCounter {
	return readCounter{
		nRead:    new(int64),
		nopsRead: new(int),
		muRead:   new(sync.RWMutex),
		Reader:   reader,
	}
}

// NewWriteCounter creates a new WriteCounter
func NewWriteCounter(writer io.Writer) WriteCounter {
	return writeCounter{
		nWrite:    new(int64),
		nopsWrite: new(int),
		muWrite:   new(sync.RWMutex),
		Writer:    writer,
	}
}

// NewReadWriteCounter creates a new ReadWriteCounter
func NewReadWriteCounter(readWriter io.ReadWriter) ReadWriteCounter {
	return readWriteCounter{
		NewReadCounter(readWriter),
		NewWriteCounter(readWriter),
	}
}
