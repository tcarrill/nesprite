# nesprite
A simple NES sprite tile viewer written in [Go](https://go.dev/).

### Building
`go build nesprite.go`
### Running
`./nesprite <path to .nes file>`

The result will be one or more .png images, depending on the size of the 
CHR ROM in the .nes file.  Each image will be 128x256 and contain 512 
sprite tiles.  Below is a sample image from Super Maio Bros.

![Super Mario Bro. Sprite Tiles](Super%20Mario%20Bros-0.png)
