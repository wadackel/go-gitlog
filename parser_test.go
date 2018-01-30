package gitlog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	assert := assert.New(t)

	commitLog := `@@__GIT_LOG_SEPARATOR__@@HASH:51064a83516c60fdffd99a7d605d168298d91464 51064a8@@__GIT_LOG_DELIMITER__@@TREE:4b825dc642cb6eb9a060e54bf8d69288fbee4904 4b825dc@@__GIT_LOG_DELIMITER__@@AUTHOR:tsuyoshiwada<mail@example.com>[1517138361]@@__GIT_LOG_DELIMITER__@@COMMITTER:tsuyoshiwada<mail@example.com>[1517138361]@@__GIT_LOG_DELIMITER__@@SUBJECT:chore(*): Has commit body@@__GIT_LOG_DELIMITER__@@BODY:This is body comment
multiline
foo
bar

@@__GIT_LOG_SEPARATOR__@@HASH:6dccb5c65f984ec8857243017f506008683342c2 6dccb5c@@__GIT_LOG_DELIMITER__@@TREE:4b825dc642cb6eb9a060e54bf8d69288fbee4904 4b825dc@@__GIT_LOG_DELIMITER__@@AUTHOR:tsuyoshiwada<mail@example.com>[1517134427]@@__GIT_LOG_DELIMITER__@@COMMITTER:tsuyoshiwada<mail@example.com>[1517134427]@@__GIT_LOG_DELIMITER__@@SUBJECT:docs(readme): Test commit@@__GIT_LOG_DELIMITER__@@BODY:
@@__GIT_LOG_SEPARATOR__@@HASH:806512fe97c9c3397b7ed30c0b4076032112f697 806512f@@__GIT_LOG_DELIMITER__@@TREE:4b825dc642cb6eb9a060e54bf8d69288fbee4904 4b825dc@@__GIT_LOG_DELIMITER__@@AUTHOR:tsuyoshiwada<mail@example.com>[1517122160]@@__GIT_LOG_DELIMITER__@@COMMITTER:tsuyoshiwada<mail@example.com>[1517122160]@@__GIT_LOG_DELIMITER__@@SUBJECT:chore(*): Initial commit@@__GIT_LOG_DELIMITER__@@BODY:`

	table := []*Commit{
		&Commit{
			Hash: &Hash{
				Long:  "51064a83516c60fdffd99a7d605d168298d91464",
				Short: "51064a8",
			},
			Tree: &Tree{
				Long:  "4b825dc642cb6eb9a060e54bf8d69288fbee4904",
				Short: "4b825dc",
			},
			Author: &Author{
				Name:  "tsuyoshiwada",
				Email: "mail@example.com",
				Date:  time.Unix(1517138361, 0),
			},
			Committer: &Committer{
				Name:  "tsuyoshiwada",
				Email: "mail@example.com",
				Date:  time.Unix(1517138361, 0),
			},
			Subject: "chore(*): Has commit body",
			Body: `This is body comment
multiline
foo
bar`,
		},
		&Commit{
			Hash: &Hash{
				Long:  "6dccb5c65f984ec8857243017f506008683342c2",
				Short: "6dccb5c",
			},
			Tree: &Tree{
				Long:  "4b825dc642cb6eb9a060e54bf8d69288fbee4904",
				Short: "4b825dc",
			},
			Author: &Author{
				Name:  "tsuyoshiwada",
				Email: "mail@example.com",
				Date:  time.Unix(1517134427, 0),
			},
			Committer: &Committer{
				Name:  "tsuyoshiwada",
				Email: "mail@example.com",
				Date:  time.Unix(1517134427, 0),
			},
			Subject: "docs(readme): Test commit",
			Body:    "",
		},
		&Commit{
			Hash: &Hash{
				Long:  "806512fe97c9c3397b7ed30c0b4076032112f697",
				Short: "806512f",
			},
			Tree: &Tree{
				Long:  "4b825dc642cb6eb9a060e54bf8d69288fbee4904",
				Short: "4b825dc",
			},
			Author: &Author{
				Name:  "tsuyoshiwada",
				Email: "mail@example.com",
				Date:  time.Unix(1517122160, 0),
			},
			Committer: &Committer{
				Name:  "tsuyoshiwada",
				Email: "mail@example.com",
				Date:  time.Unix(1517122160, 0),
			},
			Subject: "chore(*): Initial commit",
			Body:    "",
		},
	}

	parser := &parser{}
	commits, err := parser.parse(&commitLog)

	assert.Nil(err)
	assert.Equal(table, commits)
}
