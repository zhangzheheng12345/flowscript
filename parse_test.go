package parser
import (
    "testing"
    "fmt"
)

func TestRunCode(t *testing.T) {
    /* expected result: [output] 1 */
    const t1 = 
        "var a 1\n" + 
        "var b 3\n" + 
        "smlr a b > if - begin\n" +
        "\techo a\n" +
        "end"
    fmt.Println("t1----\n",t1)
    RunCode(t1)
    fmt.Println("t1----")
    /* expected result: [output] fibonacci sequence. Last one smaller tha 10 */
    const t2 =
        "def func a b begin\n" +
        "\tadd a b > echo -\n" +
        "\tsmlr a 10> if - begin\n" +
        "\t\tadd a b > func b -\n" +
        "\tend\n" +
        "end\n" +
        "func 1 2"
    fmt.Println("t2----\n",t2)
    RunCode(t2)
    fmt.Println("t2----")
    /* a small FlowScript example */
    const t3 =
        "def func a begin\n" +
        "\tmulti a a > echo -\n" +
        "\tsmlr a 1000 > if - begin\n" +
        "\t\tadd a 1 > func -\n" +
        "\tend\n" +
        "end\n" +
        "func 1\n"
    fmt.Println("t3----\n",t3)
    RunCode(t3)
    fmt.Println("t3----") 
}