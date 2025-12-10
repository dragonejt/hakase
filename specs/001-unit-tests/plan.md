# Implementation Plan: Unit Tests Implementation

**Branch**: `001-unit-tests` | **Date**: 2025-12-10 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-unit-tests/spec.md`

## Summary

Implement comprehensive unit tests for all Go files in the codebase using existing dependencies (testify) and following Go best practices. The implementation will cover core functionality, edge cases, and CI/CD integration while maintaining 80%+ code coverage.

## Technical Context

**Language/Version**: Go 1.25.0
**Primary Dependencies**: 
- github.com/stretchr/testify v1.11.1 (mocking and test suite support)
- Standard library testing package
**Testing**: testify suite and mock packages, standard Go testing
**Target Platform**: Cross-platform (Discord bot application)
**Project Type**: Single Go application with multiple packages
**Performance Goals**: Test execution under 5 minutes for full suite
**Constraints**: No new dependencies, follow existing test patterns
**Scale/Scope**: 15+ Go files across 5 packages needing unit tests

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

✅ **Testing First**: This plan implements comprehensive unit testing as required
✅ **No New Dependencies**: Using only existing testify dependency
✅ **Best Practices**: Following Go testing conventions and existing patterns
✅ **Integration Testing**: CI/CD pipeline integration included

## Project Structure

### Documentation (this feature)

```text
specs/001-unit-tests/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
.
├── clients/
│   ├── assignments.go    # Needs tests
│   ├── assignments_test.go  # New
│   ├── common.go         # Needs tests  
│   ├── common_test.go    # New
│   ├── courses.go        # Needs tests
│   ├── courses_test.go   # New
│   └── ontology.go       # Needs tests
│   └── ontology_test.go  # New
├── events/
│   ├── guild.go          # Needs tests
│   ├── guild_test.go     # Existing (expand)
│   ├── interactions.go   # Needs tests
│   ├── interactions_test.go # New
│   └── ready.go          # Needs tests
│   └── ready_test.go     # New
├── interactions/
│   ├── assignmentActions.go       # Needs tests
│   ├── assignmentActions_test.go  # New
│   ├── assignmentListActions.go   # Needs tests
│   ├── assignmentListActions_test.go # New
│   ├── configActions.go           # Needs tests
│   ├── configActions_test.go      # New
│   ├── slashAssignments.go        # Needs tests
│   ├── slashAssignments_test.go   # New
│   ├── slashHakase.go             # Needs tests
│   └── slashHakase_test.go        # New
├── settings/
│   ├── settings.go       # Needs tests
│   └── settings_test.go  # New
├── views/
│   ├── assignmentListView.go  # Needs tests
│   ├── assignmentListView_test.go # New
│   ├── assignmentView.go        # Needs tests
│   ├── assignmentView_test.go   # New
│   ├── configView.go            # Needs tests
│   └── configView_test.go       # New
└── hakase-discord.go     # Needs tests
└── hakase-discord_test.go # New
```

**Structure Decision**: Following existing Go project structure with test files colocated with source files using `_test.go` naming convention.

## Complexity Tracking

No constitution violations detected. All requirements align with existing project structure and best practices.

## Research Findings

### Go Mocking Rules and Best Practices

**Q: Can you mock an interface in an external package in Golang?**
**A: Yes, but with specific rules:**

1. **Interface Mocking Rules**:
   - You CAN mock interfaces from external packages
   - You CANNOT mock concrete types (structs) from external packages
   - Mocks must be created in the same package as the test
   - Use testify's mock package for interface mocking

2. **Best Practices for This Project**:
   - Follow existing pattern from `guild_test.go`: Create mock structs that embed both the interface and `mock.Mock`
   - Use `testify/suite` for test organization (already established pattern)
   - Keep test files colocated with source files
   - Use table-driven tests for multiple scenarios

3. **Existing Pattern Analysis**:
   - Current `guild_test.go` uses `MockHakaseClient` struct embedding `clients.HakaseClient` interface
   - Uses `testify/mock` package for method call tracking
   - Uses `testify/suite` for test suite organization
   - This pattern should be replicated across all test files

## Implementation Strategy

### Phase 1: Core Functionality Tests (P1)

**Goal**: Implement basic unit tests for all Go files
**Approach**:
1. Create test files for each Go source file
2. Use testify suite pattern for organization
3. Test main functionality paths
4. Achieve basic code coverage

**Files to Create/Update**:
- All `_test.go` files listed in structure above
- Expand existing `guild_test.go` with more test cases

### Phase 2: Edge Case and Error Handling Tests (P2)

**Goal**: Add comprehensive edge case testing
**Approach**:
1. Add test cases for error conditions
2. Test boundary values and invalid inputs
3. Test error handling paths
4. Increase coverage to 80%+

### Phase 3: CI/CD Integration (P3)

**Goal**: Integrate tests into CI/CD pipeline
**Approach**:
1. Ensure tests run on every commit
2. Configure test failure to block merges
3. Add coverage reporting
4. Optimize test execution time

## Testing Approach Details

### Mocking Strategy

For external package interfaces (like `discordgo.Session`):
```go
type MockDiscordSession struct {
    discordgo.Session  // Embed the real interface
    mock.Mock          // Embed testify mock
}

func (m *MockDiscordSession) ChannelMessageSend(channelID string, content string) (*discordgo.Message, error) {
    args := m.Called(channelID, content)
    return args.Get(0).(*discordgo.Message), args.Error(1)
}
```

### Test Organization

Use testify suite pattern:
```go
type ExampleTestSuite struct {
    suite.Suite
    // Test dependencies and mocks
}

func TestExample(t *testing.T) {
    suite.Run(t, new(ExampleTestSuite))
}

func (s *ExampleTestSuite) SetupTest() {
    // Initialize mocks and test data
}

func (s *ExampleTestSuite) TestFunctionName() {
    // Test implementation
}
```

### Code Coverage

Run coverage analysis:
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

Target: 80%+ coverage across all packages

## Risk Assessment

**Low Risk**: 
- Using existing, proven testing libraries
- Following established patterns from existing tests
- No new dependencies required
- Go's testing ecosystem is mature and stable

**Mitigation Strategies**:
- Start with core functionality tests before edge cases
- Use table-driven tests to reduce boilerplate
- Keep test execution fast (<5 minutes target)
- Document test patterns for consistency

## Success Metrics

- ✅ All Go files have corresponding test files
- ✅ 80%+ code coverage achieved
- ✅ All tests pass in local environment
- ✅ CI/CD pipeline runs tests automatically
- ✅ Test execution time < 5 minutes
- ✅ Edge cases properly tested
- ✅ Documentation complete

## Next Steps

1. **Phase 0 Complete**: Research and planning done
2. **Phase 1 Ready**: Begin implementation with `/speckit.tasks`
3. **Recommendation**: Proceed to task breakdown and implementation
