package gitlog

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func rimraf(dir string) {
	if err := os.RemoveAll(dir); err != nil {
		log.Fatalln(err)
	}
}

func mkdirp(dir string) {
	if err := os.MkdirAll(dir, 0777); err != nil {
		log.Fatalln(err)
	}
}

func git(args ...string) string {
	bytes, _ := exec.Command("git", args...).Output()
	return string(bytes)
}

func gitCommit(msg string) string {
	return git("commit", "--allow-empty", "-m", msg)
}

func setup() func() {
	cwd, _ := filepath.Abs(".")
	dir := filepath.Join(cwd, ".tmp")

	mkdirp(dir)
	os.Chdir(dir)

	// setup test repository
	git("init")
	git("config", "--local", "user.name", "authorname")
	git("config", "--local", "user.email", "mail@example.com")

	// master
	gitCommit("chore(*): Initial Commit")

	// topic
	git("checkout", "-b", "topic")
	gitCommit("docs(readme): Has body commit message\n\nThis is commit message body.\nThere are no problems on multiple lines :)")

	// master
	git("checkout", "master")
	gitCommit("feat(parser): Add foo feature")

	// merge
	git("merge", "--no-ff", "topic", "-m", "Merge pull request #12 from tsuyoshiwada/topic")

	gitCommit("fix(logger): Fix bar function")
	gitCommit("style(*): Run GoFmt")
	gitCommit("chore(release): Bump version to v0.0.0")

	os.Chdir(cwd)

	return func() {
		rimraf(dir)
	}
}

func TestNewGitLog(t *testing.T) {
	assert := assert.New(t)

	assert.NotPanics(func() {
		New(nil)
	})

	assert.NotPanics(func() {
		New(&Config{})
	})

	assert.NotPanics(func() {
		New(&Config{
			GitBin: "",
		})
	})

	assert.NotPanics(func() {
		New(&Config{
			Path: "",
		})
	})
}

func TestGitLog(t *testing.T) {
	assert := assert.New(t)

	clear := setup()
	defer clear()

	git := New(&Config{
		Path: ".tmp",
	})

	commits, err := git.Log(nil, nil)

	assert.Nil(err)
	assert.Equal(7, len(commits))

	table := [][]string{
		[]string{"chore(release): Bump version to v0.0.0", ""},
		[]string{"style(*): Run GoFmt", ""},
		[]string{"fix(logger): Fix bar function", ""},
		[]string{"Merge pull request #12 from tsuyoshiwada/topic", ""},
		[]string{"feat(parser): Add foo feature", ""},
		[]string{"docs(readme): Has body commit message", "This is commit message body.\nThere are no problems on multiple lines :)"},
		[]string{"chore(*): Initial Commit", ""},
	}

	for i, commit := range commits {
		expect := table[i]

		assert.Equal("authorname", commit.Author.Name)
		assert.Equal("mail@example.com", commit.Author.Email)

		assert.NotEmpty(commit.Hash.Long)
		assert.NotEmpty(commit.Hash.Short)

		assert.NotEmpty(commit.Tree.Long)
		assert.NotEmpty(commit.Tree.Short)

		assert.Equal(expect[0], commit.Subject)
		assert.Equal(expect[1], commit.Body)
	}
}

func TestGitLogNumber(t *testing.T) {
	assert := assert.New(t)

	clear := setup()
	defer clear()

	git := New(&Config{
		Path: ".tmp",
	})

	commits, err := git.Log(&RevNumber{2}, nil)

	assert.Nil(err)
	assert.Equal(2, len(commits))

	table := [][]string{
		[]string{"chore(release): Bump version to v0.0.0"},
		[]string{"style(*): Run GoFmt"},
	}

	for i, commit := range commits {
		assert.Equal(table[i][0], commit.Subject)
	}
}

func TestGitLogMergesOnly(t *testing.T) {
	assert := assert.New(t)

	clear := setup()
	defer clear()

	git := New(&Config{
		Path: ".tmp",
	})

	commits, err := git.Log(nil, &Params{
		MergesOnly: true,
	})

	assert.Nil(err)
	assert.Equal(1, len(commits))
}

func TestGitLogIgnoreMerges(t *testing.T) {
	assert := assert.New(t)

	clear := setup()
	defer clear()

	git := New(&Config{
		Path: ".tmp",
	})

	commits, err := git.Log(nil, &Params{
		IgnoreMerges: true,
	})

	assert.Nil(err)
	assert.Equal(6, len(commits))
}

func TestGitLogReverse(t *testing.T) {
	assert := assert.New(t)

	clear := setup()
	defer clear()

	git := New(&Config{
		Path: ".tmp",
	})

	commits, err := git.Log(nil, &Params{
		Reverse: true,
	})

	assert.Nil(err)
	assert.Equal(7, len(commits))

	table := [][]string{
		[]string{"chore(*): Initial Commit"},
		[]string{"docs(readme): Has body commit message"},
	}

	for i, commit := range commits {
		if len(table) >= i+1 {
			assert.Equal(table[i][0], commit.Subject)
		}
	}
}

// FIXME: Time tests

func TestGitLogNotFoundGitCommand(t *testing.T) {
	assert := assert.New(t)

	clear := setup()
	defer clear()

	git := New(&Config{
		GitBin: "/notfound/git/bin",
		Path:   ".tmp",
	})

	commits, err := git.Log(nil, nil)

	assert.Nil(commits)
	assert.Contains(err.Error(), "does not exists")
}

func TestGitLogNotFoundPath(t *testing.T) {
	assert := assert.New(t)

	clear := setup()
	defer clear()

	git := New(&Config{
		GitBin: "git",
		Path:   "/notfound/repo",
	})

	commits, err := git.Log(nil, nil)

	assert.Nil(commits)
	assert.Contains(err.Error(), "no such file or directory")
}
