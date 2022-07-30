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
	firstItem *ListItem
	lastItem  *ListItem
	length    int
}

func NewList() List {
	return new(list)
}

func (l *list) MoveToFront(listItem *ListItem) {
	l.remove(listItem)
	l.pushFront(listItem)
}

func (l *list) Remove(listItem *ListItem) {
	l.remove(listItem)
}

func (l *list) remove(listItem *ListItem) {
	if listItem.Next != nil {
		listItem.Next.Prev = listItem.Prev
	}

	if listItem.Prev != nil {
		listItem.Prev.Next = listItem.Next
	}

	if listItem == l.firstItem {
		l.firstItem = listItem.Next
	} else if listItem == l.lastItem {
		l.lastItem = l.lastItem.Prev
	}

	listItem.Next = nil
	listItem.Prev = nil

	l.length--
}

func (l *list) PushBack(value interface{}) *ListItem {
	item := &ListItem{
		Value: value,
	}

	if l.firstItem == nil {
		l.firstItem = item
	}

	if l.lastItem != nil {
		l.lastItem.Next = item
		item.Prev = l.lastItem
	}

	l.lastItem = item
	l.length++

	return l.lastItem
}

func (l *list) PushFront(value interface{}) *ListItem {
	item := &ListItem{
		Value: value,
	}

	l.pushFront(item)
	return item
}

func (l *list) pushFront(item *ListItem) {
	if l.lastItem == nil {
		l.lastItem = item
	}

	if l.firstItem != nil {
		l.firstItem.Prev = item
		item.Next = l.firstItem
	}

	l.firstItem = item
	l.length++
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.firstItem
}

func (l *list) Back() *ListItem {
	return l.lastItem
}
