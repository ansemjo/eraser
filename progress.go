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
	// spinner to print an updating rune
	spin *Spinner
	// avg is used to calculate average speed
	avg *Average
	// total size of disk
	total int64
	// current amount of written bytes
	current int64
	// format string for printing
	fmt string
	// start time of operation
	start time.Time
	// last drawing time
	last time.Time
}

func newProgress(total int64) *Progress {
	return &Progress{
		spin:    newSpinner(),
		avg:     newAverage(),
		total:   total,
		current: 0,
		fmt:     fmt.Sprintf("\033[2K\r %%c %%%dd / %%d", len(fmt.Sprint(total))),
		start:   time.Now(),
		last:    time.Now(),
	}
}

func (prg *Progress) estimate() time.Duration {
	if int64(time.Since(prg.start)) == 0 || prg.current == 0 {
		return time.Duration(0)
	}
	speed := float64(prg.current) / time.Since(prg.start).Seconds()
	return time.Duration(float64(prg.total-prg.current) / speed)
}

func (prg *Progress) draw() {

	fmt.Printf(prg.fmt+" (%.1f%%, %v elapsed, %.2f MiB/s, ETA %v)",
		prg.spin.next(), prg.current, prg.total, // spinner rune, bytes progress
		float64(prg.current)/float64(prg.total)*100, // progress percentage
		time.Since(prg.start).Round(time.Second),    // elapsed time
		prg.avg.current/(1024*1024),                 // current speed
		prg.estimate()*time.Second,                  // estimated time left
	)
	prg.last = time.Now()
}

func (prg *Progress) syncing() {
	fmt.Printf("\033[2K\r %c %d bytes written, syncing to disk ...", '⦁', prg.total)
}

func (prg *Progress) done() {
	fmt.Printf("\033[2K\r %c %d bytes written in %v\n", '✔', prg.total, time.Since(prg.start).Round(time.Second))
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
	avg.current = avg.current + (average-avg.current)/500
}
