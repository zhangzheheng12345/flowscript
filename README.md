# ![FlowScript!](icon/flowscript.png)

A embeded script language interpreter written by Go.

## Usage

Use:

```go
import "https://github.com/zhangzheheng12345/flowscript"
```

and add:

```
require github.com/zhangzheheng12345/flowscript v0.1.1
replace github.com/zhangzheheng12345/flowscript v0.1.1 => github.com/zhangzheheng12345/FlowScript v0.1.1
```

in go.mod to import it to your Go project.

( *I don't know why there must be a replacement on my machine. Who could give me an answer?* )

Use:

```go
parser.RunCode( <your codes(type string)> )
```

to run FlowScript code.

## Examples

### Hello World

```
echo "Hello World!\n" # output: Hello World!
```

### Fibonacci sequence

```
def fibonacci a b begin
    add a b > echo - > echo "\n"
    add a b > smlr - 10000 > if - begin
        add a b > fibonacci b -
    end else begin end # else block is a must
end
echo "fibonacci begins ...\n"
fibonacci 1 1
```

### OOP

```
def Student name li begin # This function will be used as a constructor
    echo "constructor called!\n"
    struct begin
        def getName begin expr name end
        def getScore begin expr li end
        def addScore score begin app li score > Student name - end
    end
end
list 100 96 99 > Student "Zhang" - > var A -
A.getName > echo -
A.getScore > echo -
A.addScore 100 > var B -
B.getScore > echo -`
```