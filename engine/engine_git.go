package engine

import (
	"os"
	"os/exec"
	"path"

	"github.com/Sirupsen/logrus"
)

// IsGitRepo checks if Root is a git repo.
func (e *Engine) IsGitRepo() bool {
	_, err := os.Stat(path.Join(e.Root, ".git"))
	if err != nil {
		return false
	}
	return true
}

// GitFetchCheckoutPull performs 3 commands: fetch, checkout branch, pull.
func (e *Engine) GitFetchCheckoutPull() error {
	// fetch
	cmd := exec.Command("git", "fetch", "origin")
	cmd.Path = e.Root

	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal(string(output))

		return err
	}

	logrus.Info(string(output))

	// checkout the branch
	cmd = exec.Command("git", "checkout", e.Git.Branch)
	cmd.Path = e.Root

	output, err = cmd.CombinedOutput()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal(string(output))

		return err
	}

	logrus.Info(string(output))

	// pull the branch
	cmd = exec.Command("git", "pull", "origin", e.Git.Branch)
	cmd.Path = e.Root

	output, err = cmd.CombinedOutput()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal(string(output))

		return err
	}

	logrus.Info(string(output))

	return nil
}

// GitCleanPull wipes dirty local changes and perform pull on specific git branch.
func (e *Engine) GitCleanPull() error {
	if !e.IsGitRepo() {
		return nil
	}

	// Clean all existing dirty changes
	cmd := exec.Command("git", "reset", "--hard")
	cmd.Path = e.Root

	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal(string(output))

		return err
	}

	logrus.Info(string(output))

	err = e.GitFetchCheckoutPull()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Unabled to pull from git")

		return err
	}

	logrus.Info(string(output))
	return nil
}

// GitFreshPull performs clone and pull for the first time.
func (e *Engine) GitFreshPull() error {
	if e.IsGitRepo() {
		return nil
	}

	// Make sure root dir does not exist because we will overwrite it anyway.
	os.RemoveAll(e.Root)

	// git clone from HTTPS
	cmd := exec.Command("git", "clone", e.Git.HTTPS, e.Root)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal(string(output))

		return err
	}

	err = e.GitFetchCheckoutPull()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal(string(output))

		return err
	}

	logrus.Info(string(output))
	return nil
}

// GitPull decides to clone for first time or to perform clean pull.
func (e *Engine) GitPull() error {
	if e.IsGitRepo() {
		return e.GitCleanPull()
	} else {
		return e.GitFreshPull()
	}

	return nil
}
