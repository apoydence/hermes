package datastructures

import (
	"sync/atomic"
	"unsafe"
)

type LinkedList struct {
	root unsafe.Pointer //*node
}

type node struct {
	next  unsafe.Pointer //*node
	value unsafe.Pointer
}

var NewLinkedList = func() *LinkedList {
	return new(LinkedList)
}

func (l *LinkedList) Traverse(callback func(unsafe.Pointer)) {
	root := l.loadRoot()
	if root == nil {
		return
	}

	last := l.findLast(root, callback)
	if last != nil {
		callback(l.extractValue(last))
	}
}

func (l *LinkedList) Append(value unsafe.Pointer) {
	newNode := unsafe.Pointer(&node{
		value: value,
	})

	root := l.loadRoot()
	if root == nil {
		atomic.StorePointer(&l.root, newNode)
		return
	}

	last := l.findLast(root, nil)
	l.storeNext(last, newNode)
}

//func (l *LinkedList) Remove(value unsafe.Pointer) {
//	if l.root == nil {
//		return
//	}
//
//	if (*node)(l.root).value == value {
//		l.root = unsafe.Pointer((*node)(l.root).next)
//		return
//	}
//
//	parent := l.findParent(value, l.root)
//	if parent == nil {
//		return
//	}
//
//	(*node)(parent).next = (*node)((*node)(parent).next).next
//}

func (l *LinkedList) loadRoot() unsafe.Pointer {
	return atomic.LoadPointer(&l.root)
}

func (l *LinkedList) extractValue(value unsafe.Pointer) unsafe.Pointer {
	return (*node)(value).value
}

func (l *LinkedList) loadNext(value unsafe.Pointer) unsafe.Pointer {
	return atomic.LoadPointer(&(*node)(value).next)
}

func (l *LinkedList) storeNext(value, next unsafe.Pointer) {
	atomic.StorePointer(&(*node)(value).next, next)
}

func (l *LinkedList) findLast(current unsafe.Pointer, callback func(unsafe.Pointer)) unsafe.Pointer {
	currentNode := (*node)(current)
	next := l.loadNext(current)
	if next == nil {
		return current
	}

	if callback != nil {
		callback(currentNode.value)
	}

	return l.findLast(next, callback)
}

// func (l *LinkedList) findParent(value unsafe.Pointer, current unsafe.Pointer) unsafe.Pointer {
// 	currentNode := (*node)(current)
// 	if atomic.LoadPointer(&current) == nil || atomic.LoadPointer(&currentNode.next) == nil {
// 		return nil
// 	}

// 	if (*node)(currentNode.next).value == value {
// 		return current
// 	}

// 	return l.findParent(value, currentNode.next)
// }
