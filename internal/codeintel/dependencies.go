package codeintel

import (
	"context"
	"database/sql"
	"strings"
	"sync"

	"github.com/inconshreveable/log15"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/sourcegraph/sourcegraph/cmd/frontend/backend"
	"github.com/sourcegraph/sourcegraph/internal/observation"
	"github.com/sourcegraph/sourcegraph/internal/trace"
	"github.com/sourcegraph/sourcegraph/internal/types"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"

	"github.com/sourcegraph/sourcegraph/internal/api"
	codeinteldbstore "github.com/sourcegraph/sourcegraph/internal/codeintel/stores/dbstore"
	"github.com/sourcegraph/sourcegraph/internal/database"
	"github.com/sourcegraph/sourcegraph/internal/database/basestore"
	"github.com/sourcegraph/sourcegraph/internal/gitserver"
	"github.com/sourcegraph/sourcegraph/internal/lockfiles"
)

// DependenciesServices encapsulates the resolution and persistence of dependencies at the repository
// and package levels.
type DependenciesService struct {
	db              database.DB
	sync            func(context.Context, api.RepoName) (*types.Repo, error)
	lockfileService *lockfiles.Service
	operations      *dependencyServiceOperations
}

var (
	depSvc     *DependenciesService
	depSvcOnce sync.Once
)

func GetOrCreateGlobalDependencyService(db database.DB, repoStore database.RepoStore) *DependenciesService {
	depSvcOnce.Do(func() {
		observationContext := &observation.Context{
			Logger:     log15.Root(),
			Tracer:     &trace.Tracer{Tracer: opentracing.GlobalTracer()},
			Registerer: prometheus.DefaultRegisterer,
		}

		depSvc = NewDependenciesService(db, backend.NewRepos(repoStore).GetByName, observationContext)
	})

	return depSvc
}

func NewDependenciesService(
	db database.DB,
	sync func(context.Context, api.RepoName) (*types.Repo, error),
	observationContext *observation.Context,
) *DependenciesService {
	return &DependenciesService{
		db:              db,
		sync:            sync,
		lockfileService: &lockfiles.Service{GitArchive: gitserver.DefaultClient.Archive},
		operations:      newDependencyServiceOperations(observationContext),
	}
}

// RevSpecSet is a utility type for a set of RevSpecs.
type RevSpecSet map[api.RevSpec]struct{}

// Dependencies resolves the (transitive) dependencies for a set of repository and revisions.
// Both the input repoRevs and the output dependencyRevs are a map from repository names to revspecs.
func (r *DependenciesService) Dependencies(ctx context.Context, repoRevs map[api.RepoName]RevSpecSet) (dependencyRevs map[api.RepoName]RevSpecSet, err error) {
	logFields := make([]log.Field, 0, 2)
	if len(repoRevs) == 1 {
		for repoName, revs := range repoRevs {
			revStrs := make([]string, 0, len(revs))
			for rev := range revs {
				revStrs = append(revStrs, string(rev))
			}

			logFields = append(logFields,
				log.String("repo", string(repoName)),
				log.String("revs", strings.Join(revStrs, ",")),
			)
		}
	} else {
		logFields = append(logFields, log.Int("repoRevs", len(repoRevs)))
	}

	ctx, endObservation := r.operations.dependencies.With(ctx, &err, observation.Args{LogFields: logFields})
	defer func() {
		endObservation(1, observation.Args{LogFields: []log.Field{
			log.Int("numDependencyRevs", len(dependencyRevs)),
		}})
	}()

	var mu sync.Mutex
	dependencyRevs = make(map[api.RepoName]RevSpecSet)

	depsStore := codeinteldbstore.Store{Store: basestore.NewWithDB(r.db, sql.TxOptions{})}

	sem := semaphore.NewWeighted(32)
	g, ctx := errgroup.WithContext(ctx)

	for repoName, revs := range repoRevs {
		for rev := range revs {
			repoName, rev := repoName, rev

			g.Go(func() error {
				if err := sem.Acquire(ctx, 1); err != nil {
					return err
				}
				defer sem.Release(1)

				deps, err := r.lockfileService.ListDependencies(ctx, repoName, string(rev))
				if err != nil {
					return err
				}

				for _, dep := range deps {
					if err := sem.Acquire(ctx, 1); err != nil {
						return err
					}

					dep := dep

					g.Go(func() error {
						defer sem.Release(1)

						if err := depsStore.UpsertDependencyRepo(ctx, dep); err != nil {
							return err
						}

						depName := dep.RepoName()
						if _, err := r.sync(ctx, depName); err != nil {
							return err
						}

						depRev := api.RevSpec(dep.GitTagFromVersion())

						mu.Lock()
						defer mu.Unlock()

						if _, ok := dependencyRevs[depName]; !ok {
							dependencyRevs[depName] = RevSpecSet{}
						}
						dependencyRevs[depName][depRev] = struct{}{}

						return nil
					})
				}

				return nil
			})
		}
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return dependencyRevs, nil
}
