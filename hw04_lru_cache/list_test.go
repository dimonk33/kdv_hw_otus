package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("last remove", func(t *testing.T) {
		l := NewList()

		l.PushBack(1) // [1]
		l.PushBack(2) // [1, 2]
		require.Equal(t, 2, l.Len())

		for i := l.Front(); i != nil; {
			j := i
			i = i.Next
			l.Remove(j)
		}

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("move front", func(t *testing.T) {
		l := NewList()

		item1 := l.PushBack(1) // [1]
		item2 := l.PushBack(2) // [1, 2]
		item3 := l.PushBack(3) // [1, 2, 3]

		l.MoveToFront(l.Back()) // [3, 1, 2]
		l.MoveToFront(l.Back()) // [2, 3, 1]

		require.Equal(t, l.Front(), item2)
		require.Equal(t, l.Back(), item1)

		require.Nil(t, l.Front().Prev)
		require.Equal(t, l.Front().Next, item3)

		require.Nil(t, l.Back().Next)
		require.Equal(t, l.Back().Prev, item3)

		require.Equal(t, item3.Prev, l.Front())
		require.Equal(t, item3.Next, l.Back())
	})
}
