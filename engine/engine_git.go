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

	// fetch
	cmd = exec.Command("git", "fetch")
	cmd.Path = e.Root

	output, err = cmd.CombinedOutput()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal(string(output))

		return err
	}

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
