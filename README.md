# Whack-a-mole emoji

Create a 128x128 whack-a-? emoji from any PNG.

Install the CLI with Go:

```sh
go install github.com/bushong1/whack-a/cmd/whack-a@latest
```

Or build it from a local checkout:

```sh
go build -o ./bin/whack-a ./cmd/whack-a
./bin/whack-a --left-hidden-pct 20 --right-exposed-pct 20 my-emoji.png
# Equivalent short form:
./whack-a -l 20 -r 20 my-emoji.png
```

The command writes `whack-a-my-emoji.png` in the current directory. The input
is antialiased to 45 pixels wide without changing its aspect ratio. The left
reveal area can extend from `y=27` to `y=112`, while its visible bottom always
remains at `y=111`.
By default, the lower 20% of that left copy is hidden behind its hole, and the
upper 20% grows upward from `(74, 112)` on the right. Both percentages accept
values from 0 through 100. Both portions are bottom-aligned to `y=111` and can
grow upward to the 85px reveal cap.

`fig/whack-a-blank.png` is the embedded template and `whack-a-mole.png` is an
example output.

## Credit

This project was refactored from [oncilla/old-man-yells-at](https://github.com/oncilla/old-man-yells-at), the original Abe Simpson emoji generator.
