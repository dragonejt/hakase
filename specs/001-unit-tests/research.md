# Research Findings: Go Unit Testing Implementation

## Go Mocking Rules and Best Practices

### Key Question: Can you mock an interface in an external package in Golang?

**Answer**: YES, with specific rules and patterns

### Detailed Findings

#### 1. Interface Mocking in Go

**Core Rule**: You CAN mock interfaces from external packages, but you CANNOT mock concrete types (structs) from external packages.

**Why This Works**:
- Interfaces in Go are implicit and can be implemented by any type
- Mocking creates a new type that implements the same interface
- External package interfaces are just contracts that can be satisfied

**Example Pattern (from existing guild_test.go)**:
```go
type MockHakaseClient struct {
    clients.HakaseClient  // Embed the interface from external package
    mock.Mock             // Embed testify mock for tracking
}

func (m *MockHakaseClient) CreateCourse(span *sentry.Span, course clients.Course) error {
    m.Called(span, course)  // Track the call
    return nil              // Return mock response
}
```

#### 2. Mocking External Package Interfaces

**For discordgo.Session (external package)**:
```go
type MockDiscordSession struct {
    discordgo.Session  // Embed the interface
    mock.Mock          // Embed mock functionality
}

func (m *MockDiscordSession) ChannelMessageSend(channelID string, content string) (*discordgo.Message, error) {
    args := m.Called(channelID, content)
    return args.Get(0).(*discordgo.Message), args.Error(1)
}
```

**Key Points**:
- The mock struct must be in the same package as the test
- Embed both the interface and mock.Mock
- Implement only the methods you need to mock
- Use `m.Called()` to track method calls
- Use `args.Get()` and `args.Error()` to return mock values

#### 3. Best Practices for This Project

**Follow Existing Patterns**:
- Use testify suite for test organization (already established)
- Use testify mock for mocking (already available)
- Keep test files colocated with source files
- Use table-driven tests for multiple scenarios

**Test Organization**:
```go
type ExampleTestSuite struct {
    suite.Suite
    mockDiscord *MockDiscordSession
    mockClient *MockHakaseClient
}

func (s *ExampleTestSuite) SetupTest() {
    s.mockDiscord = new(MockDiscordSession)
    s.mockClient = new(MockHakaseClient)
    // Setup mock expectations
}

func (s *ExampleTestSuite) TestFunction() {
    // Test implementation
    s.mockDiscord.AssertExpectations(s.T())
}
```

#### 4. Limitations and Workarounds

**Cannot Mock Concrete Types**:
- If you need to mock a struct from an external package, you must:
  1. Create an interface that describes the methods you need
  2. Have your code depend on the interface, not the concrete type
  3. Mock the interface

**Example Workaround**:
```go
// In your code
type DiscordSender interface {
    ChannelMessageSend(channelID string, content string) (*discordgo.Message, error)
}

// Use dependency injection
type MyService struct {
    discord DiscordSender  // Depend on interface, not concrete type
}

// In tests
type MockDiscordSender struct {
    mock.Mock
}

func (m *MockDiscordSender) ChannelMessageSend(channelID string, content string) (*discordgo.Message, error) {
    args := m.Called(channelID, content)
    return args.Get(0).(*discordgo.Message), args.Error(1)
}
```

#### 5. Testing Recommendations

**Test Coverage Targets**:
- Aim for 80%+ overall coverage
- Focus on critical paths first
- Test both happy paths and error conditions

**Test Execution**:
- Keep tests fast (<5 minutes total)
- Use `go test ./...` for running all tests
- Use `-coverprofile` for coverage analysis
- Use `-race` for race condition detection

**CI/CD Integration**:
- Run tests on every commit
- Block merges on test failures
- Report coverage metrics
- Parallelize tests where possible

## Decision Summary

**Decision**: Use testify mock package with interface embedding pattern
**Rationale**: 
- Already available in project dependencies
- Follows existing test patterns
- Provides comprehensive mocking capabilities
- Well-documented and widely used

**Alternatives Considered**:
- Manual mocks: More boilerplate, harder to maintain
- Other mocking libraries: Would require new dependencies (violates constraints)
- No mocking: Would make testing external dependencies difficult

## Implementation Checklist

- [ ] Create mock structs for all external interfaces needed
- [ ] Follow testify suite pattern for test organization
- [ ] Use interface embedding for external package mocks
- [ ] Implement table-driven tests for multiple scenarios
- [ ] Add both happy path and error condition tests
- [ ] Achieve 80%+ code coverage
- [ ] Keep test execution under 5 minutes
- [ ] Integrate with CI/CD pipeline
