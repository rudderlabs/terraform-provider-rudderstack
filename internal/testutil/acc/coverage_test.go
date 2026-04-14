package acc_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"

	_ "github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/integrations"
)

// TestAllDestinationsHaveAcceptanceTests verifies that every registered destination
// has a corresponding TestAccDestination* function. This ensures new destinations
// cannot be merged without E2E test coverage.
func TestAllDestinationsHaveAcceptanceTests(t *testing.T) {
	testFuncs, err := collectTestFunctions("rudderstack/integrations/destinations")
	if err != nil {
		t.Fatalf("failed to scan test functions: %v", err)
	}

	for name := range configs.Destinations.Entries() {
		normalized := normalizeKey(name)
		if !hasMatchingFunc(testFuncs, "TestAccDestination", normalized) {
			t.Errorf("destination %q is registered but has no acceptance test matching TestAccDestination<Name>", name)
		}
	}
}

// TestAllSourcesHaveAcceptanceTests verifies that every registered source
// has a corresponding TestAccSource* function.
func TestAllSourcesHaveAcceptanceTests(t *testing.T) {
	testFuncs, err := collectTestFunctions("rudderstack/integrations/sources")
	if err != nil {
		t.Fatalf("failed to scan test functions: %v", err)
	}

	for name := range configs.Sources.Entries() {
		normalized := normalizeKey(name)
		if !hasMatchingFunc(testFuncs, "TestAccSource", normalized) {
			t.Errorf("source %q is registered but has no acceptance test matching TestAccSource<Name>", name)
		}
	}
}

// collectTestFunctions parses *_test.go files in the given directory and returns
// a set of all top-level test function names.
func collectTestFunctions(relDir string) (map[string]bool, error) {
	rootDir, err := findRepoRoot()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(rootDir, relDir)

	fset := token.NewFileSet()
	funcs := make(map[string]bool)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(dir, entry.Name()), nil, 0)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", entry.Name(), err)
		}
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			funcs[fn.Name.Name] = true
		}
	}

	return funcs, nil
}

// findRepoRoot walks up from the working directory to find the go.mod file.
func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find repository root (go.mod)")
		}
		dir = parent
	}
}

// normalizeKey strips underscores and lowercases, so "google_pubsub" and "LINKEDIN_INSIGHT_TAG" both become simple lowercase keys.
func normalizeKey(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, "_", ""))
}

// hasMatchingFunc checks whether any collected test function starts with the given
// prefix and, after stripping the prefix, matches the normalized key (case-insensitive,
// ignoring underscores). This allows "TestAccDestinationCustomerIO" to match registry
// name "customerio" without requiring an exact PascalCase conversion.
func hasMatchingFunc(funcs map[string]bool, prefix, normalizedKey string) bool {
	lowerPrefix := strings.ToLower(prefix)
	for fn := range funcs {
		lower := strings.ToLower(fn)
		if !strings.HasPrefix(lower, lowerPrefix) {
			continue
		}
		rest := lower[len(lowerPrefix):]
		if rest == normalizedKey {
			return true
		}
	}
	return false
}
