# which-fix-version

A TUI that gives immediate insight in what release version your (merge) commit from main / master/ develop(ment) has made it for the first time

## Install binary with make
`make build`

Run `./wfv` to verify the subcommands options are showing:

```
please choose your subcommand, valid options:
- local         find fix version in local repository
- remote        scan a remote repo for the fix version
```

## Example Usage

### When git repository is locally installed

Define arguments in terminal UI:
`./wfv local`

Pass arguments with flags and skip UI:
`./wfv local --path /Users/username/LocalRepos/my-project-repo --commitHash 6c247dd --releaseBranchFormats releases/`

### When connecting to a git repo over https (in pre-alpha mode)
Currently only public repositories are supported. Private repo support is coming soon.

With terminal UI step in between:
`./wfv remote`

With skipped terminal form and arguments as flags:
`./wfv remote --commitHash 7e68ce3c6c7a7d46b86647e10c6bedafd9d8eed2 --releaseBranchFormats release- --url https://github.com/mvangent/which-fix-version.git`

