package randomname

import (
	"fmt"
	"testing"

	"github.com/wandersoulz/godes"
)

func TestGetManyNames(t *testing.T) {
	godes.SetSeed(848243241)

	for i := 0; i < 1; i++ {
		fmt.Println(GetName())
	}
}
