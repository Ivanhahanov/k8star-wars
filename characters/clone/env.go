package main

import (
	"os"
	"strconv"
)

func GetEnv[I int | float64](name string, defaultValue I) I {
	value := os.Getenv(name)
	if value != "" {

		switch any(defaultValue).(type) {
		case int:
			i, err := strconv.Atoi(value)
			if err != nil {
				panic(err)
			}
			return I(i)

		case float64:
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				panic(err)
			}
			return I(f)
		}

	}
	return defaultValue
}
