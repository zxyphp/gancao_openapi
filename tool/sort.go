package tool

import (
	"sort"
)

// MapItem 用于存储map中的键值对
type MapItem struct {
	Key   string
	Value interface{}
}

// ByKey 根据键对MapItem切片进行排序
type ByKey []MapItem

func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

// Sort 对map按照给定的键顺序进行排序，并返回一个MapItem切片
func Sort(data map[string]interface{}, order []string) []MapItem {
	items := make([]MapItem, 0, len(data))
	for _, key := range order {
		if value, ok := data[key]; ok {
			items = append(items, MapItem{Key: key, Value: value})
		}
	}
	sort.Sort(ByKey(items))
	return items
}
