package parser

import (
	"fmt"
)

/*
Global variable, Scope, contains all the variable at present.
If someone want to add global variable by Golang in a embeding way,
he should Scope.Add(string,interface{}) variable to Scope directly,
because all the codes run with RunCode(string) will be inside a independence Block_,
so you can't add global variable / function by native FlowScript.
*/
var Scope *Scope_ = &Scope_{nil, nil, make(map[string]interface{})}

/*
exp1 > exp2 > exp3
The values of the send sequence expressions will be pushed to this queue.
*/
var tmpQueue *Queue_ = nil

/* scope system ( only for variables ) */
type Scope_ struct {
	last   *Scope_
	father *Scope_
	vars   map[string]interface{}
}

func MakeScope(father *Scope_, last *Scope_) *Scope_ {
	return &Scope_{last, father, make(map[string]interface{})}
}

func (scope_ Scope_) Add(key string, value interface{}) {
	_, ok := scope_.vars[key]
	if ok {
		fmt.Println("Error: Try to add a variable that has been added. var: " + key)
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
		fmt.Println("Error: Try to read a variable that hasn't been added. var: " + key)
	}
	return 0
}

func (scope_ Scope_) Back() *Scope_ {
	return scope_.last
}

// Give a proper beginning size of a queue
const maxBufferSize = 5

type Queue_ struct {
	father  *Queue_
	data    []interface{}
	dataLen int
}

func MakeTmpQueue(father *Queue_) *Queue_ {
	return &Queue_{father, make([]interface{}, maxBufferSize), 0}
}

func (queue_ *Queue_) Add(value interface{}) {
	queue_.dataLen++
	if queue_.dataLen > maxBufferSize {
		queue_.data = append(queue_.data, value)
	} else {
		queue_.data[queue_.dataLen-1] = value
	}
}

func (queue_ Queue_) Get() interface{} {
	return queue_.data[0]
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
