package stevedore

import (
	"testing"

	"github.com/zulily/stevedore/cmd"
)

type generateNamesTestCase struct {
	registry string
	args     []string
	expected []string
}

var (
	testGenerateRepoNamesCases = []generateNamesTestCase{
		generateNamesTestCase{
			registry: "gcr.io/mydomain",
			args:     []string{"foo", "bar", "baz", "Dockerfile"},
			expected: []string{"gcr.io/mydomain/foo-bar:baz", "gcr.io/mydomain/foo-bar:latest"},
		},
		generateNamesTestCase{
			registry: "gcr.io/mydomain",
			args:     []string{"foo", "bar", "baz", "Dockerfile.api"},
			expected: []string{"gcr.io/mydomain/foo-bar-api:baz", "gcr.io/mydomain/foo-bar-api:latest"},
		},
	}
)

func TestGenerateRepoNames(t *testing.T) {
	for _, testCase := range testGenerateRepoNamesCases {
		cmd.Registry = testCase.registry
		result := generateRepoNames(testCase.args[0], testCase.args[1], testCase.args[2], testCase.args[3])
		if result[0] != testCase.expected[0] || result[1] != testCase.expected[1] {
			t.Errorf("Expected (%q), got (%q)", testCase.expected, result)
			t.FailNow()
		}
	}
}

type detectPathTagCase struct {
	repo     string
	expected string
}

var (
	testDetectPathTagCases = []detectPathTagCase{
		detectPathTagCase{
			repo:     "git@github.com:foo/bar.git",
			expected: "foo/bar",
		},
		detectPathTagCase{
			repo:     "https://github.com/foo/bar.git",
			expected: "foo/bar",
		},
	}
)

func TestDetectRepoPathAndTag(t *testing.T) {
	for _, testCase := range testDetectPathTagCases {
		actual := extractRepo(testCase.repo)

		if actual != testCase.expected {
			t.Errorf("Expected (%q), got (%q)", testCase.expected, actual)
			t.FailNow()
		}
	}
}
