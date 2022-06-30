package lockfiles

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/sourcegraph/sourcegraph/internal/api"
	"github.com/sourcegraph/sourcegraph/internal/conf/reposource"
)

func TestDependencyGraph(t *testing.T) {
	a := &testVersionedPackage{name: "a", version: "1.0.0"}
	b := &testVersionedPackage{name: "b", version: "1.0.0"}
	c := &testVersionedPackage{name: "c", version: "1.0.0"}
	d := &testVersionedPackage{name: "d", version: "1.0.0"}
	e := &testVersionedPackage{name: "e", version: "1.0.0"}
	f := &testVersionedPackage{name: "f", version: "1.0.0"}
	g := &testVersionedPackage{name: "g", version: "1.0.0"}

	dg := newDependencyGraph()
	dg.addPackage(a)
	dg.addPackage(b)
	dg.addPackage(c)
	dg.addPackage(d)
	dg.addPackage(e)
	dg.addPackage(f)
	dg.addPackage(f)

	// a -> b -> d
	//   -> c -> e -> f
	//   -> g
	dg.addDependency(a, b)
	dg.addDependency(a, c)
	dg.addDependency(a, g)

	dg.addDependency(b, d)

	dg.addDependency(c, e)
	dg.addDependency(e, f)

	want := `` +
		`test/a:
	test/b:
		test/d
	test/c:
		test/e:
			test/f
	test/g
`

	got := dg.String()
	fmt.Println(got)

	if d := cmp.Diff(want, got); d != "" {
		t.Fatalf("+want,-got\n%s", d)
	}

}

var _ reposource.VersionedPackage = &testVersionedPackage{}

type testVersionedPackage struct {
	name    string
	version string
}

func (t *testVersionedPackage) VersionedPackageSyntax() string { return t.name }
func (t *testVersionedPackage) PackageSyntax() string          { return t.name }
func (t *testVersionedPackage) RepoName() api.RepoName         { return api.RepoName("test/" + t.name) }
func (t *testVersionedPackage) PackageVersion() string         { return t.version }
func (t *testVersionedPackage) Scheme() string                 { return "test" }
func (t *testVersionedPackage) Description() string            { return "" }
func (t *testVersionedPackage) GitTagFromVersion() string      { return t.version }
func (t *testVersionedPackage) Less(other reposource.VersionedPackage) bool {
	return t.PackageSyntax() < other.PackageSyntax()
}
