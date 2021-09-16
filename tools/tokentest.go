package main
import (
    "fmt"
	"entry_task/token"
)

func main() {
	// fmt.Println(token.ParseToken("001FAD4cwTVjjmSIoqYZlzc8U5sUuoFG0x5lp+3AeG1cmBewhEGQ2EALw0A+3AeG1"))
	tokenstr := token.GenToken("wanlb1", token.Sha1("password1"))
	fmt.Println(tokenstr)

	fmt.Println(token.ParseToken(tokenstr))

	fmt.Println(token.CheckTokenByToken(tokenstr, "001FABmII+mwY/GocZBLbwoZjfViuMEWUx5lp+3AeG1V5T+4MIJQ2EALw0A"))
	fmt.Println(token.CheckTokenByToken(tokenstr, "001FABMkc+Bv91NJ8t8ZZV0MGcLiVpdjEx5lp+3AeG1cPHaqJEKQ2EALw0A"))

	fmt.Println(token.CheckToken("001FABMkc+Bv91NJ8t8ZZV0MGcLiVpdjEx5lp+3AeG1cPHaqJEKQ2EALw0A", "wanlb1", token.Sha1("password1")))
}