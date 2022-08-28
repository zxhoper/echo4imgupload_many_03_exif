// App for upload 2 image and calculate time diff
//    then write to file
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/tidwall/gjson"
)

func upload(c echo.Context) error {

	var err error
	var imgFile *os.File
	var metaData *exif.Exif
	var jsonByte []byte
	var jsonString string

	// Read form fields
	cats := c.FormValue("category")
	title := c.FormValue("title")
	tags := c.FormValue("tags")
	note := c.FormValue("note")

	text := fmt.Sprintf("- %s\n", cats)

	resStr := "<p><pre>"
	resStr += "Upload report:\n"
	resStr += fmt.Sprintf(" -> Category:%s\n", cats)
	title = strings.TrimSpace(title)
	title = strings.Replace(title, " ", "-", -1)
	resStr += fmt.Sprintf(" -> Title:%s\n", title)

	tagsori := strings.TrimSpace(tags)
	tagStringHyphen := strings.Replace(tagsori, " ", "-", -1)
	tagStringComma := strings.Replace(tagsori, " ", ",", -1)

	text += fmt.Sprintf("- %s\n", tagStringComma)

	resStr += fmt.Sprintf(" -> tagString:%s\n", tagStringHyphen)

	resStr += fmt.Sprintf(" -> Notes:%s\n\n\n", note)

	//time now
	//init the loc
	loc, _ := time.LoadLocation("Asia/Taipei")

	//set timezone,
	now := time.Now().In(loc)
	nowTimeStamp := fmt.Sprintf("%d-%02d-%02da%02d%02d-%s",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Weekday().String()[0:3])
	titleTimeStamp := fmt.Sprintf("%d-%02d-%02d %02d:%02d %s",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Weekday().String()[0:3])

	//-----------
	// Read file
	//-----------

	// Source
	// -Single file, err := c.FormFile("file")
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	var dow string
	var sYYYY string
	var sMM string
	var sDD string
	dstFolder := "up/"
	newFileName := "new"

	text += fmt.Sprintf("- IMAGES:\n")

	for i, file := range files {
		newFileName = ""
		log.Printf("\n\nPicture index: %d\n", i)

		tryGetEXIF := true

		// get file extension
		imgExt := filepath.Ext(file.Filename)

		src, err := file.Open()
		if err != nil {
			return err
		}

		dst, err := os.Create(dstFolder + file.Filename)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		// FF -- Parse-EXIF ------------------------ ___--\\

		// imgFile, err = os.Open("sample.jpg")

		imgFile, err = os.Open(dst.Name())
		if err != nil {
			log.Println("os.Open imgFile faild: ", err.Error())
		}

		metaData, err = exif.Decode(imgFile)
		if err != nil {
			tryGetEXIF = false
			log.Println("exif.Decode metaData faild: ", err.Error())
		}

		if tryGetEXIF {
			jsonByte, err = metaData.MarshalJSON()
			if err != nil {
				log.Println("jsonByte metaData err:", err.Error())
			}

			jsonString = string(jsonByte)
			fmt.Println(jsonString)

			//"DateTimeOriginal":"2022:08:20 15:25:55"

			// fmt.Println("Make: " + gjson.Get(jsonString, "Make").String())
			// fmt.Println("Model: " + gjson.Get(jsonString, "Model").String())
			// fmt.Println("Software: " + gjson.Get(jsonString, "Software").String())
			// fmt.Println("DateTimeOriginal: " + gjson.Get(jsonString, "DateTimeOriginal").String())
			dtStr := gjson.Get(jsonString, "DateTimeOriginal").String()
			fmt.Println("DateTimeOriginal: " + dtStr)
			// => time: 2022:08:20 15:25:55
			resStr += fmt.Sprintf("       => time: %s \n", dtStr)

			// LL __ Parse-EXIF ________________________ ___--//
			//Date Time Map
			dtm := parseDTString(dtStr)

			//<--MAP-->
			// ==> YYYY 2022
			// ==> MM 08
			// ==> DD 20
			// ==> hh 15
			// ==> mm 25
			// ==> ss 55

			fmt.Println("<--MAP-->")
			for k, v := range dtm {
				fmt.Println("==>", k, v)
			}

			sYYYY = dtm["YYYY"]
			sMM = dtm["MM"]
			sDD = dtm["DD"]

			iYYYY, _ := strconv.Atoi(sYYYY)
			iMM, _ := strconv.Atoi(sMM)
			iDD, _ := strconv.Atoi(sDD)
			dow = day_of_week(iDD, iMM, iYYYY)
			// <2022-08-20 Sat 17:46>

			orgTimeStamp := fmt.Sprintf("<%s-%s-%s %s %s:%s>",
				dtm["YYYY"], dtm["MM"], dtm["DD"], dow,
				dtm["hh"], dtm["mm"])
			resStr += orgTimeStamp

			resStr += fmt.Sprintf("imgExt:=========>%s", imgExt)

			newFileName = fmt.Sprintf("%s-%s-%sa%s%s-%s__%s-%03d_%s_%s%s",
				dtm["YYYY"], dtm["MM"], dtm["DD"], dtm["hh"], dtm["mm"], dow,
				title, i+1, cats, tagStringHyphen, imgExt)

			resStr += fmt.Sprintf("newFileName:====>%s", newFileName)
		}
		src.Close()

		//time.Sleep(1 * time.Second)

		if newFileName == "" {
			newFileName = fmt.Sprintf("%s__%s-%03d_%s_%s%s",
				nowTimeStamp, title, i+1,
				cats, tagStringHyphen, imgExt)
		}

		// move to new file
		err = os.Rename(dstFolder+file.Filename, dstFolder+newFileName)
		if err != nil {
			log.Println("os.Rename error")
			log.Fatal(err)
		}

		text += fmt.Sprintf("  + [[./img/%s][Pic%02d]]\n", dstFolder+newFileName, i+1)
		resStr += fmt.Sprintf("\n\nFile %d: %s upload OK!\n", i, dstFolder+newFileName)

	}

	resStr += "\n\n"

	resStr += "</pre></p>"
	resStr += "<p>"

	// FF -- Write-to-File ------------------------ ___--\\

	filename := "public/tu.org"

	text += fmt.Sprintf("- Note:\n#+BEGIN_VERSE\n%s\n#+END_VERSE\n\n", note)

	FirstLine := fmt.Sprintf("\n** %s: %s\n",
		titleTimeStamp, title)
	textWithFirstLine := FirstLine + text

	//text += fmt.Sprintf("%s %s\n", duraStr, timeUsed)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(textWithFirstLine); err != nil {
		panic(err)
	}
	// LL __ Write-to-File ________________________ ___--//

	//resStr += fmt.Sprintf("<h1><a href='/tu.org'>See The Record.</a></h1>\n\n")
	resStr += fmt.Sprintf("<h1><a href='/'>Upload again.</a></h1>\n\n")

	resStr += fmt.Sprintf("\n\n 2022-08-20 is %s <br /> \n\n", day_of_week(24, 8, 2022))

	resStr += "</p>"

	fmt.Printf("title: %s\n", title)

	fmt.Println("==============================================================================================")
	fmt.Println("==============================================================================================")
	fmt.Println("")
	fmt.Println("")

	return c.HTML(http.StatusOK, fmt.Sprintf(resStr+"<p>Uploaded total  %d files with fields Title=%s and title=%s.</p>", len(files), cats, title))

}

func parseDTString(dt string) map[string]string {
	dtmap := make(map[string]string)

	fmt.Println("in parseDTString dt:", dt)

	// FF -- RexGet-YYYY-MM-DD ------------------------ ___--\\
	// regexp get YYYY MM DD  hh mm
	// => time: 2022:08:20 15:25:55
	//rex := regexp.MustCompile(`([0-9]{4}):([0-9]{2}):([0-9]{2}) ([0-9]{2}):([0-9]{2}):([0-9]{2})`)
	rex := regexp.MustCompile(`([0-9]{4}):([0-9]{2}):([0-9]{2})\s([0-9]{2}):([0-9]{2}):([0-9]{2})`)
	rslt := rex.FindAllStringSubmatch(dt, -1)
	for i, m := range rslt {
		fmt.Println(i, "--")
		fmt.Printf(" YYYY: %s\n", m[1])
		fmt.Printf("   MM: %s\n", m[2])
		fmt.Printf("   DD: %s\n", m[3])
		fmt.Printf("   hh: %s\n", m[4])
		fmt.Printf("   mm: %s\n", m[5])
		fmt.Printf("   ss: %s\n", m[6])

		dtmap["YYYY"] = m[1]
		dtmap["MM"] = m[2]
		dtmap["DD"] = m[3]
		dtmap["hh"] = m[4]
		dtmap["mm"] = m[5]
		dtmap["ss"] = m[6]

	}
	return dtmap
	// LL __ RexGet-YYYY-MM-DD ________________________ ___--//

}

// modify from RFC3339 Date and Time on the Internet: Timestamps
// https://www.rfc-editor.org/rfc/rfc3339
// The following is a sample C subroutine loosely based on Zeller's
//   Congruence [Zeller] which may be used to obtain the day of the week
//   for dates on or after 0000-03-01:(0000-03-1 is Wednesday)
func day_of_week(day int, month int, year int) string {
	var cent int
	dayofweek := []string{
		"Sun",
		"Mon",
		"Tue",
		"Wed",
		"Thu",
		"Fri",
		"Sat",
	}
	// dayofweek := []string{
	// 	"Sunday",
	// 	"Monday",
	// 	"Tuesday",
	// 	"Wednesday",
	// 	"Thursday",
	// 	"Friday",
	// 	"Saturday",
	// }

	// adjust months so February is the last one
	month -= 2
	if month < 1 {
		month += 12
		year -= 1
	}
	// split by century
	cent = year / 100
	year %= 100

	//(  (26*month-2)/10 +  day + year<365%7=1> +
	//    year/4<leap year> + cent/4   +  5*cent )  %7

	return (dayofweek[((26*month-2)/10+day+year+year/4+cent/4+5*cent)%7])
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "public")
	e.POST("/upload", upload)

	e.Logger.Fatal(e.Start(":8715"))
}
