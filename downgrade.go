package downgrade

import (
	"crypto/tls"
	"github.com/matsuwin/siggroup/x/errcause"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"time"
)

// Request (HTTP Client)
func Request(header http.Header, method, url string, body io.Reader) (http.Header, []byte, error) {

	// 忽略证书校验
	var cli http.Client
	cli.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	var req *http.Request
	var res *http.Response
	var err error

	// 请求初始化
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return nil, nil, errors.New(err.Error())
	}
	req.Header = header

	// 发起请求
	res, err = cli.Do(req)
	if err != nil {
		return nil, nil, errors.New(err.Error())
	}
	defer res.Body.Close()

	// 读取数据
	var data []byte
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, errors.New(err.Error())
	}
	return res.Header, data, nil
}

// New (Create Work) Fuse downgrade and retry
func New(timeout time.Duration, retry int) *_Work {
	return &_Work{timeout: timeout, retry: retry}
}

// Do Start
func (work *_Work) Do() error {
	var sig = make(chan struct{})
	var err error

	// PlanA
	go func() {
		defer errcause.Recover()
		defer func() { sig <- struct{}{} }()

		if work.retry == 0 {
			work.retry = 1
		}
		for i := 0; i < work.retry; i++ {
			if err = work.PlanA(); err == nil {
				return
			}
		}
	}()

	// Wait
	select {
	case <-time.After(work.timeout):
		err = errors.New("downgrade.Work: timeout")
	case <-sig:
	}

	// PlanB
	if work.PlanB != nil && err != nil {
		err = work.PlanB(err)
	}
	return err
}

type _Work struct {
	PlanA   func() error
	PlanB   func(err error) error
	timeout time.Duration
	retry   int
}
