package parser

import errlog "github.com/zhangzheheng12345/flowscript/error_logger"

/*
Global variable, Scope, contains all the variable at present.
If someone want to add global variable by Golang in a embeding way,
he should Scope.Add(string,interface{}) variable to Scope directly,
because all the codes run with RunCode(string) will be inside a independence Block_,
so you can't add global variable / function by native FlowScript.
*/
var Scope *Scope_ = &Scope_{nil, nil, make(map[string]interface{}), 0}

/*
exp1 > exp2 > exp3
The values of the send sequence expressions will be pushed to this queue.
*/
var tmpQueue *Queue_ = MakeTmpQueue(nil, 0)

/* scope system ( only for variables ) */
type Scope_ struct {
	last        *Scope_
	father      *Scope_
	vars        map[string]interface{}
	enumCounter int
}

func MakeScope(father *Scope_, last *Scope_) *Scope_ {
	return &Scope_{last, father, make(map[string]interface{}), father.enumCounter}
}

func (scope_ Scope_) Add(key string, value interface{}) {
	_, ok := scope_.vars[key]
	if ok {
		errlog.Err("runtime", errlog.Line, "Try to add a variable that has been added. name:", key)
	} else {
		scope_.vars[key] = value
	}
}

func (scope_ Scope_) Find(key string) interface{} {
	v, ok := scope_.vars[key]
	if ok {
		return v
	} else if scope_.father != nil {
		return scope_.father.Find(key)
	} else {
		errlog.Err("runtime", errlog.Line, "Try to read a variable that hasn't been added. name:", key)
	}
	return 0
}

func (scope_ Scope_) Back() *Scope_ {
	return scope_.last
}

type Queue_ struct {
	father  *Queue_
	data    []interface{}
	dataLen int
}

func MakeTmpQueue(father *Queue_, initSize int) *Queue_ {
	return &Queue_{father, make([]interface{}, initSize), 0}
}

func (queue_ *Queue_) Add(value interface{}) {
	queue_.dataLen++
	queue_.data[queue_.dataLen-1] = value
}

func (queue_ Queue_) Get() interface{} {
	if queue_.dataLen > 0 {
		return queue_.data[0]
	}
	errlog.Err("runtime", errlog.Line, "Try to get a value from temp queue while it's empty")
	return nil
}

func (queue_ *Queue_) Pop() {
	queue_.dataLen--
	queue_.data = queue_.data[1:]
}

func (queue_ *Queue_) Clear() *Queue_ {
	queue_.data = nil
	return queue_.father
}

func (queue_ Queue_) Size() int {
	return queue_.dataLen
}
