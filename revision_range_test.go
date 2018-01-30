package gitlog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRev(t *testing.T) {
	assert := assert.New(t)

	rev := &Rev{
		Ref: "5e312d5",
	}
	assert.Equal([]string{"5e312d5"}, rev.Args())

	rev = &Rev{
		Ref: "v0.0.1",
	}
	assert.Equal([]string{"v0.0.1"}, rev.Args())
}

func TestRevRange(t *testing.T) {
	assert := assert.New(t)

	rev := &RevRange{
		New: "d01b41a",
		Old: "5e312d5",
	}
	assert.Equal([]string{"5e312d5..d01b41a"}, rev.Args())

	rev = &RevRange{
		New: "v0.1.2",
		Old: "v0.0.1",
	}
	assert.Equal([]string{"v0.0.1..v0.1.2"}, rev.Args())
}

func TestRevAll(t *testing.T) {
	assert := assert.New(t)

	rev := &RevAll{}
	assert.Equal([]string{"--all"}, rev.Args())
}

func TestRevNumber(t *testing.T) {
	assert := assert.New(t)

	rev := &RevNumber{
		Limit: 10,
	}
	assert.Equal([]string{"-n", "10"}, rev.Args())
}

func TestRevTime(t *testing.T) {
	assert := assert.New(t)
	now := time.Now()
	formattedNow := now.Format("2006-01-02 15:04:05")

	rev := &RevTime{
		Since: now,
		Until: now,
	}
	assert.Equal([]string{
		"--since", formattedNow,
		"--until", formattedNow,
	}, rev.Args())

	rev = &RevTime{
		Since: now,
	}
	assert.Equal([]string{
		"--since", formattedNow,
	}, rev.Args())

	rev = &RevTime{
		Until: now,
	}
	assert.Equal([]string{
		"--until", formattedNow,
	}, rev.Args())

	rev = &RevTime{}
	assert.Equal([]string{}, rev.Args())
}
