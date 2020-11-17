package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimMap(t *testing.T) {
	t.Run("returns false if map is the same", func(t *testing.T) {
		map1 := map[string]int{"b": 2, "a": 1}
		map2 := map[string]int{"a": 1, "b": 2}
		response := mapDiff(map1, map2)
		assert.Equal(t, false, response, "maps are equal")
	})

	t.Run("returns false if both maps are empty", func(t *testing.T) {
		map1 := map[string]int{}
		map2 := map[string]int{}
		response := mapDiff(map1, map2)
		assert.Equal(t, false, response, "maps are equal")
	})

	t.Run("returns true if new only contains original", func(t *testing.T) {
		orig := map[string]int{"a": 1}
		new := map[string]int{"a": 1, "b": 2}
		response := mapDiff(orig, new)
		assert.Equal(t, true, response, "maps are not equal")
	})

	t.Run("returns true if original only contains new", func(t *testing.T) {
		orig := map[string]int{"a": 1, "b": 2}
		new := map[string]int{"a": 1}
		response := mapDiff(orig, new)
		assert.Equal(t, true, response, "maps are not equal")
	})
}

func TestCopyMap(t *testing.T) {
	t.Run("clones map", func(t *testing.T) {
		map1 := map[string]int{
			"a": 1,
			"b": 2,
		}
		map2 := copyMap(map1)
		delete(map1, "a")

		assert.Equal(t, map[string]int{"b": 2}, map1, "map1 should be altered")
		assert.Equal(t, map[string]int{"a": 1, "b": 2}, map2, "map2 should not be altered")
	})

	t.Run("clones empty map", func(t *testing.T) {
		map1 := map[string]int{}
		map2 := copyMap(map1)
		map1["a"] = 1

		assert.Equal(t, map[string]int{"a": 1}, map1, "map1 should be altered")
		assert.Equal(t, map[string]int{}, map2, "map2 should not be altered")
	})
}

func TestDeleteExtra(t *testing.T) {
	t.Run("trims to 10", func(t *testing.T) {
		list := map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
			"d": 4,
			"e": 5,
			"f": 6,
		}
		newList := trimMap(list, 3)
		assert.Equal(t, 3, len(newList), "list should be length 10")
	})

	t.Run("no trim when smaller than length", func(t *testing.T) {
		list := map[string]int{
			"a": 1,
			"b": 2,
		}
		newList := trimMap(list, 3)
		assert.Equal(t, 2, len(newList), "returns all if under length")
	})

	t.Run("no trim when lenght is reached", func(t *testing.T) {
		list := map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		}
		newList := trimMap(list, 3)
		assert.Equal(t, 3, len(newList), "returns all if at length")
	})
}

func TestIntIndexInsert(t *testing.T) {
	t.Run("Insert value into slice", func(t *testing.T) {
		slice := []int{
			1, 2, 3,
		}
		response := intIndexInsert(slice, 1, 42)
		expected := []int{
			1, 42, 2, 3,
		}
		assert.Equal(t, expected, response, "expected 42 to be inserted")
	})

	t.Run("Insert value into emtpy slice", func(t *testing.T) {
		slice := []int{}
		response := intIndexInsert(slice, 0, 42)
		assert.Equal(t, []int{42}, response, "expected 42 to be inserted")
	})
}

func TestStrIndexInsert(t *testing.T) {
	t.Run("Insert value into slice", func(t *testing.T) {
		slice := []string{
			"a", "b", "c",
		}
		response := strIndexInsert(slice, 1, "test")
		expected := []string{
			"a", "test", "b", "c",
		}
		assert.Equal(t, response, expected, "expected 'test' to be inserted")
	})

	t.Run("Insert value into emtpy slice", func(t *testing.T) {
		slice := []string{}
		response := strIndexInsert(slice, 0, "test")
		assert.Equal(t, response, []string{"test"}, "expected 'test' to be inserted")
	})
}
