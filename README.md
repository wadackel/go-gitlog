# go-gitlog

[![godoc.org](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/tsuyoshiwada/go-gitlog)
[![Travis](https://img.shields.io/travis/tsuyoshiwada/go-gitlog.svg?style=flat-square)](https://travis-ci.org/tsuyoshiwada/go-gitlog)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/tsuyoshiwada/go-gitlog/blob/master/LICENSE)

> Go (golang) package for providing a means to handle git-log.




## Installation

To install `go-gitlog`, simply run:

```bash
$ go get -u github.com/tsuyoshiwada/go-gitlog
```




## Examples


### `$ git log master`

It is the simplest example.

```go
package main

import (
	"log"
	"github.com/tsuyoshiwada/go-gitlog"
)

func main() {
	// New gitlog
	git := gitlog.New(&gitlog.Config{})

	// List git-log
	commits, err := git.Log("master", nil, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// `commits` -> []*gitlog.Commit{...}
}
```


### `$ git log <sha1|tag|ref>`

Specification of `revisionrange` can be realized by giving a [RevArgs](https://godoc.org/github.com/tsuyoshiwada/go-gitlog#RevArgs) interface.

```go
commits, err := git.Log("master", &gitlog.Rev{"5e312d5"}, nil)
```


### `$ git log <sha1|tag|ref>..<sha1|tag|ref>`

For double dot notation use [RevRange](https://godoc.org/github.com/tsuyoshiwada/go-gitlog#RevRange).

```go
commits, err := git.Log("master", &gitlog.RevRange{
	Old: "v0.1.7",
	New: "v1.2.3",
}, nil)
```


### `$ git log -n <n>`

Use [RevNumber](https://godoc.org/github.com/tsuyoshiwada/go-gitlog#RevNumber) to get the specified number of commits.

```go
commits, err := git.Log("master", &gitlog.RevNumber{10}, nil)
```


### `$ git log --since <time> --until <time>`

By using [RevTime](https://godoc.org/github.com/tsuyoshiwada/go-gitlog#RevTime) you can specify the range by time.

**Since and Until:**

```go
commits, err := git.Log("master", &gitlog.RevTime{
	Since: time.Date(2018, 3, 4, 23, 0, 0, time.UTC),
	Until: time.Date(2018, 1, 2, 12, 0, 0, time.UTC),
}, nil)
```

**Only since:**

```go
commits, err := git.Log("master", &gitlog.RevTime{
	Since: time.Date(2018, 1, 2, 12, 0, 0, time.UTC),
}, nil)
```

**Only until:**

```go
commits, err := git.Log("master", &gitlog.RevTime{
	Until: time.Date(2018, 1, 2, 12, 0, 0, time.UTC),
}, nil)
```




## How it works

Internally we use the git command to format it with the `--pretty` option of log and parse the standard output.  
So, only local git repositories are eligible for acquisition.




## Contribute

1. Fork (https://github.com/tsuyoshiwada/go-gitlog)
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test` command and confirm that it passes
1. Create new Pull Request :muscle:

Bugs, feature requests and comments are more than welcome in the [issues](https://github.com/tsuyoshiwada/go-gitlog/issues).




## License

[MIT © tsuyoshiwada](./LICENSE)
