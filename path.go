package register

import (
	"fmt"
	"log"
	"path"
	"regexp"
	"sort"
	"strings"
)

var (
	wildcardRegexp  *regexp.Regexp
	normalizeRegexp *regexp.Regexp
)

func init() {
	var err error
	wildcardRegexp, err = regexp.Compile("{(.*?)}")
	if err != nil {
		log.Fatal(err)
	}

	normalizeRegexp, err = regexp.Compile(`/[^\w](.*)[^\w]/?`)
	if err != nil {
		log.Fatal(err)
	}
}

type pathKeys []string

func (x pathKeys) Len() int           { return len(x) }
func (x pathKeys) Less(i, j int) bool { return x[i] > x[j] } // longest first
func (x pathKeys) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func wildcard(s string) string {
	i := strings.Index(s, "{")
	if i >= 0 {
		j := strings.Index(s, "}")
		if j >= 0 {
			return s[i+1 : j]
		}
	}
	return ""
}

// Uses path.Match to match the given path to given pathKeys
// Paths are sorted by length, so the longest path (i.e. most precise) is matched first
// pathkeys have wildcards replaced with *
func findPath(keys pathKeys, ref string) string {
	// sort the paths by length, longest first
	sort.Sort(keys)

	// try to match the given path to the registered paths
	for _, k := range keys {
		if ok, err := path.Match(k, ref); ok && err == nil {
			return k
		} else if err != nil {
			Debug.Errf("error matching path %s to registered path %s: %s", ref, k, err)
		}
	}
	return ""
}

// ExtractVars extracts the variables from the path and saves them to the pathWildcards map.
// Unnamed wildcards (*) can be accessed using the index of the wildcard.
func extractVars(ref string, wildcards map[int]string) map[string]string {
	vars := make(map[string]string)
	pathParts := strings.Split(path.Clean(ref), "/")
	for idx, name := range wildcards {
		if name == "*" {
			// use k as the name of the wildcard
			vars[fmt.Sprintf("%d", idx)] = pathParts[idx]
			continue
		}
		vars[name] = pathParts[idx]
	}

	return vars
}

// breakPath takes a path and removes the starting segment as received by the CloudEvent
// in firestore and realtimeDB, the metadata.Resource.RawPath contains the following prefix
// firestore: "projects/{project-name}/databases/(default)/documents/....."
// realtimeDB: "projects/_/instances/{project-id}/refs/....."
// storage: "projects/_/buckets/{bucket}/objects/......"
func breakRef(ref string) string {
	pathParts := strings.Split(ref, "/")
	if len(pathParts) < 5 {
		return ref
	}
	if ok, _ := path.Match(fsPathBase, path.Join(pathParts[:5]...)); ok {
		return path.Join(pathParts[5:]...)
	}

	if ok, _ := path.Match(rtdbPathBase, path.Join(pathParts[:5]...)); ok {
		return path.Join(pathParts[5:]...)
	}

	if ok, _ := path.Match(gcsBasePath, path.Join(pathParts[:5]...)); ok {
		return path.Join(pathParts[5:]...)
	}
	return ref
}
