package cmd

import (
	"fmt"
	"testing"
)

func TestDirWalk(t *testing.T) {
	dir := "/Applications"
	for _, app := range dirWalk(dir) {
		fmt.Println(app.Name)
	}
}
