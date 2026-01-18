# Master checklist

- [ ] Commit 01 — Add Codex integration plan docs (commits/01_codex_plan_docs.md)
- [ ] Commit 02 — Agent selection plumbing (`--agent`) (commits/02_agent_selection_plumbing.md)
- [ ] Commit 03 — Codex CLI runner (`codex exec`) scaffold (commits/03_codex_cli_runner_scaffold.md)
- [ ] Commit 04 — Codex CLI JSONL parsing + resume support (commits/04_codex_cli_jsonl_and_resume.md)
- [ ] Commit 05 — Codex SDK runner + bash integration (commits/05_codex_sdk_runner.md)
- [ ] Commit 06 — Add Codex templates (AGENTS.md + .codex config) (commits/06_codex_templates.md)
- [ ] Commit 07 — Add Codex skills pack in templates (commits/07_codex_skills_pack.md)
- [x] Commit 08 — Charm TUI scaffold (commits/08_charm_tui_scaffold.md) ✅ COMPLETE
  - 10 sub-commits: 08a-08h (packages) + 08i (integration) + 08j (bug fixes)
  - All shell script functionality ported to Go
  - 5 subcommands: run, setup, import, status, reset-circuit
  - TUI and headless monitoring modes
  - ~90 unit tests across 8 packages
- [ ] Commit 09 — TUI polish and docs (commits/09_tui_polish_and_docs.md)
- [ ] Commit 10 — Tests & docs hardening for Codex backends (commits/10_tests_and_docs_hardening.md)
