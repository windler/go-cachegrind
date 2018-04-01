package cachegrind

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var events []string

type goCachegrind struct {
	functions       []*cgFn
	parseContext    parseContext
	functionMap     map[string]*cgFn
	functionNameMap map[string]string
	mainFn          *cgFn
}

type parseContext struct {
	currentFn      *cgFn
	lastCalledFile string
	lastFileID     string
	lastFile       string
}

type cgFn struct {
	calls        []*cgCall
	name         string
	file         string
	nameID       string
	fileID       string
	measurements []int64
}

type cgCall struct {
	fn           *cgFn
	line         int
	measurements []int64
}

func (fn cgFn) GetName() string {
	return fn.name
}

func (fn cgFn) GetFile() string {
	return fn.file
}

func (fn cgFn) GetCalls() []FunctionCall {
	res := []FunctionCall{}
	for _, fc := range fn.calls {
		res = append(res, fc)
	}
	return res
}

func (fn cgFn) GetMeasurement(part string) int64 {
	measure := fn.measurements[getMeasurementIndex(part)]

	for _, call := range fn.calls {
		measure += call.GetMeasurement(part)
	}

	return measure
}

func getMeasurementIndex(part string) int {
	i := 0
	for _, p := range events {
		if part == p {
			return i
		}
		i++
	}
	return -1
}

func (cg goCachegrind) GetMainFunction() Function {
	return cg.mainFn
}

func (c cgCall) GetFunction() Function {
	return c.fn
}

func (c cgCall) GetLine() int {
	return c.line
}

func (c cgCall) GetMeasurement(part string) int64 {
	return c.measurements[getMeasurementIndex(part)]
}

//Parse parses a callgrind file content and creates a Cachegrind object
func Parse(fileName string) (Cachegrind, error) {
	cg := &goCachegrind{
		functions:       []*cgFn{},
		functionMap:     map[string]*cgFn{},
		functionNameMap: map[string]string{},
		parseContext:    parseContext{},
	}

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cg.parseLine(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return cg, err
}

func (cg *goCachegrind) parseLine(line string) {
	if strings.HasPrefix(line, "events:") {
		cg.callParseFunction(line, "events: ", cg.parseEvents)
	} else if strings.HasPrefix(line, "fl=") {
		cg.callParseFunction(line, "fl=", cg.parseFile)
	} else if strings.HasPrefix(line, "fn=") {
		cg.callParseFunction(line, "fn=", cg.parseFunction)
	} else if strings.HasPrefix(line, "cfl=") {
		cg.callParseFunction(line, "cfl=", cg.parseCalledFile)
	} else if strings.HasPrefix(line, "cfn=") {
		cg.callParseFunction(line, "cfn=", cg.parseCalledFunction)
	}

	measurementLine := regexp.MustCompile(`^(\d+ )+`)
	if measurementLine.MatchString(line) {
		cg.parseMeasurement(line)
	}
}

func (cg *goCachegrind) callParseFunction(line, prefix string, f func(string)) {
	f(strings.Trim(line, prefix))
}

/*
Line format:
	3 15000293 32
*/
func (cg *goCachegrind) parseMeasurement(line string) {
	measurements := strings.Split(line, " ")
	callLine, _ := strconv.Atoi(measurements[0])

	for _, m := range measurements[1:] {
		measure, _ := strconv.ParseInt(m, 10, 64)
		if len(cg.parseContext.currentFn.calls) > 0 {
			cg.getLastCall().measurements = append(cg.getLastCall().measurements, measure)
			cg.getLastCall().line = callLine
		} else {
			cg.parseContext.currentFn.measurements = append(cg.parseContext.currentFn.measurements, measure)
		}
	}

}

/*
Line format (without prefix):
	Time Memory
*/
func (cg *goCachegrind) parseEvents(line string) {
	events = strings.Split(line, " ")
}

/*
Line format (without prefix):
	(2) /var/www/html/index.php
or
	(2)
*/
func (cg *goCachegrind) parseFile(line string) {
	idAndFile := strings.Split(line, " ")

	if len(idAndFile) == 2 {
		cg.parseContext.lastFile = idAndFile[1]
		cg.parseContext.lastFileID = idAndFile[0]
	}
	cg.parseContext.lastFileID = idAndFile[0]
}

/*
Line format (without prefix):
	fn=(2) fun2
or
	(2)
*/
func (cg *goCachegrind) parseFunction(line string) {
	idAndFile := strings.Split(line, " ")

	fnID := cg.parseContext.lastFileID + "_" + idAndFile[0]
	if cg.functionMap[fnID] == nil {
		var name string
		if len(idAndFile) == 2 {
			name = idAndFile[1]
		} else {
			name = cg.functionNameMap[line]
		}

		newFn := &cgFn{
			name:   name,
			nameID: idAndFile[0],
			file:   cg.parseContext.lastFile,
			fileID: cg.parseContext.lastFileID,
			calls:  []*cgCall{},
		}
		cg.functions = append(cg.functions, newFn)
		cg.parseContext.currentFn = newFn

		cg.functionMap[fnID] = newFn
		cg.functionNameMap[line] = name

		if strings.HasSuffix(line, "{main}") {
			cg.mainFn = newFn
		}
	} else {
		cg.parseContext.currentFn = cg.functionMap[fnID]
	}
}

/*
Line format (without prefix):
	cfn=(1)
*/
func (cg *goCachegrind) parseCalledFunction(line string) {
	fnID := cg.parseContext.lastCalledFile + "_" + line
	calledFn := cg.functionMap[fnID]

	cg.getLastCall().fn = calledFn
}

/*
Line format (without prefix):
	cfl=(1)
*/
func (cg *goCachegrind) parseCalledFile(line string) {
	cg.parseContext.lastCalledFile = line
	cg.parseContext.currentFn.calls = append(cg.parseContext.currentFn.calls, &cgCall{})
}

func (cg *goCachegrind) getLastCall() *cgCall {
	return cg.parseContext.currentFn.calls[len(cg.parseContext.currentFn.calls)-1]
}
