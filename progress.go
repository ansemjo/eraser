// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package main

import (
	"fmt"
	"time"
)

type Spinner struct {
	runes []rune
	spin  int
}

func newSpinner() *Spinner {
	return &Spinner{[]rune{'/', '-', '\\', '|'}, 0}
}

func (s *Spinner) next() rune {
	s.spin++
	s.spin %= len(s.runes)
	return s.runes[s.spin]
}

type Progress struct {
	spin    *Spinner
	avg     *Average
	total   int64
	current int64
	fmt     string
	start   time.Time
	last    time.Time
}

func newProgress(total int64) *Progress {
	return &Progress{
		spin:    newSpinner(),
		avg:     newAverage(),
		total:   total,
		current: 0,
		fmt:     fmt.Sprintf("\r %%c %%%dd / %%d", len(fmt.Sprint(total))),
		start:   time.Now(),
		last:    time.Now(),
	}
}

func (prg *Progress) draw() {
	fmt.Printf(prg.fmt+" (%v, %.2f MiB/s)", prg.spin.next(),
		prg.current, prg.total, time.Since(prg.start).Round(time.Millisecond), prg.avg.current/(1024*1024))
	prg.last = time.Now()
}

func (prg *Progress) done() {
	fmt.Printf(prg.fmt+"\n", '✔', prg.current, prg.total)
}

func (prg *Progress) add(bytes int64) {
	prg.avg.sample(bytes)
	prg.current += bytes
	if time.Since(prg.last) > time.Millisecond*83 {
		prg.draw()
	}
}

type Average struct {
	current float32
	last    time.Time
}

func newAverage() *Average {
	return &Average{0, time.Now()}
}

func (avg *Average) sample(added int64) {
	delta := time.Since(avg.last)
	avg.last = time.Now()
	average := float32(added) / float32(delta.Seconds())
	avg.current = avg.current + (average-avg.current)/200
}
