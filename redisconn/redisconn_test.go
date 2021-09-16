package redisconn

import (
	"testing"
	"entry_task/redisconn"
	// "fmt"
)

func TestGetToken(t *testing.T) {
	redisconn.InitPool()
	redisconn.SetToken("wanlb1", "123456")

	pass := redisconn.GetToken("wanlb1")

	if pass != "123456" {
		t.Errorf("error(-1) = %s; want 123456", pass)
	}
}