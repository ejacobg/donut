# donut

An implementation of Andy Sloane’s [donut.c](https://www.a1k0n.net/2011/07/20/donut-math.html), written in Go.

## Features

-   Renders to ASCII and GIF
-   Supports grayscale (default) and monochrome colors
-   Supports custom colors and logic
-   Properties of the scene can be controlled

## Examples

![monochrome](mono.gif)![grayscale](grayscale.gif)

[![asciicast](https://asciinema.org/a/521347.svg)](https://asciinema.org/a/521347)

## Usage

```
Usage of ./donut:
  -delay int
        Delay between frames, in milliseconds (default 50)
  -gif string
        Generate a GIF file with the given name
  -height int
        Height of the GIF file, if needed (default 100)
  -steps int
        Number of frames to generate (default 200)
  -width int
        Width of the GIF file, if needed (default 100)
```

If the `gif` flag is omitted, the program will print to the terminal, scaled proportionally to the terminal size.
