# Benchmarks

Performance benchmarks for sum.

## Running

```bash
make test-bench
```

## Writing Benchmarks

Place benchmark files here with the naming convention `*_bench_test.go`.

Example:

```go
//go:build testing

package benchmarks

import (
    "testing"
)

func BenchmarkExample(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Code to benchmark
    }
}
```

## Comparing Results

Use benchstat to compare benchmark runs:

```bash
go install golang.org/x/perf/cmd/benchstat@latest

# Run baseline
make test-bench > old.txt

# Make changes, then run again
make test-bench > new.txt

# Compare
benchstat old.txt new.txt
```
