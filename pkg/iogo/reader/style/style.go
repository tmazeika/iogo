package style

import (
	"fmt"
	"github.com/zarthus/iogo/v2/pkg/iogo"
	"regexp"
	"strings"
)

type readerStyle struct {
	writer iogo.Writer
	reader iogo.Reader

	confirmRegexp *regexp.Regexp
}

func NewReaderStyle(writer iogo.Writer, reader iogo.Reader) iogo.ReaderStyle {
	return &readerStyle{
		writer:        writer,
		reader:        reader,
		confirmRegexp: regexp.MustCompile("^[yY]"),
	}
}

func (style readerStyle) Prompt(prompt string, options iogo.Options) (string, error) {
	style.writer.WriteLine(prompt)
	return style.reader.ReadLine(options)
}

func (style readerStyle) RequirePrompt(prompt string, options iogo.Options) (string, error) {
	style.writer.WriteLine(prompt)
	return style.reader.ReadLine(options)
}

func (style readerStyle) Confirm(prompt string, options iogo.Options) (bool, error) {
	defaultYes := &options.Default == nil || options.Default == "" || style.confirmRegexp.MatchString(options.Default)

	var yes string
	var no string

	if defaultYes {
		yes = "Y"
		no = "n"
	} else {
		yes = "y"
		no = "N"
	}

	style.writer.WriteLine(prompt + fmt.Sprintf(" (%s/%s)", yes, no))
	result, err := style.reader.ReadLine(options)

	if &result == nil || result == "" {
		return defaultYes, err
	}

	return style.confirmRegexp.MatchString(result), err
}

func (style readerStyle) Select(prompt string, valid []string, options iogo.Options) (string, error) {
	var safeValid []string
	for _, value := range valid {
		safeValid = append(safeValid, value)
	}
	selectRegexp := regexp.MustCompile("^" + strings.Join(safeValid, "|") + "$")

	style.writer.WriteLine(prompt)
	style.writer.WriteLine("Valid options: " + strings.Join(valid, ", "))

	for {
		result, err := style.reader.ReadLine(options)
		if err != nil {
			continue
		}

		if selectRegexp.MatchString(result) {
			return result, err
		} else {
			style.writer.WriteLine("Your input did not match the valid selection.")
			style.writer.WriteLine(prompt)
			style.writer.WriteLine("Valid options: " + strings.Join(valid, ", "))
		}
	}
}
