#
# Ultra-fast statusline script for Claude Code hooks (PowerShell)
# Extracts git branch from JSON input in <10ms
#
# Requirements:
# - PowerShell 5.1+ compatible
# - Reads Claude Code JSON from stdin
# - Outputs git branch name or empty string
# - Performance: <10ms execution time
#
# Fast Path: Direct .git/HEAD file read (fastest)
# Fallback: git symbolic-ref command (slower but reliable)
#
# [ref: PRD; lines: 213-223]
# [ref: SDD; lines: 799-820, 235-240, 1153]

# Set strict mode for better error handling
$ErrorActionPreference = 'SilentlyContinue'

# Parse JSON from stdin
function Get-JsonCwd {
    param([string]$JsonInput)

    try {
        $parsed = $JsonInput | ConvertFrom-Json
        if ($parsed.cwd) {
            return $parsed.cwd
        }
    }
    catch {
        # JSON parsing failed, return empty
    }

    return ""
}

# Expand tilde to home directory (PowerShell equivalent)
function Expand-TildePath {
    param([string]$Path)

    if ($Path -match '^~') {
        return $Path -replace '^~', $env:USERPROFILE
    }

    return $Path
}

# Fast path: Read .git/HEAD directly
function Get-BranchFast {
    param([string]$Directory)

    $gitHeadPath = Join-Path $Directory '.git\HEAD'

    if (Test-Path $gitHeadPath -PathType Leaf) {
        try {
            $headContent = Get-Content $gitHeadPath -Raw -ErrorAction Stop

            # Extract branch from "ref: refs/heads/branch-name"
            if ($headContent -match '^ref:\s*refs/heads/(.+)$') {
                return $Matches[1].Trim()
            }
        }
        catch {
            # File read error
        }
    }

    return $null
}

# Fallback: Use git command
function Get-BranchFallback {
    param([string]$Directory)

    # Check if git command exists
    $gitExists = Get-Command git -ErrorAction SilentlyContinue
    if (-not $gitExists) {
        return $null
    }

    # Check if directory is a git repo
    $gitDir = Join-Path $Directory '.git'
    if (-not (Test-Path $gitDir -PathType Container)) {
        return $null
    }

    try {
        Push-Location $Directory
        $branch = git symbolic-ref --short HEAD 2>$null
        Pop-Location

        if ($branch) {
            return $branch.Trim()
        }
    }
    catch {
        Pop-Location
    }

    return $null
}

# Main execution
function Main {
    # Read JSON from stdin (all at once for performance)
    $jsonInput = [Console]::In.ReadToEnd()

    # Extract cwd from JSON
    $cwd = Get-JsonCwd -JsonInput $jsonInput

    # Handle empty cwd (use current directory as fallback)
    if ([string]::IsNullOrEmpty($cwd)) {
        $cwd = Get-Location
    }

    # Expand tilde in path
    $cwd = Expand-TildePath -Path $cwd

    # Verify directory exists
    if (-not (Test-Path $cwd -PathType Container)) {
        Write-Output ""
        exit 0
    }

    # Try fast path first
    $branch = Get-BranchFast -Directory $cwd
    if ($branch) {
        Write-Output $branch
        exit 0
    }

    # Try fallback with git command
    $branch = Get-BranchFallback -Directory $cwd
    if ($branch) {
        Write-Output $branch
        exit 0
    }

    # No git repo or error - return empty string
    Write-Output ""
    exit 0
}

# Execute main function
Main
