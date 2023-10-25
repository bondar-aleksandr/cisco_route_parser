package main

import (
	"io"
)

// buildRoutesCache parses all routes from specified source
func buildRoutesCache(source io.Reader) {
	allRoutes = ParseRoute(source) 
}

// buildNHCache builds NH cache as map, where keys are hashes and values are *nextHop
func addNhToCache(nh *nextHop) {
	if _, ok := allNH[nh.GetHash()]; ok {
		return
	}
	allNH[nh.GetHash()] = nh
}