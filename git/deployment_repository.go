package git

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/artdarek/go-unzip"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/otiai10/copy"
	"gosh/log"
	"gosh/util"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Repository interface {
	Initialize() error
	Clone() error
	Pull() error
	Push() error
	CommitChanges(msg string) error
}

const (
	defaultDeploymentRepoTemplateUrl = "https://github.com/ndriessen/gosh-git-template/archive/refs/heads/master.zip"
	defaultUnzipDirectory            = "gosh-git-template-master"
)

var (
	WorkingDirNotEmptyErr    = errors.New("working dir is not empty")
	WorkingDirEmptyErr       = errors.New("your working directory is empty, please initialize it first using gosh init")
	InvalidDeploymentRepoErr = errors.New("working dir does not point to configured deployment repo or has an invalid structure")
)

type DeploymentRepository struct {
	url        string
	publicKeys *ssh.PublicKeys
	git        *git.Repository
}

func isValid(repo *DeploymentRepository) bool {
	return repo != nil && &repo.publicKeys != nil && &repo.git != nil
}

func (repo *DeploymentRepository) OpenWorkingDir() error {
	if !isValid(repo) {
		return errors.New("invalid struct, please use NewDeploymentRepo to create one")
	}
	if isDirectoryEmpty(util.Context.WorkingDir) {
		return WorkingDirEmptyErr
	}
	if gitRepo, err := git.PlainOpen(util.Context.WorkingDir); err == nil {
		repo.git = gitRepo
		return nil
	} else {
		return log.Errf(err, "Error opening working dir as git repository")
	}
}

func (repo *DeploymentRepository) InitFromTemplate() error {
	if file, err := downloadTemplate(); err == nil {
		if err = unzipTemplate(file); err == nil {
			_ = os.Remove(file)
			return nil
		} else {
			return log.Err(err, "Could not extract template repository contents")
		}
	} else {
		return log.Errf(err, "Could not download template repository contents")
	}
}

func NewDeploymentRepository(url string, cloneIfEmpty bool) (*DeploymentRepository, error) {
	var repo *DeploymentRepository
	if publicKeys, err := initSsh(util.Config); err == nil {
		repo = &DeploymentRepository{
			url:        url,
			publicKeys: publicKeys,
		}
	} else {
		return nil, err
	}
	if isDirectoryEmpty(util.Context.WorkingDir) {
		if cloneIfEmpty {
			if err := repo.Clone(); err != nil {
				return nil, err
			}
		} else {
			log.Fatal(errors.New("your working directory is empty, please initialize it first using gosh init"), "Empty working directory")
		}
	} else {
		if err := repo.OpenWorkingDir(); err != nil {
			return nil, err
		}
		//is a valid deployment repo, pull changes
		if err := repo.Pull(); err != nil && err != git.NoErrAlreadyUpToDate {
			return nil, err
		}
	}
	return repo, nil
}

func initSsh(config *util.GoshConfig) (*ssh.PublicKeys, error) {
	sshKey := os.ExpandEnv(config.SshKey)
	_, err := os.Stat(sshKey)
	if err != nil {
		return nil, log.Errf(err, "SSH key %s could not be read", config.SshKey)
	}
	var pwd = ""
	if config.SshPrivateKeyPass != "" {
		decoded, err := base64.StdEncoding.DecodeString(config.SshPrivateKeyPass)
		if err != nil {
			return nil, log.Errf(err, "Unable to decode Base64 private key password from config for key %s", config.SshKey)
		}
		//for some reason it decodes a newline at the end, or the MacOS base64 encoding adds one...
		pwd = strings.TrimSuffix(string(decoded), "\n")
	}
	publicKeys, err := ssh.NewPublicKeysFromFile("git", sshKey, pwd)
	if err != nil {
		return nil, log.Errf(err, "Unable to load and decrypt SSH keys for key %s", config.SshKey)
	}
	return publicKeys, nil
}

func (repo *DeploymentRepository) openWorkingDir() error {
	gitRepo, err := git.PlainOpen(util.Context.WorkingDir)
	repo.git = gitRepo
	return err
}

func (repo *DeploymentRepository) Initialize() error {
	//create new repo from template
	return errors.New("not yet implemented")
}

func (repo *DeploymentRepository) isValidRepository() bool {
	if _, err := os.Stat(filepath.Join(util.Context.WorkingDir, ".git")); err == nil {
		if _, err = os.Stat(filepath.Join(util.Context.WorkingDir, "inventory", "classes")); err == nil {
			if _, err = os.Stat(filepath.Join(util.Context.WorkingDir, "inventory", "targets")); err == nil {
				if repo.git != nil {
					if remotes, err := repo.git.Remotes(); err == nil {
						for _, remote := range remotes {
							for _, url := range remote.Config().URLs {
								if strings.ToLower(url) == strings.ToLower(repo.url) {
									return true
								}
							}
						}
					}
				}
			}
		}
	}
	return false
}

func isDirectoryEmpty(path string) bool {
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		if dir, err := os.Open(path); err == nil {
			if files, err := dir.Readdirnames(1); (err == nil || err == io.EOF) && (files == nil || len(files) == 0) {
				return true
			}
		}
	}
	return false
}

func (repo *DeploymentRepository) Clone() error {
	if !isDirectoryEmpty(util.Context.WorkingDir) {
		return WorkingDirNotEmptyErr
	}
	log.Infof("Cloning deployment repo %s into %s", repo.url, util.Context.WorkingDir)
	if gitRepo, err := git.PlainClone(util.Context.WorkingDir, false, &git.CloneOptions{
		URL:               repo.url,
		Auth:              repo.publicKeys,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Depth:             1,
	}); err == nil {
		repo.git = gitRepo
		return nil
	} else {
		return log.Errf(err, "Error cloning deployment repository in %s", util.Context.WorkingDir)
	}
}

func (repo *DeploymentRepository) Pull() error {
	if !isValid(repo) {
		return errors.New("invalid DeploymentRepository struct, please use NewDeploymentRepository() to create one")
	}
	if !repo.isValidRepository() {
		return InvalidDeploymentRepoErr
	}
	if worktree, err := repo.git.Worktree(); err == nil {
		if err = worktree.Pull(&git.PullOptions{Depth: 1, Auth: repo.publicKeys}); err != nil && err != git.NoErrAlreadyUpToDate {
			return log.Errf(err, "Error updating working dir with remote")
		}
	} else {
		return log.Errf(err, "Error accessing working tree in working dir")
	}
	return nil
}

func (repo *DeploymentRepository) Push(msg string) error {
	if w, err := repo.git.Worktree(); err == nil {
		//err = w.AddGlob(filepath.Join("inventory", "classes", "*"))
		err = w.AddWithOptions(&git.AddOptions{
			All: true,
		})
		if err != nil {
			return err
		}
		if msg == "" {
			msg = "chore: gosh version changes"
		}
		commit, err := w.Commit(msg, &git.CommitOptions{Committer: &object.Signature{
			Name:  "gosh",
			Email: "gosh@github.com",
			When:  time.Now(),
		}})
		if err != nil {
			return err
		}
		commitObject, _ := repo.git.CommitObject(commit)
		log.Debugf("commit: %+v", commitObject)

		err = repo.git.Push(&git.PushOptions{Auth: repo.publicKeys})
		return err
	} else {
		return log.Errf(err, "error pushing changes")
	}
}

func (repo *DeploymentRepository) Commit(msg string) error {
	log.Fatal(errors.New("not implemented"), "not implemented")
	return nil
}

func downloadTemplate() (file string, err error) {
	var client http.Client
	if resp, err := client.Get(defaultDeploymentRepoTemplateUrl); err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			if out, err := os.CreateTemp("", "gosh_git_template_*.zip"); err == nil {
				defer out.Close()
				_, err = io.Copy(out, resp.Body)
				return out.Name(), err
			} else {
				return "", err
			}
		} else {
			return "", errors.New("received " + fmt.Sprint(resp.StatusCode) + " response")
		}
	} else {
		return "", err
	}
}

func unzipTemplate(src string) error {
	if tmpDir, err := os.MkdirTemp("", "*"); err == nil {
		uz := unzip.New(src, tmpDir)
		if err = uz.Extract(); err == nil {
			err = copy.Copy(filepath.Join(tmpDir, defaultUnzipDirectory), util.Context.WorkingDir)
			_ = os.RemoveAll(tmpDir)
		}
		return err
	} else {
		return err
	}
}
