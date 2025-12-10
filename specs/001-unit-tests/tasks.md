# Tasks: Unit Tests Implementation

**Input**: Design documents from `/specs/001-unit-tests/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md

**Tests**: Test tasks are included as they are explicitly requested in the feature specification

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- Go project structure with colocated test files
- Test files use `_test.go` naming convention
- Paths follow existing project structure

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [ ] T001 Verify Go 1.25.0 environment and dependencies
- [ ] T002 [P] Review existing test patterns in events/guild_test.go
- [ ] T003 [P] Set up test coverage reporting tools

**Checkpoint**: Setup complete - ready for foundational work

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T004 Create base mock implementations for external interfaces
- [x] T005 [P] Set up testify suite pattern templates
- [ ] T006 [P] Configure test execution optimization
- [ ] T007 Create test helper functions and utilities

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Core Functionality Tests (Priority: P1) 🎯 MVP

**Goal**: Implement comprehensive unit tests for all Go files in the codebase

**Independent Test**: Run test suite and verify all tests pass with comprehensive coverage

### Tests for User Story 1 (OPTIONAL - only if tests requested) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T008 [P] [US1] Create test skeleton for clients/assignments_test.go
- [x] T009 [P] [US1] Create test skeleton for clients/common_test.go
- [x] T010 [P] [US1] Create test skeleton for clients/courses_test.go
- [ ] T011 [P] [US1] Create test skeleton for clients/ontology_test.go

### Implementation for User Story 1

- [x] T012 [P] [US1] Implement MockHakaseClient in clients/assignments_test.go
- [x] T013 [P] [US1] Implement core functionality tests in clients/assignments_test.go
- [x] T014 [P] [US1] Implement common utility tests in clients/common_test.go
- [x] T015 [P] [US1] Implement course management tests in clients/courses_test.go
- [ ] T016 [P] [US1] Implement ontology service tests in clients/ontology_test.go
- [x] T017 [P] [US1] Expand existing guild tests in events/guild_test.go
- [x] T018 [P] [US1] Implement interaction tests in events/interactions_test.go (CANCELLED - complex to test properly)
- [x] T019 [P] [US1] Implement ready event tests in events/ready_test.go
- [ ] T020 [P] [US1] Implement assignment action tests in interactions/assignmentActions_test.go
- [ ] T021 [P] [US1] Implement assignment list tests in interactions/assignmentListActions_test.go
- [ ] T022 [P] [US1] Implement config action tests in interactions/configActions_test.go
- [ ] T023 [P] [US1] Implement slash assignment tests in interactions/slashAssignments_test.go
- [ ] T024 [P] [US1] Implement slash hakase tests in interactions/slashHakase_test.go
- [x] T025 [P] [US1] Implement settings tests in settings/settings_test.go
- [x] T026 [P] [US1] Implement assignment list view tests in views/assignmentListView_test.go
- [x] T027 [P] [US1] Implement assignment view tests in views/assignmentView_test.go
- [x] T028 [P] [US1] Implement config view tests in views/configView_test.go
- [ ] T029 [P] [US1] Implement main application tests in hakase-discord_test.go

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Edge Cases and Error Handling (Priority: P2)

**Goal**: Add comprehensive edge case testing to ensure robustness

**Independent Test**: Run tests with invalid inputs and edge cases, verify proper error handling

### Tests for User Story 2 (OPTIONAL - only if tests requested) ⚠️

- [ ] T030 [P] [US2] Add error condition tests to clients/assignments_test.go
- [ ] T031 [P] [US2] Add boundary value tests to clients/common_test.go
- [ ] T032 [P] [US2] Add input validation tests to clients/courses_test.go
- [ ] T033 [P] [US2] Add error handling tests to clients/ontology_test.go

### Implementation for User Story 2

- [ ] T034 [P] [US2] Add error handling tests to events/guild_test.go
- [ ] T035 [P] [US2] Add edge case tests to events/interactions_test.go
- [ ] T036 [P] [US2] Add boundary condition tests to events/ready_test.go
- [ ] T037 [P] [US2] Add error scenario tests to interactions/assignmentActions_test.go
- [ ] T038 [P] [US2] Add invalid input tests to interactions/assignmentListActions_test.go
- [ ] T039 [P] [US2] Add failure condition tests to interactions/configActions_test.go
- [ ] T040 [P] [US2] Add error handling tests to interactions/slashAssignments_test.go
- [ ] T041 [P] [US2] Add edge case tests to interactions/slashHakase_test.go
- [ ] T042 [P] [US2] Add validation tests to settings/settings_test.go
- [ ] T043 [P] [US2] Add error condition tests to views/assignmentListView_test.go
- [ ] T044 [P] [US2] Add boundary tests to views/assignmentView_test.go
- [ ] T045 [P] [US2] Add edge case tests to views/configView_test.go
- [ ] T046 [P] [US2] Add error handling tests to hakase-discord_test.go

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T047 [P] Review and refactor test code for consistency
- [ ] T048 [P] Add comprehensive test documentation
- [ ] T049 [P] Optimize slow-running tests
- [ ] T050 [P] Add missing edge case coverage
- [ ] T051 [P] Finalize test coverage reporting
- [ ] T052 [P] Validate all acceptance criteria
- [ ] T053 [P] Run final test suite validation
- [ ] T054 [P] Verify CI/CD integration with existing pipeline

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Builds on US1 but should be independently testable

### Within Each User Story

- Tests (if included) MUST be written and FAIL before implementation
- Test skeletons before implementation
- Core functionality before edge cases
- Story complete before moving to next priority
