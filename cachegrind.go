package cachegrind

//Cachegrind represents a cachegrind-file
type Cachegrind interface {
	GetMainFunction() Function
}

//Function represents a function within a cachegrind-file
type Function interface {
	GetName() string
	GetFile() string
	GetCalls() []FunctionCall
	GetMeasurement(part string) int64
}

//FunctionCall represents a function call within a cachegrind-file
type FunctionCall interface {
	GetFunction() Function
	GetLine() int
	GetMeasurement(part string) int64
}
