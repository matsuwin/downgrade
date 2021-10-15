# downgrade
Go 的熔断降级和重试功能库

<br>

## Installing

```
go get github.com/matsuwin/downgrade
```

## Quick Start

```go
// 新建任务实例，期望一秒内完成，三次重试机会
work := downgrade.New(time.Second, 3)

// 定义主要任务内容
work.Plan1 = func() error {
  fmt.Println("run plan1")
  return errors.New("plan1 error")
}

// 定义备选任务内容
work.Plan2 = func(err error) error {
  fmt.Println(err)
  fmt.Println("run plan2")
  return nil
}

// 开始任务
if err := work.Do(); err != nil {
  panic(err)
}
```
```
run plan1
run plan1
run plan1
plan1 error
run plan2
```
