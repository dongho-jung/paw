// Package service provides business logic services for PAW.
package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/dongho-jung/paw/internal/logging"
)

// CommandResult captures the result of running a shell command.
type CommandResult struct {
	Command    string
	Success    bool
	ExitCode   int
	Duration   time.Duration
	Output     string
	TimeoutHit bool
}

// RunCommand executes a shell command with timeout and captures output.
func RunCommand(command string, workDir string, env []string, timeout time.Duration) (*CommandResult, error) {
	start := time.Now()
	ctx := context.Background()
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	if workDir != "" {
		cmd.Dir = workDir
	}
	if len(env) > 0 {
		cmd.Env = env
	}

	output, err := cmd.CombinedOutput()
	result := &CommandResult{
		Command:  command,
		Duration: time.Since(start),
		Output:   string(output),
	}

	if err == nil {
		result.Success = true
		result.ExitCode = 0
		return result, nil
	}

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		result.TimeoutHit = true
		result.ExitCode = -1
		result.Success = false
		return result, fmt.Errorf("command timed out after %s", timeout)
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		result.ExitCode = exitErr.ExitCode()
		result.Success = false
		return result, err
	}

	result.ExitCode = -1
	result.Success = false
	return result, err
}

// RunHook executes a named hook command, writing output and metadata.
func RunHook(name, command, workDir string, env []string, outputPath, metaPath string, timeout time.Duration) (*HookMetadata, error) {
	logging.Debug("Hook: start name=%s cmd=%s", name, command)
	result, err := RunCommand(command, workDir, env, timeout)

	status := "success"
	if result.TimeoutHit {
		status = "timeout"
	} else if !result.Success {
		status = "failed"
	}

	meta := &HookMetadata{
		Name:       name,
		Command:    command,
		Status:     status,
		ExitCode:   result.ExitCode,
		DurationMs: result.Duration.Milliseconds(),
		OutputFile: outputPath,
	}

	if outputPath != "" {
		if writeErr := os.WriteFile(outputPath, []byte(result.Output), 0644); writeErr != nil {
			logging.Warn("Hook: failed to write output name=%s err=%v", name, writeErr)
		}
	}

	if metaPath != "" {
		if data, marshalErr := json.MarshalIndent(meta, "", "  "); marshalErr == nil {
			if writeErr := os.WriteFile(metaPath, data, 0644); writeErr != nil {
				logging.Warn("Hook: failed to write metadata name=%s err=%v", name, writeErr)
			}
		} else {
			logging.Warn("Hook: failed to marshal metadata name=%s err=%v", name, marshalErr)
		}
	}

	if result.Success {
		logging.Debug("Hook: success name=%s duration=%s", name, result.Duration)
	} else {
		logging.Warn("Hook: failed name=%s status=%s err=%v", name, status, err)
	}

	return meta, err
}

// RunVerification executes the verification command and writes output metadata.
func RunVerification(command, workDir string, env []string, outputPath, metaPath string, timeout time.Duration) (*VerificationMetadata, error) {
	logging.Debug("Verify: start cmd=%s", command)
	result, err := RunCommand(command, workDir, env, timeout)

	meta := &VerificationMetadata{
		Command:    command,
		Success:    result.Success,
		ExitCode:   result.ExitCode,
		DurationMs: result.Duration.Milliseconds(),
		OutputFile: outputPath,
	}

	if outputPath != "" {
		if writeErr := os.WriteFile(outputPath, []byte(result.Output), 0644); writeErr != nil {
			logging.Warn("Verify: failed to write output err=%v", writeErr)
		}
	}

	if metaPath != "" {
		if data, marshalErr := json.MarshalIndent(meta, "", "  "); marshalErr == nil {
			if writeErr := os.WriteFile(metaPath, data, 0644); writeErr != nil {
				logging.Warn("Verify: failed to write metadata err=%v", writeErr)
			}
		} else {
			logging.Warn("Verify: failed to marshal metadata err=%v", marshalErr)
		}
	}

	if result.Success {
		logging.Debug("Verify: success duration=%s", result.Duration)
	} else {
		logging.Warn("Verify: failed exit=%d err=%v", result.ExitCode, err)
	}

	return meta, err
}
