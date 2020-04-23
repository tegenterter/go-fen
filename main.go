package main

import (
	"bytes"
	"github.com/gorilla/mux"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const dimension = 45

func drawBoard(squareDimension int, lightColor color.RGBA, darkColor color.RGBA) *image.RGBA {
	board := image.NewRGBA(image.Rect(0, 0, squareDimension * 8, squareDimension * 8))

	colors := make(map[int]color.RGBA, 2)
	colors[0] = darkColor
	colors[1] = lightColor

	c := 0

	for i := 0; i < 8; i++ {
		rank := i * dimension

		for j := 0; j < 8; j++ {
			file := j * dimension
			draw.Draw(board, image.Rect(rank, file, rank + dimension, file + dimension), &image.Uniform{colors[c]}, image.Point{}, draw.Src)
			c = 1 - c
		}

		c = 1 - c
	}

	return board
}

func diagram(w http.ResponseWriter, r *http.Request) {
	fen := mux.Vars(r)["fen"]

	board := drawBoard(dimension, color.RGBA{209, 139, 71, 255}, color.RGBA{255, 206, 158, 255})
	rgxp := regexp.MustCompile(`(?i)^[a-z\d+\/]+`)

	for i, row := range strings.Split(rgxp.FindString(fen), "/") {
		k := 0

		for j := 0; j < len(row); j++ {
			character := string(row[j])
			emptySquares, err := strconv.Atoi(character)

			if err == nil {
				k += emptySquares
			} else {
				draw.Draw(board, image.Rect(0, 0, dimension, dimension).Add(image.Point{k * dimension, dimension * i}), getPiece(row[j]), image.Point{}, draw.Over)
				k++
			}
		}
	}

	buffer := new(bytes.Buffer)
	png.Encode(buffer, board)

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	w.Write(buffer.Bytes())
}

func getPiece(character uint8) image.Image {
	var piece string
	if unicode.IsUpper(rune(character)) {
		piece = "w" + string(unicode.ToLower(rune(character)))
	} else {
		piece = "b" + string(character)
	}

	image, err := os.Open("assets/" + piece + ".png")
	second, err := png.Decode(image)
	if err != nil {
		panic(err.Error())
	}
	defer image.Close()

	return second
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/{fen:.*}", diagram)
	http.ListenAndServe(":8080", router)
}
