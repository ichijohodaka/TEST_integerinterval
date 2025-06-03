# Interval Algebra for Integer Half-Open Intervals

This Go package provides algebraic operations on **half-open integer intervals** and their sets.

- In Go, the slice expression `s[a:b]` represents the substring from index `a` (inclusive) to `b` (exclusive).
- Mathematically, this corresponds to the half-open interval `[a, b)` — assuming `a`, `b` ∈ ℤ.
- This package adopts that model: all intervals are treated as half-open, and operations on them are defined accordingly.
- Comments and documentation use mathematical notation wherever possible, for clarity and brevity.

---

## Notation and Style

All intervals in this package follow the **half-open interval** model: [a, b) := { x ∈ ℤ | a ≤ x < b }

We use concise mathematical expressions in function comments and documentation.  
This helps make the logic of interval operations transparent and formally verifiable.

For example:

```go
// Contains(n) ⇔ n ∈ [Start, End)
func (iv IntegerInterval) Contains(n int) bool

// Overlaps(other) ⇔ [Start, End) ∩ [other.Start, other.End) ≠ ∅
func (iv IntegerInterval) Overlaps(other IntegerInterval) bool
```

Likewise, interval sets are represented as {[a, b), [c, d), ...}
and operations are defined in terms of set theory:

- Union: A ∪ B
- Intersection: A ∩ B
- Difference (subtract): A − B
- Complement within base: base − A

## Example

Given:

```
A     = {[0,1), [2,3)}
base  = [0,3)
Aᶜ    = base − A = {[1,2)}
```

This matches Go string slicing precisely:

```
s := "abc"
s[1:2] == "b"
```



