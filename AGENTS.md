# 微服务业务规范
```text
首先遵从agent_read_v2.md
```
# Repository Guidelines

## Project Structure & Module Organization
The Go 1.25.1 backend entrypoint is `cmd/main.go`; the data-import CLI lives in `cmd/migrate`. Core layers are grouped under `internal/` for configuration (`conf`), data access (`dao`), HTTP wiring (`router`, `middleware`), and domain logic (`service`, `models`, `dto`, `utils`). Shared helpers sit in `pkg/pke`, API contracts in `api`, UI assets in `web/templates` and `web/static`, and configuration plus vocab sources in `configs/` (notably `configs/config.yaml` and `configs/words/IELTS.xlsx`). Docker build assets remain in `build/`.

## Build, Test, and Development Commands
Use the Makefile to keep workflows consistent:
- `make localbuild` – tidies modules and compiles the web binary into `word-hero`.
- `make up` / `make down` – start or stop the stack defined in `docker-compose.yaml`.
- `go run ./cmd/main.go` – launch the API server against your active Postgres instance.
- `go run ./cmd/migrate --excel configs/words/IELTS.xlsx --force` – import the bundled vocabulary; omit `--force` to guard against overwrites.

## Coding Style & Naming Conventions
Format Go code with `gofmt` (tabs for indentation) and order imports via `goimports` or `go fmt ./...`. Package names stay lowercase (`internal/router`), exported types use PascalCase, and internal helpers use camelCase. Keep handlers focused on translation and delegate validation or persistence to services and DAOs. HTML templates belong in `web/templates`; static assets belong in `web/static`.

## Testing Guidelines
Author table-driven tests beside the code as `*_test.go` files (see `pkg/pke/example_test.go`). Run `go test ./...` before submitting and add lightweight benchmarks when performance is a concern. Focus coverage on DAO error paths and service validation; document any intentional gaps in the PR.

## Commit & Pull Request Guidelines
Follow the short, imperative commit style visible in history (`Refactor pagination system…`). Group related changes, expand on motivation in the body when the subject is not enough, and reference issue IDs when relevant. PRs should include a brief change summary, test evidence (`go test ./...` output or UI screenshots), config or migration callouts, and rollback notes. Request review from at least one maintainer.

## Configuration & Environment
Defaults load from `configs/config.yaml`, then environment overrides (`WORD_HERO_PORT`, `WORD_HERO_EXCEL_FILE`, `WORD_HERO_PAGE_SIZE`). Keep secrets out of source control—supply them via environment or a mounted config file. The app expects a reachable Postgres instance; update the `Database` block before running migrations. Mount or replace `configs/words` when working with alternate Excel dictionaries.
