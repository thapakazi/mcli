# User Profiles via SSH Public Key

## Goal

Add per-user profiles identified by SSH public key, enabling persistent user-specific
preferences: location, bookmarks, read/unread state, and extensible filters.

---

## Current State

- Wish SSH server (`main.go:27`) creates a fresh `model` per session — no user identity
- `ssh.Session` is available in `teaHandler` but the public key is not extracted
- No persistence layer exists — all state is ephemeral
- Events are fetched from `mcli.d` API, filtered/sorted client-side

---

## Plan

### Phase 1: User Identity from SSH Session

- [ ] Extract the SSH public key fingerprint from `ssh.Session` in `teaHandler`
  - `s.PublicKey()` returns `ssh.PublicKey`; derive a fingerprint (SHA256) as the user ID
  - Handle the case where no key is provided (anonymous/password-only) — assign a guest profile or reject
- [ ] Pass the user identity into the `model` so it's available throughout the session
  - Add a `UserID string` (fingerprint) field to the `model` struct

### Phase 2: Profile Data Model

- [ ] Create a `profile` package (`profile/profile.go`) with types:
  ```go
  type UserProfile struct {
      UserID      string            // SSH key fingerprint
      Location    string            // preferred default location
      Bookmarks   []types.EventId   // bookmarked event IDs
      ReadEvents  []types.EventId   // events marked as read
      Filters     map[string]string // extensible key-value filters (e.g. "source": "meetup")
      CreatedAt   time.Time
      UpdatedAt   time.Time
  }
  ```

### Phase 3: Persistence (File-based, JSON)

- [ ] Store profiles as JSON files under a configurable data dir (default: `./data/profiles/`)
  - Filename: `<sha256-fingerprint>.json`
- [ ] Implement CRUD operations in the `profile` package:
  - `Load(userID string) (*UserProfile, error)` — loads or creates a new default profile
  - `Save(profile *UserProfile) error` — writes profile to disk
- [ ] Profile is loaded at session start, saved on meaningful changes (bookmark, mark read, etc.)

### Phase 4: Location Preference

- [ ] On first connect (new profile), prompt user to set a default location via `:set-location <city>`
- [ ] Store location in `UserProfile.Location`
- [ ] Use stored location as default for `:fetch` — no need to type it each time
- [ ] Allow overriding with `:fetch <other-location>` without changing the saved default
- [ ] Add `:set-location <city>` command to update and persist the preference

### Phase 5: Bookmarks

- [ ] Add `b` keybinding to toggle bookmark on the currently selected event
- [ ] Persist bookmarked event IDs in `UserProfile.Bookmarks`
- [ ] Visual indicator in the table for bookmarked events (e.g. a star/marker column)
- [ ] Add `:bookmarks` command to filter view to only bookmarked events

### Phase 6: Read/Unread

- [ ] Track which events a user has "opened" (viewed in sidebar or opened URL)
- [ ] Persist read event IDs in `UserProfile.ReadEvents`
- [ ] Visual indicator in the table — bold/highlight for unread, dimmed for read
- [ ] Add `:unread` command to filter to only unread events

### Phase 7: Extensible Filters

- [ ] Allow saving named filter presets in `UserProfile.Filters`
  - e.g. `:save-filter go` saves the current `/go` filter text
  - `:filter go` restores it
- [ ] These are nice-to-have and can be iterated on

### Phase 8: Config File (Local Mode)

- [ ] When running in CLI mode (not wish), support a local config file (`~/.config/mcli/config.json`)
  - Contains the same fields as `UserProfile` but loaded from disk instead of SSH key
- [ ] Merge with existing `.env` loading in `utils/api.go`

---

## File Changes Summary

| Action | Path | Description |
|--------|------|-------------|
| New    | `profile/profile.go` | UserProfile type + Load/Save |
| Edit   | `model.go` | Add UserID field, wire profile into model |
| Edit   | `main.go` | Extract SSH pubkey in teaHandler, pass to model |
| Edit   | `model.go` | Handle bookmark/read keybindings + commands |
| Edit   | `tui/table.go` | Add bookmark/read visual indicators |
| Edit   | `tui/sidebar.go` | Mark events as read when viewed |
| Edit   | `types/types.go` | (minor) ensure EventId is comparable |
| Edit   | `cmdprompt/cmdprompt.go` | New commands: set-location, bookmarks, unread |

---

## Open Questions

1. Should anonymous SSH users (no public key) get a temporary guest profile, or should we require key-based auth?
2. Should we eventually move persistence to the `mcli.d` API server (server-side profiles) instead of local JSON files?
3. Bookmark storage by event ID — do event IDs remain stable across API fetches?

---

## Suggested Implementation Order

Start with Phase 1-3 (identity + data model + persistence) as the foundation,
then Phase 4 (location) as first visible feature, then Phase 5-6 (bookmarks, read/unread).
Phase 7-8 can be deferred.
