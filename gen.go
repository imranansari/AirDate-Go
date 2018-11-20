//+build generate

package main

import "github.com/zserge/lorca"

func main() {
	// You can also run "npm build" or webpack here, or compress www, or
	// generate manifests, or do other preparations for your www.
	lorca.Embed("www.go", "www")
}
