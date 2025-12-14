# "Line" Go Generic based Validation Library

**Read this README in:**
[English](README.md) · [Русский](README.ru.md)

## Overview

This project is a lightweight, expressive validation library for Go designed to describe and execute complex validation rules in a declarative and composable way. It focuses on validating values, structures, and collections while producing structured, context-aware validation errors.

The library is built around a small set of core concepts: **constraints**, **arguments**, **validators**, and **violations**. These components work together to allow validation logic to remain readable, reusable, and easy to extend.

## Core Concepts

### Validator

The validator is the execution engine. It collects validation arguments, applies them in the correct context, and aggregates all violations into a single result. Validation is performed explicitly and always returns either `nil` or a structured list of violations.

Validation is context-aware and supports nested property paths, indexed paths, and validation groups.

### Constraints

Constraints describe *what* rule should be applied to a value. Each constraint is responsible for validating a specific type of data and reporting a violation when the rule is not satisfied.

The library provides a rich set of built-in constraints, including:

* Presence constraints (blank, not blank, nil, not nil)
* String constraints (length limits, exact length, regular expressions)
* Numeric constraints
* Comparable constraints (choices, equality, ranges)
* Date and time format validation
* Collection constraints (minimum, maximum, exact count, divisibility)
* Predicate-based constraints (JSON, integer, numeric checks)

Constraints are immutable and chainable. You can customize them with:

* Conditional execution
* Validation groups
* Custom errors
* Custom message templates
* Message parameters

### Arguments

Arguments bind values to constraints and define *where* and *how* validation happens. An argument represents a single validation unit that can be attached to a property path.

The library provides argument helpers for:

* Scalars (strings, numbers, booleans, time values)
* Nullable values
* Collections and iterables
* Struct-like objects that implement validation
* Per-element validation for slices
* Conditional validation flows

Arguments can be nested, combined, and reused across validators.

### Property Paths

Every violation is associated with a property path that precisely identifies where the error occurred. Paths support:

* Dot notation for properties
* Indexed access for arrays and slices
* Escaped and quoted property names

This makes the output suitable for APIs, form validation, and structured error reporting.

### Violations and Errors

When a constraint fails, a violation is created. A violation contains:

* A machine-readable error identifier
* A human-readable message
* A fully resolved property path
* Rendered message parameters

Multiple violations are collected into a single error object, allowing consumers to handle all validation issues at once.

Message templates support placeholders that are replaced at runtime with actual values, limits, or contextual data.

## Conditional and Group-Based Validation

Validation can be controlled dynamically using conditions and groups:

* Conditional execution based on runtime boolean expressions
* Group-based validation to enable or disable rules depending on context
* Branching logic with “then / else” semantics
* Sequential validation that stops on first failure
* Parallel validation for independent rules
* Logical composition such as “all” or “at least one” rules

This makes the library suitable for complex business rules and multi-step validation workflows.

## Extensibility

The library is designed to be extended without modifying its core:

* Custom constraints can be implemented by satisfying the appropriate constraint interfaces
* Custom predicates can be wrapped into reusable constraints
* Custom error types and message templates can be introduced easily
* Domain objects can become self-validating by implementing a single interface

## Design Goals

* Clear separation between validation definition and execution
* Strong typing without sacrificing flexibility
* Minimal reflection and predictable runtime behavior
* Structured error reporting suitable for APIs and user interfaces
* Composable building blocks instead of monolithic validation rules

## Typical Use Cases

* Validating request payloads in HTTP APIs
* Domain-level invariants in business logic
* Configuration and input validation
* Nested object and collection validation
* Conditional validation scenarios

## Summary

This library provides a declarative, composable approach to validation in Go. By modeling validation as a combination of constraints and arguments executed by a validator, it enables expressive validation rules, precise error reporting, and flexible control flow without sacrificing clarity or type safety.
