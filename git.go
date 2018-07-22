package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/src-d/go-billy.v4/memfs"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func branchChange(info *asnInfo) error {
	r, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL: fmt.Sprintf("https://%s:%s@%s", os.Getenv("GH_USERNAME"), os.Getenv("GH_TOKEN"), "github.com/innovate-technologies/mobile-asn"),
	})
	if err != nil {
		return err
	}

	// used to check git branches etc
	bare, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: fmt.Sprintf("https://%s:%s@%s", os.Getenv("GH_USERNAME"), os.Getenv("GH_TOKEN"), "github.com/innovate-technologies/mobile-asn"),
	})
	if err != nil {
		return err
	}

	headRef, err := r.Head()
	if err != nil {
		return err
	}

	err = bare.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	})
	if err != nil {
		return err
	}

	branchName := plumbing.ReferenceName(fmt.Sprintf("refs/heads/add-%d-%s", info.AsNumber, info.FirstIP))
	branches, err := bare.Branches()
	if err != nil {
		return err
	}
	err = branches.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name() == branchName {
			return errors.New("Branch exists")
		}
		return nil
	})
	if err != nil {
		return err // asn add bracnh exists
	}
	ref := plumbing.NewHashReference(branchName, headRef.Hash())

	// The created reference is saved in the storage.
	err = r.Storer.SetReference(ref)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}
	err = w.Checkout(&git.CheckoutOptions{Branch: branchName})
	if err != nil {
		return err
	}

	f, err := w.Filesystem.OpenFile("mobile-asn.tsv", os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		return err
	}

	_, err = f.Write([]byte(fmt.Sprintf("%d\t%s\t%s\t%s\t1\n", info.AsNumber, info.AsDescription, info.FirstIP, info.LastIP)))
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	_, err = w.Add("mobile-asn.tsv")
	if err != nil {
		return err
	}
	commit, err := w.Commit(fmt.Sprintf("Add ASN %d", info.AsNumber), &git.CommitOptions{
		Parents: []plumbing.Hash{ref.Hash()},
		Author: &object.Signature{
			Name:  "I Am Mobile",
			Email: "iammobile@shoutca.st",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	_, err = r.CommitObject(commit)
	if err != nil {
		return err
	}

	return r.Push(&git.PushOptions{})
}
