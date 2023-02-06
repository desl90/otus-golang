package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

type Pipeline struct {
	in, done In
	stages   []Stage
	out      Bi
}

func NewPipeline(in In, done In, stages []Stage) *Pipeline {
	return &Pipeline{
		in:     in,
		done:   done,
		stages: stages,
		out:    make(Bi),
	}
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	pipeline := NewPipeline(in, done, stages)

	for _, stage := range pipeline.stages {
		pipeline.in = stage(pipeline.in)
	}

	return pipeline.stream()
}

func (p *Pipeline) stream() Out {
	go func() {
		defer close(p.out)

		for {
			select {
			case <-p.done:
				return

			case v, ok := <-p.in:
				if !ok {
					return
				}

				select {
				case <-p.done:
					return

				case p.out <- v:
				}
			}
		}
	}()

	return p.out
}
