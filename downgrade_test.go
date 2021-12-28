package downgrade

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {

	// 新建任务实例，期望一秒内完成，三次重试机会
	work := New(time.Second, 3)

	// 定义主要任务内容
	work.PlanA = func() error {
		fmt.Println("run plan1")
		return errors.New("plan1 error")
	}

	// 定义备选任务内容
	work.PlanB = func(err error) error {
		fmt.Println(err)
		fmt.Println("run plan2")
		return nil
	}

	// 开始任务
	if err := work.Do(); err != nil {
		panic(err)
	}
}
