package jobutil

import (
	"context"

	"github.com/opentracing/opentracing-go/log"
	sglog "github.com/sourcegraph/log"

	"github.com/sourcegraph/sourcegraph/internal/search"
	"github.com/sourcegraph/sourcegraph/internal/search/job"
	"github.com/sourcegraph/sourcegraph/internal/search/query"
	"github.com/sourcegraph/sourcegraph/internal/search/repos"
	"github.com/sourcegraph/sourcegraph/internal/search/searcher"
	"github.com/sourcegraph/sourcegraph/internal/search/streaming"
	"github.com/sourcegraph/sourcegraph/internal/search/zoekt"
	"github.com/sourcegraph/sourcegraph/internal/trace"
)

type repoPagerJob struct {
	repoOpts         search.RepoOptions
	useIndex         query.YesNoOnly // whether to include indexed repos
	containsRefGlobs bool            // whether to include repositories with refs
	child            job.Job         // child job tree that need populating a repos field to run
	log              sglog.Logger
}

// setRepos populates the repos field for all jobs that need repos. Jobs are
// copied, ensuring this function is side-effect free.
func setRepos(logger sglog.Logger, job job.Job, indexed *zoekt.IndexedRepoRevs, unindexed []*search.RepositoryRevisions) job.Job {
	setZoektRepos := func(job *zoekt.RepoSubsetTextSearchJob) *zoekt.RepoSubsetTextSearchJob {
		jobCopy := *job
		jobCopy.Repos = indexed
		return &jobCopy
	}

	setSearcherRepos := func(job *searcher.TextSearchJob) *searcher.TextSearchJob {
		jobCopy := *job
		jobCopy.Repos = unindexed
		job.Log = logger
		return &jobCopy
	}

	setZoektSymbolRepos := func(job *zoekt.SymbolSearchJob) *zoekt.SymbolSearchJob {
		jobCopy := *job
		jobCopy.Repos = indexed
		return &jobCopy
	}

	setSymbolSearcherRepos := func(job *searcher.SymbolSearchJob) *searcher.SymbolSearchJob {
		jobCopy := *job
		jobCopy.Repos = unindexed
		return &jobCopy
	}

	setRepos := Mapper{
		Log:                             logger,
		MapZoektRepoSubsetTextSearchJob: setZoektRepos,
		MapZoektSymbolSearchJob:         setZoektSymbolRepos,
		MapSearcherTextSearchJob:        setSearcherRepos,
		MapSymbolSearcherJob:            setSymbolSearcherRepos,
	}

	return setRepos.Map(job)
}

func (p *repoPagerJob) Run(ctx context.Context, clients job.RuntimeClients, stream streaming.Sender) (alert *search.Alert, err error) {
	_, ctx, stream, finish := job.StartSpan(ctx, stream, p)
	defer func() { finish(alert, err) }()

	var maxAlerter search.MaxAlerter

	repoResolver := &repos.Resolver{DB: clients.DB, Opts: p.repoOpts, Log: p.log}
	pager := func(page *repos.Resolved) error {
		indexed, unindexed, err := zoekt.PartitionRepos(
			ctx,
			p.log,
			page.RepoRevs,
			clients.Zoekt,
			search.TextRequest,
			p.useIndex,
			p.containsRefGlobs,
		)
		if err != nil {
			return err
		}

		job := setRepos(p.log, p.child, indexed, unindexed)
		alert, err := job.Run(ctx, clients, stream)
		maxAlerter.Add(alert)
		return err
	}

	return maxAlerter.Alert, repoResolver.Paginate(ctx, pager)
}

func (p *repoPagerJob) Name() string {
	return "RepoPagerJob"
}

func (p *repoPagerJob) Tags() []log.Field {
	return []log.Field{
		trace.Stringer("repoOpts", &p.repoOpts),
		log.String("useIndex", string(p.useIndex)),
		log.Bool("containsRefGlobs", p.containsRefGlobs),
	}
}
