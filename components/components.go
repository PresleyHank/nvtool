package components

import (
	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	theme "github.com/Nicify/nvtool/theme"
)

func InputTextVWithFont(label string, width float32, value *string, flags g.InputTextFlags, font imgui.Font, cb imgui.InputTextCallback, onChange func()) g.Layout {
	useFont := theme.UseFont(font)
	return g.Layout{
		g.Custom(useFont.Push),
		g.InputTextV(label, width, value, flags, cb, onChange),
		g.Custom(useFont.Pop),
	}
}

func InputTextMultilineWithFont(label string, text *string, width float32, height float32, flags g.InputTextFlags, font imgui.Font, cb imgui.InputTextCallback, onChange func()) g.Layout {
	useFont := theme.UseFont(font)
	return g.Layout{
		g.Custom(useFont.Push),
		g.InputTextMultiline(label, text, width, height, flags, cb, onChange),
		g.Custom(useFont.Pop),
	}
}
