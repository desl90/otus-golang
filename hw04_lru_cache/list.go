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
	List
	len   int
	front *ListItem
	back  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(value interface{}) *ListItem {
	return l.pushToFront(&ListItem{
		Value: value,
		Next:  l.front,
	})
}

func (l *list) PushBack(value interface{}) *ListItem {
	return l.pushToBack(&ListItem{
		Value: value,
		Prev:  l.back,
	})
}

func (l *list) Remove(item *ListItem) {
	if l.len < 1 {
		return
	}

	next := item.Next
	prev := item.Prev
	isDone := false

	if next != nil {
		next.Prev = prev
		isDone = true
	}

	if prev != nil {
		prev.Next = next
		isDone = true
	}

	if l.front == item {
		l.front = item.Next
		isDone = true
	}

	if l.back == item {
		l.back = item.Prev
		isDone = true
	}

	if isDone {
		item.Next = nil
		item.Prev = nil
		l.len--
	}
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	i.Next = l.front
	l.pushToFront(i)
}

func (l *list) pushToFront(item *ListItem) *ListItem {
	l.len++

	if l.front != nil {
		l.front.Prev = item
		l.front = item

		return item
	}

	l.front = item
	l.back = item

	return item
}

func (l *list) pushToBack(item *ListItem) *ListItem {
	l.len++

	if l.back != nil {
		l.back.Next = item
		l.back = item

		return item
	}

	l.front = item
	l.back = item

	return item
}
