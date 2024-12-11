package captchaimage

import (
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"os"
	"unicode/utf8"

	"github.com/golang/freetype"
	"golang.org/x/image/math/fixed"
)

type CaptchaGenerator struct {
	Width  int
	Height int
	Font   string
}

// 默认构造函数，返回一个CaptchaGenerator实例
func NewCaptchaGenerator() *CaptchaGenerator {
	return &CaptchaGenerator{
		Width:  150,
		Height: 50,
		Font:   "resource/HarmonyOS_Sans_SC_Medium.ttf", // 默认字体路径
	}
}

// 设置宽度
func (cg *CaptchaGenerator) SetWidth(width int) {
	cg.Width = width
}

// 设置高度
func (cg *CaptchaGenerator) SetHeight(height int) {
	cg.Height = height
}

// 设置字体路径
func (cg *CaptchaGenerator) SetFont(fontPath string) {
	cg.Font = fontPath
}

// 根据图像宽度和内容长度动态计算字体大小
func (cg *CaptchaGenerator) calculateFontSize(content string) float64 {

	// 字符长度的影响系数
	charLengthFactor := 0.95 // 每个字符占宽度的比例，取决于字体大小

	// 最大字体大小
	maxFontSize := float64(cg.Height) * charLengthFactor

	// 根据图像的宽度和验证码的内容长度来动态计算字体大小
	// 目标是：字体大小 = (图像宽度 / 字符数) * 字符长度因子，确保字体适应图像宽度
	fontSize := float64(cg.Width) / float64(utf8.RuneCountInString(content)) * charLengthFactor

	// 限制字体大小范围
	return min(maxFontSize, fontSize)
}

// 生成验证码图像
func (cg *CaptchaGenerator) GenerateCaptcha(captcha string) (*image.RGBA, error) {
	// 创建一个新的空白图片
	img := image.NewRGBA(image.Rect(0, 0, cg.Width, cg.Height))

	// 设置背景颜色为白色
	bgColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// 加载字体
	fontBytes, err := os.ReadFile(cg.Font)
	if err != nil {
		return nil, err
	}

	// 解析字体
	fontParsed, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	c := freetype.NewContext()
	c.SetFont(fontParsed)   // 字体
	c.SetDPI(72)            // 屏幕每英寸的分辨率
	c.SetClip(img.Bounds()) // 设置绘制区域的裁剪矩形。img.Bounds() 返回图像的矩形区域，SetClip 方法告诉 freetype.Context 只在这个矩形区域内绘制文本。防止文本超出图像边界
	c.SetDst(img)           // 设置目标图像 (img) 为绘制的目的地。img 是一个 RGBA 图像，freetype.Context 将会在这个图像上进行文本的绘制。

	// 根据图像宽度和内容长度动态调整字体大小
	fontSize := cg.calculateFontSize(captcha)

	// 设置字体大小
	c.SetFontSize(fontSize)

	// 宽度大约的中间值
	midH := cg.Height/2 + int(fontSize/2)

	pt := freetype.Pt(0, midH)

	heightDiff := cg.Height - int(fontSize)
	verticalOffset := int(heightDiff / 4) // 差值的1/4

	xOffset := cg.Width / utf8.RuneCountInString(captcha)
	// 绘制文本
	for i, v := range captcha {
		c.SetSrc(image.NewUniform(randomColor())) // 随机字体颜色
		// 让字体上下动
		if i%2 == 0 {
			pt.Y -= fixed.Int26_6(verticalOffset) * 64 // 上移
		} else {
			pt.Y += fixed.Int26_6(verticalOffset) * 64 // 下移
		}

		_, err = c.DrawString(string(v), pt)
		if err != nil {
			return nil, err
		}
		pt.X += fixed.Int26_6(xOffset) * 64
	}

	// 添加噪点
	addNoise(img, 0.05)

	// 添加干扰线
	addInterferenceLines(img, 5)

	// 添加波浪线
	addWavyLines(img, 3, 5.0, 0.2)

	//
	// // 保存图像到文件
	// outFile, err := os.Create("imgPath.png")
	// if err != nil {
	// 	return err
	// }
	// defer outFile.Close()

	// // 将图片保存为 PNG 格式
	// err = png.Encode(outFile, img)
	// if err != nil {
	// 	return err
	// }

	return img, nil
}

// 生成随机颜色，避免接近白色
func randomColor() color.RGBA {
	// 定义一个函数来计算颜色的亮度
	// 亮度的计算方式是：亮度 = 0.2126 * R + 0.7152 * G + 0.0722 * B
	// 这是基于人眼对各色光的敏感度加权后的计算公式（ITU-R BT.709标准）
	// 亮度值越高，颜色越接近白色
	isTooBright := func(r, g, b uint8) bool {
		brightness := float64(0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b))
		return brightness > 200 // 如果亮度大于200，则认为颜色接近白色
	}

	// 生成 RGB 的随机值（0 到 255 之间）
	var r, g, b uint8
	for {
		// 限制随机数的范围，避免过亮的颜色
		r = uint8(rand.Intn(206)) + 50 // 范围从 50 到 255
		g = uint8(rand.Intn(206)) + 50 // 范围从 50 到 255
		b = uint8(rand.Intn(206)) + 50 // 范围从 50 到 255

		// 检查生成的颜色是否接近白色
		if !isTooBright(r, g, b) {
			break
		}
	}

	// 返回 RGBA 格式的颜色（A 透明度：255，完全不透明）
	return color.RGBA{R: r, G: g, B: b, A: 255}
}
