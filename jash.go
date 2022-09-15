package jash

import (
	"hash"
	"hash/fnv"
)

type data[T any] struct {
	key   string
	value T
}

type Jash[T any] struct {
	buckets [][]data[T]
	size    int
	hash    hash.Hash32
	resize  bool
}

func New[T any](size int) *Jash[T] {
	return &Jash[T]{
		buckets: make([][]data[T], size),
		size:    size,
		hash:    fnv.New32(),
		resize:  true,
	}
}

func (j *Jash[T]) Set(key string, value T) {
	b := j.bucket(key)
	j.buckets[b] = append(j.buckets[b], data[T]{
		key:   key,
		value: value,
	})

	if j.resize && len(j.buckets[b]) > 256 {
		j.grow()
	}
}

func (j *Jash[T]) GetExists(key string) (T, bool) {
	b := j.bucket(key)
	data := j.buckets[b]

	for _, d := range data {
		if d.key == key {
			return d.value, true
		}
	}

	var z T
	return z, false
}

func (j *Jash[T]) grow() {
	newSize := j.size * j.size
	j.size = newSize
	buckets := make([][]data[T], newSize)

	for _, b := range j.buckets {
		for _, d := range b {
			k := j.bucket(d.key)
			buckets[k] = append(buckets[k], d)
		}
	}

	j.buckets = buckets
}

func (j *Jash[T]) bucket(key string) int {
	j.hash.Reset()
	j.hash.Write([]byte(key))
	h := j.hash.Sum32()
	k := int(h) % (j.size)
	return k
}
