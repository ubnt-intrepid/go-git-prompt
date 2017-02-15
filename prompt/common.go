package prompt

import "github.com/ubnt-intrepid/go-git-prompt/color"

type Status interface {
	Prompt(color color.Colored) string
}
