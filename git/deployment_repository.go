package git

import (
	"encoding/base64"
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"gosh/log"
	"gosh/util"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Repository interface {
	Initialize() error
	Clone() error
	Pull() error
	Push() error
	CommitChanges(msg string) error
}

var (
	WorkingDirNotEmptyErr    = errors.New("working dir is not empty")
	WorkingDirEmptyErr       = errors.New("your working directory is empty, please initialize it first using gosh init")
	InvalidDeploymentRepoErr = errors.New("working dir does not point to configured deployment repo or has an invalid structure")
	ValidationErr            = errors.New("invalid struct, please use NewDeploymentRepository() to create one")
	DeploymentRepo           *DeploymentRepository
)

type DeploymentRepository struct {
	url        string
	publicKeys *ssh.PublicKeys
	git        *git.Repository
}

func InitializeGit(cloneIfEmpty bool) {
	if repo, err := NewDeploymentRepository(cloneIfEmpty); err != nil {
		log.Fatal(err, "Unable to initialize deployment repository in working directory")
	} else {
		DeploymentRepo = repo
	}
}

func isValid(repo *DeploymentRepository) bool {
	return repo != nil && repo.url != "" && &repo.publicKeys != nil && &repo.git != nil
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

func NewDeploymentRepository(cloneIfEmpty bool) (*DeploymentRepository, error) {
	var repo *DeploymentRepository
	if publicKeys, err := initSsh(util.Config); err == nil {
		repo = &DeploymentRepository{
			url:        util.Config.Url,
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
		if err := repo.openWorkingDirAsGitRepo(); err != nil {
			return nil, err
		}
		if repo.isValidRepository() {
			//is a valid deployment repo, pull changes
			if err := repo.Pull(); err != nil && err != git.NoErrAlreadyUpToDate {
				return nil, err
			}
		} else {
			return nil, log.Errf(InvalidDeploymentRepoErr, "The working dir %s is not empty and is not pointing to the deployment repo %s", util.Context.WorkingDir, util.Config.Url)
		}
	}
	return repo, nil
}

func (repo *DeploymentRepository) openWorkingDirAsGitRepo() error {
	if gitRepo, err := git.PlainOpen(util.Context.WorkingDir); err == nil {
		repo.git = gitRepo
		return nil
	} else {
		return log.Errf(err, "Error opening working dir as git repository")
	}
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
								if strings.ToLower(url) == strings.ToLower(util.Config.Url) {
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
	log.Infof("Cloning deployment repo %s into %s", util.Config.Url, util.Context.WorkingDir)
	if gitRepo, err := git.PlainClone(util.Context.WorkingDir, false, &git.CloneOptions{
		URL:               util.Config.Url,
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

func (repo *DeploymentRepository) Push() error {
	return nil
}

func (repo *DeploymentRepository) Commit(msg string) error {
	return nil
}
