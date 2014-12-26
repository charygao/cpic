package cpic

//语法制导的图片制作工具

import (
	"strings"
)

type Shape interface {
	Content(content string) []string
	Color() int
}

type Circle struct {
}

func (c *Circle) Content(word string) []string {
	upLine := "‾"     //U+203E
	underScore := "_" //
	l1 := " " + strings.Repeat(underScore, len(word)) + " "
	l2 := "(" + word + ")"
	l3 := " " + strings.Repeat(upLine, len(word)) + " "
	return []string{l1, l2, l3}
}
