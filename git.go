package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func getLatestTag(repo git.Repository) (*object.Commit, string, error) {
	var latestTagCommit *object.Commit
	var latestTagName string

	tagRefs, err := repo.Tags()
	if err != nil { // standard error check
		return latestTagCommit, latestTagName, err // bail
	}

	err = tagRefs.ForEach(func(tagRef *plumbing.Reference) error {
		// fmt.Println(tagRef)
		revision := plumbing.Revision(tagRef.Name().String())
		tagCommitHash, err := repo.ResolveRevision(revision)
		if err != nil {
			return err
		}

		commit, err := repo.CommitObject(*tagCommitHash)
		if err != nil {
			return err
		}

		if latestTagCommit == nil {
			latestTagCommit = commit
			latestTagName = tagRef.Name().String()
		}

		if commit.Committer.When.After(latestTagCommit.Committer.When) {
			latestTagCommit = commit
			latestTagName = tagRef.Name().String()
		}

		return nil
	})
	if err != nil { // standard error check
		return latestTagCommit, latestTagName, err // bail
	}
	return latestTagCommit, latestTagName, nil // success
}

func getTag(repo git.Repository, tag string) (*object.Commit, error) {
	var commit *object.Commit

	tagRef, err := repo.Tag(tag)
	if err != nil { // standard error check
		return commit, err // bail
	}

	fmt.Println(tagRef)

	return commit, nil // success
}

func checkoutRef(repo git.Repository, tag string) (string, error) {
	tree, err := repo.Worktree()
	if err != nil { // standard error check
		return tag, err // bail
	}

	opts := git.CheckoutOptions{
		Force: true,
	}

	if tag == magicLatest {
		_, latestTagName, err := getLatestTag(repo)
		if err != nil { // standard error check
			return tag, err // bail
		}
		if len(latestTagName) < 1 {
			return tag, errors.New("no tags to inspect")
		}
		tag = strings.ReplaceAll(latestTagName, "refs/tags/", "") // remove refs/tags/ so that it matches what the user would input
	}

	opts.Branch = plumbing.ReferenceName("refs/tags/" + tag) // add refs/tags/ to our tag so go-git can understand it

	// Checkout our tag
	err = tree.Checkout(&opts)
	if err != nil { // standard error check
		return tag, err // bail
	}

	return tag, nil // return human-readable tag (in case it was discovered)
}

func clone(remote string, path string) (git.Repository, error) {
	// setup cloning options, including authentication (if any)
	cloneOpts := git.CloneOptions{
		URL:               remote,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Auth: &http.BasicAuth{
			Username: *cmd_user,
			Password: *cmd_token,
		},
		// Progress: os.Stdout,
	}

	// attempt to clone
	repo, err := git.PlainClone(path, false, &cloneOpts)
	// if the clone failed and the directory already exists, remove it and try again
	if err == git.ErrRepositoryAlreadyExists {
		// deletion
		err = os.RemoveAll(path) // delete everything in the path(!)
		if err != nil {          // couldn't delete
			return *repo, err // bail
		} // successfully deleted

		// recursive function call to try again
		_repo, err := clone(remote, path) // pass in the same data
		if err == nil {                   // see if we made it all of the way through this time
			return _repo, nil // we made it!
		} else { // something bad happened
			return _repo, err // bail
		}
	} else if err != nil { // an error other than the directory already existing
		return *repo, err // bail
	} // we've cloned successfully

	return *repo, nil // we made it!
}
