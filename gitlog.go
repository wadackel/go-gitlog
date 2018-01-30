package gitlog

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	hashField      = "HASH"
	treeField      = "TREE"
	authorField    = "AUTHOR"
	committerField = "COMMITTER"
	subjectField   = "SUBJECT"
	bodyField      = "BODY"

	hashFormat      = hashField + ":%H %h"
	treeFormat      = treeField + ":%T %t"
	authorFormat    = authorField + ":%an<%ae>[%at]"
	committerFormat = committerField + ":%cn<%ce>[%ct]"
	subjectFormat   = subjectField + ":%s"
	bodyFormat      = bodyField + ":%b"

	separator = "@@__GIT_LOG_SEPARATOR__@@"
	delimiter = "@@__GIT_LOG_DELIMITER__@@"

	logFormat = separator +
		hashFormat + delimiter +
		treeFormat + delimiter +
		authorFormat + delimiter +
		committerFormat + delimiter +
		subjectFormat + delimiter +
		bodyFormat
)

// Config for getting git-log
type Config struct {
	GitBin string // default "git"
	Path   string // default "."
}

// Params for getting git-log
type Params struct {
	MergesOnly   bool
	IgnoreMerges bool
}

// GitLog is an interface for git-log acquisition
type GitLog interface {
	Log(string, RevArgs, *Params) ([]*Commit, error)
}

type gitLogImpl struct {
	parser *parser
	config *Config
}

// New GitLog interface
func New(config *Config) GitLog {
	bin := config.GitBin
	if bin == "" {
		bin = "git"
	}

	path := config.Path
	if path == "" {
		path = "."
	}

	return &gitLogImpl{
		parser: &parser{},
		config: &Config{
			GitBin: bin,
			Path:   path,
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

// Check if the `git` command can be executed
func (gitLog *gitLogImpl) canExecuteGit() error {
	_, err := exec.LookPath(gitLog.config.GitBin)
	if err != nil {
		return fmt.Errorf("\"%s\" does not exists", gitLog.config.GitBin)
	}
	return nil
}

// git command execute
func (gitLog *gitLogImpl) git(subcmd string, args ...string) (string, error) {
	gitArgs := append([]string{subcmd}, args...)

	var out bytes.Buffer
	cmd := exec.Command(gitLog.config.GitBin, gitArgs...)
	cmd.Stdout = &out
	cmd.Stderr = ioutil.Discard

	err := cmd.Run()
	if exitError, ok := err.(*exec.ExitError); ok {
		if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
			if waitStatus.ExitStatus() != 0 {
				return "", err
			}
		}
	}

	return strings.TrimRight(strings.TrimSpace(out.String()), "\000"), nil
}

func (gitLog *gitLogImpl) isInsideWorkTree() error {
	out, err := gitLog.git("rev-parse", "--is-inside-work-tree")
	if err != nil {
		return err
	}

	if out != "true" {
		return fmt.Errorf("\"%s\" is no git repository", gitLog.config.Path)
	}

	return nil
}

// Build command line args
func (gitLog *gitLogImpl) buildArgs(ref string, rev RevArgs, params *Params) []string {
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
	}

	args = append(args, ref)

	if rev != nil {
		revisions := rev.Args()
		args = append(args, revisions...)
	}

	return args
}

// Log internally uses the git command to get a list of git-logs
// func (gitLog *gitLogImpl) Log(ref string, rev RevArgs) ([]*Commit, error) {
func (gitLog *gitLogImpl) Log(ref string, rev RevArgs, params *Params) ([]*Commit, error) {
	// Can execute the git command?
	if err := gitLog.canExecuteGit(); err != nil {
		return nil, err
	}

	// To repository path
	back, err := gitLog.workdir()
	if err != nil {
		return nil, err
	}
	defer back()

	// Check inside work tree
	err = gitLog.isInsideWorkTree()
	if err != nil {
		return nil, err
	}

	// Dump git-log
	args := gitLog.buildArgs(ref, rev, params)

	out, err := gitLog.git("log", args...)
	if err != nil {
		return nil, err
	}

	commits, err := gitLog.parser.parse(&out)
	if err != nil {
		return nil, err
	}

	return commits, nil
}
