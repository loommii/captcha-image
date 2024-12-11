package captchaimage

import (
	"image/png"
	"os"
	"testing"
)

func TestCaptchaGenerator_GenerateCaptcha(t *testing.T) {
	t.Run("测试生成英文验证码", func(t *testing.T) {
		captcha := "LOOMMII"
		c := NewCaptchaGenerator()
		img, err := c.GenerateCaptcha(captcha)
		if err != nil {
			t.Error(err.Error())
		}
		outFile, err := os.Create(captcha + ".png")
		if err != nil {
			t.Error(err.Error())
		}
		defer outFile.Close()

		// 将图片保存为 PNG 格式
		err = png.Encode(outFile, img)
		if err != nil {
			t.Error(err.Error())
		}
	})
	t.Run("测试生成中文验证码", func(t *testing.T) {
		captcha := "测试验证码"
		c := NewCaptchaGenerator()
		img, err := c.GenerateCaptcha(captcha)
		if err != nil {
			t.Error(err.Error())
		}
		outFile, err := os.Create(captcha + ".png")
		if err != nil {
			t.Error(err.Error())
		}
		defer outFile.Close()

		// 将图片保存为 PNG 格式
		err = png.Encode(outFile, img)
		if err != nil {
			t.Error(err.Error())
		}
	})
}
