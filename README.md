# PaaS IO

![workflow badge](https://github.com/macintoshpie/paasio/actions/workflows/ci.yaml/badge.svg)

`passio` provides an io wrapper which reports:
- The total number of bytes read/written.
- The total number of read/write operations.

## Development

```bash
git clone git@github.com:macintoshpie/paasio.git

go mod download
```

### Running tests

```bash
go test -v
```

## Implementation Notes
I used embedded types to implement the required interfaces. See the docs in paasio.go for more details.

### Alternatives considered
We might not need a full read/write mutex for this use case, a single mutex might be sufficient depending on the frequency of reading.

Also, the current implementation includes the Read and Write calls inside of the critical zone (ie locked section). This ensures a more accurate "reporting" of bytes written/read. However locking the mutex before these calls will result in a slower implementation, so I'd suggest excluding them if we're ok with possibly being slightly less accurate.

In addition, if performance was of *much* greater importance than accuracy, I'd recommend we use the atomic package's AtomicAdd functions instead of mutexes. For example, when calling read, we could have done something like this:
```
func (rc readCounter) Read(p []byte) (n int, err error) {
	n, err = rc.Reader.Read(p)
	atomic.AddInt64(rc.nRead, int64(n))
	atomic.AddInt32(rc.nopsRead, 1)
	return n, err
}
```
This will be much faster than using locks. But callers of `*Count` beware, you might be reading a partially complete report if a goroutine is suspended right after the first increment.
