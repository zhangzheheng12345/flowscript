package parser

import (
	"fmt"
	"strconv"
	"strings"

	errlog "github.com/zhangzheheng12345/flowscript/error_logger"
)

func Add_(args []interface{}) interface{} {
	v1 := args[0]
	v2 := args[1]
	switch v := v1.(type) {
	case int:
		return v + WantInt(v2)
	case byte:
		return int(v) + WantInt(v2)
	case []interface{}:
		/* Connect two lists */
		return append(v, WantList(v2)...)
	case string:
		return v + WantString(v2)
	default:
		errlog.Err("runtime", errlog.Line, "Try to add a unkown type value with another")
		return nil
	}
}

func Sub_(args []interface{}) interface{} {
	return WantInt(args[0]) - WantInt(args[1])
}

func Multi_(args []interface{}) interface{} {
	return WantInt(args[0]) * WantInt(args[1])
}

func Div_(args []interface{}) interface{} {
	return WantInt(args[0]) / WantInt(args[1])
}

func Mod_(args []interface{}) interface{} {
	return WantInt(args[0]) % WantInt(args[1])
}

func Bigr_(args []interface{}) interface{} {
	return BoolToInt(WantInt(args[0]) > WantInt(args[1]))
}

func Smlr_(args []interface{}) interface{} {
	return BoolToInt(WantInt(args[0]) < WantInt(args[1]))
}

func Equal_(args []interface{}) interface{} {
	v1 := args[0]
	v2 := args[1]
	switch v := v1.(type) {
	case int:
		return BoolToInt(v == WantInt(v2))
	case byte:
		return BoolToInt(int(v) == WantInt(v2))
	case []interface{}:
		/* Compare two lists */
		li2 := WantList(v2)
		for i, value := range v {
			if value != li2[i] {
				return 0
			}
		}
		return 1
	case string:
		return BoolToInt(v == WantString(v2))
	case Func_:
		errlog.Err("runtime", errlog.Line, "Try to compare a function to another. Operation: equal")
		return nil
	default:
		errlog.Err("runtime", errlog.Line, "Try to compare a unkown type value to another. Operation: equal")
		return nil
	}
}

func And_(args []interface{}) interface{} {
	return WantInt(args[0]) & WantInt(args[1])
}

func Or_(args []interface{}) interface{} {
	return WantInt(args[0]) | WantInt(args[1])
}

func Xor_(args []interface{}) interface{} {
	return WantInt(args[0]) ^ WantInt(args[1])
}

func Not_(args []interface{}) interface{} {
	return BoolToInt(WantInt(args[0]) == 0)
}

func Expr_(args []interface{}) interface{} {
	return args[0]
}

/* type cmd returns a string which tells the operand's type */
func Type_(args []interface{}) interface{} {
	switch v := args[0].(type) {
	case int:
		_ = v // FIXME: why must i write like this
		return "int"
	case byte:
		return "char"
	case []interface{}:
		return "list"
	case string:
		return "string"
	case Func_:
		return "function"
	case Struct:
		return "struct"
	default:
		return "?unknown_type?"
	}
}

func Toint_(args []interface{}) interface{} {
	switch v := args[0].(type) {
	case int:
		return v
	case byte:
		return int(v)
	case []interface{}:
		errlog.Err("runtime", errlog.Line, "Cannot convert list to int")
		return 0
	case string:
		res, err := strconv.Atoi(v)
		if err != nil {
			errlog.Err("runtime", errlog.Line, "Cannot convert string to int. string:", v)
			return 0
		}
		return res
	case Func_:
		errlog.Err("runtime", errlog.Line, "Cannot convert function to int")
		return 0
	case Struct:
		errlog.Err("runtime", errlog.Line, "Cannot convert structure to int")
		return 0
	default:
		return 0
	}
}

func Tochar_(args []interface{}) interface{} {
	switch v := args[0].(type) {
	case int:
		return byte(v)
	case byte:
		return v
	case []interface{}:
		errlog.Err("runtime", errlog.Line, "Cannot convert list to char")
		return 0
	case string:
		errlog.Err("runtime", errlog.Line, "Cannot convert string to char")
		return 0
	case Func_:
		errlog.Err("runtime", errlog.Line, "Cannot convert function to char")
		return 0
	case Struct:
		errlog.Err("runtime", errlog.Line, "Cannot convert structure to char")
		return 0
	default:
		return 0
	}
}

func Tostr_(args []interface{}) interface{} {
	switch v := args[0].(type) {
	case int:
		return strconv.Itoa(v)
	case byte:
		return string(v)
	case []interface{}:
		return fmt.Sprint(v)
	case string:
		return v
	case Func_:
		errlog.Err("runtime", errlog.Line, "Cannot convert function to string")
		return 0
	case Struct:
		errlog.Err("runtime", errlog.Line, "Cannot convert structure to string")
		return 0
	default:
		return 0
	}
}

/* output to stdout */
func Echo_(args []interface{}) interface{} {
	for k, op := range args {
		if k != 0 {
			fmt.Print(" ")
		}
		v, ok := op.(byte)
		if ok {
			fmt.Print(string([]byte{v}))
		} else {
			fmt.Print(op)
		}
	}
	return 0
}

/* Output to stdout */
func Echoln_(args []interface{}) interface{} {
	for _, op := range args {
		v, ok := op.(byte)
		if ok {
			fmt.Println(string([]byte{v}))
		} else {
			fmt.Println(op)
		}
	}
	return 0
}

/* Read from stdin */
func Input_(args []interface{}) interface{} {
	Echo_([]interface{}{args[0]})
	var res string
	_, err := fmt.Scan(&res)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return res
}

func List_(args []interface{}) interface{} {
	return args
}

func Len_(args []interface{}) interface{} {
	switch v := args[0].(type) {
	case int, byte:
		errlog.Err("runtime", errlog.Line, "Try to get the length of a int or char value.")
		return 0
	case []interface{}:
		return len(v)
	case string:
		return len(strings.Split(v, ""))
	default:
		errlog.Err("runtime", errlog.Line, "Try to get the length of a unknown type value.")
		return 0
	}
}

func Index_(args []interface{}) interface{} {
	index := WantInt(args[1])
	switch v := args[0].(type) {
	case int, byte:
		errlog.Err("runtime", errlog.Line, "Try to index a int or char value.")
		return 0
	case []interface{}:
		abs, err := AbsIndex(len(v), index)
		if err {
			return 0
		} else {
			return v[abs]
		}
	case string:
		abs, err := AbsIndex(len(strings.Split(v, "")), index)
		if err {
			return ""
		} else {
			return strings.Split(v, "")[abs]
		}
	default:
		errlog.Err("runtime", errlog.Line, "Try to index a unknown type value.")
		return 0
	}
}

func App_(args []interface{}) interface{} {
	switch v := args[0].(type) {
	case int, byte:
		errlog.Err("runtime", errlog.Line, "Try to append a value after a int or char value.")
		return 0
	case []interface{}:
		return append(v, args[1])
	case string:
		return v + string([]byte{byte(WantInt(args[1]))})
	default:
		errlog.Err("runtime", errlog.Line, "Try to append sth after a unknown type value.")
		return 0
	}
}

func Slice_(args []interface{}) interface{} {
	begin := WantInt(args[1])
	end := WantInt(args[2])
	var err1, err2 bool
	switch v := args[0].(type) {
	case int, byte:
		errlog.Err("runtime", errlog.Line, "Try to slice a int or char value.")
		return 0
	case []interface{}:
		begin, err1 = AbsIndex(len(v)+1, begin)
		end, err2 = AbsIndex(len(v)+1, end)
		if err1 || err2 {
			return []interface{}{}
		} else {
			return v[begin:end]
		}
	case string:
		splited := strings.Split(v, "")
		begin, err1 = AbsIndex(len(splited), begin)
		end, err2 = AbsIndex(len(splited), end)
		if err1 || err2 {
			return 0
		} else {
			return strings.Join(splited[begin:end], "")
		}
	default:
		errlog.Err("runtime", errlog.Line, "Try to slice a unknown type value.")
		return 0
	}
}

/* Split the string by words */
func Words_(args []interface{}) interface{} {
	str := WantString(args[0])
	return strings.Split(str, " ")
}

/* Split the string by lines*/
func Lines_(args []interface{}) interface{} {
	str := WantString(args[0])
	return strings.Split(str, "\n")
}

func Fmap_(args []interface{}) interface{} {
	v := WantList(args[0])
	res := make([]interface{}, len(v))
	f := WantFunc(args[1])
	for k, value := range v {
		res[k] = f.run([]interface{}{value})
	}
	return res
}

func Reduce_(args []interface{}) interface{} {
	f := WantFunc(args[1])
	var res interface{}
	for k, v := range WantList(args[0]) {
		if k == 0 {
			res = v
		} else {
			res = f.run([]interface{}{res, v})
		}
	}
	return res
}

func Filter_(args []interface{}) interface{} {
	v := WantList(args[0])
	res := make([]interface{}, 0)
	f := WantFunc(args[1])
	for _, value := range v {
		if f.run([]interface{}{value}) != 0 {
			res = append(res, value)
		}
	}
	return res
}

func AddBuildinFuncs() {
	AddGoFunc("add", Add_, 2)
	AddGoFunc("sub", Sub_, 2)
	AddGoFunc("multi", Multi_, 2)
	AddGoFunc("div", Div_, 2)
	AddGoFunc("mod", Mod_, 2)
	AddGoFunc("bigr", Bigr_, 2)
	AddGoFunc("smlr", Smlr_, 2)
	AddGoFunc("equal", Equal_, 2)
	AddGoFunc("and", And_, 2)
	AddGoFunc("or", Or_, 2)
	AddGoFunc("xor", Xor_, 2)
	AddGoFunc("not", Not_, 1)
	AddGoFunc("expr", Expr_, 1)
	AddGoFunc("type", Type_, 1)
	AddGoFunc("toint", Toint_, 1)
	AddGoFunc("tochar", Tochar_, 1)
	AddGoFunc("tostr", Tostr_, 1)
	AddGoFunc("echo", Echo_, -1) // args num = -1 means the arg num is limitless
	AddGoFunc("echoln", Echoln_, -1)
	AddGoFunc("input", Input_, 1)
	AddGoFunc("list", List_, -1)
	AddGoFunc("len", Len_, 1)
	AddGoFunc("index", Index_, 2)
	AddGoFunc("app", App_, 2)
	AddGoFunc("slice", Slice_, 3)
	AddGoFunc("words", Words_, 1)
	AddGoFunc("lines", Lines_, 1)
	AddGoFunc("fmap", Fmap_, 2)
	AddGoFunc("reduce", Reduce_, 2)
    AddGoFunc("filter", Filter_, 2)
}
