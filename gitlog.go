// Package gitlog is providing a means to handle git-log.
package gitlog

import (
	"os"
	"path/filepath"

	gitcmd "github.com/tsuyoshiwada/go-gitcmd"
)

const (
	hashField      = "HASH"
	treeField      = "TREE"
	authorField    = "AUTHOR"
	committerField = "COMMITTER"
	subjectField   = "SUBJECT"
	bodyField      = "BODY"
	tagField       = "TAG"

	hashFormat      = hashField + ":%H %h"
	treeFormat      = treeField + ":%T %t"
	authorFormat    = authorField + ":%an<%ae>[%at]"
	committerFormat = committerField + ":%cn<%ce>[%ct]"
	subjectFormat   = subjectField + ":%s"
	bodyFormat      = bodyField + ":%b"
	tagFormat       = tagField + ":%D"

	separator = "@@__GIT_LOG_SEPARATOR__@@"
	delimiter = "@@__GIT_LOG_DELIMITER__@@"

	logFormat = separator +
		hashFormat + delimiter +
		treeFormat + delimiter +
		authorFormat + delimiter +
		committerFormat + delimiter +
		tagFormat + delimiter +
		subjectFormat + delimiter +
		bodyFormat
)

// Config for getting git-log
type Config struct {
	Bin  string // default "git"
	Path string // default "."
}

// Params for getting git-log
type Params struct {
	MergesOnly   bool
	IgnoreMerges bool
	Reverse      bool
}

// GitLog is an interface for git-log acquisition
type GitLog interface {
	Log(RevArgs, *Params) ([]*Commit, error)
}

type gitLogImpl struct {
	client gitcmd.Client
	parser *parser
	config *Config
}

// New GitLog interface
func New(config *Config) GitLog {
	bin := "git"
	path := "path"

	if config != nil {
		if config.Bin != "" {
			bin = config.Bin
		}

		if config.Path != "" {
			path = config.Path
		}
	}

	return &gitLogImpl{
		client: gitcmd.New(&gitcmd.Config{
			Bin: bin,
		}),
		parser: &parser{},
		config: &Config{
			Bin:  bin,
			Path: path,
		},
	}
}

// workingdir is go to the temporary working directory
func (gitLog *gitLogImpl) workdir() (func() error, error) {
	cwd, err := filepath.Abs(".")

	back := func() error {
		return os.Chdir(cwd)
	}

	if err != nil {
		return back, err
	}

	dir, err := filepath.Abs(gitLog.config.Path)
	if err != nil {
		return back, err
	}

	err = os.Chdir(dir)
	if err != nil {
		return back, err
	}

	return back, nil
}

// Build command line args
func (gitLog *gitLogImpl) buildArgs(rev RevArgs, params *Params) []string {
	args := []string{
		"--no-decorate",
		"--pretty=\"" + logFormat + "\"",
	}

	if params != nil {
		if params.MergesOnly {
			args = append(args, "--merges")
		}

		if params.IgnoreMerges {
			args = append(args, "--no-merges")
		}

		if params.Reverse {
			args = append(args, "--reverse")
		}
	}

	if rev != nil {
		revisions := rev.Args()
		args = append(args, revisions...)
	}

	return args
}

// Log internally uses the git command to get a list of git-logs
// func (gitLog *gitLogImpl) Log(ref string, rev RevArgs) ([]*Commit, error) {
func (gitLog *gitLogImpl) Log(rev RevArgs, params *Params) ([]*Commit, error) {
	// Can execute the git command?
	if err := gitLog.client.CanExec(); err != nil {
		return nil, err
	}

	// To repository path
	back, err := gitLog.workdir()
	if err != nil {
		return nil, err
	}
	defer back()

	// Check inside work tree
	err = gitLog.client.InsideWorkTree()
	if err != nil {
		return nil, err
	}

	// Dump git-log
	args := gitLog.buildArgs(rev, params)

	out, err := gitLog.client.Exec("log", args...)
	if err != nil {
		return nil, err
	}

	commits, err := gitLog.parser.parse(&out)
	if err != nil {
		return nil, err
	}

	return commits, nil
}
