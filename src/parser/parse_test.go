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
        "\techo a\n"
        "end"
    fmt.Println("t1----\n",t1)
    RunCode(t1)
    fmt.Println("t1----")
    /* expected result: [output] fibonacci sequence. Last one smaller tha 100 */
    const t2 =
        "def func a b begin\n" +
        "\tadd a b > echo -\n" +
        "\tsmlr a 100 > if - begin\n" +
        "\t\tadd a b > func b -\n" +
        "\tend\n" +
        "end\n" +
        "func 1 2"
    fmt.Println("t2----\n",t2)
    RunCode(t1)
    fmt.Println("t2----")
}