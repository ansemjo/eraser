package main

import (
	"fmt"
	"time"
)

var spinRunes = []rune{'/', '-', '\\', '|'}

const megabyte float32 = 1024 * 1024

type progressSpinner struct {
	total   int64
	current int64
	fmt     string
	spin    int
	last    time.Time
	avg     *welfordAverage
}

func newProgressSpinner(total int64) *progressSpinner {
	format := fmt.Sprintf("\r %%c %%%d.2f / %%.2f MiB", len(fmt.Sprintf("%.2f", float32(total)/megabyte)))
	return &progressSpinner{
		total:   total,
		current: 0,
		fmt:     format,
		spin:    0,
		last:    time.Now(),
		avg:     newWelfordAverage(200),
	}
}

func (sp *progressSpinner) draw() {
	fmt.Printf(sp.fmt+" (%.2f MiB/s)",
		spinRunes[sp.spin], float32(sp.current)/megabyte, float32(sp.total)/megabyte, sp.avg.avg/megabyte)
	sp.spin = (sp.spin + 1) % len(spinRunes)
	sp.last = time.Now()
}

func (sp *progressSpinner) add(b int64) {
	sp.avg.sample(b)
	sp.current += b
	if time.Since(sp.last) > 80*time.Millisecond {
		sp.draw()
	}
}

func (sp *progressSpinner) done() {
	fmt.Printf(sp.fmt+"\n", '✔', float32(sp.current)/megabyte, float32(sp.total)/megabyte)
}

type welfordAverage struct {
	avg     float32
	samples float32
	last    time.Time
}

func newWelfordAverage(samples float32) *welfordAverage {
	return &welfordAverage{0, samples, time.Now()}
}

func (w *welfordAverage) sample(copied int64) {
	delta := time.Since(w.last)
	w.last = time.Now()
	sample := float32(copied) / float32(delta.Seconds())
	w.avg = w.avg + (sample-w.avg)/w.samples
}
