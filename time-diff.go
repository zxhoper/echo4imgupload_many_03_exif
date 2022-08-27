package main

import (
	"fmt"
	// -"os"
	"regexp"
	"strconv"
	//"strings"
	"time"
)

type dateTime struct {
	Year, Month, Date, Hour, Minute int
}

func timeDiff(dura string) (tdiff string, mustswitch bool) {

	//argsWithProg := os.Args
	//argsWithoutProg := os.Args[1:]

	//fmt.Println(argsWithProg)
	//fmt.Println(argsWithoutProg)

	//for i, v := range os.Args[1:] {
	//	fmt.Println(i, v)
	//}
	fmt.Println("input string", dura)

	// if len(os.Args) == 1 {
	// 	fmt.Println("No Argument! ")
	// 	fmt.Println("Usage: (to calculate time difference)")
	// 	fmt.Println("      ./goGetArguSplitString \"<2022-07-03 Sun 11:33>--<2022-07-03 Sun 11:50>\"")
	// 	fmt.Println("      ./goGetArguSplitString \"2022-07-03 Sun 11:33\" \"2022-07-03 Sun 11:50\"")
	// 	os.Exit(1)
	// }

	//input := strings.Join(os.Args[1:], " ")
	input := dura
	//	fmt.Println(input)

	isMatched, _ := regexp.MatchString(`.*[0-9]{4}-[0-9]{2}-[0-9]{2} [a-zA-Z]{3} [0-9]{2}:[0-9]{2}[ <>-]{1,4}[0-9]{4}-[0-9]{2}-[0-9]{2} [a-zA-Z]{3} [0-9]{2}:[0-9]{2}`, input)
	//	2022-07-03 Sun 11:33 2022-07-03 Sun 11:50
	if !isMatched {
		fmt.Println("Input String is NOT OK !!!!!!!!!!!!!!!!!!!!!!!! ")
		//os.Exit(1)
		input = "<1978-01-05 Sun 11:33>--<2022-08-20 Sun 21:50>"
	}

	//fmt.Println("Input String is good! ")

	rexp := regexp.MustCompile(`.*([0-9]{4})-([0-9]{2})-([0-9]{2}) [a-zA-Z]{3} ([0-9]{2}):([0-9]{2})[ <>-]{1,4}([0-9]{4})-([0-9]{2})-([0-9]{2}) [a-zA-Z]{3} ([0-9]{2}):([0-9]{2})`)
	result := rexp.FindAllStringSubmatch(string(input), -1)

	var o, n dateTime
	for _, m := range result {
		//		fmt.Printf("Nth:%d  Y=%v M=%v D=%v %v:%v  to Y=%v M=%v D=%v %v:%v  \n", i+1, m[1], m[2], m[3], m[4], m[5], m[6], m[7], m[8], m[9], m[10])

		o.Year, _ = strconv.Atoi(m[1])
		o.Month, _ = strconv.Atoi(m[2])
		o.Date, _ = strconv.Atoi(m[3])
		o.Hour, _ = strconv.Atoi(m[4])
		o.Minute, _ = strconv.Atoi(m[5])

		n.Year, _ = strconv.Atoi(m[6])
		n.Month, _ = strconv.Atoi(m[7])
		n.Date, _ = strconv.Atoi(m[8])
		n.Hour, _ = strconv.Atoi(m[9])
		n.Minute, _ = strconv.Atoi(m[10])

	}

	//fmt.Println(o)
	//fmt.Println(n)

	oldDate := time.Date(o.Year, time.Month(o.Month), o.Date, o.Hour, o.Minute, 0, 0, time.UTC)
	newDate := time.Date(n.Year, time.Month(n.Month), n.Date, n.Hour, n.Minute, 0, 0, time.UTC)

	// Using time.Before() method
	//    g1 := today.Before(tomorrow)
	//fmt.Println("today before tomorrow:", g1)

	if newDate.Before(oldDate) {
		mustswitch = true
		tmpDate := newDate
		newDate = oldDate
		oldDate = tmpDate
	}

	//newDate := time.Date(2022, 4, 13, 1, 0, 0, 0, time.UTC)
	difference := newDate.Sub(oldDate)

	// fmt.Printf("Years: %d\n", int64(difference.Hours()/24/365))
	// fmt.Printf("Months: %d\n", int64(difference.Hours()/24/30))
	// fmt.Printf("Weeks: %d\n", int64(difference.Hours()/24/7))
	// fmt.Printf("Days: %d\n", int64(difference.Hours()/24))
	// fmt.Printf("Hours: %.f\n", difference.Hours())
	// fmt.Printf("Minutes: %.f\n", difference.Minutes())
	// fmt.Printf("Seconds: %.f\n", difference.Seconds())
	// fmt.Printf("Milliseconds: %d\n", difference.Milliseconds())
	// fmt.Printf("Microseconds: %d\n", difference.Microseconds())
	// fmt.Printf("Nanoseconds: %d\n", difference.Nanoseconds())

	tdiff = fmt.Sprintf("%v", FormatSince(difference))
	return

}

//https://stackoverflow.com/questions/42391869/time-since-in-days-hours-minutes-seconds-format
func FormatSince(ts time.Duration) string {
	const (
		Decisecond = 100 * time.Millisecond
		Day        = 24 * time.Hour
		Year       = 365 * Day
	)
	//ts := time.Since(t)
	y := ts / Year
	ts = ts % Year
	d := ts / Day
	ts = ts % Day
	h := ts / time.Hour
	ts = ts % time.Hour
	m := ts / time.Minute
	ts = ts % time.Minute
	//s := ts / time.Second
	//ts = ts % time.Second
	//f := ts / Decisecond
	//return fmt.Sprintf("%dd%dh%dm%d.%ds", d, h, m, s, f)
	rtn := "return string"
	if y != 0 {
		rtn = fmt.Sprintf("%d Year, ", y)
	} else {
		rtn = ""
	}

	if d != 0 {
		rtn += fmt.Sprintf("%d Day, ", d)
	}
	if h != 0 {
		if h == 1 {
			rtn += fmt.Sprintf("%d Hour, ", h)
		} else {
			rtn += fmt.Sprintf("%d Hours, ", h)
		}
	}
	if m == 0 || m == 1 {
		rtn += fmt.Sprintf("%d Minute", m)
	} else {
		rtn += fmt.Sprintf("%d Minutes", m)
	}

	//	return fmt.Sprintf("%dy%dd%dh%dm", y, d, h, m)
	return rtn
}
