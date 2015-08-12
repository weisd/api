package polling

import (
	"testing"
)

func TestPolling(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	p := NewPolling(len(arr))

	for i := 0; i < 100; i++ {
		t.Log(arr[p.Index()])
	}
}
