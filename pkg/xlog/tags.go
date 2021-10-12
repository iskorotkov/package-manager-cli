package xlog

import (
	"log"
	"strings"
)

const TagSeparator = "."

func Push(tag string) {
	trimmed := strings.TrimRight(log.Prefix(), " ")

	var tags []string
	if trimmed != "" {
		tags = strings.Split(trimmed, TagSeparator)
	}

	tags = append(tags, tag)

	prefix := strings.Join(tags, TagSeparator)

	log.SetPrefix(prefix + " ")
}

func Pop() {
	trimmed := strings.TrimRight(log.Prefix(), " ")
	if trimmed == "" {
		return
	}

	tags := strings.Split(trimmed, TagSeparator)
	prefix := strings.Join(tags[:len(tags)-1], TagSeparator)

	log.SetPrefix(prefix + " ")
}

func PushN(tags ...string) {
	for _, tag := range tags {
		Push(tag)
	}
}

func PopN(n int) {
	for i := 0; i < n; i++ {
		Pop()
	}
}
