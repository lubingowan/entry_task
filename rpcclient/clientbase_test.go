package rpcclient

import (
	"testing"
	"entry_task/rpcclient"
	"strconv"
)

func BenchmarkGetAuth(b *testing.B) {
	username := "test"
	password := "password"
	for i:= 0; i< b.N; i++ {
		istr := strconv.Itoa(i)
		rpcclient.GetAuth(username + istr, password + istr)
	}
}