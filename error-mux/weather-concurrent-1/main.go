package main

import (
	"fmt"
	"sync"
	"time"
)

type Weather struct {
	Name  string
	TempC float64
}

func getName(zipCode string) (n string, err error) {
	names := map[string]string{
		"19123": "Philadelphia, PA",
		"90210": "Beverly Hills, CA",
	}

	if name, ok := names[zipCode]; ok {
		n = name
	} else {
		err = fmt.Errorf("getName: %d not found", zipCode)
	}

	time.Sleep(time.Second)

	return
}

func getTemp(zipCode string) (t float64, err error) {
	temps := map[string]float64{
		"19123": -5.0,
		"90210": 27.3,
	}

	if temp, ok := temps[zipCode]; ok {
		t = temp
	} else {
		err = fmt.Errorf("getTemp: %d not found", zipCode)
	}

	time.Sleep(time.Second)

	return
}

func getWeather(zipCode string) (w *Weather, err error) {

	var (
		name             string
		tempC            float64
		nameErr, tempErr error
	)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		name, nameErr = getName(zipCode)
		wg.Done()
	}()

	go func() {
		tempC, tempErr = getTemp(zipCode)
		wg.Done()
	}()

	wg.Wait()

	switch {
	case nameErr != nil:
		err = nameErr
	case tempErr != nil:
		err = tempErr
	default:
		w = &Weather{Name: name, TempC: tempC}
	}

	return
}

func main() {
	w, err := getWeather("19123")

	fmt.Println(w, err)
}
