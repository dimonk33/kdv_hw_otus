package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	count     int
	frontItem *ListItem
	backItem  *ListItem
}

func (l *list) Len() int {
	return l.count
}

func (l *list) Front() *ListItem {
	return l.frontItem
}

func (l *list) Back() *ListItem {
	return l.backItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := ListItem{}
	item.Value = v
	if l.frontItem == nil {
		item.Next = nil
		l.frontItem = &item
		l.backItem = &item
	} else {
		item.Next = l.frontItem
		l.frontItem.Prev = &item
		l.frontItem = &item
	}
	l.count++

	return &item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := ListItem{}
	item.Value = v
	item.Next = nil
	if l.backItem == nil {
		item.Prev = nil
		l.frontItem = &item
		l.backItem = &item
	} else {
		item.Prev = l.backItem
		l.backItem.Next = &item
		l.backItem = &item
	}
	l.count++

	return &item
}

func (l *list) Remove(i *ListItem) {
	if l.count == 1 {
		l.frontItem = nil
		l.backItem = nil
	} else {
		switch i {
		case l.frontItem:
			i.Next.Prev = nil
			l.frontItem = i.Next
		case l.backItem:
			i.Prev.Next = nil
			l.backItem = i.Prev
		default:
			i.Prev.Next = i.Next
			i.Next.Prev = i.Prev
		}
	}
	l.count--
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev == nil {
		return
	}
	i.Prev.Next = i.Next
	if i.Next == nil {
		l.backItem = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	l.frontItem.Prev = i
	i.Prev = nil
	i.Next = l.frontItem
	l.frontItem = i
}

func NewList() List {
	return new(list)
}
