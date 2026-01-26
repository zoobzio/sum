# Testing

This directory contains testing infrastructure for sum.

## Structure

```
testing/
├── README.md           # This file
├── helpers.go          # Domain-specific test helpers
├── helpers_test.go     # Tests for helpers themselves
├── integration/        # Integration tests
│   └── README.md
└── benchmarks/         # Performance benchmarks
    └── README.md
```

## Running Tests

```bash
# All tests
make test

# Unit tests only (short mode)
make test-unit

# Integration tests
make test-integration

# Benchmarks
make test-bench
```

## Writing Tests

- Place unit tests alongside source files (`foo.go` -> `foo_test.go`)
- Place integration tests in `testing/integration/`
- Place benchmarks in `testing/benchmarks/`
- Use helpers from this package for common test operations

## Coverage

Coverage targets:
- Project: 70% minimum
- New code (patch): 80% minimum

Generate coverage report:
```bash
make coverage
```
