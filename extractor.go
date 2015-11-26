/*
ColorEx extracts dominant color palette from an image
*/
package colorex

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"sort"

	"io"

	"image/color"

	"errors"

	gwc "github.com/jyotiska/go-webcolors"
	"github.com/nfnt/resize"
)

type Result struct {
	Hex   string
	Match int
}

type rgbTriplet struct {
	r, g, b uint32
}

type hexWeight struct {
	hex    string
	weight int
}

type byWeight []hexWeight

func (bw byWeight) Len() int           { return len(bw) }
func (bw byWeight) Swap(i, j int)      { bw[i], bw[j] = bw[j], bw[i] }
func (bw byWeight) Less(i, j int) bool { return bw[i].weight > bw[j].weight }

// Default limit = 5
// Default palette = HTML4 palette
func ExtractColors(imgReader io.Reader, limit int, hexPalette []string) ([]Result, error) {
	if limit == 0 {
		limit = 5 // Default limit
	}
	if len(hexPalette) == 0 {
		hexPalette = HTML4Palette // Default palette
	}
	image, _, err := image.Decode(imgReader)
	if err != nil {
		return nil, err
	}

	// Resize the image
	image = resize.Resize(100, 0, image, resize.NearestNeighbor)
	bounds := image.Bounds()
	totalPixels := bounds.Max.X * bounds.Max.Y

	// Prepare palette
	paletteMap := make(map[string]rgbTriplet)
	for _, hexcode := range hexPalette {
		triplet := gwc.HexToRGB(hexcode)
		paletteMap[hexcode] = rgbTriplet{uint32(triplet[0]), uint32(triplet[1]), uint32(triplet[2])}
	}

	weightMap := make(map[string]hexWeight)

	var pixel color.Color
	var red, green, blue uint32
	var minDist uint32
	var curHex string
	var dist uint32

	for i := 0; i <= bounds.Max.X; i++ {
		for j := 0; j <= bounds.Max.Y; j++ {
			pixel = image.At(i, j)
			red, green, blue, _ = pixel.RGBA()
			red /= 255
			green /= 255
			blue /= 255

			minDist = 0
			curHex = ""
			for hex, triplet := range paletteMap {
				dist = distance(red, green, blue, triplet.r, triplet.g, triplet.b)
				if dist < minDist || minDist == 0 {
					minDist = dist
					curHex = hex
				}
			}
			_, exists := weightMap[curHex]
			if exists {
				hw := weightMap[curHex]
				hw.weight++
				weightMap[curHex] = hw
			} else {
				weightMap[curHex] = hexWeight{hex: curHex, weight: 1}
			}
		}
	}

	weights := make([]hexWeight, 0, len(weightMap))
	for _, w := range weightMap {
		weights = append(weights, w)
	}

	if len(weights) == 0 {
		return nil, errors.New("No result")
	}

	sort.Sort(byWeight(weights))

	// Compute match in percents
	res := make([]Result, 0, len(weights))
	for i, w := range weights {
		if i >= limit {
			break
		}
		r := Result{
			Hex:   w.hex,
			Match: w.weight * 100 / totalPixels,
		}
		if r.Match > 0 {
			res = append(res, r)
		}
	}
	return res, nil
}

// Relative distance between two points in 3D space
func distance(x1, y1, z1, x2, y2, z2 uint32) uint32 {
	xd := x2 - x1
	yd := y2 - y1
	zd := z2 - z1
	return xd*xd + yd*yd + zd*zd
}
