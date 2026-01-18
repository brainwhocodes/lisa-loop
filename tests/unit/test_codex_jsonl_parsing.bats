#!/usr/bin/env bats
# Unit Tests for Codex JSONL Parsing and Resume Support

load '../helpers/test_helper'

# Test: parse_codex_jsonl extracts thread_id from thread.started event
@test "parse_codex_jsonl extracts thread_id" {
    source "$(dirname "$BATS_TEST_FILENAME")/../../ralph_loop.sh"
    
    # Create test JSONL output file
    local test_file="$(mktemp /tmp/codex-test.XXXXXX.log)"
    cat > "$test_file" << 'JSONL'
{"event": "thread.started", "thread_id": "thread-abc123-def456", "timestamp": "2026-01-16T12:00:00Z"}
{"event": "message", "type": "text", "content": "Final output"}
JSONL
    
    # Parse JSONL
    run parse_codex_jsonl "$test_file"
    
    # Check thread_id extracted
    [[ "$CODEX_THREAD_ID" == "thread-abc123-def456" ]]
    assert_success
    
    # Cleanup
    rm -f "$test_file"
}

# Test: parse_codex_jsonl accumulates messages
@test "parse_codex_jsonl accumulates message content" {
    source "$(dirname "$BATS_TEST_FILENAME")/../../ralph_loop.sh"
    
    local test_file="$(mktemp /tmp/codex-test.XXXXXX.log)"
    cat > "$test_file" << 'JSONL'
{"event": "thread.started", "thread_id": "thread-xyz789", "timestamp": "2026-01-16T12:00:00Z"}
{"event": "message", "type": "text", "content": "First message."}
{"event": "message", "type": "text", "content": "Second message."}
JSONL
    
    run parse_codex_jsonl "$test_file"
    
    # Check that messages were accumulated
    [[ "$CODEX_FINAL_MESSAGE" == *"First message."* ]]
    [[ "$CODEX_FINAL_MESSAGE" == *"Second message."* ]]
    assert_success
    
    rm -f "$test_file"
}

# Test: parse_codex_jsonl handles single JSON format (fallback)
@test "parse_codex_jsonl handles single JSON" {
    source "$(dirname "$BATS_TEST_FILENAME")/../../ralph_loop.sh"
    
    local test_file="$(mktemp /tmp/codex-test.XXXXXX.log)"
    echo '{"result": "Output", "sessionId": "test-123"}' > "$test_file"
    
    run parse_codex_jsonl "$test_file"
    
    # Check that session_id was extracted from single JSON
    [[ "$CODEX_THREAD_ID" == "test-123" ]]
    assert_success
    
    rm -f "$test_file"
}

# Test: parse_codex_jsonl handles plain text (fallback)
@test "parse_codex_jsonl handles plain text" {
    source "$(dirname "$BATS_TEST_FILENAME")/../../ralph_loop.sh"
    
    local test_file="$(mktemp /tmp/codex-test.XXXXXX.log)"
    echo "Plain text output from Codex" > "$test_file"
    
    run parse_codex_jsonl "$test_file"
    
    # Check that plain text is treated as final message
    [[ "$CODEX_FINAL_MESSAGE" == "Plain text output from Codex" ]]
    assert_success
    
    rm -f "$test_file"
}

# Test: save_codex_session persists thread_id
@test "save_codex_session saves thread_id to file" {
    source "$(dirname "$BATS_TEST_FILENAME")/../../ralph_loop.sh"
    source "$(dirname "$BATS_TEST_FILENAME")/../helpers/mocks.bash"
    setup_mocks
    
    # Create mock JSONL output
    local test_file="$(mktemp /tmp/codex-test.XXXXXX.log)"
    cat > "$test_file" << 'JSONL'
{"event": "thread.started", "thread_id": "saved-thread-456", "timestamp": "2026-01-16T12:00:00Z"}
{"event": "message", "type": "text", "content": "Test output"}
JSONL
    
    # Save session
    run save_codex_session "$test_file"
    teardown_mocks
    
    # Check thread_id was saved
    if [[ -f ".codex_session_id" ]]; then
        local saved_id=$(cat ".codex_session_id")
        [[ "$saved_id" == "saved-thread-456" ]]
    else
        fail "Session file was not created"
    fi
    
    assert_success
    
    rm -f "$test_file" ".codex_session_id"
}

# Test: init_codex_session returns existing session_id
@test "init_codex_session resumes existing session" {
    source "$(dirname "$BATS_TEST_FILENAME")/../../ralph_loop.sh"
    source "$(dirname "$BATS_TEST_FILENAME")/../helpers/mocks.bash"
    setup_mocks
    
    # Create existing session file
    echo "existing-session-789" > ".codex_session_id"
    
    # Init session
    run init_codex_session
    teardown_mocks
    
    # Check that existing session was returned
    [[ "$output" == "existing-session-789" ]]
    assert_success
    
    rm -f ".codex_session_id"
}

# Test: init_codex_session starts new session when expired
@test "init_codex_session starts new session when expired" {
    source "$(dirname "$BATS_TEST_FILENAME")/../../ralph_loop.sh"
    source "$(dirname "$BATS_TEST_FILENAME")/../helpers/mocks.bash"
    setup_mocks
    
    # Create old session file (>24 hours old)
    # Use macOS compatible date command
    local old_timestamp=$(date -v -48H +"%Y%m%d%H%M.%S" 2>/dev/null || date -u -d "48 hours ago" +"%Y%m%d%H%M.%S" 2>/dev/null)
    if [[ -n "$old_timestamp" ]]; then
        touch -t "$old_timestamp" ".codex_session_id"
    fi
    echo "old-session-123" > ".codex_session_id"
    
    # Init session
    CODEX_SESSION_EXPIRY_HOURS=24
    run init_codex_session
    teardown_mocks
    
    # Check that new session was started (no output)
    [[ "$output" == "" ]]
    
    # Check old session was deleted
    [[ ! -f ".codex_session_id" ]]
    
    assert_success
}

# Test: build_codex_command adds resume flags with session_id
@test "build_codex_command uses resume flags with session_id" {
    source "$(dirname "$BATS_TEST_FILENAME")/../../ralph_loop.sh"
    
    # Create test prompt
    local test_prompt_file="$(mktemp /tmp/prompt.XXXXXX.md)"
    echo "Test prompt" > "$test_prompt_file"
    
    # Build command with session_id
    CODEX_USE_CONTINUE=true
    run build_codex_command "$test_prompt_file" "" "resume-session-xyz"
    
    # Check resume flags are present
    [[ "${CODEX_CMD_ARGS[*]}" == *"exec"* ]]
    [[ "${CODEX_CMD_ARGS[*]}" == *"--resume"* ]]
    [[ "${CODEX_CMD_ARGS[*]}" == *"--thread-id"* ]]
    [[ "${CODEX_CMD_ARGS[*]}" == *"resume-session-xyz"* ]]
    
    assert_success
    
    rm -f "$test_prompt_file"
}
