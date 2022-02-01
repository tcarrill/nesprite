package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	rom, err :=  RetrieveROM(os.Args[1])

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	prgRomSize := uint(rom[0x0004]) * 16384
	chrRomSize := uint(rom[0x0005]) * 8192

	if chrRomSize == 0 {
		fmt.Printf("CHR RAM currently unsupported!\n")
		os.Exit(1)
	}
	numTiles := int(chrRomSize / 16)
	pages := numTiles / 512 // 512 tiles per page

	chrRom := rom[prgRomSize + 16 : prgRomSize + chrRomSize + 16]
	tiles := ConvertChrRom(chrRom)

	fmt.Printf("%s\n", os.Args[1])
	fmt.Printf("PRG ROM Size: %d\n", prgRomSize)
	fmt.Printf("CHR ROM Size: %d\n", chrRomSize)
	fmt.Printf("Number of tiles: %d\n", numTiles)
	fmt.Printf("Number of pages: %d\n", pages)
	fmt.Printf("len(tiles): %d\n", len(tiles))

	for i := 0; i < pages; i++ {
		filename := fmt.Sprintf("%s-%d.png", strings.TrimSuffix(os.Args[1], filepath.Ext(os.Args[1])), i)
		start := i * 32768 // 32768 pixels per page
		end := start + 32768
		CreatePNG(tiles[start:end],  filename)
	}
}

func DrawTile(tile []byte, img *image.RGBA, x int, y int) {
	ox := x
	white := color.RGBA{255, 255,255,0xff}
	grey1 := color.RGBA{128, 128, 128, 0xff}
	grey2 := color.RGBA{200, 200, 200, 0xff}
	grey3 := color.RGBA{160, 160, 160, 0xff}

	for i := 0; i < len(tile); i++ {
		pixColor := white
		if tile[i] == 1 {
			pixColor = grey1
		} else if tile[i] == 2 {
			pixColor = grey2
		} else if tile[i] == 3 {
			pixColor = grey3
		}

		img.Set(x, y, pixColor)
		x += 1
		if x % 8 == 0 {
			x = ox
			y += 1
		}
	}
}

func CreatePNG(sprites []byte, filename string) {
	fmt.Printf("Creating %s\n", filename)
	width :=  128
	height := 256

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	x := 0
	y := 0
	for i := 0; i < len(sprites); i++ {
		start := i * 64
		end := start + 64
		if end > len(sprites) {
			break
		}

		DrawTile(sprites[start:end], img, x, y)

		x += 8
		if x >= width {
			x = 0
			y += 8
		}
	}

	f, _ := os.Create(filename)
	png.Encode(f, img)
}

func ConvertChrRom(spriteBytes []byte) []byte {
	pixels := make([]byte, len(spriteBytes) * 4)
	index := 0
	mask := byte(0x1)

	for i := 0; i < len(spriteBytes); i += 16 {
		for j := 0; j < 8; j++ {
			for k := 7; k >= 0; k-- {
				channel1 := spriteBytes[i + j]
				channel2 := spriteBytes[i + j + 8]
				pixels[index] = ((channel1 >> byte(k)) & mask) + (((channel2 >> byte(k)) & mask) << 1)
				index++
			}
		}
	}

	return pixels
}

func RetrieveROM(filename string) ([]byte, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	var size = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_,err = bufr.Read(bytes)

	return bytes, err
}