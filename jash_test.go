package jash

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	j := New[int](256)

	j.Set("one", 1)
	j.Set("two", 2)

	v, ok := j.GetExists("one")
	require.True(t, ok)
	require.Equal(t, 1, v)

	v, ok = j.GetExists("two")
	require.True(t, ok)
	require.Equal(t, 2, v)

	v, ok = j.GetExists("three")
	require.False(t, ok)
	require.Equal(t, 0, v)
}

func TestMultiple(t *testing.T) {
	j := New[string](256)

	var keys []string
	for i := 0; i < 1_000_000; i++ {
		k := RandStringBytes(10)
		j.Set(k, k)
		keys = append(keys, k)
	}

	for _, k := range keys {
		v, ok := j.GetExists(k)
		require.True(t, ok)
		require.Equal(t, k, v)
	}
}

func bench(j *Jash[string]) {
	var keys []string
	for i := 0; i < 1_000_000; i++ {
		k := RandStringBytes(10)
		j.Set(k, k)
		keys = append(keys, k)
	}

	for _, k := range keys {
		j.GetExists(k)
	}
}

func BenchmarkNoGrow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		j := New[string](256)
		j.resize = false
		bench(j)
	}
}

func BenchmarkGrow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		j := New[string](256)
		bench(j)
	}
}

func BenchmarkMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := map[string]string{}
		var keys []string
		for i := 0; i < 1_000_000; i++ {
			k := RandStringBytes(10)
			h[k] = k
			keys = append(keys, k)
		}

		for _, k := range keys {
			_ = h[k]
		}
	}
}

// copied from https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
