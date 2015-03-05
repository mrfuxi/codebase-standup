package main

import (
    "fmt"

    "github.com/mrfuxi/go-codebase/codebase"
)

type ChangeMapping struct {
    Status map[string]string
}

func (c ChangeMapping) MapChange(field, before, after string) (description string) {
    change := fmt.Sprintf("%s -> %s", before, after)

    switch {
    case field == codebase.CHANGE_STATUS:
        change = c.Status[change]
    }

    return change
}
