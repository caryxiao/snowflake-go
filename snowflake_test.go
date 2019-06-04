package snowflake_go

import (
	"fmt"
	"testing"
)

func TestSnowflake_GetId(t *testing.T) {
	sf, err := New(1)
	if err != nil {
		fmt.Println(err)
		return
	}

	idCh := make(chan int64)
	defer close(idCh)
	for i := 0; i < 10000; i++ {
		go func(sf *Snowflake) {
			nid := sf.GetId()
			idCh <- nid
		}(sf)
	}

	idMap := make(map[int64]int)
	for i := 0; i < 10000; i++ {
		id := <-idCh
		if _, ok := idMap[id]; ok {
			t.Error("ID is not unique!\n")
			return
		}

		idMap[id] = i
	}

	fmt.Printf("total: %d", len(idMap))
}
