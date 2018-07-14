// Package wkhtmltoimage provides a simple wrapper around wkhtmltoimage (http://wkhtmltopdf.org/) binary.
package htmlToimage

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"image/png"
	"log"
	"os/exec"
	"strconv"
	"strings"

	proto "github.com/zale144/instagram-bot/htmlToimage/proto"
)

// GenerateImage creates an image from an input.
// It returns the image ([]byte) and any error encountered.
func GenerateImage(options *proto.ImageRequest) ([]byte, error) {
	arr, err := buildParams(options)
	if err != nil {
		return []byte{}, err
	}

	cmd := exec.Command("wkhtmltoimage", arr...)

	if options.Html != "" {
		cmd.Stdin = strings.NewReader(options.Html)
	}
	fmt.Println(options.Input)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Println("can't generate file ", err)
		return nil, err
	}
	if options.Output == "" && len(output) > 0 {
		trimmed := cleanupOutput(output, options.Format)
		return trimmed, err
	}
	return output, err
}

// buildParams takes the image options set by the user and turns them into command flags for wkhtmltoimage
// It returns an array of command flags.
func buildParams(options *proto.ImageRequest) ([]string, error) {
	a := []string{}

	if options.Input == "" {
		return []string{}, errors.New("Must provide input")
	}

	// silence extra wkhtmltoimage output
	// might want to add --javascript-delay too?
	a = append(a, "-q")
	a = append(a, "--disable-plugins")

	a = append(a, "--format")
	if options.Format != "" {
		a = append(a, options.Format)
	} else {
		a = append(a, "png")
	}

	if options.Height != 0 {
		a = append(a, "--height")
		a = append(a, strconv.Itoa(int(options.Height)))
	}

	if options.Width != 0 {
		a = append(a, "--width")
		a = append(a, strconv.Itoa(int(options.Width)))
	}

	if options.Zoom != 0 {
		a = append(a, "--zoom")
		a = append(a, strconv.FormatFloat(float64(options.Zoom), 'E', -1, 64))
	}

	if options.Quality != 0 {
		a = append(a, "--quality")
		a = append(a, strconv.Itoa(int(options.Quality)))
	}

	if options.CropX != 0 {
		a = append(a, "--crop-x")
		a = append(a, strconv.Itoa(int(options.CropX)))
	}

	if options.CropY != 0 {
		a = append(a, "--crop-y")
		a = append(a, strconv.Itoa(int(options.CropX)))
	}

	if options.CropW != 0 {
		a = append(a, "--crop-w")
		a = append(a, strconv.Itoa(int(options.CropX)))
	}

	if options.CropH != 0 {
		a = append(a, "--crop-h")
		a = append(a, strconv.Itoa(int(options.CropX)))
	}

	// url and output come last
	if options.Input != "-" {
		// make sure we dont pass stdin if we aren't expecting it
		options.Html = ""
	}

	a = append(a, options.Input)

	if options.Output == "" {
		a = append(a, "-")
	} else {
		a = append(a, options.Output)
	}

	return a, nil
}

func cleanupOutput(img []byte, format string) []byte {
	buf := new(bytes.Buffer)
	switch {
	case format == "png":
		decoded, err := png.Decode(bytes.NewReader(img))
		for err != nil {
			img = img[1:]
			if len(img) == 0 {
				break
			}
			decoded, err = png.Decode(bytes.NewReader(img))
		}
		png.Encode(buf, decoded)
		return buf.Bytes()
	case format == "jpg":
		decoded, err := jpeg.Decode(bytes.NewReader(img))
		for err != nil {
			img = img[1:]
			if len(img) == 0 {
				break
			}
			decoded, err = jpeg.Decode(bytes.NewReader(img))
		}
		jpeg.Encode(buf, decoded, nil)
		return buf.Bytes()
		// case format == "svg":
		// 	return img
	}
	return img
}
