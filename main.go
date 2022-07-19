package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// model of output file

type RandomIntegers struct {
	Stddev float64 `json:"stddev"`
	Data   []int   `json:"data"`
}

const API_SERVER string = "https://www.random.org/integers/"
const PARAMS string = "min=1&max=100&col=1&base=10&format=plain&rnd=new"

func main() {

	fmt.Println("Random integer API")

	http.HandleFunc("/random/mean", getParameters)

	http.ListenAndServe(":80", nil)

}

func getParameters(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusBadRequest)
		return
	}

	length, requests, err := parametersHandler(r.URL.Query().Get("length"), r.URL.Query().Get("requests"), w)
	if err != nil {
		return
	}

	calculation, err := calculations(requests, length, w)
	if err != nil {
		http.Error(w, "Service internal error", http.StatusInternalServerError)
		return
	}
	convertCalculations, err := json.Marshal(calculation)
	if err != nil {
		http.Error(w, "json encoding error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(convertCalculations))

}

func calculations(request int, length int, w http.ResponseWriter) ([]RandomIntegers, error) {
	var output []RandomIntegers
	var reqUrl = fmt.Sprintf("%s?%s&num=%v", API_SERVER, PARAMS, length)

	for i := 0; i < request; i++ {
		buf := new(bytes.Buffer)
		response, err := http.Get(reqUrl)

		if err != nil {
			fmt.Print(err.Error())
			http.Error(w, err.Error(), http.StatusBadGateway)
			return nil, errors.New(err.Error())
		}
		if response.StatusCode != 200 {
			buf.ReadFrom(response.Body)
			newStr := buf.String()
			http.Error(w, newStr, response.StatusCode)
			return nil, errors.New(newStr)
		} else {
			data, _ := io.ReadAll(response.Body)
			lines := strings.Split(string(data), "\n")
			numbers := []int{}

			for i := range lines {
				line := strings.TrimSpace(lines[i])
				if line == "" {
					continue
				}

				num, err := strconv.Atoi(line)
				if err != nil {
					fmt.Println(err)
					return nil, err
				}

				numbers = append(numbers, num)
			}
			dev := stdDev(numbers)
			resoults := RandomIntegers{dev, numbers}
			output = append(output, resoults)

		}
	}

	sumResoults := sumCalculations(output)
	output = append(output, sumResoults)
	return output, nil
}

func stdDev(numbers []int) float64 {
	var average float64 = 0
	for _, v := range numbers {
		average += float64(v)
	}
	average = average / float64(len(numbers))

	var dev float64 = 0
	for _, v := range numbers {
		dev += ((float64(v) - average) * (float64(v) - average))
	}
	dev = dev / float64(len(numbers))
	dev = math.Sqrt(dev)
	return math.Round(dev*100) / 100
}

func sumCalculations(values []RandomIntegers) RandomIntegers {
	var sum []int
	for _, data := range values {
		sum = append(sum, data.Data...)
	}
	dev := stdDev(sum)
	sumResoults := RandomIntegers{dev, sum}
	return sumResoults
}

func parametersHandler(length string, requests string, w http.ResponseWriter) (int, int, error) {
	if length == "" || requests == "" {
		http.Error(w, "You have to provide 2 parameters: request(int) and length(int)", http.StatusBadRequest)
		return 0, 0, errors.New("error")
	}

	response, err := strconv.Atoi(requests)
	if err != nil || response > 10 || response < 1 {
		http.Error(w, "Request parameter must to be integer in range 1 - 10!", http.StatusBadRequest)
		return 0, 0, errors.New("error")
	}

	len, err := strconv.Atoi(length)
	if err != nil || len > 10000 || len < 1 {
		http.Error(w, "Length parameter must to be integer in range 1 - 10.000!", http.StatusBadRequest)
		return 0, 0, errors.New("error")
	}
	return len, response, nil
}
