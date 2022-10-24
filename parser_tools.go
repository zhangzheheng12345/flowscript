package parser

import errlog "github.com/zhangzheheng12345/flowscript/error_logger"

/*
A tool function, for cannot convert bool to int directly
*/
func BoolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func WantInt(value interface{}, ctx *Context = Global) int {
	v, ok := value.(int)
	if !ok {
		v, ok := value.(byte)
		if !ok {
			errlog.Err("runtime", ctx.Line, "int or char value wanted, but got other type value")
		}
		return int(v)
	}
	return v
}

func WantList(value interface{}, ctx *Context) []interface{} {
	v, ok := value.([]interface{})
	if !ok {
		errlog.Err("runtime", ctx.Line, "int list wanted, but got other type value")
	}
	return v
}

func WantString(value interface{}, ctx *Context) string {
	v, ok := value.(string)
	if !ok {
		errlog.Err("runtime", ctx.Line, "string wanted, but got other type value")
	}
	return v
}

func AbsIndex(length int, index int, ctx *Context) (int, bool) {
	if index >= 0 && index < length {
		return index, false
	} else if index < 0 && -index <= length {
		return length + index, false
	} else {
		errlog.Err("runtime", ctx.Line, "index out of range. index:", index, "length of the list:", length)
		return 0, true
	}
}

func WantStruct(value interface{}, ctx *Context) Struct {
	v, ok := value.(Struct)
	if !ok {
		errlog.Err("runtime", ctx.Line, "struct wanted, but got other type value")
	}
	return v
}

func WantFunc(value interface{}, ctx *Context) Func_ {
	v, ok := value.(Func_)
	if !ok {
		errlog.Err("runtime", ctx.Line, "function wanted, but got other type value")
		return GoFunc{
			func(i []interface{}, ctx *Context) interface{} { return nil },
			0, make([]interface{}, 0)}
	}
	return v
}
