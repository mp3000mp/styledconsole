package internal

import (
	"testing"

	"github.com/coreoas/styledconsole/internal/style"
	"github.com/stretchr/testify/assert"
)

// TestGetSubstring checks we can extract a substring safely
func TestGetSubstring(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("aaa", getSubstring("zaaaz", 1, 4))
	assert.Equal("aaaz", getSubstring("zaaaz", 1, 50))
	assert.Equal("", getSubstring("zaaaz", 50, 51))
}

// TestExtractStyle checks we can extract styles from a tag
func TestExtractStyle(t *testing.T) {
	assert := assert.New(t)

	// Test a simple, valid style
	assert.Equal(&style.OutputStyle{Foreground: "red", Background: "green"}, extractStyle("bg=green;fg=red"))

	// Test a style with options and href
	assert.Equal(
		&style.OutputStyle{Foreground: "ieua", Background: "aie", Href: "http://github.com", Options: []string{"bold", "italic"}},
		extractStyle("bg=aie;fg=ieua;href=http://github.com;options=bold,italic"),
	)

	// Test an invalid style
	assert.Equal((*style.OutputStyle)(nil), extractStyle("toto=titi;fg=red"))
}

// TestEscapeTrailingBackslash checks we can remove trailing "\"" from texts
func TestEscapeTrailingBackslash(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("noop", EscapeTrailingBackslash("noop"))
	assert.Equal("super super\x00\x00", EscapeTrailingBackslash("super super\\\\"))
	assert.Equal("super \x00 soupaire \x00 awesome \x00", EscapeTrailingBackslash("super \x00 soupaire \x00 awesome \x00"))
	assert.Equal("super \x00 soupaire \x00 awesome \x00\x00\x00", EscapeTrailingBackslash("super \x00 soupaire \x00 awesome \x00\\\\"))
}

// TestFormatWithStyleWithoutTextBefore checks there are no errors of newline
func TestFormatWithStyleWithoutTextBefore(t *testing.T) {
	assert := assert.New(t)
	width := 20
	stack := style.OutputStyleStack{}

	// Test line-cutting
	lastLineLength := 0
	assert.Equal("abc", formatStringWithStyle("abc", width, &lastLineLength, stack))
	assert.Equal(3, lastLineLength)

	lastLineLength = 0
	assert.Equal("supertoto           \nabc                 \n", formatStringWithStyle("supertoto\nabc\n", width, &lastLineLength, stack))
	assert.Equal(0, lastLineLength)

	lastLineLength = 0
	assert.Equal("abc                 \n", formatStringWithStyle("abc\n\n\n", width, &lastLineLength, stack))
	assert.Equal(0, lastLineLength)

	lastLineLength = 5
	assert.Equal(
		"super super sup\ner super super super\n super super super s\nuper",
		formatStringWithStyle("super super super super super super super super super super", width, &lastLineLength, stack),
	)
	assert.Equal(4, lastLineLength)

	// Test with style
	stack.Push(*extractStyle("bg=green;fg=red"))
	lastLineLength = 5
	assert.Equal(
		"\x1b[31;42msuper super sup\x1b[39;49m\n\x1b[31;42mer super super super\x1b[39;49m\n",
		formatStringWithStyle("super super super super super super", width, &lastLineLength, stack),
	)
	assert.Equal(0, lastLineLength)
	stack.PopCurrent()

	// Test edge-cases
	lastLineLength = 0
	assert.Equal(" ", formatStringWithStyle(" ", width, &lastLineLength, stack))
	assert.Equal(1, lastLineLength)

	lastLineLength = 10
	assert.Equal("", formatStringWithStyle("", width, &lastLineLength, stack))
	assert.Equal(10, lastLineLength)

	lastLineLength = 150
	assert.Equal("\nabc", formatStringWithStyle("abc", width, &lastLineLength, stack))
	assert.Equal(3, lastLineLength)

	lastLineLength = -10
	assert.Equal("abc", formatStringWithStyle("abc", width, &lastLineLength, stack))
	assert.Equal(3, lastLineLength)
}

// TestFormatText checks we can render a full text using style tags
func TestFormatText(t *testing.T) {
	assert := assert.New(t)
	width := 20

	assert.Equal(
		"great text",
		FormatText("great text", width),
	)
	assert.Equal(
		"awesome text        \non                  \nmultiple lines.",
		FormatText("awesome text\non\nmultiple lines.", width),
	)
	assert.Equal(
		"\x1b[31mawesome text\x1b[39m",
		FormatText("<fg=red>awesome text</>", width),
	)
	assert.Equal(
		"awesome text \x1b[31mwith st\x1b[39m\n\x1b[31myle and on          \x1b[39m\n\x1b[31mmultiple\x1b[39m lines.",
		FormatText("awesome text <fg=red>with style and on\nmultiple</> lines.", width),
	)
	assert.Equal(
		"awesome text \x1b[31mwith \x1b[39m\x1b[44mim\x1b[49m\n\x1b[44mbricated styles\x1b[49m\x1b[31m and \x1b[39m\n\x1b[31mon                  \x1b[39m\n\x1b[31mmultiple\x1b[39m lines.",
		FormatText("awesome text <fg=red>with <bg=blue>imbricated styles</> and on\nmultiple</> lines.", width),
	)

	// Test edge-cases
	assert.Equal("", FormatText("", width))
	assert.Equal("", FormatText("<fg=red></>", width))
	assert.Equal("qsdf", FormatText("<fg=wrong>qsdf</>", width))
	assert.Equal("<toto=titi>qsdf", FormatText("<toto=titi>qsdf</fg=blue>", width))
	assert.Equal("\x1b[34mtesttest\x1b[39m", FormatText("<fg=blue>testtest", width))
	assert.Equal("testtest", FormatText("testt</fg=blue>est", width))
}
