package main

import (
	"errors"
	"fmt"
	"time"

	"gopkg.in/src-d/go-billy.v3/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func main() {
	err := getCommitTree()
	if err != nil {
		panic(err)
	}
}

// Test commit.Tree() to see if it returns a reference copy of the tree,
// and not the reference of the underlying one
func getCommitTree() error {
	fs := memfs.New()
	storer := memory.NewStorage()

	// Init new in-mem repo
	repo, err := git.Init(storer, fs)
	if err != nil {
		return err
	}

	// Create file
	file, err := fs.Create("example-git-file")
	if err != nil {
		return err
	}

	_, err = file.Write([]byte("hello world!"))
	if err != nil {
		return err
	}
	file.Close()

	// Add to commit
	wt, err := repo.Worktree()
	if err != nil {
		return err
	}
	wt.Add("example-git-file")

	// Commit changes
	_, err = wt.Commit("example go-git commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "John Doe",
			Email: "john@doe.org",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	itr, err := repo.CommitObjects()
	if err != nil {
		return err
	}

	var tree []*object.Tree
	err = itr.ForEach(func(commit *object.Commit) error {
		t, err := commit.Tree()
		if err != nil {
			return err
		}
		tree = append(tree, t)

		commit = nil
		return nil
	})
	if err != nil {
		return err
	}

	err = fs.Remove("example-git-file")
	if err != nil {
		return err
	}

	for _, t := range tree {
		fmt.Println(t.Hash.String())
		f, err := t.File("example-git-file")
		if err != nil {
			return err
		}
		fmt.Println(f.Size)
	}

	_, err = fs.Open("example-git-file")
	if err == nil {
		return errors.New("expected error: file does not exist")
	}

	return nil
}
