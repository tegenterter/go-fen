package main

import (
	"bytes"
	"fmt"
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

func drawDiagram(fen string) image.Image {
	// Draw an empty board
	board := drawBoard(dimension, color.RGBA{209, 139, 71, 255}, color.RGBA{255, 206, 158, 255})

	rgxp := regexp.MustCompile(`(?i)^[a-z\d+\/]+`)

	// Iterate the separate ranks as defined in the FEN string
	for i, row := range strings.Split(rgxp.FindString(fen), "/") {
		k := 0

		for j := 0; j < len(row); j++ {
			character := string(row[j])
			// Check if the character is numeric
			emptySquares, err := strconv.Atoi(character)

			if err == nil {
				// Pad the number of empty squares
				k += emptySquares
			} else {
				// Draw the piece onto the board
				draw.Draw(board, image.Rect(0, 0, dimension, dimension).Add(image.Point{k * dimension, i * dimension}), getPiece(row[j]), image.Point{}, draw.Over)
				k++
			}
		}
	}

	return board
}

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

			// Draw an individual square
			draw.Draw(board, image.Rect(rank, file, rank + dimension, file + dimension), &image.Uniform{colors[c]}, image.Point{}, draw.Src)

			c = 1 - c
		}

		c = 1 - c
	}

	return board
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fen := r.URL.Path[1:]
		rgxp := regexp.MustCompile(`(?i)^[rnbqk1-7]{1,8}\/[rnbqkp1-7]{1,8}\/[rnbqkp1-7]{1,8}\/[rnbqkp1-7]{1,8}\/[rnbqkp1-7]{1,8}\/[rnbqkp1-7]{1,8}\/[rnbqkp1-7]{1,8}\/[rnbqk1-7]{1,8}\sw|b\sk|q|-\s\d+\s\d+$`)

		if rgxp.MatchString(fen) == false {
			http.Error(w, "Invalid FEN string", http.StatusBadRequest)

			fmt.Println("GET", http.StatusBadRequest, "/" + fen)
			return
		}

		buffer := new(bytes.Buffer)
		// Generate the diagram
		png.Encode(buffer, drawDiagram(fen))

		fmt.Println("GET", http.StatusOK, "/" + fen)

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
		w.Write(buffer.Bytes())
	})

	http.ListenAndServe(":8080", nil)
}
