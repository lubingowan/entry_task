功能测试：




性能测试：

cd goprofile/rpcclient;

go test -bench clientbase_test.go -parallel 200 -count 10000
PASS
ok  	rpcclient	0.230s

go test -bench clientbase_test.go -parallel 10000  -count 3000
PASS
ok  	rpcclient	0.076s
