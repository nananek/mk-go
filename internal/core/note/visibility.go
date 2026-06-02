package note

import (
	"slices"

	"github.com/shiroha-a/mk/internal/model"
	"github.com/shiroha-a/mk/internal/repository"
)

// CanSeeNote reports whether viewer is allowed to see the given note based on
// its visibility level. viewer may be nil for unauthenticated requests.
//
// followingChecker is invoked only when the visibility level requires a follow
// relationship check; pass nil to skip the check (in which case "followers"
// notes are treated as invisible to non-author viewers).
//
// NOTE: internal/stream/note_publisher.go の noteVisibleToNotifiee に同条件の
// mirror あり (#1471)。stream package が internal/repository / internal/core/note
// を import しない方針 (narrow interface pattern) のため inline 再実装している。
// 本関数の条件 (visibility branch / fail-closed semantics) を変更する際は、
// stream 側の mirror と internal/stream/note_publisher_test.go の equivalence
// test (TestNoteVisibleToNotifiee_EquivalentToCanSeeNote) を必ず追従させること。
func CanSeeNote(viewer *model.User, n *model.Note, followingChecker repository.FollowingRepository) bool {
	if n == nil {
		return false
	}
	switch n.Visibility {
	case model.NoteVisibilityPublic, model.NoteVisibilityHome:
		return true
	case model.NoteVisibilityFollowers:
		if viewer == nil {
			return false
		}
		// 投稿者本人は常に閲覧可
		if viewer.ID == n.UserID {
			return true
		}
		if followingChecker == nil {
			return false
		}
		ok, err := followingChecker.Exists(viewer.ID, n.UserID)
		if err != nil {
			return false
		}
		return ok
	case model.NoteVisibilitySpecified:
		if viewer == nil {
			return false
		}
		if viewer.ID == n.UserID {
			return true
		}
		return slices.Contains(n.VisibleUserIDs, viewer.ID)
	}
	return false
}
