package usecase

type FanIn struct {
	input  chan string
	input2 chan string
	input3 chan string
	output chan string
}

func NewFanIn(numSources int) *FanIn {
	return &FanIn{
		input:  make(chan string),
		input2: make(chan string),
		input3: make(chan string),
		output: make(chan string, 10),
	}
}

func (f *FanIn) Merge() chan string {
	go func() {
		for val := range f.input {
			f.output <- val
		}
	}()
	go func() {
		for val := range f.input2 {
			f.output <- val
		}
	}()
	go func() {
		for val := range f.input3 {
			f.output <- val
		}
	}()
	return f.output
}

func (f *FanIn) Send(sourceID int, data string) {
	switch sourceID {
	case 1:
		f.input <- data
	case 2:
		f.input2 <- data
	case 3:
		f.input3 <- data
	default:
	}
}

func (f *FanIn) Close() {
	close(f.input)
	close(f.input2)
	close(f.input3)
	close(f.output)
}
