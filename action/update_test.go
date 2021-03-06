package action

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-github/github"
	"github.com/helm/helm/log"
	"github.com/helm/helm/release"
)

func TestEnsurePrereqs(t *testing.T) {
	pp := os.Getenv("PATH")
	defer os.Setenv("PATH", pp)

	os.Setenv("PATH", filepath.Join(helmRoot, "testdata")+":"+pp)
	ensurePrereqs()
}

func TestEnsureHome(t *testing.T) {
	tmpHome := createTmpHome()
	ensureHome(tmpHome)
}

func TestCheckLatest(t *testing.T) {
	setupTestCheckLatest()
	var b bytes.Buffer
	defer func() {
		log.Stdout = os.Stdout
		release.RepoService = nil
	}()

	log.IsDebugging = true
	log.Stdout = &b

	CheckLatest("0.0.1")

	if !strings.Contains(b.String(), "A new version of Helm") {
		t.Error("Expected notification of a new release.")
	}
}

type MockGHRepoService struct {
	Release *github.RepositoryRelease
}

func setupTestCheckLatest() {
	v := "9.8.7"
	u := "http://example.com/latest/release"
	i := 987
	r := &github.RepositoryRelease{
		TagName: &v,
		HTMLURL: &u,
		ID:      &i,
	}
	release.RepoService = &MockGHRepoService{Release: r}
}

func (m *MockGHRepoService) GetLatestRelease(o, p string) (*github.RepositoryRelease, *github.Response, error) {
	return m.Release, nil, nil
}
