package main

import (
    "fmt"

    "github.com/mrfuxi/go-codebase/codebase"
)

type ChangeMapping struct {
    Status    map[string]string
    NewTicket string
}

func (c ChangeMapping) MapChange(field, before, after string) (description string) {
    raw_change := field

    if before != "" && after != "" {
        raw_change = fmt.Sprintf("%s -> %s", before, after)
    } else if before != "" || after != "" {
        raw_change = fmt.Sprintf("%s%s", before, after)
    }

    change := ""
    switch {
    case field == codebase.CHANGE_STATUS:
        change = c.Status[raw_change]
    case field == codebase.CHANGE_MILESTONE:
        change = fmt.Sprintf("Moved ticket to %s", after)
    case field == codebase.CHANGE_NEW_TICKET:
        change = c.NewTicket
    }

    if includeRawChange && change == "" {
        change = raw_change
    }

    return change
}
