package store

// List is a double-ended queue backed by a ring buffer.
//
// It gives O(1) amortized push/pop at both ends, unlike a Go slice where
// prepending is O(N). This was the #1 performance blocker: Engine.LPush used
// `append([]string{v}, list...)`, making redis-benchmark's single-growing-list
// LPUSH an O(N²) allocation bomb (98.7% of all allocs in profiling, ~4.2K
// ops/sec collapsing toward zero). With this deque, LPUSH is O(1) per element.
type List struct {
	buf   []string
	head  int // index of the first (logical head) element
	count int // number of elements
}

// NewList returns an empty list.
func NewList() *List {
	return &List{}
}

// Len returns the number of elements.
func (l *List) Len() int { return l.count }

// at returns the physical buffer index for logical index i (0 = head).
func (l *List) at(i int) int {
	if len(l.buf) == 0 {
		return 0
	}
	return (l.head + i) % len(l.buf)
}

// PushFront prepends an element (O(1) amortized).
func (l *List) PushFront(v string) {
	if l.count == len(l.buf) {
		l.grow()
	}
	l.head = (l.head - 1 + len(l.buf)) % len(l.buf)
	l.buf[l.head] = v
	l.count++
}

// PushBack appends an element (O(1) amortized).
func (l *List) PushBack(v string) {
	if l.count == len(l.buf) {
		l.grow()
	}
	l.buf[l.at(l.count)] = v
	l.count++
}

// PopFront removes and returns the head element.
func (l *List) PopFront() (string, bool) {
	if l.count == 0 {
		return "", false
	}
	v := l.buf[l.head]
	l.head = (l.head + 1) % len(l.buf)
	l.count--
	return v, true
}

// PopBack removes and returns the tail element.
func (l *List) PopBack() (string, bool) {
	if l.count == 0 {
		return "", false
	}
	idx := l.at(l.count - 1)
	v := l.buf[idx]
	l.count--
	return v, true
}

// Get returns the element at logical index i (0 = head). Negative index
// counts from the tail, matching Redis semantics (-1 = last element).
func (l *List) Get(i int) (string, bool) {
	if i < 0 {
		i = l.count + i
	}
	if i < 0 || i >= l.count {
		return "", false
	}
	return l.buf[l.at(i)], true
}

// Copy returns all elements in head-to-tail order as a new slice.
func (l *List) Copy() []string {
	out := make([]string, l.count)
	for i := 0; i < l.count; i++ {
		out[i] = l.buf[l.at(i)]
	}
	return out
}

// grow doubles the buffer, re-linearizing elements so head becomes 0.
func (l *List) grow() {
	newCap := 8
	if len(l.buf) > 0 {
		newCap = len(l.buf) * 2
	}
	newBuf := make([]string, newCap)
	for i := 0; i < l.count; i++ {
		newBuf[i] = l.buf[l.at(i)]
	}
	l.buf = newBuf
	l.head = 0
}
