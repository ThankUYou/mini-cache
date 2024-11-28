.PHONY: test bench

# 运行所有测试
test:
	go test ./...

# 运行基准测试并显示内存分配信息
bench:
	go test -bench=. -benchmem ./benchmark

