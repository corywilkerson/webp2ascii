package main

import (
    "flag"
    "fmt"
    "image"
    "image/color"
    _ "image/jpeg"
    _ "image/png"
    "math"
    "os"
    "strings"

    _ "golang.org/x/image/webp"
)

const (
    // ASCII characters from darkest to lightest
    asciiChars         = " .:-=+*#%@"
    asciiCharsDetailed = " .'`^\",:;Il!i><~+_-?][}{1)(|/\\tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
)

type Config struct {
    imagePath  string
    width      int
    invert     bool
    detailed   bool
    output     string
    contrast   float64
    gamma      float64
    brightness float64
}

func main() {
    config := parseFlags()

    // Open the image file
    file, err := os.Open(config.imagePath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
        os.Exit(1)
    }
    defer file.Close()

    // Decode the image (supports WebP, JPEG, PNG)
    img, _, err := image.Decode(file)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error decoding image: %v\n", err)
        os.Exit(1)
    }

    // Convert to ASCII
    ascii := imageToASCII(img, config)

    // Output the result
    if config.output != "" {
        err := os.WriteFile(config.output, []byte(ascii), 0644)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("ASCII art saved to %s\n", config.output)
    } else {
        fmt.Println(ascii)
    }
}

func parseFlags() Config {
    var config Config

    flag.IntVar(&config.width, "w", 80, "Width of ASCII output in characters")
    flag.IntVar(&config.width, "width", 80, "Width of ASCII output in characters")
    flag.BoolVar(&config.invert, "i", false, "Invert brightness (for dark terminals)")
    flag.BoolVar(&config.invert, "invert", false, "Invert brightness (for dark terminals)")
    flag.BoolVar(&config.detailed, "d", false, "Use detailed character set")
    flag.BoolVar(&config.detailed, "detailed", false, "Use detailed character set")
    flag.StringVar(&config.output, "o", "", "Save output to file")
    flag.StringVar(&config.output, "output", "", "Save output to file")
    flag.Float64Var(&config.contrast, "c", 1.0, "Contrast adjustment (0.5-3.0, default 1.0)")
    flag.Float64Var(&config.contrast, "contrast", 1.0, "Contrast adjustment (0.5-3.0, default 1.0)")
    flag.Float64Var(&config.gamma, "g", 1.0, "Gamma correction (0.5-2.0, default 1.0)")
    flag.Float64Var(&config.gamma, "gamma", 1.0, "Gamma correction (0.5-2.0, default 1.0)")
    flag.Float64Var(&config.brightness, "b", 0.0, "Brightness adjustment (-0.5 to 0.5, default 0.0)")
    flag.Float64Var(&config.brightness, "brightness", 0.0, "Brightness adjustment (-0.5 to 0.5, default 0.0)")

    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Image to ASCII Art Converter\n\n")
        fmt.Fprintf(os.Stderr, "Usage: %s [options] image.(webp|jpg|jpeg|png)\n\n", os.Args[0])
        fmt.Fprintf(os.Stderr, "Options:\n")
        flag.PrintDefaults()
        fmt.Fprintf(os.Stderr, "\nExamples:\n")
        fmt.Fprintf(os.Stderr, "  %s image.webp\n", os.Args[0])
        fmt.Fprintf(os.Stderr, "  %s image.jpg\n", os.Args[0])
        fmt.Fprintf(os.Stderr, "  %s -w 120 photo.jpeg\n", os.Args[0])
        fmt.Fprintf(os.Stderr, "  %s -i -d image.png\n", os.Args[0])
        fmt.Fprintf(os.Stderr, "  %s -c 1.5 -g 0.8 image.webp  # Enhanced contrast and gamma\n", os.Args[0])
        fmt.Fprintf(os.Stderr, "  %s -o output.txt image.webp\n", os.Args[0])
    }

    flag.Parse()

    // Get the image path from remaining args
    args := flag.Args()
    if len(args) != 1 {
        flag.Usage()
        os.Exit(1)
    }
    config.imagePath = args[0]

    // Validate width
    if config.width < 10 || config.width > 300 {
        fmt.Fprintf(os.Stderr, "Error: Width should be between 10 and 300 characters\n")
        os.Exit(1)
    }

    // Validate enhancement parameters
    if config.contrast < 0.5 || config.contrast > 3.0 {
        fmt.Fprintf(os.Stderr, "Error: Contrast should be between 0.5 and 3.0\n")
        os.Exit(1)
    }
    if config.gamma < 0.5 || config.gamma > 2.0 {
        fmt.Fprintf(os.Stderr, "Error: Gamma should be between 0.5 and 2.0\n")
        os.Exit(1)
    }
    if config.brightness < -0.5 || config.brightness > 0.5 {
        fmt.Fprintf(os.Stderr, "Error: Brightness should be between -0.5 and 0.5\n")
        os.Exit(1)
    }

    return config
}

func imageToASCII(img image.Image, config Config) string {
    // Get image bounds
    bounds := img.Bounds()
    imgWidth := bounds.Max.X - bounds.Min.X
    imgHeight := bounds.Max.Y - bounds.Min.Y

    // Calculate new dimensions
    // Terminal characters are typically ~2.2x taller than wide (more accurate)
    newWidth := config.width
    aspectRatio := float64(imgHeight) / float64(imgWidth)
    newHeight := int(aspectRatio * float64(newWidth) * 0.45)

    // Choose character set
    chars := asciiChars
    if config.detailed {
        chars = asciiCharsDetailed
    }
    if config.invert {
        chars = reverseString(chars)
    }

    // Build ASCII art
    var builder strings.Builder
    
    for y := 0; y < newHeight; y++ {
        for x := 0; x < newWidth; x++ {
            // Map ASCII coordinates to image coordinates
            imgX := bounds.Min.X + x*imgWidth/newWidth
            imgY := bounds.Min.Y + y*imgHeight/newHeight
            
            // Get pixel color and convert to grayscale
            c := img.At(imgX, imgY)
            gray := rgbToGray(c, config)
            
            // Map grayscale value to ASCII character
            charIndex := int(float64(gray) * float64(len(chars)-1) / 65535.0)
            builder.WriteByte(chars[charIndex])
        }
        builder.WriteByte('\n')
    }

    return builder.String()
}

func rgbToGray(c color.Color, config Config) uint16 {
    // Convert to RGBA
    r, g, b, a := c.RGBA()
    
    // Handle transparency by blending with white
    if a < 65535 {
        alpha := float64(a) / 65535.0
        r = uint32(float64(r)*alpha + 65535*(1-alpha))
        g = uint32(float64(g)*alpha + 65535*(1-alpha))
        b = uint32(float64(b)*alpha + 65535*(1-alpha))
    }
    
    // Convert to grayscale using luminance formula
    // Y = 0.299*R + 0.587*G + 0.114*B
    gray := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
    
    // Normalize to 0-1 range
    normalized := gray / 65535.0
    
    // Apply gamma correction
    if config.gamma != 1.0 {
        normalized = math.Pow(normalized, 1.0/config.gamma)
    }
    
    // Apply brightness adjustment
    normalized += config.brightness
    
    // Apply contrast enhancement
    if config.contrast != 1.0 {
        normalized = (normalized - 0.5) * config.contrast + 0.5
    }
    
    // Clamp to 0-1 range
    if normalized < 0 {
        normalized = 0
    } else if normalized > 1 {
        normalized = 1
    }
    
    // Convert back to uint16
    return uint16(normalized * 65535.0)
}


func reverseString(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}