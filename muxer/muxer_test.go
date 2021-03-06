package muxer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMux(t *testing.T) {
	for _, tcase := range []TestCase{
		{
			name:            "NoDataOneSource",
			expectedOutputs: []int{},
			sourceChannels:  1,
		},
		{
			name:            "NoDataMultipleSources",
			expectedOutputs: []int{},
			sourceChannels:  10,
		},
		{
			name:            "OneOutputOneSource",
			expectedOutputs: []int{666},
			sourceChannels:  1,
		},
		{
			name:            "MultipleOutputsOneSource",
			expectedOutputs: []int{666, 777, 10, 0, 1},
			sourceChannels:  1,
		},
		{
			name:            "SameOutputsAsSources",
			expectedOutputs: []int{666, 777, 10},
			sourceChannels:  3,
		},
		{
			name:            "LessOutputsThanSources",
			expectedOutputs: []int{666, 777},
			sourceChannels:  10,
		},
		{
			name:            "MoreOutputsThanSources",
			expectedOutputs: []int{666, 777, 234, 1, 0},
			sourceChannels:  2,
		},
	} {
		testMux(t, tcase)
	}
}

func TestMuxClosedChannels(t *testing.T) {
	sink := make(chan int)
	source1 := make(chan int)
	source2 := make(chan int)

	close(source1)
	close(source2)

	assert.NoError(t, Do(sink, source1, source2))

	v, ok := <-sink
	if ok {
		t.Fatalf("expected sink to be closed, instead got val[%d]", v)
	}
}

func TestMuxCloseFirstSource(t *testing.T) {
	sink := make(chan int)
	source1 := make(chan int)
	source2 := make(chan int)
	source3 := make(chan int)

	assert.NoError(t, Do(sink, source1, source2, source3))

	close(source1)

	want1 := 666
	want2 := 777
	go func() {
		source3 <- want1
		source2 <- want2
	}()

	assert.Equal(t, want1, <-sink)
	assert.Equal(t, want2, <-sink)
}

func TestMuxCloseMiddleSource(t *testing.T) {
	sink := make(chan int)
	source1 := make(chan int)
	source2 := make(chan int)
	source3 := make(chan int)

	assert.NoError(t, Do(sink, source1, source2, source3))

	close(source2)

	want1 := 100
	want2 := 7
	go func() {
		source3 <- want1
		source1 <- want2
	}()

	assert.Equal(t, want1, <-sink)
	assert.Equal(t, want2, <-sink)
}

func TestMuxCloseLastSource(t *testing.T) {
	sink := make(chan int)
	source1 := make(chan int)
	source2 := make(chan int)
	source3 := make(chan int)

	assert.NoError(t, Do(sink, source1, source2, source3))

	close(source3)

	want1 := 100
	want2 := 7
	go func() {
		source2 <- want1
		source1 <- want2
	}()

	assert.Equal(t, want1, <-sink)
	assert.Equal(t, want2, <-sink)
}

func TestMuxDirectionedChannels(t *testing.T) {
	sink := make(chan string)
	source := make(chan string)

	var sinkSendOnly chan<- string = sink
	var sourceReadOnly <-chan string = source

	const expectedVal = "lambda"

	go func() {
		source <- expectedVal
		close(source)
	}()

	assert.NoError(t, Do(sinkSendOnly, sourceReadOnly))
	assert.Equal(t, expectedVal, <-sink)
}

func TestFailsOnWrongSourceDirection(t *testing.T) {
	sink := make(chan string)
	source := make(chan string)
	var sourceSendOnly chan<- string = source

	assert.Error(t, Do(sink, sourceSendOnly))
}

func TestFailsOnWrongSinkDirection(t *testing.T) {
	sink := make(chan string)
	source := make(chan string)

	var sinkReadOnly <-chan string = sink
	assert.Error(t, Do(sinkReadOnly, source))
}

func TestErrorOnInvalidSink(t *testing.T) {
	for name, invalidSink := range invalidCases() {
		t.Run(name, func(t *testing.T) {
			source := make(chan int)
			assert.Error(t, Do(invalidSink, source))
		})
	}
}

func TestErrorOnInvalidSource(t *testing.T) {
	for name, invalidSource := range invalidCases() {
		t.Run(name, func(t *testing.T) {
			sink := make(chan int)
			validSource := make(chan int)
			assert.Error(t, Do(sink, validSource, invalidSource))
			assert.Error(t, Do(sink, invalidSource, validSource))
		})
	}
}

type TestCase struct {
	name            string
	expectedOutputs []int
	sourceChannels  int
}

func invalidCases() map[string]interface{} {
	valid := make(chan int)
	var nilChannel chan int

	return map[string]interface{}{
		"nil":              nil,
		"nilChannel":       nilChannel,
		"notChannel":       1,
		"pointerToChannel": &valid,
		"wrongType":        make(chan uint),
	}
}

func testMux(t *testing.T, tcase TestCase) {
	t.Run(tcase.name, func(t *testing.T) {
		sources := []chan int{}
		sourcesgen := []interface{}{}
		for i := 0; i < tcase.sourceChannels; i++ {
			source := make(chan int)
			sources = append(sources, source)
			sourcesgen = append(sourcesgen, source)
		}
		sink := make(chan int)
		assert.NoError(t, Do(sink, sourcesgen...))

		go func() {
			for i, v := range tcase.expectedOutputs {
				srcindex := i % len(sources)
				sources[srcindex] <- v
			}
			for _, source := range sources {
				close(source)
			}
		}()

		gotOutputs := []int{}
		for got := range sink {
			gotOutputs = append(gotOutputs, got)
		}

		if len(gotOutputs) != len(tcase.expectedOutputs) {
			t.Fatalf("got %v != want %v", gotOutputs, tcase.expectedOutputs)
		}

		for i, want := range tcase.expectedOutputs {
			got := gotOutputs[i]
			assert.Equal(t, want, got)
		}
	})
}