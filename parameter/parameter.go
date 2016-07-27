package parameter

import (
	"sort"
)

func Escape(s string) string {
	t := make([]byte, 0, 3*len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if IsEscapable(c) {
			t = append(t, '%')
			t = append(t, "0123456789ABCDEF"[c>>4])
			t = append(t, "0123456789ABCDEF"[c&15])
		} else {
			t = append(t, s[i])
		}
	}
	return string(t)
}

func IsEscapable(b byte) bool {
	return !('A' <= b && b <= 'Z' || 'a' <= b && b <= 'z' || '0' <= b && b <= '9' || b == '-' || b == '.' || b == '_' || b == '~')
}

type ByValue []string

func (a ByValue) Len() int {
	return len(a)
}

func (a ByValue) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByValue) Less(i, j int) bool {
	return a[i] < a[j]
}

type OrderedParams struct {
	allParams   map[string][]string
	keyOrdering []string
}

func NewOrderedParams() *OrderedParams {
	return &OrderedParams{
		allParams:   make(map[string][]string),
		keyOrdering: make([]string, 0),
	}
}

func (o *OrderedParams) Get(key string) []string {
	sort.Sort(ByValue(o.allParams[key]))
	return o.allParams[key]
}

func (o *OrderedParams) Keys() []string {
	sort.Sort(o)
	return o.keyOrdering
}

func (o *OrderedParams) Add(key, value string) {
	o.AddUnescaped(key, Escape(value))
}

func (o *OrderedParams) AddUnescaped(key, value string) {
	if _, exists := o.allParams[key]; !exists {
		o.keyOrdering = append(o.keyOrdering, key)
		o.allParams[key] = make([]string, 1)
		o.allParams[key][0] = value
	} else {
		o.allParams[key] = append(o.allParams[key], value)
	}
}

func (o *OrderedParams) Len() int {
	return len(o.keyOrdering)
}

func (o *OrderedParams) Less(i int, j int) bool {
	return o.keyOrdering[i] < o.keyOrdering[j]
}

func (o *OrderedParams) Swap(i int, j int) {
	o.keyOrdering[i], o.keyOrdering[j] = o.keyOrdering[j], o.keyOrdering[i]
}

func (o *OrderedParams) Clone() *OrderedParams {
	clone := NewOrderedParams()
	for _, key := range o.Keys() {
		for _, value := range o.Get(key) {
			clone.AddUnescaped(key, value)
		}
	}
	return clone
}

type pair struct {
	Key   string
	Value string
}

type Pairs []pair

func (p Pairs) Len() int           { return len(p) }
func (p Pairs) Less(i, j int) bool { return p[i].Key < p[j].Key }
func (p Pairs) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func ParamsToSortedPairs(params map[string]string) Pairs {
	// Sort parameters alphabetically
	paramPairs := make(Pairs, len(params))
	i := 0
	for key, value := range params {
		paramPairs[i] = pair{Key: key, Value: value}
		i++
	}
	sort.Sort(paramPairs)

	return paramPairs
}
