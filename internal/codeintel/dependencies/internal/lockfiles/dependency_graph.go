package lockfiles

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/sourcegraph/sourcegraph/internal/conf/reposource"
)

type Edge struct {
	Source, Target reposource.VersionedPackage
}

func newDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		dependencies: make(map[reposource.VersionedPackage][]reposource.VersionedPackage),
		edges:        map[Edge]struct{}{},
	}
}

type DependencyGraph struct {
	dependencies map[reposource.VersionedPackage][]reposource.VersionedPackage
	edges        map[Edge]struct{}
}

func (dg *DependencyGraph) addPackage(pkg reposource.VersionedPackage) {
	if _, ok := dg.dependencies[pkg]; !ok {
		dg.dependencies[pkg] = []reposource.VersionedPackage{}
	}
}
func (dg *DependencyGraph) addDependency(a, b reposource.VersionedPackage) {
	dg.dependencies[a] = append(dg.dependencies[a], b)
	dg.edges[Edge{a, b}] = struct{}{}
}

func (dg *DependencyGraph) Roots() (roots []reposource.VersionedPackage) {
	set := make(map[reposource.VersionedPackage]struct{}, len(dg.dependencies))
	for pkg := range dg.dependencies {
		set[pkg] = struct{}{}
	}

	for edge := range dg.edges {
		delete(set, edge.Target)
	}

	for k := range set {
		roots = append(roots, k)
	}

	return roots
}

func (dg *DependencyGraph) AllEdges() (edges []Edge) {
	for edge := range dg.edges {
		edges = append(edges, edge)
	}
	return edges
}

func (dg *DependencyGraph) String() string {
	var out strings.Builder

	roots := dg.Roots()
	sort.Slice(roots, func(i, j int) bool { return roots[i].Less(roots[j]) })

	for _, root := range roots {
		printDependencies(&out, dg, 0, root)
	}

	return out.String()
}

func printDependencies(out io.Writer, graph *DependencyGraph, level int, node reposource.VersionedPackage) {
	deps, ok := graph.dependencies[node]
	if !ok || len(deps) == 0 {
		fmt.Fprintf(out, "%s%s\n", strings.Repeat("\t", level), node.RepoName())
		return
	}

	fmt.Fprintf(out, "%s%s:\n", strings.Repeat("\t", level), node.RepoName())

	sortedDeps := deps
	sort.Slice(sortedDeps, func(i, j int) bool { return sortedDeps[i].Less(sortedDeps[j]) })

	for _, dep := range sortedDeps {
		printDependencies(out, graph, level+1, dep)
	}
}
