// Modified in 2026 from the original old-man-yells-at for whack-a purposes.
// Copyright 2021 oncilla
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package whacka

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"

	xdraw "golang.org/x/image/draw"
)

const (
	canvasSize       = 128
	emojiWidth       = 45
	leftWindowHeight = emojiWidth + 40
	leftWindowBottom = 112

	defaultLeftHiddenPct   = 20
	defaultRightExposedPct = 20
)

var (
	//go:embed fig/whack-a-blank.png
	blankPNG []byte

	blank image.Image = func() image.Image {
		m, err := png.Decode(bytes.NewReader(blankPNG))
		if err != nil {
			panic(err)
		}
		return m
	}()
)

// Options controls how much of each emoji appearance is visible.
type Options struct {
	// LeftHiddenPct is the percent of the lower-left emoji hidden behind its hole.
	LeftHiddenPct int
	// RightExposedPct is the percent of the right emoji exposed above its hole.
	RightExposedPct int
}

// DefaultOptions returns the standard 20% left-hidden and right-exposed layout.
func DefaultOptions() Options {
	return Options{
		LeftHiddenPct:   defaultLeftHiddenPct,
		RightExposedPct: defaultRightExposedPct,
	}
}

// WhackA creates a 128x128 whack-a-mole emoji using the standard layout.
func WhackA(target image.Image) image.Image {
	result, err := WhackAWithOptions(target, DefaultOptions())
	if err != nil {
		panic(err)
	}
	return result
}

// WhackAWithOptions creates a 128x128 whack-a-mole emoji containing target.
// The target is scaled to 45 pixels wide with alpha-aware antialiasing.
func WhackAWithOptions(target image.Image, options Options) (image.Image, error) {
	if err := options.validate(); err != nil {
		return nil, err
	}
	result := image.NewRGBA(image.Rect(0, 0, canvasSize, canvasSize))
	draw.Draw(result, result.Bounds(), blank, blank.Bounds().Min, draw.Src)

	scaled := scaleEmoji(target)
	source := scaled.Bounds().Min

	leftVisibleHeight := scaled.Bounds().Dy() - percentOf(options.LeftHiddenPct, scaled.Bounds().Dy())
	if leftVisibleHeight > leftWindowHeight {
		leftVisibleHeight = leftWindowHeight
	}
	leftTop := leftWindowBottom - leftVisibleHeight
	draw.Draw(result, image.Rect(9, leftTop, 54, leftWindowBottom), scaled, source, draw.Over)

	rightVisibleHeight := percentOf(options.RightExposedPct, scaled.Bounds().Dy())
	if rightVisibleHeight > leftWindowHeight {
		rightVisibleHeight = leftWindowHeight
	}
	rightTop := leftWindowBottom - rightVisibleHeight
	draw.Draw(result, image.Rect(74, rightTop, 119, leftWindowBottom), scaled, source, draw.Over)

	return result, nil
}

func (options Options) validate() error {
	if options.LeftHiddenPct < 0 || options.LeftHiddenPct > 100 {
		return fmt.Errorf("left hidden percentage must be between 0 and 100")
	}
	if options.RightExposedPct < 0 || options.RightExposedPct > 100 {
		return fmt.Errorf("right exposed percentage must be between 0 and 100")
	}
	return nil
}

func percentOf(percent, total int) int { return total * percent / 100 }

// scaleEmoji preserves the source aspect ratio while scaling to 45 pixels
// wide. The reveal windows crop the result during compositing; no vertical
// distortion or separately prescribed source-height cutoff is applied.
func scaleEmoji(target image.Image) image.Image {
	bounds := target.Bounds()
	size := bounds.Size()
	if size.X <= 0 || size.Y <= 0 {
		return image.NewRGBA(image.Rect(0, 0, emojiWidth, 1))
	}
	height := int(math.Ceil(float64(size.Y) * float64(emojiWidth) / float64(size.X)))
	scaled := image.NewRGBA(image.Rect(0, 0, emojiWidth, height))
	if bounds.Empty() {
		return scaled
	}
	xdraw.CatmullRom.Scale(scaled, scaled.Bounds(), target, bounds, draw.Src, nil)
	return scaled
}
