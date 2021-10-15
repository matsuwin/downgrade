package downgrade

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	work := New(time.Second, 3)
	work.Plan1 = func() error {
		fmt.Println("run plan1")
		return errors.New("plan1 error")
	}
	work.Plan2 = func(err error) error {
		fmt.Println(err)
		fmt.Println("run plan2")
		return nil
	}
	if err := work.Do(); err != nil {
		panic(err)
	}
}
