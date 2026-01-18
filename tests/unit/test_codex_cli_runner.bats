#!/usr/bin/env bats
# Unit Tests for Codex CLI Runner

load '../helpers/test_helper'

# Test: Codex dependency check function exists
@test "check_codex_available function exists" {
    source "$(dirname "$BATS_TEST_FILENAME")/../../ralph_loop.sh"
    run type check_codex_available
    assert_success
}

# Test: Codex CLI command uses 'exec' subcommand
@test "build_codex_command uses exec subcommand" {
    source "$(dirname "$BATS_TEST_FILENAME")/../../ralph_loop.sh"
    
    # Create temporary test environment
    export TEST_TEMP_DIR="$(mktemp -d /tmp/ralph-test.XXXXXX)"
    cd "$TEST_TEMP_DIR"
    echo "Test prompt" > "PROMPT.md"
    
    run build_codex_command "PROMPT.md" "" ""
    
    # Check that 'exec' is in the command
    [[ "${CODEX_CMD_ARGS[*]}" =~ "exec" ]]
    assert_success
    
    # Cleanup
    cd /
    rm -rf "$TEST_TEMP_DIR"
}

# Test: build_codex_command adds --json flag for JSON output
@test "build_codex_command adds --json flag" {
    source "$(dirname "$BATS_TEST_FILENAME")/../../ralph_loop.sh"
    
    export TEST_TEMP_DIR="$(mktemp -d /tmp/ralph-test.XXXXXX)"
    cd "$TEST_TEMP_DIR"
    export CODEX_OUTPUT_FORMAT="json"
    echo "Test prompt" > "PROMPT.md"
    
    run build_codex_command "PROMPT.md" "" ""
    
    [[ "${CODEX_CMD_ARGS[*]}" =~ "--json" ]]
    assert_success
    
    cd /
    rm -rf "$TEST_TEMP_DIR"
}

# Test: build_codex_command does not add --json for text output
@test "build_codex_command omits --json flag for text" {
    source "$(dirname "$BATS_TEST_FILENAME")/../../ralph_loop.sh"
    
    export TEST_TEMP_DIR="$(mktemp -d /tmp/ralph-test.XXXXXX)"
    cd "$TEST_TEMP_DIR"
    export CODEX_OUTPUT_FORMAT="text"
    echo "Test prompt" > "PROMPT.md"
    
    run build_codex_command "PROMPT.md" "" ""
    
    ! [[ "${CODEX_CMD_ARGS[*]}" =~ "--json" ]]
    assert_success
    
    cd /
    rm -rf "$TEST_TEMP_DIR"
}

# Test: build_codex_command adds --resume flags with session
@test "build_codex_command adds resume with session ID" {
    source "$(dirname "$BATS_TEST_FILENAME")/../../ralph_loop.sh"
    
    export TEST_TEMP_DIR="$(mktemp -d /tmp/ralph-test.XXXXXX)"
    cd "$TEST_TEMP_DIR"
    export CODEX_USE_CONTINUE=true
    echo "Test prompt" > "PROMPT.md"
    
    run build_codex_command "PROMPT.md" "" "test-session-123"
    
    [[ "${CODEX_CMD_ARGS[*]}" =~ "--resume" ]]
    [[ "${CODEX_CMD_ARGS[*]}" =~ "--thread-id" ]]
    [[ "${CODEX_CMD_ARGS[*]}" =~ "test-session-123" ]]
    assert_success
    
    cd /
    rm -rf "$TEST_TEMP_DIR"
}
