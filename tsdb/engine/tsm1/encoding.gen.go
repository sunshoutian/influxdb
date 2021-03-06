// Generated by tmpl
// https://github.com/benbjohnson/tmpl
//
// DO NOT EDIT!
// Source: encoding.gen.go.tmpl

package tsm1

import (
	"fmt"
	"sort"
)

// Values represents a slice of  values.
type Values []Value

func (a Values) MinTime() int64 {
	return a[0].UnixNano()
}

func (a Values) MaxTime() int64 {
	return a[len(a)-1].UnixNano()
}

func (a Values) Size() int {
	sz := 0
	for _, v := range a {
		sz += v.Size()
	}
	return sz
}

func (a Values) ordered() bool {
	if len(a) <= 1 {
		return true
	}
	for i := 1; i < len(a); i++ {
		if av, ab := a[i-1].UnixNano(), a[i].UnixNano(); av >= ab {
			return false
		}
	}
	return true
}

func (a Values) assertOrdered() {
	if len(a) <= 1 {
		return
	}
	for i := 1; i < len(a); i++ {
		if av, ab := a[i-1].UnixNano(), a[i].UnixNano(); av >= ab {
			panic(fmt.Sprintf("not ordered: %d %d >= %d", i, av, ab))
		}
	}
}

// Deduplicate returns a new slice with any values that have the same timestamp removed.
// The Value that appears last in the slice is the one that is kept.
func (a Values) Deduplicate() Values {
	if len(a) == 0 {
		return a
	}

	// See if we're already sorted and deduped
	var needSort bool
	for i := 1; i < len(a); i++ {
		if a[i-1].UnixNano() >= a[i].UnixNano() {
			needSort = true
			break
		}
	}

	if !needSort {
		return a
	}

	sort.Stable(a)
	var i int
	for j := 1; j < len(a); j++ {
		v := a[j]
		if v.UnixNano() != a[i].UnixNano() {
			i++
		}
		a[i] = v

	}
	return a[:i+1]
}

//  Exclude returns the subset of values not in [min, max]
func (a Values) Exclude(min, max int64) Values {
	var i int
	for j := 0; j < len(a); j++ {
		if a[j].UnixNano() >= min && a[j].UnixNano() <= max {
			continue
		}

		a[i] = a[j]
		i++
	}
	return a[:i]
}

// Include returns the subset values between min and max inclusive.
func (a Values) Include(min, max int64) Values {
	var i int
	for j := 0; j < len(a); j++ {
		if a[j].UnixNano() < min || a[j].UnixNano() > max {
			continue
		}

		a[i] = a[j]
		i++
	}
	return a[:i]
}

// Merge overlays b to top of a.  If two values conflict with
// the same timestamp, b is used.  Both a and b must be sorted
// in ascending order.
func (a Values) Merge(b Values) Values {
	if len(a) == 0 {
		return b
	}

	if len(b) == 0 {
		return a
	}

	// Normally, both a and b should not contain duplicates.  Due to a bug in older versions, it's
	// possible stored blocks might contain duplicate values.  Remove them if they exists before
	// merging.
	a = a.Deduplicate()
	b = b.Deduplicate()

	if a[len(a)-1].UnixNano() < b[0].UnixNano() {
		return append(a, b...)
	}

	if b[len(b)-1].UnixNano() < a[0].UnixNano() {
		return append(b, a...)
	}

	for i := 0; i < len(a) && len(b) > 0; i++ {
		av, bv := a[i].UnixNano(), b[0].UnixNano()
		// Value in a is greater than B, we need to merge
		if av > bv {
			// Save value in a
			temp := a[i]

			// Overwrite a with b
			a[i] = b[0]

			// Slide all values of b down 1
			copy(b, b[1:])
			b = b[:len(b)-1]

			var k int
			if len(b) > 0 && av > b[len(b)-1].UnixNano() {
				// Fast path where a is after b, we skip the search
				k = len(b)
			} else {
				// See where value we save from a should be inserted in b to keep b sorted
				k = sort.Search(len(b), func(i int) bool { return b[i].UnixNano() >= temp.UnixNano() })
			}

			if k == len(b) {
				// Last position?
				b = append(b, temp)
			} else if b[k].UnixNano() != temp.UnixNano() {
				// Save the last element, since it will get overwritten
				last := b[len(b)-1]
				// Somewhere in the middle of b, insert it only if it's not a duplicate
				copy(b[k+1:], b[k:])
				// Add the last vale to the end
				b = append(b, last)
				b[k] = temp
			}
		} else if av == bv {
			// Value in a an b are the same, use b
			a[i] = b[0]
			b = b[1:]
		}
	}

	if len(b) > 0 {
		return append(a, b...)
	}
	return a
}

// Sort methods
func (a Values) Len() int           { return len(a) }
func (a Values) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Values) Less(i, j int) bool { return a[i].UnixNano() < a[j].UnixNano() }

// FloatValues represents a slice of Float values.
type FloatValues []FloatValue

func (a FloatValues) MinTime() int64 {
	return a[0].UnixNano()
}

func (a FloatValues) MaxTime() int64 {
	return a[len(a)-1].UnixNano()
}

func (a FloatValues) Size() int {
	sz := 0
	for _, v := range a {
		sz += v.Size()
	}
	return sz
}

func (a FloatValues) ordered() bool {
	if len(a) <= 1 {
		return true
	}
	for i := 1; i < len(a); i++ {
		if av, ab := a[i-1].UnixNano(), a[i].UnixNano(); av >= ab {
			return false
		}
	}
	return true
}

func (a FloatValues) assertOrdered() {
	if len(a) <= 1 {
		return
	}
	for i := 1; i < len(a); i++ {
		if av, ab := a[i-1].UnixNano(), a[i].UnixNano(); av >= ab {
			panic(fmt.Sprintf("not ordered: %d %d >= %d", i, av, ab))
		}
	}
}

// Deduplicate returns a new slice with any values that have the same timestamp removed.
// The Value that appears last in the slice is the one that is kept.
func (a FloatValues) Deduplicate() FloatValues {
	if len(a) == 0 {
		return a
	}

	// See if we're already sorted and deduped
	var needSort bool
	for i := 1; i < len(a); i++ {
		if a[i-1].UnixNano() >= a[i].UnixNano() {
			needSort = true
			break
		}
	}

	if !needSort {
		return a
	}

	sort.Stable(a)
	var i int
	for j := 1; j < len(a); j++ {
		v := a[j]
		if v.UnixNano() != a[i].UnixNano() {
			i++
		}
		a[i] = v

	}
	return a[:i+1]
}

//  Exclude returns the subset of values not in [min, max]
func (a FloatValues) Exclude(min, max int64) FloatValues {
	var i int
	for j := 0; j < len(a); j++ {
		if a[j].UnixNano() >= min && a[j].UnixNano() <= max {
			continue
		}

		a[i] = a[j]
		i++
	}
	return a[:i]
}

// Include returns the subset values between min and max inclusive.
func (a FloatValues) Include(min, max int64) FloatValues {
	var i int
	for j := 0; j < len(a); j++ {
		if a[j].UnixNano() < min || a[j].UnixNano() > max {
			continue
		}

		a[i] = a[j]
		i++
	}
	return a[:i]
}

// Merge overlays b to top of a.  If two values conflict with
// the same timestamp, b is used.  Both a and b must be sorted
// in ascending order.
func (a FloatValues) Merge(b FloatValues) FloatValues {
	if len(a) == 0 {
		return b
	}

	if len(b) == 0 {
		return a
	}

	// Normally, both a and b should not contain duplicates.  Due to a bug in older versions, it's
	// possible stored blocks might contain duplicate values.  Remove them if they exists before
	// merging.
	a = a.Deduplicate()
	b = b.Deduplicate()

	if a[len(a)-1].UnixNano() < b[0].UnixNano() {
		return append(a, b...)
	}

	if b[len(b)-1].UnixNano() < a[0].UnixNano() {
		return append(b, a...)
	}

	for i := 0; i < len(a) && len(b) > 0; i++ {
		av, bv := a[i].UnixNano(), b[0].UnixNano()
		// Value in a is greater than B, we need to merge
		if av > bv {
			// Save value in a
			temp := a[i]

			// Overwrite a with b
			a[i] = b[0]

			// Slide all values of b down 1
			copy(b, b[1:])
			b = b[:len(b)-1]

			var k int
			if len(b) > 0 && av > b[len(b)-1].UnixNano() {
				// Fast path where a is after b, we skip the search
				k = len(b)
			} else {
				// See where value we save from a should be inserted in b to keep b sorted
				k = sort.Search(len(b), func(i int) bool { return b[i].UnixNano() >= temp.UnixNano() })
			}

			if k == len(b) {
				// Last position?
				b = append(b, temp)
			} else if b[k].UnixNano() != temp.UnixNano() {
				// Save the last element, since it will get overwritten
				last := b[len(b)-1]
				// Somewhere in the middle of b, insert it only if it's not a duplicate
				copy(b[k+1:], b[k:])
				// Add the last vale to the end
				b = append(b, last)
				b[k] = temp
			}
		} else if av == bv {
			// Value in a an b are the same, use b
			a[i] = b[0]
			b = b[1:]
		}
	}

	if len(b) > 0 {
		return append(a, b...)
	}
	return a
}

func (a FloatValues) Encode(buf []byte) ([]byte, error) {
	return encodeFloatValuesBlock(buf, a)
}

func encodeFloatValuesBlock(buf []byte, values []FloatValue) ([]byte, error) {
	if len(values) == 0 {
		return nil, nil
	}

	venc := getFloatEncoder(len(values))
	tsenc := getTimeEncoder(len(values))

	var b []byte
	err := func() error {
		for _, v := range values {
			tsenc.Write(v.unixnano)
			venc.Write(v.value)
		}
		venc.Flush()

		// Encoded timestamp values
		tb, err := tsenc.Bytes()
		if err != nil {
			return err
		}
		// Encoded values
		vb, err := venc.Bytes()
		if err != nil {
			return err
		}

		// Prepend the first timestamp of the block in the first 8 bytes and the block
		// in the next byte, followed by the block
		b = packBlock(buf, BlockFloat64, tb, vb)

		return nil
	}()

	putTimeEncoder(tsenc)
	putFloatEncoder(venc)

	return b, err
}

// Sort methods
func (a FloatValues) Len() int           { return len(a) }
func (a FloatValues) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a FloatValues) Less(i, j int) bool { return a[i].UnixNano() < a[j].UnixNano() }

// IntegerValues represents a slice of Integer values.
type IntegerValues []IntegerValue

func (a IntegerValues) MinTime() int64 {
	return a[0].UnixNano()
}

func (a IntegerValues) MaxTime() int64 {
	return a[len(a)-1].UnixNano()
}

func (a IntegerValues) Size() int {
	sz := 0
	for _, v := range a {
		sz += v.Size()
	}
	return sz
}

func (a IntegerValues) ordered() bool {
	if len(a) <= 1 {
		return true
	}
	for i := 1; i < len(a); i++ {
		if av, ab := a[i-1].UnixNano(), a[i].UnixNano(); av >= ab {
			return false
		}
	}
	return true
}

func (a IntegerValues) assertOrdered() {
	if len(a) <= 1 {
		return
	}
	for i := 1; i < len(a); i++ {
		if av, ab := a[i-1].UnixNano(), a[i].UnixNano(); av >= ab {
			panic(fmt.Sprintf("not ordered: %d %d >= %d", i, av, ab))
		}
	}
}

// Deduplicate returns a new slice with any values that have the same timestamp removed.
// The Value that appears last in the slice is the one that is kept.
func (a IntegerValues) Deduplicate() IntegerValues {
	if len(a) == 0 {
		return a
	}

	// See if we're already sorted and deduped
	var needSort bool
	for i := 1; i < len(a); i++ {
		if a[i-1].UnixNano() >= a[i].UnixNano() {
			needSort = true
			break
		}
	}

	if !needSort {
		return a
	}

	sort.Stable(a)
	var i int
	for j := 1; j < len(a); j++ {
		v := a[j]
		if v.UnixNano() != a[i].UnixNano() {
			i++
		}
		a[i] = v

	}
	return a[:i+1]
}

//  Exclude returns the subset of values not in [min, max]
func (a IntegerValues) Exclude(min, max int64) IntegerValues {
	var i int
	for j := 0; j < len(a); j++ {
		if a[j].UnixNano() >= min && a[j].UnixNano() <= max {
			continue
		}

		a[i] = a[j]
		i++
	}
	return a[:i]
}

// Include returns the subset values between min and max inclusive.
func (a IntegerValues) Include(min, max int64) IntegerValues {
	var i int
	for j := 0; j < len(a); j++ {
		if a[j].UnixNano() < min || a[j].UnixNano() > max {
			continue
		}

		a[i] = a[j]
		i++
	}
	return a[:i]
}

// Merge overlays b to top of a.  If two values conflict with
// the same timestamp, b is used.  Both a and b must be sorted
// in ascending order.
func (a IntegerValues) Merge(b IntegerValues) IntegerValues {
	if len(a) == 0 {
		return b
	}

	if len(b) == 0 {
		return a
	}

	// Normally, both a and b should not contain duplicates.  Due to a bug in older versions, it's
	// possible stored blocks might contain duplicate values.  Remove them if they exists before
	// merging.
	a = a.Deduplicate()
	b = b.Deduplicate()

	if a[len(a)-1].UnixNano() < b[0].UnixNano() {
		return append(a, b...)
	}

	if b[len(b)-1].UnixNano() < a[0].UnixNano() {
		return append(b, a...)
	}

	for i := 0; i < len(a) && len(b) > 0; i++ {
		av, bv := a[i].UnixNano(), b[0].UnixNano()
		// Value in a is greater than B, we need to merge
		if av > bv {
			// Save value in a
			temp := a[i]

			// Overwrite a with b
			a[i] = b[0]

			// Slide all values of b down 1
			copy(b, b[1:])
			b = b[:len(b)-1]

			var k int
			if len(b) > 0 && av > b[len(b)-1].UnixNano() {
				// Fast path where a is after b, we skip the search
				k = len(b)
			} else {
				// See where value we save from a should be inserted in b to keep b sorted
				k = sort.Search(len(b), func(i int) bool { return b[i].UnixNano() >= temp.UnixNano() })
			}

			if k == len(b) {
				// Last position?
				b = append(b, temp)
			} else if b[k].UnixNano() != temp.UnixNano() {
				// Save the last element, since it will get overwritten
				last := b[len(b)-1]
				// Somewhere in the middle of b, insert it only if it's not a duplicate
				copy(b[k+1:], b[k:])
				// Add the last vale to the end
				b = append(b, last)
				b[k] = temp
			}
		} else if av == bv {
			// Value in a an b are the same, use b
			a[i] = b[0]
			b = b[1:]
		}
	}

	if len(b) > 0 {
		return append(a, b...)
	}
	return a
}

func (a IntegerValues) Encode(buf []byte) ([]byte, error) {
	return encodeIntegerValuesBlock(buf, a)
}

func encodeIntegerValuesBlock(buf []byte, values []IntegerValue) ([]byte, error) {
	if len(values) == 0 {
		return nil, nil
	}

	venc := getIntegerEncoder(len(values))
	tsenc := getTimeEncoder(len(values))

	var b []byte
	err := func() error {
		for _, v := range values {
			tsenc.Write(v.unixnano)
			venc.Write(v.value)
		}
		venc.Flush()

		// Encoded timestamp values
		tb, err := tsenc.Bytes()
		if err != nil {
			return err
		}
		// Encoded values
		vb, err := venc.Bytes()
		if err != nil {
			return err
		}

		// Prepend the first timestamp of the block in the first 8 bytes and the block
		// in the next byte, followed by the block
		b = packBlock(buf, BlockInteger, tb, vb)

		return nil
	}()

	putTimeEncoder(tsenc)
	putIntegerEncoder(venc)

	return b, err
}

// Sort methods
func (a IntegerValues) Len() int           { return len(a) }
func (a IntegerValues) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a IntegerValues) Less(i, j int) bool { return a[i].UnixNano() < a[j].UnixNano() }

// StringValues represents a slice of String values.
type StringValues []StringValue

func (a StringValues) MinTime() int64 {
	return a[0].UnixNano()
}

func (a StringValues) MaxTime() int64 {
	return a[len(a)-1].UnixNano()
}

func (a StringValues) Size() int {
	sz := 0
	for _, v := range a {
		sz += v.Size()
	}
	return sz
}

func (a StringValues) ordered() bool {
	if len(a) <= 1 {
		return true
	}
	for i := 1; i < len(a); i++ {
		if av, ab := a[i-1].UnixNano(), a[i].UnixNano(); av >= ab {
			return false
		}
	}
	return true
}

func (a StringValues) assertOrdered() {
	if len(a) <= 1 {
		return
	}
	for i := 1; i < len(a); i++ {
		if av, ab := a[i-1].UnixNano(), a[i].UnixNano(); av >= ab {
			panic(fmt.Sprintf("not ordered: %d %d >= %d", i, av, ab))
		}
	}
}

// Deduplicate returns a new slice with any values that have the same timestamp removed.
// The Value that appears last in the slice is the one that is kept.
func (a StringValues) Deduplicate() StringValues {
	if len(a) == 0 {
		return a
	}

	// See if we're already sorted and deduped
	var needSort bool
	for i := 1; i < len(a); i++ {
		if a[i-1].UnixNano() >= a[i].UnixNano() {
			needSort = true
			break
		}
	}

	if !needSort {
		return a
	}

	sort.Stable(a)
	var i int
	for j := 1; j < len(a); j++ {
		v := a[j]
		if v.UnixNano() != a[i].UnixNano() {
			i++
		}
		a[i] = v

	}
	return a[:i+1]
}

//  Exclude returns the subset of values not in [min, max]
func (a StringValues) Exclude(min, max int64) StringValues {
	var i int
	for j := 0; j < len(a); j++ {
		if a[j].UnixNano() >= min && a[j].UnixNano() <= max {
			continue
		}

		a[i] = a[j]
		i++
	}
	return a[:i]
}

// Include returns the subset values between min and max inclusive.
func (a StringValues) Include(min, max int64) StringValues {
	var i int
	for j := 0; j < len(a); j++ {
		if a[j].UnixNano() < min || a[j].UnixNano() > max {
			continue
		}

		a[i] = a[j]
		i++
	}
	return a[:i]
}

// Merge overlays b to top of a.  If two values conflict with
// the same timestamp, b is used.  Both a and b must be sorted
// in ascending order.
func (a StringValues) Merge(b StringValues) StringValues {
	if len(a) == 0 {
		return b
	}

	if len(b) == 0 {
		return a
	}

	// Normally, both a and b should not contain duplicates.  Due to a bug in older versions, it's
	// possible stored blocks might contain duplicate values.  Remove them if they exists before
	// merging.
	a = a.Deduplicate()
	b = b.Deduplicate()

	if a[len(a)-1].UnixNano() < b[0].UnixNano() {
		return append(a, b...)
	}

	if b[len(b)-1].UnixNano() < a[0].UnixNano() {
		return append(b, a...)
	}

	for i := 0; i < len(a) && len(b) > 0; i++ {
		av, bv := a[i].UnixNano(), b[0].UnixNano()
		// Value in a is greater than B, we need to merge
		if av > bv {
			// Save value in a
			temp := a[i]

			// Overwrite a with b
			a[i] = b[0]

			// Slide all values of b down 1
			copy(b, b[1:])
			b = b[:len(b)-1]

			var k int
			if len(b) > 0 && av > b[len(b)-1].UnixNano() {
				// Fast path where a is after b, we skip the search
				k = len(b)
			} else {
				// See where value we save from a should be inserted in b to keep b sorted
				k = sort.Search(len(b), func(i int) bool { return b[i].UnixNano() >= temp.UnixNano() })
			}

			if k == len(b) {
				// Last position?
				b = append(b, temp)
			} else if b[k].UnixNano() != temp.UnixNano() {
				// Save the last element, since it will get overwritten
				last := b[len(b)-1]
				// Somewhere in the middle of b, insert it only if it's not a duplicate
				copy(b[k+1:], b[k:])
				// Add the last vale to the end
				b = append(b, last)
				b[k] = temp
			}
		} else if av == bv {
			// Value in a an b are the same, use b
			a[i] = b[0]
			b = b[1:]
		}
	}

	if len(b) > 0 {
		return append(a, b...)
	}
	return a
}

func (a StringValues) Encode(buf []byte) ([]byte, error) {
	return encodeStringValuesBlock(buf, a)
}

func encodeStringValuesBlock(buf []byte, values []StringValue) ([]byte, error) {
	if len(values) == 0 {
		return nil, nil
	}

	venc := getStringEncoder(len(values))
	tsenc := getTimeEncoder(len(values))

	var b []byte
	err := func() error {
		for _, v := range values {
			tsenc.Write(v.unixnano)
			venc.Write(v.value)
		}
		venc.Flush()

		// Encoded timestamp values
		tb, err := tsenc.Bytes()
		if err != nil {
			return err
		}
		// Encoded values
		vb, err := venc.Bytes()
		if err != nil {
			return err
		}

		// Prepend the first timestamp of the block in the first 8 bytes and the block
		// in the next byte, followed by the block
		b = packBlock(buf, BlockString, tb, vb)

		return nil
	}()

	putTimeEncoder(tsenc)
	putStringEncoder(venc)

	return b, err
}

// Sort methods
func (a StringValues) Len() int           { return len(a) }
func (a StringValues) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a StringValues) Less(i, j int) bool { return a[i].UnixNano() < a[j].UnixNano() }

// BooleanValues represents a slice of Boolean values.
type BooleanValues []BooleanValue

func (a BooleanValues) MinTime() int64 {
	return a[0].UnixNano()
}

func (a BooleanValues) MaxTime() int64 {
	return a[len(a)-1].UnixNano()
}

func (a BooleanValues) Size() int {
	sz := 0
	for _, v := range a {
		sz += v.Size()
	}
	return sz
}

func (a BooleanValues) ordered() bool {
	if len(a) <= 1 {
		return true
	}
	for i := 1; i < len(a); i++ {
		if av, ab := a[i-1].UnixNano(), a[i].UnixNano(); av >= ab {
			return false
		}
	}
	return true
}

func (a BooleanValues) assertOrdered() {
	if len(a) <= 1 {
		return
	}
	for i := 1; i < len(a); i++ {
		if av, ab := a[i-1].UnixNano(), a[i].UnixNano(); av >= ab {
			panic(fmt.Sprintf("not ordered: %d %d >= %d", i, av, ab))
		}
	}
}

// Deduplicate returns a new slice with any values that have the same timestamp removed.
// The Value that appears last in the slice is the one that is kept.
func (a BooleanValues) Deduplicate() BooleanValues {
	if len(a) == 0 {
		return a
	}

	// See if we're already sorted and deduped
	var needSort bool
	for i := 1; i < len(a); i++ {
		if a[i-1].UnixNano() >= a[i].UnixNano() {
			needSort = true
			break
		}
	}

	if !needSort {
		return a
	}

	sort.Stable(a)
	var i int
	for j := 1; j < len(a); j++ {
		v := a[j]
		if v.UnixNano() != a[i].UnixNano() {
			i++
		}
		a[i] = v

	}
	return a[:i+1]
}

//  Exclude returns the subset of values not in [min, max]
func (a BooleanValues) Exclude(min, max int64) BooleanValues {
	var i int
	for j := 0; j < len(a); j++ {
		if a[j].UnixNano() >= min && a[j].UnixNano() <= max {
			continue
		}

		a[i] = a[j]
		i++
	}
	return a[:i]
}

// Include returns the subset values between min and max inclusive.
func (a BooleanValues) Include(min, max int64) BooleanValues {
	var i int
	for j := 0; j < len(a); j++ {
		if a[j].UnixNano() < min || a[j].UnixNano() > max {
			continue
		}

		a[i] = a[j]
		i++
	}
	return a[:i]
}

// Merge overlays b to top of a.  If two values conflict with
// the same timestamp, b is used.  Both a and b must be sorted
// in ascending order.
func (a BooleanValues) Merge(b BooleanValues) BooleanValues {
	if len(a) == 0 {
		return b
	}

	if len(b) == 0 {
		return a
	}

	// Normally, both a and b should not contain duplicates.  Due to a bug in older versions, it's
	// possible stored blocks might contain duplicate values.  Remove them if they exists before
	// merging.
	a = a.Deduplicate()
	b = b.Deduplicate()

	if a[len(a)-1].UnixNano() < b[0].UnixNano() {
		return append(a, b...)
	}

	if b[len(b)-1].UnixNano() < a[0].UnixNano() {
		return append(b, a...)
	}

	for i := 0; i < len(a) && len(b) > 0; i++ {
		av, bv := a[i].UnixNano(), b[0].UnixNano()
		// Value in a is greater than B, we need to merge
		if av > bv {
			// Save value in a
			temp := a[i]

			// Overwrite a with b
			a[i] = b[0]

			// Slide all values of b down 1
			copy(b, b[1:])
			b = b[:len(b)-1]

			var k int
			if len(b) > 0 && av > b[len(b)-1].UnixNano() {
				// Fast path where a is after b, we skip the search
				k = len(b)
			} else {
				// See where value we save from a should be inserted in b to keep b sorted
				k = sort.Search(len(b), func(i int) bool { return b[i].UnixNano() >= temp.UnixNano() })
			}

			if k == len(b) {
				// Last position?
				b = append(b, temp)
			} else if b[k].UnixNano() != temp.UnixNano() {
				// Save the last element, since it will get overwritten
				last := b[len(b)-1]
				// Somewhere in the middle of b, insert it only if it's not a duplicate
				copy(b[k+1:], b[k:])
				// Add the last vale to the end
				b = append(b, last)
				b[k] = temp
			}
		} else if av == bv {
			// Value in a an b are the same, use b
			a[i] = b[0]
			b = b[1:]
		}
	}

	if len(b) > 0 {
		return append(a, b...)
	}
	return a
}

func (a BooleanValues) Encode(buf []byte) ([]byte, error) {
	return encodeBooleanValuesBlock(buf, a)
}

func encodeBooleanValuesBlock(buf []byte, values []BooleanValue) ([]byte, error) {
	if len(values) == 0 {
		return nil, nil
	}

	venc := getBooleanEncoder(len(values))
	tsenc := getTimeEncoder(len(values))

	var b []byte
	err := func() error {
		for _, v := range values {
			tsenc.Write(v.unixnano)
			venc.Write(v.value)
		}
		venc.Flush()

		// Encoded timestamp values
		tb, err := tsenc.Bytes()
		if err != nil {
			return err
		}
		// Encoded values
		vb, err := venc.Bytes()
		if err != nil {
			return err
		}

		// Prepend the first timestamp of the block in the first 8 bytes and the block
		// in the next byte, followed by the block
		b = packBlock(buf, BlockBoolean, tb, vb)

		return nil
	}()

	putTimeEncoder(tsenc)
	putBooleanEncoder(venc)

	return b, err
}

// Sort methods
func (a BooleanValues) Len() int           { return len(a) }
func (a BooleanValues) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BooleanValues) Less(i, j int) bool { return a[i].UnixNano() < a[j].UnixNano() }
