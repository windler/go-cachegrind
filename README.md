[![CircleCI](https://circleci.com/gh/windler/go-cachegrind.svg?style=svg)](https://circleci.com/gh/windler/go-cachegrind) [![Go Report Card](https://goreportcard.com/badge/github.com/windler/go-cachegrind)](https://goreportcard.com/report/github.com/windler/go-cachegrind) [![codebeat badge](https://codebeat.co/badges/79fe429b-0f54-4c35-a359-4526a4294647)](https://codebeat.co/projects/github-com-windler-go-cachegrind-master)
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
