package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type GitService struct {
	workDir string
}

func NewGitService(workDir string) *GitService {
	return &GitService{
		workDir: workDir,
	}
}

func (g *GitService) CloneOrPull(ctx context.Context, repoURL, tag string) (string, error) {
	repoName := g.extractRepoName(repoURL)
	repoPath := filepath.Join(g.workDir, repoName)

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		if err := os.MkdirAll(repoPath, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}

		_, err := git.PlainClone(repoPath, false, &git.CloneOptions{
			URL:      repoURL,
			Progress: os.Stdout,
		})
		if err != nil {
			return "", fmt.Errorf("failed to clone repository: %w", err)
		}
	} else {
		repo, err := git.PlainOpen(repoPath)
		if err != nil {
			return "", fmt.Errorf("failed to open repository: %w", err)
		}

		err = repo.Fetch(&git.FetchOptions{
			RemoteName: "origin",
			Tags:       git.AllTags,
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return "", fmt.Errorf("failed to fetch: %w", err)
		}
	}

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	w, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("failed to get worktree: %w", err)
	}

	tagRef, err := repo.Tag(tag)
	if err != nil {
		return "", fmt.Errorf("failed to get tag %s: %w", tag, err)
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewTagReferenceName(tag),
		Force:  true,
	})
	if err != nil {
		hash := tagRef.Hash()
		err = w.Checkout(&git.CheckoutOptions{
			Hash:  hash,
			Force: true,
		})
		if err != nil {
			return "", fmt.Errorf("failed to checkout tag %s: %w", tag, err)
		}
	}

	return repoPath, nil
}

func (g *GitService) GetPreviousTag(ctx context.Context, repoPath, currentTag string) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	tags, err := repo.Tags()
	if err != nil {
		return "", fmt.Errorf("failed to get tags: %w", err)
	}

	var tagList []string
	err = tags.ForEach(func(ref *plumbing.Reference) error {
		tagName := ref.Name().Short()
		if tagName != currentTag {
			tagList = append(tagList, tagName)
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to iterate tags: %w", err)
	}

	if len(tagList) == 0 {
		return "", nil
	}

	return tagList[len(tagList)-1], nil
}

func (g *GitService) GetChangedApps(ctx context.Context, repoPath, fromTag, toTag string) ([]string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	var fromHash plumbing.Hash
	if fromTag != "" {
		fromRef, err := repo.Tag(fromTag)
		if err != nil {
			return nil, fmt.Errorf("failed to get tag %s: %w", fromTag, err)
		}
		fromHash = fromRef.Hash()
	} else {
		iter, err := repo.CommitObjects()
		if err != nil {
			return nil, fmt.Errorf("failed to get commits: %w", err)
		}
		var firstCommit *object.Commit
		err = iter.ForEach(func(c *object.Commit) error {
			firstCommit = c
			return nil
		})
		if err != nil {
			return nil, err
		}
		if firstCommit != nil {
			fromHash = firstCommit.Hash
		}
	}

	toRef, err := repo.Tag(toTag)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag %s: %w", toTag, err)
	}
	toHash := toRef.Hash()

	fromCommit, err := repo.CommitObject(fromHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get from commit: %w", err)
	}

	toCommit, err := repo.CommitObject(toHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get to commit: %w", err)
	}

	fromTree, err := fromCommit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get from tree: %w", err)
	}

	toTree, err := toCommit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get to tree: %w", err)
	}

	changes, err := fromTree.Diff(toTree)
	if err != nil {
		return nil, fmt.Errorf("failed to diff trees: %w", err)
	}

	appSet := make(map[string]bool)
	for _, change := range changes {
		path := change.To.Name
		if change.To.Name == "" {
			path = change.From.Name
		}

		parts := strings.Split(path, "/")
		if len(parts) > 0 && parts[0] != "" {
			appSet[parts[0]] = true
		}
	}

	apps := make([]string, 0, len(appSet))
	for app := range appSet {
		apps = append(apps, app)
	}

	return apps, nil
}

func (g *GitService) extractRepoName(repoURL string) string {
	parts := strings.Split(repoURL, "/")
	if len(parts) == 0 {
		return "repo"
	}
	name := parts[len(parts)-1]
	name = strings.TrimSuffix(name, ".git")
	return name
}
