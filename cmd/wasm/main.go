// Modified in 2026 from the original old-man-yells-at for whack-a purposes.
//go:build js && wasm

package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"strings"
	"syscall/js"

	whacka "github.com/bushong1/whack-a"
)

func newError(err error) any {
	return map[string]any{
		"error": err.Error(),
	}
}

func whackAWrapper() js.Func {
	whackAFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return newError(fmt.Errorf("expected 1 argument, got %d", len(args)))
		}
		input := args[0].String()
		_, input, _ = strings.Cut(input, ";base64,")
		decoded, err := base64.StdEncoding.DecodeString(input)
		if err != nil {
			return newError(fmt.Errorf("decoding base64: %w", err))
		}

		im, _, err := image.Decode(bytes.NewReader(decoded))
		if err != nil {
			return newError(fmt.Errorf("decoding image: %w", err))
		}

		whacked := whacka.WhackA(im)

		var buf bytes.Buffer
		enc := base64.NewEncoder(base64.StdEncoding, &buf)
		if err := png.Encode(enc, whacked); err != nil {
			return newError(fmt.Errorf("encoding image: %w", err))
		}
		if err := enc.Close(); err != nil {
			return newError(fmt.Errorf("closing encoder: %w", err))
		}

		return map[string]any{
			"result": buf.String(),
		}
	})
	return whackAFunc
}

func main() {
	js.Global().Set("whackA", whackAWrapper())
	<-make(chan bool)
}
