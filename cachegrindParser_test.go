package cachegrind

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
PHP-CODE:
<?php

fun1();

function fun1() {
    for($i = 0; $i < 2; $i++)
        fun2();
    echo "test";
    fun2();
}

function fun2() {
    sleep(5);
}
*/
func TestParseExample1(t *testing.T) {
	cg, err := Parse("internal/testfiles/s1.cachegrind")
	assert.Nil(t, err)

	main := cg.GetMainFunction()
	assert.NotNil(t, main)
	assert.Equal(t, int64(15000303), main.GetMeasurement("Time"))
	assert.Equal(t, int64(32), main.GetMeasurement("Memory"))

	calls := main.GetCalls()
	assert.Equal(t, 1, len(calls))

	call1 := calls[0]
	assert.Equal(t, int64(15000293), call1.GetMeasurement("Time"))
	assert.Equal(t, int64(32), call1.GetMeasurement("Memory"))

	calledFunction := call1.GetFunction()
	assert.Equal(t, "/var/www/html/index.php", calledFunction.GetFile())
	assert.Equal(t, "fun1", calledFunction.GetName())
	assert.Equal(t, 3, len(calledFunction.GetCalls()))
	assert.Equal(t, int64(15000290), calledFunction.GetMeasurement("Time"))
	assert.Equal(t, int64(32), calledFunction.GetMeasurement("Memory"))

	for _, call := range calledFunction.GetCalls() {
		assert.Equal(t, "/var/www/html/index.php", call.GetFunction().GetFile())
		assert.Equal(t, "fun2", call.GetFunction().GetName())
	}
}
