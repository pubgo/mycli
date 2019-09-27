# errors
error handle for go

## export func
```go
package errors

import "github.com/pubgo/errors/internal"

// Err
type Err = internal.Err

// error assert
var Panic = internal.Panic
var Wrap = internal.Wrap
var WrapM = func(err interface{}, fn func(err *Err)) {
	internal.WrapM(err, fn)
}
var TT = func(b bool, fn func(err *Err)) {
	internal.TT(b, fn)
}
var T = internal.T

// error handle
var Throw = internal.Throw
var Assert = internal.Assert
var Resp = func(fn func(err *Err)) {
	internal.Resp(fn)
}
var RespErr = internal.RespErr
var ErrLog = internal.ErrLog
var ErrHandle = internal.ErrHandle
var Debug = internal.Debug

// test
type Test = internal.Test

var TestRun = func(fn interface{}, desc func(desc func(string) *Test)) {
	internal.TestRun(fn, desc)
}

// config
var Cfg = &internal.Cfg

// err tag
var ErrTagRegistry = internal.ErrTagRegistry
var ErrTags = internal.ErrTags
var ErrTagsMatch = internal.ErrTagsMatch

// utils
var AssertFn = internal.AssertFn
var If = internal.If
var IsZero = internal.IsZero
var IsNone = internal.IsNone
var P = internal.P
var FuncCaller = internal.FuncCaller
var GetCallerFromFn = internal.GetCallerFromFn
var LoadEnvFile = internal.LoadEnvFile
var InitDebug = internal.InitDebug

// try
var Try = internal.Try
var Retry = internal.Retry
var RetryAt = internal.RetryAt
var Ticker = internal.Ticker
```

## test
```go
package tests_test

import (
	es "errors"
	"fmt"
	"github.com/pubgo/errors"
	"github.com/pubgo/errors/internal"
	"io/ioutil"
	"os"
	"reflect"
	"runtime/trace"
	"strings"
	"testing"
	"time"
)

func init() {
	errors.InitDebug()
	//internal.InitSkipErrorFile()
}

func TestCfg(t *testing.T) {
	errors.P("errors.Cfg", errors.Cfg)
}

func TestT(t *testing.T) {
	errors.TestRun(errors.T, func(desc func(string) *errors.Test) {
		desc("params is true").In(true, "test t").IsErr()
		desc("params is false").In(false, "test t").IsNil()
	})
}

func TestErrLog2(t *testing.T) {
	errors.TestRun(errors.ErrLog, func(desc func(string) *errors.Test) {
		desc("err log params").In(es.New("sss")).IsNil()
		desc("nil params").In(es.New("sss")).IsNil()
	})
}

func TestRetry(t *testing.T) {
	defer errors.Assert()

	errors.TestRun(errors.Retry, func(desc func(string) *errors.Test) {
		desc("retry(3)").In(3, func() {
			errors.T(true, "test t")
		}).IsErr(func(err error) {
			errors.Wrap(err, "test Retry error")
		})
	})
}

func TestIf(t *testing.T) {
	defer errors.Assert()

	errors.T(errors.If(true, "test true", "test false").(string) != "test true", "")
}

func TestTT(t *testing.T) {
	defer errors.Assert()

	_fn := func(b bool) {
		errors.TT(b, func(err *internal.Err) {
			err.Msg("test tt")
			err.M("k", "v")
			err.SetTag("12")
		})
	}

	errors.TestRun(_fn, func(desc func(string) *errors.Test) {
		desc("true params 1").In(true).IsErr()
		desc("true params 2").In(true).IsErr()
		desc("true params 3").In(true).IsErr()
		desc("false params").In(false).IsNil()
	})
}

func TestWrap(t *testing.T) {
	defer errors.Assert()

	errors.Wrap(es.New("test"), "test")
}

func TestWrapM(t *testing.T) {
	defer errors.Assert()

	errors.Wrap(es.New("dd"), "test")
}

func testFunc_2() {
	errors.WrapM(es.New("testFunc_1"), func(err *internal.Err) {
		err.Msg("test shhh")
		err.M("ss", 1)
		err.M("input", 2)
	})
}

func testFunc_1() {
	testFunc_2()
}

func testFunc() {
	errors.Wrap(errors.Try(testFunc_1), "errors.Wrap")
}

func TestErrLog(t *testing.T) {
	defer errors.Assert()

	errors.TestRun(testFunc, func(desc func(string) *errors.Test) {
		desc("test func").In().IsErr()
	})
}

func init11() {
	errors.T(true, "test tt")
}

func TestT2(t *testing.T) {
	defer errors.Assert()

	errors.TestRun(init11, func(desc func(string) *errors.Test) {
		desc("simple test").In().IsErr()
	})
}

func TestTry(t *testing.T) {
	defer errors.Assert()

	errors.Panic(errors.Try(errors.T)(true, "sss"))
}

func TestTask(t *testing.T) {
	defer errors.Assert()

	errors.Wrap(errors.Try(func() {
		errors.Wrap(es.New("dd"), "err ")
	}), "test wrap")
}

func TestHandle(t *testing.T) {
	defer errors.Assert()

	func() {
		errors.Wrap(es.New("hello error"), "sss")
	}()

}

func TestErrHandle(t *testing.T) {
	defer errors.Assert()

	errors.ErrHandle(errors.Try(func() {
		errors.T(true, "test T")
	}), func(err *errors.Err) {
		err.P()
	})

	errors.ErrHandle("ttt", func(err *errors.Err) {
		err.P()
	})
	errors.ErrHandle(es.New("eee"), func(err *errors.Err) {
		err.P()
	})
	errors.ErrHandle([]string{"dd"}, func(err *errors.Err) {
		err.P()
	})
}

func TestIsZero(t *testing.T) {
	//defer errors.Log()

	var ss = func() map[string]interface{} {
		return make(map[string]interface{})
	}

	var ss1 = func() map[string]interface{} {
		return nil
	}

	var s = 1
	var ss2 map[string]interface{}
	errors.T(errors.IsZero(reflect.ValueOf(1)), "")
	errors.T(errors.IsZero(reflect.ValueOf(1.2)), "")
	errors.T(!errors.IsZero(reflect.ValueOf(nil)), "")
	errors.T(errors.IsZero(reflect.ValueOf("ss")), "")
	errors.T(errors.IsZero(reflect.ValueOf(map[string]interface{}{})), "")
	errors.T(errors.IsZero(reflect.ValueOf(ss())), "")
	errors.T(!errors.IsZero(reflect.ValueOf(ss1())), "")
	errors.T(errors.IsZero(reflect.ValueOf(&s)), "")
	errors.T(!errors.IsZero(reflect.ValueOf(ss2)), "")
}

func TestResp(t *testing.T) {
	defer errors.Assert()

	errors.TestRun(errors.Resp, func(desc func(string) *errors.Test) {
		desc("resp ok").In(func(err *errors.Err) {
			err.Caller(errors.FuncCaller(2))
		}).IsNil()
	})

}

func TestTicker(t *testing.T) {
	defer errors.Assert()

	errors.Ticker(func(dur time.Time) time.Duration {
		fmt.Println(dur)
		return time.Second
	})
}

func TestRetryAt(t *testing.T) {
	errors.RetryAt(time.Second*2, func(dur time.Duration) {
		fmt.Println(dur.String())

		errors.T(true, "test RetryAt")
	})
}

func TestErr(t *testing.T) {
	errors.ErrHandle(errors.Try(func() {
		errors.ErrHandle(errors.Try(func() {
			errors.T(true, "90999 error")
		}), func(err *errors.Err) {
			errors.Wrap(err, "wrap")
		})
	}), func(err *errors.Err) {
		fmt.Println(err.P())
	})
}

func _GetCallerFromFn2() {
	errors.WrapM(es.New("test 123"), func(err *internal.Err) {
		err.Msg("test GetCallerFromFn")
		err.M("ss", "dd")
	})
}

func _GetCallerFromFn1(fn func()) {
	errors.Panic(errors.AssertFn(reflect.ValueOf(fn)))
	fn()
}

func TestGetCallerFromFn(t *testing.T) {
	defer errors.Assert()

	fmt.Println(errors.GetCallerFromFn(reflect.ValueOf(_GetCallerFromFn1)))

	errors.TestRun(_GetCallerFromFn1, func(desc func(string) *errors.Test) {
		desc("GetCallerFromFn ok").In(_GetCallerFromFn2).IsErr()
		desc("GetCallerFromFn nil").In(nil).IsErr()
	})
}

func TestErrTagRegistry(t *testing.T) {
	defer errors.Assert()

	errors.ErrTagRegistry("errors_1")
	errors.ErrTagRegistry("errors_2")
	fmt.Printf("%#v\n", errors.ErrTags())

	errors.T(errors.ErrTagsMatch("errors") == true, "errors match error")
	errors.T(errors.ErrTagsMatch("errors_1") == false, "errors_1 not match")
}

func TestTest(t *testing.T) {
	defer errors.Assert()

	errors.TestRun(errors.AssertFn, func(desc func(string) *errors.Test) {
		desc("params is func 1").
			In(reflect.ValueOf(func() {})).
			IsNil(func(err error) {
				errors.Wrap(err, "check error")
			})

		desc("params is func 2").
			In(reflect.ValueOf(func() {})).
			IsNil(func(err error) {
				errors.Wrap(err, "check error")
			})

		desc("params is func 3").
			In(reflect.ValueOf(func() {})).
			IsNil(func(err error) {
				errors.Wrap(err, "check error")
			})

		desc("params is nil").
			In(reflect.ValueOf(nil)).
			IsErr(func(err error) {
				errors.Wrap(err, "check error ok")
			})
	})
}

func TestThrow(t *testing.T) {
	defer errors.Assert()

	errors.TestRun(errors.Throw, func(desc func(string) *errors.Test) {
		desc("not func type params").In(es.New("ss")).IsErr()
		desc("func type params").In(func() {}).IsNil()
		desc("nil type params").In(nil).IsErr()
	})
}

func TestLoadEnv(t *testing.T) {
	errors.LoadEnvFile("../.env")
	errors.T(os.Getenv("a") != "1", "env error")
}

func init2() (err error) {
	defer errors.RespErr(&err)

	errors.TT(true, func(err *errors.Err) {
		err.Msg("ok sss %d", 23)
	})
	return
}

func TestSig(t *testing.T) {
	defer errors.Assert()
	errors.Panic(init2())
}

func TestIsNone(t *testing.T) {
	defer errors.Debug()

	buf := &strings.Builder{}
	trace.Start(buf)
	defer func() {
		ioutil.WriteFile("trace.log", []byte(buf.String()), 0666)
	}()
	defer trace.Stop()

	errors.TestRun(errors.IsNone, func(desc func(string) *errors.Test) {
		desc("is null").In(nil).IsNil(func(b bool) {
			errors.T(b != true, "error")
		})
		desc("is ok").In("ok").IsNil(func(b bool) {
			errors.T(b == false, "error")
		})
	})
}
```