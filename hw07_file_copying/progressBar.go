package main

import (
	"fmt"
	"strings"
	"sync"
)

type ProgressBar interface {
	Add(value int64)
	Finish()
}

type progressBar struct {
	sync.Mutex
	maxValue     int64
	currentValue int64
	width        int
	finish       bool
}

func StartNewProgressBar(count int64) ProgressBar {
	return &progressBar{
		maxValue: count,
		width:    15,
	}
}

func (pb *progressBar) Add(value int64) {
	pb.Lock()
	defer pb.Unlock()

	if pb.finish {
		return
	}

	pb.currentValue += value

	if pb.currentValue >= pb.maxValue {
		pb.finish = true
		pb.currentValue = pb.maxValue
		fmt.Println("\r", pb.createOutput())
	} else {
		fmt.Print("\r", pb.createOutput())
	}
}

func (pb *progressBar) Finish() {
	pb.Lock()
	defer pb.Unlock()

	if pb.finish {
		return
	}

	pb.finish = true
	pb.currentValue = pb.maxValue

	fmt.Println("\r", pb.createOutput())
}

func (pb *progressBar) createOutput() string {
	counter := fmt.Sprint(pb.currentValue, "/", pb.maxValue)

	currentWidth := int(float64(pb.currentValue) / float64(pb.maxValue) * float64(pb.width))
	bar := fmt.Sprint("[", strings.Repeat("-", currentWidth), strings.Repeat("_", pb.width-currentWidth), "]")

	percent := int(float64(pb.currentValue) / float64(pb.maxValue) * 100)

	return fmt.Sprint("\r", counter, " ", bar, " ", percent, "%")
}
