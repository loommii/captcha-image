package captchaimage

import (
	"image"
	"image/color"
	"math/rand"
)

// 添加噪点干扰
func addNoise(img *image.RGBA, noiseDensity float64) {
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	for i := 0; i < int(noiseDensity*float64(width*height)); i++ {
		// 随机选择一个像素位置
		x := rand.Intn(width)
		y := rand.Intn(height)

		// 随机生成一个颜色
		noiseColor := color.RGBA{
			R: uint8(rand.Intn(256)),
			G: uint8(rand.Intn(256)),
			B: uint8(rand.Intn(256)),
			A: 255,
		}

		// 修改该像素点的颜色
		img.Set(x, y, noiseColor)
	}
}
