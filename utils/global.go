package utils

import (
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func IsNumeric(s string) bool {
	for _, c := range s {
		if !(c >= 48 && c <= 57) {
			return false
		}
	}
	return true
}

func RemoveAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output
}

func FormatArtistName(artist string) string {
	artist = RemoveAccents(artist)
	artist = strings.Join(strings.Split(artist, " "), "%20")
	return artist
}

func TransformDateToText(dateTime string) string {
	d := strings.Split(dateTime, "-")
	day, _ := strconv.Atoi(d[0])
	month, _ := strconv.Atoi(d[1])
	year, _ := strconv.Atoi(d[2])

	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Now().Location())

	str := date.Weekday().String() + "-" + date.Month().String() + "-" + strconv.Itoa(date.Year())
	return str
}
