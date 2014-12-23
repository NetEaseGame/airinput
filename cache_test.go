package main

import (
	"image"
	"testing"
)

var cache = NewRGBACache(2)

func TestCache(t *testing.T) {
	im1 := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{10, 10}})
	cache.Put("1234", im1)
	cache.Put("4234", im1)
	cache.Put("5234", im1)
	if _, exists := cache.Get("1234"); exists {
		t.Fatal("should not exists 1234")
	}
	if im, exists := cache.Get("4234"); !exists || im == nil {
		t.Fatal("get cache failed")
	}
}
