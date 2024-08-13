package runtime

import "io"

type Banner interface {
	PrintBanner(writer io.Writer) error
}
