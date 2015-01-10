package main

import (
	"fmt"
	"sync"
	"time"
)

func mergeErrors(chs ...<-chan error) <-chan error {
	out := make(chan error)
	wg := sync.WaitGroup{}

	for _, ch := range chs {
		wg.Add(1)
		go func(c <-chan error) {
			for e := range c {
				out <- e
			}
			wg.Done()
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

type Weather struct {
	Name  string
	TempC float64
}

func getName(zipCode string) (<-chan string, <-chan error) {
	names := map[string]string{
		"19123": "Philadelphia, PA",
		"90210": "Beverly Hills, CA",
	}

	out := make(chan string, 1)
	errs := make(chan error, 1)

	go func() {
		time.Sleep(time.Second)
		if name, ok := names[zipCode]; ok {
			out <- name
		} else {
			errs <- fmt.Errorf("getName: %d not found", zipCode)
		}

		close(out)
		close(errs)
	}()

	return out, errs
}

func getTemp(zipCode string) (<-chan float64, <-chan error) {
	temps := map[string]float64{
		"19123": -5.0,
		"90210": 27.3,
	}

	out := make(chan float64, 1)
	errs := make(chan error, 1)

	go func() {
		time.Sleep(time.Second)
		if temp, ok := temps[zipCode]; ok {
			out <- temp
		} else {
			errs <- fmt.Errorf("getTemp: %d not found", zipCode)
		}

		close(out)
		close(errs)
	}()

	return out, errs
}

func getWeather(zipCode string) (w *Weather, err error) {

	nameOut, nameErr := getName(zipCode)
	tempOut, tempErr := getTemp(zipCode)

	errs := mergeErrors(nameErr, tempErr)

	for e := range errs {
		if err == nil {
			err = e
		}
	}

	if err != nil {
		return
	}

	w = &Weather{Name: <-nameOut, TempC: <-tempOut}

	return
}

func main() {
	w, err := getWeather("19123")

	fmt.Println(w, err)
}
