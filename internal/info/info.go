package info

import "fmt"

type Info struct {
	from    string
	content string
}

func NewInfo(from, content string) *Info {
	return &Info{
		from:    from,
		content: content,
	}
}

func (i *Info) Format() string {
	return fmt.Sprintf(`From:%s,Content:%s`, i.from, "\n"+i.content)
}
