package downgrade

import (
	"crypto/tls"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/matsuwin/siggroup/x/errcause"
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
	data, err = ioutil.ReadAll(res.Body)
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

	// Plan1
	go func() {
		defer errcause.Recover()
		defer func() { sig <- struct{}{} }()

		if work.retry == 0 {
			work.retry = 1
		}
		for i := 0; i < work.retry; i++ {
			if err = work.Plan1(); err == nil {
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

	// Plan2
	if work.Plan2 != nil && err != nil {
		err = work.Plan2(err)
	}
	return err
}

type _Work struct {
	Plan1   func() error
	Plan2   func(err error) error
	timeout time.Duration
	retry   int
}
