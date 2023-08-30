package gitops

import (
	"github.com/fluxcd/go-git-providers/gitprovider"
	"github.com/save-abandoned-projects/ignite/pkg/apis/ignite/scheme"
	"github.com/save-abandoned-projects/ignite/pkg/constants"
	"github.com/save-abandoned-projects/ignite/pkg/operations/reconcile"
	"github.com/save-abandoned-projects/ignite/pkg/providers/manifeststorage"
	"github.com/save-abandoned-projects/libgitops/pkg/gitdir"
	"github.com/save-abandoned-projects/libgitops/pkg/storage/sync"
	log "github.com/sirupsen/logrus"
)

func RunGitOps(url string, opts gitdir.GitDirectoryOptions) error {
	log.Infof("Starting GitOps loop for repo at %q\n", url)
	log.Info("Whenever changes are pushed to the target branch, Ignite will apply the desired state locally\n")

	// Construct the GitDirectory implementation which backs the storage
	repoRef, err := gitprovider.ParseOrgRepositoryURL(url)
	if err != nil {
		return err
	}

	gitDir, err := gitdir.NewGitDirectory(repoRef, opts)
	if err != nil {
		return err
	}
	// TODO: Run gitDir.Cleanup() on SIGINT

	// Wait for the repo to be cloned
	if err := gitDir.StartCheckoutLoop(); err != nil {
		return err
	}

	// Construct a manifest storage for the path backed by git
	s, err := manifeststorage.NewTwoWayManifestStorage(gitDir.Dir(), constants.DATA_DIR, scheme.Serializer)
	if err != nil {
		return err
	}

	// TODO: Make the reconcile function signal-aware
	reconcile.ReconcileManifests(s.(*sync.SyncStorage))
	return nil
}
