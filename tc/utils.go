package tc

import(
	"errors"
	"strings"
	"fmt"
)

type ChannelPermission struct {
	SendHTMLMessage bool // wether or not to allow <html> tags in the message
	SendMessages bool // wether or not to allow user to send messages in channel
}

func MakeError(msg ...any) error {
	if len(msg) == 0 {
		return nil
	}

	var b strings.Builder
	for i, m := range msg {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(fmt.Sprint(m))
	}

	return errors.New(b.String())
}
