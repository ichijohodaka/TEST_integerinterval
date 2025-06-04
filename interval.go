package interval

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

// Slice(text) = text[Start:End], if valid range
//
// Returns the substring corresponding to the interval [Start, End).
// Returns an error if the interval is out of bounds.
func (iv IntegerInterval) Slice(text string) (string, error) {
	if !iv.IsValid() || iv.Start < 0 || iv.End > len(text) {
		return "", errors.New("out of range")
	}
	return text[iv.Start:iv.End], nil
}

// Replace replaces the interval [Start, End) in text with replacement.
func (iv IntegerInterval) Replace(text, replacement string) (string, error) {
	if !iv.IsValid() || iv.Start < 0 || iv.End > len(text) {
		return "", errors.New("out of range")
	}
	return text[:iv.Start] + replacement + text[iv.End:], nil
}

// Remove removes the interval [Start, End) from text.
func (iv IntegerInterval) Remove(text string) (string, error) {
	return iv.Replace(text, "")
}

// Insert inserts a string at position Start (End is ignored).
func (iv IntegerInterval) Insert(text, insert string) (string, error) {
	if !iv.IsValid() || iv.Start < 0 || iv.Start > len(text) {
		return "", errors.New("invalid insert position")
	}
	return text[:iv.Start] + insert + text[iv.Start:], nil
}

// ExtractSlices returns a slice of substrings from `text`
// corresponding to each interval in the set.
// Returns an error if any interval is out of range.
func (set IntervalSet) ExtractSlices(text string) ([]string, error) {
	result := make([]string, 0, len(set))
	for _, iv := range set {
		if !iv.IsValid() || iv.Start < 0 || iv.End > len(text) {
			return nil, fmt.Errorf("interval %v out of range", iv)
		}
		part := text[iv.Start:iv.End]
		result = append(result, part)
	}
	return result, nil
}

// IntegerInterval represents a [start, end) interval of byte]()

// 数学的には[Start, End)と表される。文字列を扱うときのindexに適合する。
// "abc" 全体 → [0,3)
// "a" → [0,1)、補集合 → "bc" = [1,3)
// "c" → [2,3)、補集合 → "ab" = [0,2)
type IntegerInterval struct {
	Start int
	End   int
}

type IntervalSet []IntegerInterval

// IsValid ⇔ Start ≤ End
func (iv IntegerInterval) IsValid() bool {
	return iv.Start <= iv.End
}

// Length() = End − Start
func (iv IntegerInterval) Length() int {
	return iv.End - iv.Start
}

// Contains(n) ⇔ n ∈ [Start, End)
func (iv IntegerInterval) Contains(n int) bool {
	return iv.Start <= n && n < iv.End
}

// Overlaps(other) ⇔ [Start, End) ∩ [other.Start, other.End) ≠ ∅
func (iv IntegerInterval) Overlaps(other IntegerInterval) bool {
	return iv.Start < other.End && other.Start < iv.End
}

// Intersect(other) = iv ∩ other, if non-empty
func (iv IntegerInterval) Intersect(other IntegerInterval) (IntegerInterval, bool) {
	start := max(iv.Start, other.Start)
	end := min(iv.End, other.End)
	if start < end {
		return IntegerInterval{Start: start, End: end}, true
	}
	return IntegerInterval{}, false
}

// 連続または重複していればマージ可能
// → [0,2) + [2,5) → [0,5)
// Merge(other) = [Start, End) ∪ [other.Start, other.End), if Overlaps or IsAdjacent
func (iv IntegerInterval) Merge(other IntegerInterval) (IntegerInterval, bool) {
	if iv.End < other.Start || other.End < iv.Start {
		// 完全に離れていればマージ不可（接してない）。離れている場合はスライスにすべき。
		return IntegerInterval{}, false
	}
	start := min(iv.Start, other.Start)
	end := max(iv.End, other.End)
	return IntegerInterval{Start: start, End: end}, true
}

// 自分から other を引く
// → 1つまたは2つの区間に分かれる可能性あり（または空）
// Subtract(other) = [Start, End) − [other.Start, other.End)
// (may return 0, 1, or 2 intervals)
func (iv IntegerInterval) Subtract(other IntegerInterval) []IntegerInterval {
	intersection, ok := iv.Intersect(other)
	if !ok {
		return []IntegerInterval{iv}
	}
	result := []IntegerInterval{}
	if iv.Start < intersection.Start {
		result = append(result, IntegerInterval{iv.Start, intersection.Start})
	}
	if intersection.End < iv.End {
		result = append(result, IntegerInterval{intersection.End, iv.End})
	}
	return result
}

// Equal ⇔ Start = other.Start ∧ End = other.End
func (iv IntegerInterval) Equal(other IntegerInterval) bool {
	return iv.Start == other.Start && iv.End == other.End
}

// IsEmpty ⇔ Length() = 0 ⇔ Start = End
func (iv IntegerInterval) IsEmpty() bool {
	return iv.Start == iv.End
}

// 自分の直後または直前に他の区間が続いているか
// つまり iv.End == other.Start または iv.Start == other.End
// IsAdjacent(other) ⇔ End = other.Start ∨ Start = other.End
func (iv IntegerInterval) IsAdjacent(other IntegerInterval) bool {
	return iv.End == other.Start || other.End == iv.Start
}

// ソートのための比較関数（Start優先、Endはタイブレーク）
// Compare by Start, then End
func (iv IntegerInterval) Compare(other IntegerInterval) int {
	if iv.Start != other.Start {
		return iv.Start - other.Start
	}
	return iv.End - other.End
}

// Covers(other) ⇔ [Start, End) ⊇ [other.Start, other.End)
func (iv IntegerInterval) Covers(other IntegerInterval) bool {
	return iv.Start <= other.Start && iv.End >= other.End
}

// Normalize returns a new IntervalSet where all overlapping or adjacent intervals are merged.
//
// All intervals are assumed to be half-open: [start, end).
// The result is sorted and non-overlapping.
//
// For example:
//
//	input  = {[0,2), [1,4), [5,6)}
//	result = {[0,4), [5,6)}
//
// Normalize merges all overlapping or adjacent intervals.
// Result: disjoint, sorted, minimal form.
func (set IntervalSet) Normalize() IntervalSet {
	if len(set) == 0 {
		return nil
	}

	// まずコピーしてソート
	sorted := slices.Clone(set)
	slices.SortFunc(sorted, func(a, b IntegerInterval) int {
		return a.Compare(b)
	})

	result := make(IntervalSet, 0, len(sorted))
	current := sorted[0]

	for _, next := range sorted[1:] {
		if merged, ok := current.Merge(next); ok {
			current = merged
		} else {
			result = append(result, current)
			current = next
		}
	}
	result = append(result, current)
	return result
}

// ContainsPoint reports whether the given integer n is contained in any of the intervals in the set.
//
// For example:
//
//	set = {[0,3), [5,7)}
//	n = 2   → true
//	n = 4   → false
//
// ContainsPoint(n) ⇔ ∃ iv ∈ set, n ∈ iv
func (set IntervalSet) ContainsPoint(n int) bool {
	for _, iv := range set {
		if iv.Contains(n) {
			return true
		}
	}
	return false
}

// ContainsInterval reports whether the given interval is entirely covered by any interval in the set.
//
// For example:
//
//	set = {[0,3), [5,7)}
//	iv  = [1,2)   → true
//	iv  = [2,5)   → false
//
// ContainsInterval(iv) ⇔ ∃ s ∈ set, iv ⊆ s
func (set IntervalSet) ContainsInterval(iv IntegerInterval) bool {
	for _, current := range set {
		if current.Covers(iv) {
			return true
		}
	}
	return false
}

// Subtract returns a new IntervalSet where the given interval has been removed from all intervals in the set.
//
// Each interval in the set is individually subtracted by iv, and the results are collected and normalized.
//
// For example:
//
//	set = {[0,5)}
//	iv  = [2,4)
//	result = {[0,2), [4,5)}
//
// Subtract(iv) = ⋃ (s − iv) for s ∈ set
func (set IntervalSet) Subtract(iv IntegerInterval) IntervalSet {
	result := make(IntervalSet, 0, len(set)*2) // 最悪ケース：全区間が2つに分かれる
	for _, current := range set {
		subtracted := current.Subtract(iv) // []IntegerInterval
		result = append(result, subtracted...)
	}
	return result
}

// Union returns the union of the set and another IntervalSet, merging overlapping or adjacent intervals.
//
// All intervals are treated as half-open: [start, end).
//
// For example:
//
//	a = {[0,2), [5,6)}
//	b = {[1,4), [6,8)}
//	result = {[0,4), [5,8)} → Normalize ⇒ {[0,8)}
//
// Union(set') = Normalize(set ∪ set')
func (set IntervalSet) Union(other IntervalSet) IntervalSet {
	combined := make(IntervalSet, 0, len(set)+len(other))
	combined = append(combined, set...)
	combined = append(combined, other...)
	return combined.Normalize()
}

// Intersect returns a new IntervalSet consisting of all intersections between intervals in the set and another set.
//
// Each pair of intervals is intersected, and all non-empty intersections are collected and normalized.
//
// For example:
//
//	a = {[0,5), [6,8)}
//	b = {[3,7)}
//	result = {[3,5), [6,7)}
//
// Intersect(set') = Normalize({ s ∩ t | s ∈ set, t ∈ set', s ∩ t ≠ ∅ })
func (set IntervalSet) Intersect(other IntervalSet) IntervalSet {
	result := make(IntervalSet, 0)

	for _, iv1 := range set {
		for _, iv2 := range other {
			if intersection, ok := iv1.Intersect(iv2); ok {
				result = append(result, intersection)
			}
		}
	}

	return result.Normalize()
}

// Complement returns the complement of the interval set within the given base interval.
//
// All intervals are treated as half-open: [start, end).
// This method computes: base - set
//
// For example:
//
//	set  = {[0,1), [2,3)}
//	base = [0,3)
//	result = {[1,2)}
//
// Complement(base) = base − ⋃(iv ∈ set)
func (set IntervalSet) Complement(base IntegerInterval) IntervalSet {
	if base.IsEmpty() {
		return nil
	}
	subtracted := IntervalSet{base}
	for _, iv := range set {
		next := make(IntervalSet, 0)
		for _, s := range subtracted {
			next = append(next, s.Subtract(iv)...)
		}
		subtracted = next
	}
	return subtracted.Normalize()
}

func (iv IntegerInterval) String() string {
	return fmt.Sprintf("[%d,%d)", iv.Start, iv.End)
}

func (set IntervalSet) String() string {
	if len(set) == 0 {
		return "{}"
	}
	parts := make([]string, len(set))
	for i, iv := range set {
		parts[i] = iv.String()
	}
	return "{" + strings.Join(parts, ", ") + "}"
}
