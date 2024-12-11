package captchaimage

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"
)

// 添加干扰线
func addInterferenceLines(img *image.RGBA, lineCount int) {
	rand.Seed(time.Now().UnixNano())
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	for i := 0; i < lineCount; i++ {
		// 随机生成线条的起点和终点
		x1 := rand.Intn(width)
		y1 := rand.Intn(height)
		x2 := rand.Intn(width)
		y2 := rand.Intn(height)

		// 随机生成线条颜色
		lineColor := color.RGBA{
			R: uint8(rand.Intn(256)),
			G: uint8(rand.Intn(256)),
			B: uint8(rand.Intn(256)),
			A: 255,
		}

		// 绘制线条
		drawLine(img, x1, y1, x2, y2, lineColor)
	}
}

// 绘制线条
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, col color.Color) {
	// 使用 Bresenham 算法绘制线条
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx := 1
	if x1 > x2 {
		sx = -1
	}
	sy := 1
	if y1 > y2 {
		sy = -1
	}
	err := dx - dy
	for {
		img.Set(x1, y1, col)
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := err * 2
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

// 计算绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// 添加波浪干扰线
func addWavyLines(img *image.RGBA, waveCount int, amplitude, frequency float64) {
	rand.Seed(time.Now().UnixNano())
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	for i := 0; i < waveCount; i++ {
		// 计算波浪的起点和终点
		startX := rand.Intn(width)
		endX := startX + rand.Intn(width/2)
		yOffset := rand.Intn(height / 2)

		// 随机颜色
		lineColor := color.RGBA{
			R: uint8(rand.Intn(256)),
			G: uint8(rand.Intn(256)),
			B: uint8(rand.Intn(256)),
			A: 255,
		}

		// 绘制波浪线
		for x := startX; x < endX; x++ {
			y := int(math.Sin(float64(x)*frequency+float64(rand.Intn(10)))*amplitude) + yOffset
			if y >= 0 && y < height {
				img.Set(x, y, lineColor)
			}
		}
	}
}
