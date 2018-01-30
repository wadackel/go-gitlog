package gitlog

import (
	"fmt"
	"strconv"
	"time"
)

// RevArgs is used to specify the range of git-log
type RevArgs interface {
	Args() []string
}

// Rev is useful for specifying refname
type Rev struct {
	Ref string
}

// Args ...
func (rev *Rev) Args() []string {
	return []string{rev.Ref}
}

// RevRange is useful for specifying refname range
// alias for `<ref>..<ref>`
type RevRange struct {
	New string
	Old string
}

// Args ...
func (rev *RevRange) Args() []string {
	return []string{fmt.Sprintf("%s..%s", rev.Old, rev.New)}
}

// RevAll alias for `--all`
type RevAll struct{}

// Args ...
func (*RevAll) Args() []string {
	return []string{"--all"}
}

// RevNumber alias for `-n <number>`
type RevNumber struct {
	Limit int
}

// Args ...
func (rev *RevNumber) Args() []string {
	return []string{"-n", strconv.Itoa(rev.Limit)}
}

// RevTime alias for `--since <date> --until <date>`
type RevTime struct {
	Since time.Time
	Until time.Time
}

// Args ...
func (rev *RevTime) Args() []string {
	since := rev.Since.Format("2006-01-02 15:04:05")
	until := rev.Until.Format("2006-01-02 15:04:05")

	if !rev.Since.IsZero() && !rev.Until.IsZero() {
		return []string{
			"--since", since,
			"--until", until,
		}
	} else if !rev.Since.IsZero() {
		return []string{"--since", since}
	} else if !rev.Until.IsZero() {
		return []string{"--until", until}
	}

	return []string{}
}
