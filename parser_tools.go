package parser

import errlog "github.com/zhangzheheng12345/flowscript/error_logger"

/* Tool function, for cannot convert bool to int directly */
func BoolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func WantInt(value interface{}) int {
	v, ok := value.(int)
	if !ok {
		v, ok := value.(byte)
		if !ok {
			errlog.Err("runtime", errlog.Line, "int or char value wanted, but got other type value")
		}
		return int(v)
	}
	return v
}

func WantList(value interface{}) []interface{} {
	v, ok := value.([]interface{})
	if !ok {
		errlog.Err("runtime", errlog.Line, "int list wanted, but got other type value")
	}
	return v
}

func WantString(value interface{}) string {
	v, ok := value.(string)
	if !ok {
		errlog.Err("runtime", errlog.Line, "string wanted, but got other type value")
	}
	return v
}

func AbsIndex(length int, index int) (int, bool) {
	if index >= 0 && index < length {
		return index, false
	} else if index < 0 && -index <= length {
		return length + index, false
	} else {
		errlog.Err("runtime", errlog.Line, "index out of range. index:", index, "length of the list:", length)
		return 0, true
	}
}

func WantStruct(value interface{}) Struct {
	v, ok := value.(Struct)
	if !ok {
		errlog.Err("runtime", errlog.Line, "struct wanted, but got other type value")
	}
	return v
}

func WantFunc(value interface{}) Func_ {
	v, ok := value.(Func_)
	if !ok {
		errlog.Err("runtime", errlog.Line, "function wanted, but got other type value")
		return GoFunc{
			func(i []interface{}) interface{} { return nil },
			0, make([]interface{}, 0)}
	}
	return v
}
