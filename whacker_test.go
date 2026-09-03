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

package whacka_test

import (
	"image"
	"image/color"
	"testing"

	whacka "github.com/bushong1/whack-a"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhackAPlacesAndCropsEmoji(t *testing.T) {
	target := image.NewRGBA(image.Rect(0, 0, 45, 90))
	blue := color.RGBA{B: 255, A: 255}
	for y := 0; y < 90; y++ {
		for x := 0; x < 45; x++ {
			target.SetRGBA(x, y, blue)
		}
	}

	result := whacka.WhackA(target)
	require.Equal(t, image.Rect(0, 0, 128, 128), result.Bounds())
	assert.NotEqual(t, blue, color.RGBAModel.Convert(result.At(9, 39)))
	assert.Equal(t, blue, color.RGBAModel.Convert(result.At(9, 40)))
	assert.Equal(t, blue, color.RGBAModel.Convert(result.At(53, 111)))
	assert.NotEqual(t, blue, color.RGBAModel.Convert(result.At(74, 93)))
	assert.Equal(t, blue, color.RGBAModel.Convert(result.At(74, 94)))
	assert.Equal(t, blue, color.RGBAModel.Convert(result.At(118, 111)))
}

func TestWhackAWithOptionsControlsVisibleHeights(t *testing.T) {
	target := image.NewRGBA(image.Rect(0, 0, 45, 90))
	blue := color.RGBA{B: 255, A: 255}
	for y := 0; y < 90; y++ {
		for x := 0; x < 45; x++ {
			target.SetRGBA(x, y, blue)
		}
	}

	result, err := whacka.WhackAWithOptions(target, whacka.Options{
		LeftHiddenPct:   0,
		RightExposedPct: 100,
	})
	require.NoError(t, err)
	assert.Equal(t, blue, color.RGBAModel.Convert(result.At(9, 27)))
	assert.Equal(t, blue, color.RGBAModel.Convert(result.At(53, 111)))
	assert.Equal(t, blue, color.RGBAModel.Convert(result.At(118, 27)))
	assert.Equal(t, blue, color.RGBAModel.Convert(result.At(118, 67)))
	assert.Equal(t, blue, color.RGBAModel.Convert(result.At(118, 111)))

	_, err = whacka.WhackAWithOptions(target, whacka.Options{LeftHiddenPct: 101})
	assert.Error(t, err)
}

func TestWhackABottomAlignsSquareTargets(t *testing.T) {
	target := image.NewRGBA(image.Rect(0, 0, 45, 45))
	blue := color.RGBA{B: 255, A: 255}
	for y := 0; y < 45; y++ {
		for x := 0; x < 45; x++ {
			target.SetRGBA(x, y, blue)
		}
	}

	result := whacka.WhackA(target)
	assert.NotEqual(t, blue, color.RGBAModel.Convert(result.At(9, 75)))
	assert.Equal(t, blue, color.RGBAModel.Convert(result.At(9, 76)))
	assert.Equal(t, blue, color.RGBAModel.Convert(result.At(53, 111)))
	assert.NotEqual(t, blue, color.RGBAModel.Convert(result.At(74, 102)))
	assert.Equal(t, blue, color.RGBAModel.Convert(result.At(74, 103)))
	assert.Equal(t, blue, color.RGBAModel.Convert(result.At(118, 111)))
}

func TestWhackAPreservesNonSquareTargetAspectRatio(t *testing.T) {
	target := image.NewRGBA(image.Rect(0, 0, 45, 90))
	red := color.RGBA{R: 255, A: 255}
	blue := color.RGBA{B: 255, A: 255}
	for y := 0; y < 90; y++ {
		for x := 0; x < 45; x++ {
			if y < 45 {
				target.SetRGBA(x, y, red)
			} else {
				target.SetRGBA(x, y, blue)
			}
		}
	}

	result, err := whacka.WhackAWithOptions(target, whacka.Options{RightExposedPct: 0})
	require.NoError(t, err)
	assert.Equal(t, red, color.RGBAModel.Convert(result.At(9, 67)))
	assert.Equal(t, blue, color.RGBAModel.Convert(result.At(9, 111)))
}
