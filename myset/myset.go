// package main
package myset

import (
//	"fmt"
	"sort"
	"sync"
)

type Set struct {
	m map[string]bool
	sync.RWMutex
}

func New() *Set {
	return &Set{
		m: map[string]bool{},
	}
}

func (s *Set) Add(item string) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = true
}

func (s *Set) Remove(item string) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, item)
}

func (s *Set) Has(item string) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

func (s *Set) Len() int {
	return len(s.List())
}

func (s *Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[string]bool{}
}

func (s *Set) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}

func (s *Set) Equal(c *Set) bool {
	s.RLock()
	defer s.RUnlock()
    if s.Len() != c.Len() {
        return false
    }
	for item := range s.m {
        if ! c.Has(item) {
            return false
        }
	}
	return true
}

func (s *Set) List() []string {
	s.RLock()
	defer s.RUnlock()
	list := []string{}
	for item := range s.m {
		list = append(list, item)
	}
	return list
}

func (s *Set) SortList() []string {
	s.RLock()
	defer s.RUnlock()
	list := []string{}
	for item := range s.m {
		list = append(list, item)
	}
	sort.Strings(list)
	return list
}

/*
func main() {
	//初始化
	s := New()

	s.Add("1")
	s.Add("1")
	s.Add("0")
	s.Add("2")

	if s.Has("2") {
		fmt.Println("2 does exist")
	}

	s.Remove("2")
	fmt.Println("无序的切片", s.List())
	fmt.Println("有序的切片", s.SortList())

    c1 := New()
    c2 := New()

	c1.Add("0")
	c1.Add("1")
	c1.Add("2")

	c2.Add("0")
	c2.Add("1")

    fmt.Printf("Equal1 \ns: %v\nc: %v\nresult: %v\n", s.List(), c1.List(), s.Equal(c1))
    fmt.Printf("Equal2 \ns: %v\nc: %v\nresult: %v\n", s.List(), c2.List(), s.Equal(c2))

	s.Clear()
	if s.IsEmpty() {
		fmt.Println("0 item")
	}
}
*/
