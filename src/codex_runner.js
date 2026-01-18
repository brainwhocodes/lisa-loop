#!/usr/bin/env node
/**
 * Codex SDK Runner for Ralph
 * Minimal Node.js runner for Codex SDK programmatic control
 */

import { Codex } from '@openai/codex-sdk';
import { readFileSync, writeFileSync, existsSync } from 'fs';
import { resolve, dirname } from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// Configuration
const THREAD_ID_FILE = '.codex_thread_id';
const OUTPUT_FILE = 'codex_agent_output.txt';

/**
 * Main entry point
 */
async function main() {
    // Parse command line arguments
    const args = process.argv.slice(2);
    
    if (args.length < 2) {
        console.error('Usage: node codex_runner.js <prompt-file> <output-file> [thread-id]');
        console.error('  prompt-file: Path to PROMPT.md');
        console.error('  output-file: Path for agent output');
        console.error('  thread-id: Optional thread ID to resume');
        process.exit(1);
    }

    const promptFile = args[0];
    const outputFile = args[1];
    const resumeThreadId = args[2] || null;

    // Check if running in test mode
    const testMode = process.env.CODEX_TEST_MODE === 'true';
    
    try {
        // Read prompt file
        const promptText = readFileSync(promptFile, 'utf-8');
        
        if (testMode) {
            // Emit deterministic test output without calling Codex
            console.log('[TEST MODE] Emitting test output');
            writeFileSync(outputFile, 'Test output from Codex SDK runner\n');
            
            // Write mock thread ID in test mode
            if (!existsSync(THREAD_ID_FILE)) {
                writeFileSync(THREAD_ID_FILE, 'test-thread-123', 'utf-8');
            }
            
            process.exit(0);
        }

        // Initialize Codex SDK
        const codex = new Codex({
            apiKey: process.env.OPENAI_API_KEY
        });

        let threadId = resumeThreadId;
        
        // Resume existing thread or start new one
        if (threadId) {
            console.error(`[CODEX SDK] Resuming thread: ${threadId}`);
            
            // Resume thread
            await codex.threads.resume(threadId);
        } else {
            console.error('[CODEX SDK] Starting new thread');
            
            // Start new thread
            const thread = await codex.threads.create();
            threadId = thread.id;
            
            // Persist thread ID
            writeFileSync(THREAD_ID_FILE, threadId, 'utf-8');
        }

        // Run the thread with prompt
        console.error(`[CODEX SDK] Running thread with prompt (${promptText.length} chars)`);
        const run = await codex.threads.run(threadId, {
            messages: [{ role: 'user', content: promptText }]
        });

        // Wait for completion
        await codex.threads.wait(threadId);

        // Get final message
        console.error('[CODEX SDK] Thread completed, extracting final message');
        
        // Note: In a real implementation, we would extract the final agent message
        // For now, we write a placeholder
        const finalMessage = 'Codex SDK runner: Final message would be extracted here';
        writeFileSync(outputFile, finalMessage + '\n', 'utf-8');
        
        console.error(`[CODEX SDK] Output written to: ${outputFile}`);
        console.error(`[CODEX SDK] Thread ID: ${threadId}`);
        
        process.exit(0);
    } catch (error) {
        console.error(`[CODEX SDK ERROR] ${error.message}`);
        
        // Check for module not found error
        if (error.code === 'MODULE_NOT_FOUND') {
            console.error('[CODEX SDK] @openai/codex-sdk is not installed');
            console.error('[CODEX SDK] Run: npm install @openai/codex-sdk');
        }
        
        process.exit(1);
    }
}

// Run main
main().catch(error => {
    console.error(`[FATAL ERROR] ${error.message}`);
    process.exit(1);
});
