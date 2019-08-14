// +build integration

package change_points

// Collection of slow tests.
// Run with:
// go test -v  --tags="* integration" -parallel 8

import (
	"testing"
)

func TestLarge(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow test")
		return
	}
	ChangePointTestHelper(t)
}

func TestHuge(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow test")
	}
	ChangePointTestHelper(t)
}

func TestHugeEDM(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow test")
		return
	}
	EdmTestHelper(t)
}

func TestHumungousPlusEDM(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow test")
		return
	}
	EdmTestHelper(t)
}

func TestHumungousEDM(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow test")
		return
	}
	EdmTestHelper(t)
}