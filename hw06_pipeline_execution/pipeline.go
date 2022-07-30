package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var out Out
	bi := make(Bi)

	exitStage := func(in In) Out {
		out := make(Bi)

		go func() {
			defer close(out)

			exit := false
			for v := range in {
				select {
				case <-done:
					exit = true
				default:
					out <- v
				}

				if exit {
					return
				}
			}
		}()

		return out
	}

	for index, stage := range stages {
		if index == 0 {
			out = stage(bi)
		} else {
			out = exitStage(stage(out))
		}
	}

	go func() {
		for {
			select {
			case <-done:
				return
			case value, ok := <-in:
				if ok {
					bi <- value
				} else {
					close(bi)
					return
				}
			}
		}
	}()

	return out
}
