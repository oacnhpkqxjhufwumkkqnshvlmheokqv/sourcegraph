package perforce

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/sourcegraph/sourcegraph/internal/authz"
	"github.com/sourcegraph/sourcegraph/internal/extsvc"
)

func TestConvertToPostgresMatch(t *testing.T) {
	// Only needs to implement directory-level perforce protects
	tests := []struct {
		name  string
		match string
		want  string
	}{{
		name:  "*",
		match: "//Sourcegraph/Engineering/*/Frontend/",
		want:  "//Sourcegraph/Engineering/[^/]+/Frontend/",
	}, {
		name:  "...",
		match: "//Sourcegraph/Engineering/.../Frontend/",
		want:  "//Sourcegraph/Engineering/%/Frontend/",
	}, {
		name:  "* and ...",
		match: "//Sourcegraph/*/Src/.../Frontend/",
		want:  "//Sourcegraph/[^/]+/Src/%/Frontend/",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertToPostgresMatch(tt.match)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestConvertToGlobMatch(t *testing.T) {
	// Should fully implement perforce protects
	// Some cases taken directly from https://www.perforce.com/manuals/cmdref/Content/CmdRef/filespecs.html
	// Useful for debugging:
	//
	//   go run github.com/gobwas/glob/cmd/globdraw -p '{//gra*/dep*/,//gra*/dep*}' -s '/' | dot -Tpng -o pattern.png
	//
	tests := []struct {
		name  string
		match string
		want  string

		shouldMatch    []string
		shouldNotMatch []string
	}{{
		name:  "*",
		match: "//Sourcegraph/Engineering/*/Frontend/",
		want:  "//Sourcegraph/Engineering/*/Frontend/",
	}, {
		name:  "...",
		match: "//Sourcegraph/Engineering/.../Frontend/",
		want:  "//Sourcegraph/Engineering/**/Frontend/",
	}, {
		name:           "* and ...",
		match:          "//Sourcegraph/*/Src/.../Frontend/",
		want:           "//Sourcegraph/*/Src/**/Frontend/",
		shouldMatch:    []string{"//Sourcegraph/Path/Src/One/Two/Frontend/"},
		shouldNotMatch: []string{"//Sourcegraph/One/Two/Src/Path/Frontend/"},
	}, {
		name:  "./....c",
		match: "./....c",
		want:  "./**.c",
		shouldMatch: []string{
			"./file.c", "./dir/file.c",
		},
	}, {
		name:  "//gra*/dep*",
		match: "//gra*/dep*",
		want:  `//gra*/dep*{/,}`,
		shouldMatch: []string{
			"//graph/depot/", "//graphs/depots",
		},
		shouldNotMatch: []string{"//graph/depot/release1/"},
	}, {
		name:        "//depot/main/rel...",
		match:       "//depot/main/rel...",
		want:        "//depot/main/rel**",
		shouldMatch: []string{"//depot/main/rel/", "//depot/main/releases/", "//depot/main/release-note.txt", "//depot/main/rel1/product1"},
	}, {
		name:        "//depot/*",
		match:       "//depot/*",
		want:        "//depot/*{/,}",
		shouldMatch: []string{"//depot/main", "//depot/main/"},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToGlobMatch(tt.match)
			if err != nil {
				t.Fatal(fmt.Sprintf("unexpected error: %+v", err))
			}
			if diff := cmp.Diff(tt.want, got.pattern); diff != "" {
				t.Fatal(diff)
			}
			if len(tt.shouldMatch) > 0 {
				for _, m := range tt.shouldMatch {
					if !got.Match(m) {
						t.Errorf("%q should have matched %q", got.pattern, m)
					}
				}
			}
			if len(tt.shouldNotMatch) > 0 {
				for _, m := range tt.shouldNotMatch {
					if got.Match(m) {
						t.Errorf("%q should not have matched %q", got.pattern, m)
					}
				}
			}
		})
	}
}

func mustGlob(t *testing.T, match string) globMatch {
	m, err := convertToGlobMatch(match)
	if err != nil {
		t.Error(err)
	}
	return m
}

// mustGlobPattern gets the glob pattern for a given p4 match for use in testing
func mustGlobPattern(t *testing.T, match string) string {
	return mustGlob(t, match).pattern
}

func TestMatchesAgainstDepot(t *testing.T) {
	type args struct {
		match globMatch
		depot string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{{
		name: "simple match",
		args: args{
			match: mustGlob(t, "//depot/main/..."),
			depot: "//depot/main/",
		},
		want: true,
	}, {
		name: "no wildcard in match",
		args: args{
			match: mustGlob(t, "//depot/"),
			depot: "//depot/main/",
		},
		want: false,
	}, {
		name: "match parent path",
		args: args{
			match: mustGlob(t, "//depot/..."),
			depot: "//depot/main/",
		},
		want: true,
	}, {
		name: "match sub path with all wildcard",
		args: args{
			match: mustGlob(t, "//depot/.../file"),
			depot: "//depot/main/",
		},
		want: true,
	}, {
		name: "match sub path with dir wildcard",
		args: args{
			match: mustGlob(t, "//depot/*/file"),
			depot: "//depot/main/",
		},
		want: true,
	}, {
		name: "match sub path with dir and all wildcards",
		args: args{
			match: mustGlob(t, "//depot/*/file/.../path"),
			depot: "//depot/main/",
		},
		want: true,
	}, {
		name: "match sub path with dir wildcard that's deeply nested",
		args: args{
			match: mustGlob(t, "//depot/*/file/*/another-file/path/"),
			depot: "//depot/main/",
		},
		want: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchesAgainstDepot(tt.args.match, tt.args.depot); got != tt.want {
				t.Errorf("matchesAgainstDepot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanFullRepoPermissions(t *testing.T) {
	f, err := os.Open("testdata/sample-protects-u.txt")
	if err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	rc := io.NopCloser(bytes.NewReader(data))

	execer := p4ExecFunc(func(ctx context.Context, host, user, password string, args ...string) (io.ReadCloser, http.Header, error) {
		return rc, nil, nil
	})

	p := NewTestProvider("", "ssl:111.222.333.444:1666", "admin", "password", execer)
	p.depots = []extsvc.RepoID{
		"//depot/main/",
		"//depot/training/",
		"//depot/test/",
		"//depot/rickroll/",
		"//not-depot/not-main/", // no rules exist
	}
	perms := &authz.ExternalUserPermissions{
		SubRepoPermissions: make(map[extsvc.RepoID]*authz.SubRepoPermissions),
	}
	if err := scanProtects(rc, fullRepoPermsScanner(perms, p.depots)); err != nil {
		t.Fatal(err)
	}

	// See sample-protects-u.txt for notes
	want := &authz.ExternalUserPermissions{
		Exacts: []extsvc.RepoID{
			"//depot/main/",
			"//depot/training/",
			"//depot/test/",
		},
		SubRepoPermissions: map[extsvc.RepoID]*authz.SubRepoPermissions{
			"//depot/main/": {
				PathIncludes: []string{
					mustGlobPattern(t, "base/..."),
					mustGlobPattern(t, "*/stuff/..."),
					mustGlobPattern(t, "frontend/.../stuff/*"),
					mustGlobPattern(t, "config.yaml"),
					mustGlobPattern(t, "subdir/**"),
					mustGlobPattern(t, ".../README.md"),
					mustGlobPattern(t, "dir.yaml"),
				},
				PathExcludes: []string{
					mustGlobPattern(t, "subdir/remove/"),
					mustGlobPattern(t, "subdir/*/also-remove/..."),
					mustGlobPattern(t, ".../.secrets.env"),
				},
			},
			"//depot/test/": {
				PathIncludes: []string{
					mustGlobPattern(t, "..."),
					mustGlobPattern(t, ".../README.md"),
					mustGlobPattern(t, "dir.yaml"),
				},
				PathExcludes: []string{
					mustGlobPattern(t, ".../.secrets.env"),
				},
			},
			"//depot/training/": {
				PathIncludes: []string{
					mustGlobPattern(t, "..."),
					mustGlobPattern(t, ".../README.md"),
					mustGlobPattern(t, "dir.yaml"),
				},
				PathExcludes: []string{
					mustGlobPattern(t, "secrets/..."),
					mustGlobPattern(t, ".env"),
					mustGlobPattern(t, ".../.secrets.env"),
				},
			},
		},
	}
	if diff := cmp.Diff(want, perms); diff != "" {
		t.Fatal(diff)
	}
}

func TestScanFullRepoPermissionsWithWildcardMatchingDepot(t *testing.T) {
	f, err := os.Open("testdata/sample-protects-m.txt")
	if err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	rc := io.NopCloser(bytes.NewReader(data))

	execer := p4ExecFunc(func(ctx context.Context, host, user, password string, args ...string) (io.ReadCloser, http.Header, error) {
		return rc, nil, nil
	})

	p := NewTestProvider("", "ssl:111.222.333.444:1666", "admin", "password", execer)
	p.depots = []extsvc.RepoID{
		"//depot/main/base/",
	}
	perms := &authz.ExternalUserPermissions{
		SubRepoPermissions: make(map[extsvc.RepoID]*authz.SubRepoPermissions),
	}
	if err := scanProtects(rc, fullRepoPermsScanner(perms, p.depots)); err != nil {
		t.Fatal(err)
	}

	want := &authz.ExternalUserPermissions{
		Exacts: []extsvc.RepoID{
			"//depot/main/base/",
		},
		SubRepoPermissions: map[extsvc.RepoID]*authz.SubRepoPermissions{
			"//depot/main/base/": {
				PathIncludes: []string{
					mustGlobPattern(t, "**"),
				},
				PathExcludes: []string{
					mustGlobPattern(t, "**"),
					mustGlobPattern(t, "**/base/build/deleteorgs.txt"),
					mustGlobPattern(t, "build/deleteorgs.txt"),
					mustGlobPattern(t, "**/base/build/**/asdf.txt"),
					mustGlobPattern(t, "build/**/asdf.txt"),
				},
			},
		},
	}

	if diff := cmp.Diff(want, perms); diff != "" {
		t.Fatal(diff)
	}
}

func TestFullScanWildcardDepotMatching(t *testing.T) {
	f, err := os.Open("testdata/sample-protects-x.txt")
	if err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	rc := io.NopCloser(bytes.NewReader(data))

	execer := p4ExecFunc(func(ctx context.Context, host, user, password string, args ...string) (io.ReadCloser, http.Header, error) {
		return rc, nil, nil
	})

	p := NewTestProvider("", "ssl:111.222.333.444:1666", "admin", "password", execer)
	p.depots = []extsvc.RepoID{
		"//depot/654/deploy/base/",
	}
	perms := &authz.ExternalUserPermissions{
		SubRepoPermissions: make(map[extsvc.RepoID]*authz.SubRepoPermissions),
	}
	if err := scanProtects(rc, fullRepoPermsScanner(perms, p.depots)); err != nil {
		t.Fatal(err)
	}

	want := &authz.ExternalUserPermissions{
		Exacts: []extsvc.RepoID{
			"//depot/654/deploy/base/",
		},
		SubRepoPermissions: map[extsvc.RepoID]*authz.SubRepoPermissions{
			"//depot/654/deploy/base/": {
				PathExcludes: []string{
					mustGlobPattern(t, "**/base/build/deleteorgs.txt"),
					mustGlobPattern(t, "build/deleteorgs.txt"),
					mustGlobPattern(t, "asdf/plsql/base/cCustomSchema*.sql"),
				},
				PathIncludes: []string{
					mustGlobPattern(t, "db/upgrade-scripts/**"),
					mustGlobPattern(t, "db/my_db/upgrade-scripts/**"),
					mustGlobPattern(t, "asdf/config/my_schema.xml"),
					mustGlobPattern(t, "db/plpgsql/**"),
				},
			},
		},
	}

	if diff := cmp.Diff(want, perms); diff != "" {
		t.Fatal(diff)
	}
}

func TestCheckWildcardDepotMatch(t *testing.T) {
	testDepot := extsvc.RepoID("//depot/main/base/")
	testCases := []struct {
		name               string
		pattern            string
		original           string
		expectedNewRules   []string
		expectedFoundMatch bool
		depot              extsvc.RepoID
	}{
		{
			name:             "depot match ends with double wildcard",
			pattern:          "//depot/**/README.md",
			original:         "//depot/.../README.md",
			expectedNewRules: []string{"**/README.md"},
			depot:            "//depot/test/",
		},
		{
			name:             "single wildcard",
			pattern:          "//depot/*/dir.yaml",
			original:         "//depot/*/dir.yaml",
			expectedNewRules: []string{"dir.yaml"},
			depot:            "//depot/test/",
		},
		{
			name:             "single wildcard in depot match",
			pattern:          "//depot/**/base/build/deleteorgs.txt",
			original:         "//depot/.../base/build/deleteorgs.txt",
			expectedNewRules: []string{"**/base/build/deleteorgs.txt", "build/deleteorgs.txt"},
			depot:            testDepot,
		},
		{
			name:             "ends with wildcard",
			pattern:          "//depot/**",
			original:         "//depot/...",
			expectedNewRules: []string{"**"},
			depot:            testDepot,
		},
		{
			name:             "two wildcards",
			pattern:          "//depot/**/tests/**/my_test",
			original:         "//depot/.../test/.../my_test",
			expectedNewRules: []string{"**/tests/**/my_test"},
			depot:            testDepot,
		},
		{
			name:             "no match no effect",
			pattern:          "//foo/**/base/build/asdf.txt",
			original:         "//foo/.../base/build/asdf.txt",
			expectedNewRules: []string{"//foo/**/base/build/asdf.txt"},
			depot:            testDepot,
		},
		{
			name:             "original rule is fine, no changes needed",
			pattern:          "//**/.secrets.env",
			original:         "//.../.secrets.env",
			expectedNewRules: []string{"//**/.secrets.env"},
			depot:            testDepot,
		},
		{
			name:             "single wildcard match",
			pattern:          "//depot/6*/*/base/schema/submodules**",
			original:         "//depot/6*/*/base/schema/submodules**",
			expectedNewRules: []string{"schema/submodules**"},
			depot:            "//depot/654/deploy/base/",
		},
		{
			name:             "single wildcard match no double wildcard",
			pattern:          "//depot/6*/*/base/asdf/java/resources/foo.xml",
			original:         "//depot/6*/*/base/asdf/java/resources/foo.xml",
			expectedNewRules: []string{"asdf/java/resources/foo.xml"},
			depot:            "//depot/654/deploy/base/",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pattern := tc.pattern
			glob := mustGlob(t, pattern)
			rule := globMatch{
				glob,
				pattern,
				tc.original,
			}
			newRules := convertRulesForWildcardDepotMatch(rule, tc.depot, map[string]globMatch{})
			if diff := cmp.Diff(newRules, tc.expectedNewRules); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestScanAllUsers(t *testing.T) {
	ctx := context.Background()
	f, err := os.Open("testdata/sample-protects-a.txt")
	if err != nil {
		t.Fatal(err)
	}

	data, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	rc := io.NopCloser(bytes.NewReader(data))

	execer := p4ExecFunc(func(ctx context.Context, host, user, password string, args ...string) (io.ReadCloser, http.Header, error) {
		return rc, nil, nil
	})

	p := NewTestProvider("", "ssl:111.222.333.444:1666", "admin", "password", execer)
	p.cachedGroupMembers = map[string][]string{
		"dev": {"user1", "user2"},
	}
	p.cachedAllUserEmails = map[string]string{
		"user1": "user1@example.com",
		"user2": "user2@example.com",
	}

	users := make(map[string]struct{})
	if err := scanProtects(rc, allUsersScanner(ctx, p, users)); err != nil {
		t.Fatal(err)
	}
	want := map[string]struct{}{
		"user1": {},
		"user2": {},
	}
	if diff := cmp.Diff(want, users); diff != "" {
		t.Fatal(diff)
	}
}
