package gotry

import (
	"fmt"
	"github.com/pubgo/g/errors"
	"reflect"
	"strconv"
	"sync"
	"time"
)

// TryRaw errors
func TryRaw(fn reflect.Value) func(...reflect.Value) func(...reflect.Value) (err error) {
	errors.PanicM(errors.AssertFn(fn), "func error")

	var variadicType reflect.Value
	var isVariadic = fn.Type().IsVariadic()
	if isVariadic {
		variadicType = reflect.New(fn.Type().In(fn.Type().NumIn() - 1).Elem()).Elem()
	}

	_NumIn := fn.Type().NumIn()
	return func(args ...reflect.Value) func(...reflect.Value) (err error) {
		errors.PanicT(isVariadic && len(args) < _NumIn-1, "func %s input params is error,func(%d,%d)", fn.Type(), _NumIn, len(args))
		errors.PanicT(!isVariadic && _NumIn != len(args), "func %s input params is not match,func(%d,%d)", fn.Type(), _NumIn, len(args))

		for i, k := range args {
			if _isZero(k) {
				args[i] = reflect.New(fn.Type().In(i)).Elem()
				continue
			}

			if isVariadic {
				args[i] = variadicType
			}

			args[i] = k
		}

		return func(cfn ...reflect.Value) (err error) {
			defer func() {
				errors.ErrHandle(recover(), func(_err *errors.Err) {
					_fn := fn
					if len(cfn) > 0 && !_isZero(cfn[0]) {
						_fn = cfn[0]
					}
					_err.Caller(_caller.FromFunc(_fn))
					err = _err
				})
			}()

			_c := fn.Call(args)
			if len(cfn) > 0 && !_isZero(cfn[0]) {
				errors.PanicM(errors.AssertFn(cfn[0]), "func type error")
				errors.PanicTT(cfn[0].Type().NumIn() != fn.Type().NumOut(), func(err *errors.Err) {
					err.Msg("callback func input num and output num not match[%d]<->[%d]", cfn[0].Type().NumIn(), fn.Type().NumOut())
				})

				if cfn[0].Type().NumIn() != 0 && cfn[0].Type().In(0) != fn.Type().Out(0) {
					errors.PanicTT(true, func(err *errors.Err) {
						err.Msg("callback func out type error [%s]<->[%s]", cfn[0].Type().In(0), fn.Type().Out(0))
					})
				}

				cfn[0].Call(_c)
			}
			return
		}
	}
}

// Try errors
func Try(fn interface{}) func(...interface{}) func(...interface{}) (err error) {
	_tr := TryRaw(reflect.ValueOf(fn))

	return func(args ...interface{}) func(...interface{}) (err error) {
		var _args = valueGet()
		defer valuePut(_args)

		for _, k := range args {
			_args = append(_args, reflect.ValueOf(k))
		}
		_tr1 := _tr(_args...)

		return func(cfn ...interface{}) (err error) {
			var _cfn = valueGet()
			defer valuePut(_cfn)

			for _, k := range cfn {
				_cfn = append(_cfn, reflect.ValueOf(k))
			}
			return _tr1(_cfn...)
		}
	}
}

// Retry errors
func Retry(num int, fn func()) (err error) {
	errors.PanicT(num < 1, "the num is less than 0")

	var all = 0
	var _fn = TryRaw(reflect.ValueOf(fn))
	for i := 0; i < num; i++ {
		if err = _fn()(); err == nil {
			return
		}

		all += i
		if Cfg.Debug {
			fmt.Printf("Retry current state, cur_sleep_time: %d, all_sleep_time: %d\n", i, all)
		}
		time.Sleep(time.Second * time.Duration(i))
	}

	errors.PanicMM(err, func(err *errors.Err) {
		err.Msg("retry error,retry_num: " + strconv.Itoa(num))
	})
	return
}

// RetryAt errors
func RetryAt(t time.Duration, fn func(at time.Duration)) {
	var err error
	var all = time.Duration(0)
	var _fn = TryRaw(reflect.ValueOf(fn))
	for {
		if err = _fn(reflect.ValueOf(all))(); err == nil {
			return
		}

		all += t
		errors.PanicTT(all > errors.Cfg.MaxRetryDur, func(err *errors.Err) {
			err.Msg("more than the max(%s) retry duration", errors.Cfg.MaxRetryDur.String())
		})

		if Cfg.Debug {
			fmt.Printf("cur_retry_time: %f, all_retry_time: %f", t.Seconds(), all.Seconds())
			errors.ErrLog(err)
		}
		time.Sleep(t)
	}
}

// Ticker errors
func Ticker(fn func(dur time.Time) time.Duration) {
	var _err error
	var _dur = time.Duration(0)
	var _all = time.Duration(0)
	var _fn = TryRaw(reflect.ValueOf(fn))
	var _cfn = reflect.ValueOf(func(t time.Duration) {
		_dur = t
	})

	for i := 0; ; i++ {
		_err = _fn(reflect.ValueOf(time.Now()))(_cfn)
		if _dur < 0 {
			return
		}

		if _dur == 0 {
			_dur = time.Second
		}

		_all += _dur
		errors.PanicT(_all > errors.Cfg.MaxRetryDur, "more than the max ticker time")
		if Cfg.Debug {
			fmt.Printf("retry_count: %d, retry_all_time: %f", i, _all.Seconds())
			errors.ErrLog(_err)
		}
		time.Sleep(_dur)
	}
}

var _valuePool = sync.Pool{
	New: func() interface{} {
		return []reflect.Value{}
	},
}

func valueGet() []reflect.Value {
	return _valuePool.Get().([]reflect.Value)
}

func valuePut(v []reflect.Value) {
	v = v[:0]
	_valuePool.Put(v)
}