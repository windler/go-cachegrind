# go-cachegrind
`go-cachegrind` is a GO library that parses files in [cachegrind compatible format](http://valgrind.org/docs/manual/cg-manual.html), used e.g. by [xdebug profiler](https://xdebug.org/docs/profiler). This package is mainly developed to parse xdebug profiling files.

# Installation 
```bash
go get github.com/windler/go-cachegrind
```

# Usage

```go
    cg := cachegrind.Parse("path/to/file.cachegrind")
    main := cg.GetMainFunction()
    totalTime := main.GetMeasurement("Time")

    for _, call := range main.GetCalls() {
        calledFn := call.GetFunction()
        //... 
    }
```