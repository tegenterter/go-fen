package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	port = os.Getenv("PORT")
	dimension, _ = strconv.Atoi(os.Getenv("DIMENSION"))
	assetsPath = os.Getenv("ASSETS_PATH")
	lightSquareRGB = strings.Split(os.Getenv("LIGHT_SQUARE_RGB"), ",")
	darkSquareRGB = strings.Split(os.Getenv("DARK_SQUARE_RGB"), ",")
)

func drawDiagram(fen string) image.Image {
	// Parse square colors from configuration
	lr, _ := strconv.Atoi(lightSquareRGB[0])
	lg, _ := strconv.Atoi(lightSquareRGB[1])
	lb, _ := strconv.Atoi(lightSquareRGB[2])
	dr, _ := strconv.Atoi(darkSquareRGB[0])
	dg, _ := strconv.Atoi(darkSquareRGB[1])
	db, _ := strconv.Atoi(darkSquareRGB[2])

	// Draw an empty board
	board := drawBoard(dimension, color.RGBA{R: uint8(lr), G: uint8(lg), B: uint8(lb), A: 255}, color.RGBA{R: uint8(dr), G: uint8(dg), B: uint8(db), A: 255})
	
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
				draw.Draw(board, image.Rect(0, 0, dimension, dimension).Add(image.Point{X: k * dimension, Y: i * dimension}), getPiece(row[j]), image.Point{}, draw.Over)
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
			draw.Draw(board, image.Rect(rank, file, rank + dimension, file + dimension), &image.Uniform{C: colors[c]}, image.Point{}, draw.Src)

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

	image, err := os.Open(assetsPath + "/" + piece + ".png")
	second, err := png.Decode(image)
	defer image.Close()

	if err != nil {
		log.Fatal(err.Error())
	}

	return second
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fen := r.URL.Path[1:]
		rgxp := regexp.MustCompile(`(?i)^[rnbqk1-8]{1,8}\/[rnbqkp1-8]{1,8}\/[rnbqkp1-8]{1,8}\/[rnbqkp1-8]{1,8}\/[rnbqkp1-8]{1,8}\/[rnbqkp1-8]{1,8}\/[rnbqkp1-8]{1,8}\/[rnbqk1-8]{1,8}\s[wb]{1}\s[kq-]{1,4}\s[a-h1-8-]{1,2}\s\d+\s\d+$`)

		if rgxp.MatchString(fen) == false {
			http.Error(w, "Invalid FEN string", http.StatusBadRequest)
			log.Println("GET", http.StatusBadRequest, "/" + fen)
			return
		}

		// Generate the diagram
		buffer := new(bytes.Buffer)
		err := png.Encode(buffer, drawDiagram(fen))

		if err != nil {
			http.Error(w, "Could not generate diagram", http.StatusInternalServerError)
			log.Println("GET", http.StatusBadRequest, "/" + fen)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
		_, err = w.Write(buffer.Bytes())

		if err != nil {
			log.Fatal(err.Error())
		}

		log.Println("GET", http.StatusOK, "/" + fen)
	})

	log.Fatal(http.ListenAndServe(":" + port, nil))
}
