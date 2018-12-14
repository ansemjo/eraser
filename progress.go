package main

import (
	"fmt"
)

var spinRunes = []rune{'/', '-', '\\', '|'}

type progressSpinner struct {
	total   int64
	current int64
	fmt     string
	spin    int
}

func newProgressSpinner(total int64) (sp *progressSpinner) {
	format := fmt.Sprintf("\r %%c %%%dd / %%d ", len(fmt.Sprint(total)))
	return &progressSpinner{total, 0, format, 0}
}

func (sp *progressSpinner) draw() {
	fmt.Printf(sp.fmt, spinRunes[sp.spin], sp.current, sp.total)
	sp.spin = (sp.spin + 1) % len(spinRunes)
}

func (sp *progressSpinner) add(b int64) {
	sp.current += b
	sp.draw()
}

func (sp *progressSpinner) done() {
	fmt.Printf(sp.fmt+"\n", '✔', sp.current, sp.total)
}
