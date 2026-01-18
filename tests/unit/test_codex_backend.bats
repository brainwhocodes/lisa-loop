#!/usr/bin/env bats
# Unit Tests for Codex Backend

load '../helpers/test_helper'

# Test: Codex command is used by default
@test "CODEX_CMD is defined as codex" {
    source "$(dirname "$BATS_TEST_FILENAME")/../../ralph_loop.sh"
    [ "$CODEX_CMD" = "codex" ]
}

@test "CODEX_TIMEOUT_MINUTES defaults to 15" {
    source "$(dirname "$BATS_TEST_FILENAME")/../../ralph_loop.sh"
    [ "$CODEX_TIMEOUT_MINUTES" = "15" ]
}

@test "CODEX_USE_CONTINUE is true by default" {
    source "$(dirname "$BATS_TEST_FILENAME")/../../ralph_loop.sh"
    [ "$CODEX_USE_CONTINUE" = "true" ]
}

@test "CODEX_OUTPUT_FORMAT defaults to json" {
    source "$(dirname "$BATS_TEST_FILENAME")/../../ralph_loop.sh"
    [ "$CODEX_OUTPUT_FORMAT" = "json" ]
}

@test "CODEX_SESSION_FILE points to .codex_session_id" {
    source "$(dirname "$BATS_TEST_FILENAME")/../../ralph_loop.sh"
    [ "$CODEX_SESSION_FILE" = ".codex_session_id" ]
}

@test "CODEX_SESSION_EXPIRY_HOURS defaults to 24" {
    source "$(dirname "$BATS_TEST_FILENAME")/../../ralph_loop.sh"
    [ "$CODEX_SESSION_EXPIRY_HOURS" = "24" ]
}
