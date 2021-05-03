package model

// ThreadName a thread name.
type ThreadName string

// String ...
func (threadName ThreadName) String() string {
	return string(threadName)
}
