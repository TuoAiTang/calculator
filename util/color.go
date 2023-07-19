package util

import "github.com/TuoAiTang/gotable/color"

func Color(clr color.Color, s string) string {
	return color.ColorfulString(clr, s).Val()
}
