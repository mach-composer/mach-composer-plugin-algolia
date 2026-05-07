package internal

import "fmt"

type InvalidSiteConfigError struct {
	msg string
}

func NewInvalidSiteConfigError(format string, a ...any) *InvalidSiteConfigError {
	return &InvalidSiteConfigError{msg: fmt.Sprintf(format, a...)}
}

func (n *InvalidSiteConfigError) Error() string {
	return n.msg
}

type NoApplicationConfigError struct {
	msg string
}

func NewNoApplicationConfigError(format string, a ...any) *NoApplicationConfigError {
	return &NoApplicationConfigError{msg: fmt.Sprintf(format, a...)}
}

func (n *NoApplicationConfigError) Error() string {
	return n.msg
}
