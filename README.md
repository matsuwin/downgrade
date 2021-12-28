# downgrade
Go 的熔断降级和重试功能库

<br>

## Quick Start

```go
// 新建任务实例，期望一秒内完成，三次重试机会
work := New(time.Second, 3)

// 定义主要任务内容
work.PlanA = func() error {
    fmt.Println("run planA")
    return errors.New("planA error")
}

// 定义备选任务内容
work.PlanB = func(err error) error {
    fmt.Println(err)
    fmt.Println("run planB")
    return nil
}

// 开始任务
if err := work.Do(); err != nil {
    panic(err)
}
```
```
run planA
run planA
run planA
planA error
run planB
```

## Installing

```
go get github.com/matsuwin/downgrade
```
