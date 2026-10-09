# Night Swipe (holiday-tinder)

A Tinder-style way for a group to pick tonight's spot while on holiday.
One person starts a session (city + restaurants / bars / clubs) and shares
the 6-character code. Everyone swipes through TripAdvisor-style venue cards
(bubble ratings, review counts, price level, neighbourhood). When **everyone**
likes a place, it's a match.

## Stack

- **Backend:** Go (standard library `net/http` + `pgx`), serves the API and
  the embedded mobile-first frontend (`web/`).
- **Database:** PostgreSQL via a one-instance CloudNativePG cluster
  (`k8s/db/`). Schema and sample venues (Lisbon, Barcelona, Amsterdam) are
  created on startup.

## API

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/api/cities` | Cities with venues |
| POST | `/api/sessions` | `{name, city, categories}` → `{code, participantId}` |
| POST | `/api/sessions/{code}/join` | `{name}` → `{code, participantId}` |
| GET | `/api/sessions/{code}/deck?participant=` | Venues not yet swiped |
| POST | `/api/sessions/{code}/swipes` | `{participantId, venueId, liked}` → `{match}` |
| GET | `/api/sessions/{code}/matches` | Participants and liked venues |

## Local development

```sh
DATABASE_URL=postgres://user:pass@localhost:5432/app go run .
```
