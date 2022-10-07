package git

import (
	"fmt"
	// "log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	// "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
)

func getPublicKeys() *ssh.PublicKeys {
	privateKeyFile := fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME"))

	/* sshKey, err := ioutil.ReadFile(s)
	   signer, err := ssh.ParsePrivateKey([]byte(sshKey))
	   auth := &gitssh.PublicKeys{User: "git", Signer: signer}

	*/
	password := ""

	_, err := os.Stat(privateKeyFile)

	if err != nil {
		fmt.Printf("read file %s failed %s\n", privateKeyFile, err.Error())
		os.Exit(1)
	}

	publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKeyFile, password)

	if err != nil {
		fmt.Printf("generate publickeys failed: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Print(publicKeys)

	return publicKeys
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

func IsCommitPresentOnBranch(repoUrl string, rootCommit *object.Commit, branch string, remoteName string) bool {
	result := false

	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:           repoUrl,
		ReferenceName: plumbing.ReferenceName(branch),
		RemoteName:    remoteName,
		Auth:          getPublicKeys(),
	})

	CheckIfError(err)

	// Gets the HEAD history from HEAD, just like this command:

	// ... retrieves the branch pointed by HEAD
	ref, err := r.Head()
	CheckIfError(err)

	// ... retrieves the commit history
	// FIXME: needs to be configurable
	since := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2099, 7, 30, 0, 0, 0, 0, time.UTC)
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash(), Since: &since, Until: &until})
	CheckIfError(err)

	// ... just iterates over the commits, printing it
	err = cIter.ForEach(func(c *object.Commit) error {
		// FIXME: get to the bottom of isAncestor logic
		isAncestor, parseErr := rootCommit.IsAncestor(c)

		CheckIfError(parseErr)
		if isAncestor {
			result = true
			return nil
		}

		return nil
	})

	CheckIfError(err)

	return result
}

func GetSortedReleases(releases map[string]string) []string {
	versions := make([]string, 0)
	for k := range releases {
		versions = append(versions, k)
	}

	// FIXME: do real semver number sort instead of string alphabetical sort
	sort.Strings(versions)

	for i := len(versions)/2 - 1; i >= 0; i-- {
		opp := len(versions) - 1 - i
		versions[i], versions[opp] = versions[opp], versions[i]
	}

	return versions
}

// SelectRoot is some kind of property where we are not yet sure how it should be impl
func SelectRoot(rootCandidates []string) string {
	// TODO: this should come as default from a flag, lets have, main, master, development fallback
	return "main" // rootCandidates[0]
}

func GetRootCommit(repoUrl string, hash string, rootBranch string) *object.Commit {
	// Clones the given repository, creating the remote, the local branches
	// and fetching the objects, everything in memory:
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:  repoUrl,
		Auth: getPublicKeys(),
	})

	CheckIfError(err)

	// Gets the HEAD history from HEAD, just like this command:
	// ... retrieves the branch pointed by HEAD
	ref, err := r.Head()
	CheckIfError(err)

	// ... retrieves the commit history
	since := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2099, 7, 30, 0, 0, 0, 0, time.UTC)
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash(), Since: &since, Until: &until})
	CheckIfError(err)

	var commit *object.Commit
	// ... just iterates over the commits, printing it
	err = cIter.ForEach(func(c *object.Commit) error {
		if c.Hash.String() == hash {
			commit = c
			return nil
		}
		return nil
	})

	CheckIfError(err)

	return commit
}

// RemoteRemoteBranches fetches remote branches from the repo origin and filters out the root and release branches
func FormatRemoteBranches(repoUrl string, developBranchName string, releaseBranchIdentifiers []string, remoteName string) ([]string, map[string]string) {
	/*	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
			Name: remoteName,
			URLs: []string{repoUrl},
		})

		refs, err := remote.List(&git.ListOptions{})

		if err != nil {
			log.Fatal(err)
			panic(err)
		}

	*/

	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:        repoUrl,
		RemoteName: remoteName,
		Auth:       getPublicKeys(),
	})

	if err != nil {
		panic(err)
	}

	sIter, err := r.Branches()

	if err != nil {
		panic(err)
	}

	releases := make(map[string]string)
	rootCandidates := make([]string, 0)

	err = sIter.ForEach(func(r *plumbing.Reference) error {
		fmt.Println(r.String())
		s := r.String()
		if strings.Contains(s, "refs/heads/") {
			branchName := strings.SplitAfter(s, " ")[1]

			var branchVersion string

			for _, releaseIdentifier := range releaseBranchIdentifiers {
				if strings.Contains(branchName, releaseIdentifier) {
					branchVersion = strings.SplitAfter(branchName, releaseIdentifier)[1]
					releases[branchVersion] = branchName
				}
			}

			if branchName == developBranchName {
				rootCandidates = append(rootCandidates, developBranchName)
			}
		}

		return nil
	})

	/* for _, ref := range refs {
	} */

	return rootCandidates, releases
}
