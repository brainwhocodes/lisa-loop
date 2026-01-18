#!/usr/bin/env bats
# Unit Tests for Codex SDK Runner

load '../helpers/test_helper'

# Test: CODEX_BACKEND defaults to cli
@test "CODEX_BACKEND defaults to cli" {
    source "$(dirname "$BATS_TEST_FILENAME")/../../ralph_loop.sh"
    [ "$CODEX_BACKEND" = "cli" ]
}

# Test: check_node_available function exists
@test "check_node_available function exists" {
    source "$(dirname "$BATS_TEST_FILENAME")/../../ralph_loop.sh"
    run type check_node_available
    assert_success
}

# Test: execute_codex_sdk function exists
@test "execute_codex_sdk function exists" {
    source "$(dirname "$BATS_TEST_FILENAME")/../../ralph_loop.sh"
    run type execute_codex_sdk
    assert_success
}
