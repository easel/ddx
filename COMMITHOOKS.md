# Git Hooks Framework Analysis and Recommendations

## Evaluation Criteria (from prompts/common/commit-hooks.md)
- Broken into individual units
- Minimal dependencies
- Cross platform (Mac, Linux, Windows)
- Fast execution, optimized for staged files only
- Follow best practices for 2025
- Native tools preferred (C, C++, Rust, Go)
- Support for: formatting, types, compilation, unit tests, integration tests, security checks

## Framework Comparison

### 1. Lefthook (RECOMMENDED) ⭐
**Written in:** Go  
**Dependencies:** None (single binary)  
**Cross-platform:** Excellent (native Windows, Mac, Linux support)  
**Performance:** ~5x faster than pre-commit (1s vs 5s in benchmarks)  

**Pros:**
- Native Go implementation aligns with DDx project
- Parallel execution by default
- Zero runtime dependencies
- Built-in support for `{staged_files}` variable
- Simple YAML configuration
- Can integrate custom scripts
- Actively maintained with strong community

**Cons:**
- Smaller ecosystem than pre-commit
- Less extensive documentation

**Example Configuration:**
```yaml
pre-commit:
  parallel: true
  commands:
    format:
      glob: "*.go"
      run: gofmt -w {staged_files}
      stage_fixed: true
    lint:
      glob: "*.go"
      run: golangci-lint run --fix
      stage_fixed: true
    test:
      glob: "*.go"
      run: go test -short ./...
    secrets:
      run: ./scripts/hooks/checks/check-secrets.sh
    conflicts:
      run: ./scripts/hooks/checks/check-conflicts.sh
```

### 2. Pre-commit
**Written in:** Python  
**Dependencies:** Python runtime required  
**Cross-platform:** Good (requires Python on all platforms)  
**Performance:** Slower (sequential execution by default)  

**Pros:**
- Massive ecosystem of pre-built hooks
- Extensive documentation
- Industry standard with wide adoption
- Remote hooks capability
- Multi-language support

**Cons:**
- Requires Python runtime (not native)
- Sequential execution impacts performance
- Heavier resource usage
- Not aligned with "minimal dependencies" criterion

### 3. Prek (Future Option)
**Written in:** Rust  
**Dependencies:** None (single binary)  
**Cross-platform:** Excellent  
**Performance:** ~10x faster than pre-commit  

**Pros:**
- Drop-in replacement for pre-commit
- Excellent performance
- Single binary, no dependencies
- Compatible with pre-commit configs

**Cons:**
- Not production-ready yet
- Limited documentation
- Small community

### 4. Custom Bash Scripts (Current Implementation)
**Written in:** Bash  
**Dependencies:** Bash, standard Unix tools  
**Cross-platform:** Poor (doesn't work natively on Windows)  
**Performance:** Good (optimized for staged files)  

**Pros:**
- Full control and customization
- Already implemented and working
- No external dependencies
- Modular design

**Cons:**
- Not cross-platform (Windows requires WSL/Git Bash)
- Maintenance burden
- No community support
- Reinventing the wheel

## Recommendation: Hybrid Approach with Lefthook

### Why Lefthook?
1. **Native Go** - Aligns perfectly with DDx being a Go project
2. **Zero dependencies** - Single binary distribution
3. **True cross-platform** - Works natively on Windows without WSL
4. **Performance** - 5x faster with parallel execution
5. **Flexibility** - Can call custom scripts where needed

### Implementation Strategy

#### Phase 1: Core Lefthook Setup
- Install lefthook as the primary hook manager
- Configure basic Go-specific hooks (format, lint, test)
- Maintain parallel execution for speed

#### Phase 2: Custom Script Integration
- Keep modular bash scripts in `scripts/hooks/checks/`
- Call them from lefthook for specialized checks
- Provides fallback and customization options

#### Phase 3: Cross-Platform Enhancements
- For Windows compatibility, consider rewriting critical bash scripts in Go
- Leverage lefthook's native Windows support
- Ensure all team members can contribute regardless of OS

### Sample Lefthook Configuration for DDx

```yaml
# lefthook.yml
min_version: 1.5.0

pre-commit:
  parallel: true
  commands:
    # Go-specific checks
    go-format:
      glob: "*.go"
      run: gofmt -l -w {staged_files}
      stage_fixed: true
    
    go-imports:
      glob: "*.go"
      run: goimports -l -w {staged_files}
      stage_fixed: true
    
    go-lint:
      glob: "*.go"
      run: golangci-lint run --fix --timeout=2m
      stage_fixed: true
    
    # Security and quality checks
    secrets-check:
      run: |
        if command -v gitleaks &> /dev/null; then
          gitleaks detect --source . --verbose
        else
          ./scripts/hooks/checks/check-secrets.sh
        fi
    
    conflicts-check:
      glob: "*"
      run: grep -H -n "^<<<<<<< \|^=======$\|^>>>>>>> " {staged_files} && exit 1 || exit 0
    
    debug-check:
      glob: "*.go"
      run: |
        grep -H -n "fmt\.Printf\|fmt\.Println\|log\.Printf\|log\.Println" {staged_files} | 
        grep -v "^[^:]*_test\.go:" && exit 1 || exit 0

pre-push:
  parallel: true
  commands:
    go-build:
      root: "cli/"
      run: go build ./...
    
    go-test:
      root: "cli/"
      run: go test -race -timeout=5m ./...
    
    go-mod-tidy:
      root: "cli/"
      run: |
        cp go.mod go.mod.backup
        cp go.sum go.sum.backup
        go mod tidy
        if ! diff -q go.mod go.mod.backup || ! diff -q go.sum go.sum.backup; then
          mv go.mod.backup go.mod
          mv go.sum.backup go.sum
          echo "Please run 'go mod tidy' before pushing"
          exit 1
        fi
        rm -f go.mod.backup go.sum.backup

# Skip hooks with: LEFTHOOK_EXCLUDE=go-test,go-lint git commit
# Or: git commit --no-verify
```

### Migration Path

1. **Keep existing scripts** - Don't delete the modular bash scripts yet
2. **Install lefthook** - `go install github.com/evilmartians/lefthook@latest`
3. **Add configuration** - Create lefthook.yml with basic hooks
4. **Test thoroughly** - Ensure all checks work on all platforms
5. **Document** - Update README with setup instructions
6. **Gradual rollout** - Team members can opt-in initially

### Alternative Tools to Consider

For specific checks, consider these native tools:
- **gitleaks** (Go) - Secret detection
- **typos** (Rust) - Spell checking
- **golangci-lint** (Go) - Comprehensive Go linting
- **gofumpt** (Go) - Stricter gofmt
- **gosec** (Go) - Security analysis

## Conclusion

Lefthook provides the best balance of performance, cross-platform support, and minimal dependencies while meeting all specified criteria. Its Go implementation makes it a natural fit for the DDx project, and its ability to integrate custom scripts provides flexibility for specialized checks.

The hybrid approach allows us to leverage the best of both worlds: a robust, fast framework with the ability to customize as needed.