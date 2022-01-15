package register

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPath(t *testing.T) {
	t.Run("Wildcards", func(t *testing.T) {
		u := "{uid}"
		w := wildcard(u)
		assert.Equal(t, w, "uid", "Wildcard should be uid")

		e := "{emails}"
		w = wildcard(e)
		assert.Equal(t, w, "emails", "Wildcard should be emails")

		x := "not-wildcard"
		w = wildcard(x)
		assert.Equal(t, w, "", "Wildcard should be empty")
	})

	t.Run("FindPaths", func(t *testing.T) {
		p := pathKeys{
			"users/*",
			"users/*/profile/*",
			"msg/*/user/*",
		}

		assert.Equal(t, "users/*", findPath(p, "users/123"), "Should find users/*")
		assert.Equal(t, "users/*/profile/*", findPath(p, "users/123/profile/456"), "Should find users/*/profile/*")
		assert.Equal(t, "msg/*/user/*", findPath(p, "msg/123/user/456"), "Should find msg/*/user/*")

		assert.Equal(t, "", findPath(p, "user/123"), "Should not find, found: %s", findPath(p, "user/123"))
		assert.Equal(t, "", findPath(p, "users/123/"), "Should not find, found %s", findPath(p, "users/123/"))
		assert.Equal(t, "", findPath(p, "users/123/profile"), "Should not find, found %s", findPath(p, "users/123/profile"))
		assert.Equal(t, "", findPath(p, "/msg/*/user/*"), "Should not find, found %s", findPath(p, "/msg/*/user/*"))
	})

	t.Run("ExtractVars", func(t *testing.T) {
		cards := map[int]string{
			1: "uid",
			3: "emails",
		}

		p := "users/12345/contact/faketest"

		vars := extractVars(p, cards)
		assert.Equal(t, "12345", vars["uid"], "Should extract uid")
		assert.Equal(t, "faketest", vars["emails"], "Should extract emails")
	})

	t.Run("BreakPath", func(t *testing.T) {
		fspath := "projects/{project-name}/databases/(default)/documents/users/12345/contact/faketest"
		assert.Equal(t, "users/12345/contact/faketest", breakRef(fspath), "Should break path")

		dbpath := "projects/_/instances/{project-id}/refs/msg/12345/user/34567"
		assert.Equal(t, "msg/12345/user/34567", breakRef(dbpath), "Should break path")

		nonpath := "/non-path/for/testing"
		assert.Equal(t, nonpath, breakRef(nonpath), "Should not break path")

		gtpath := "abc/12345/def/project/users/45678/contact/faketest"
		assert.Equal(t, gtpath, breakRef(gtpath), "Should not break path")
	})
}
