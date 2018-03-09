package anitogo

import "github.com/juju/errors"

// HasFlags returns whether the token has the specified flags or not
func (t token) HasFlags(flags tokenFlags) bool {
	return (t.Flags & flags) == flags
}

// HasCategories returns whether tho token has the specified categories
func (t token) HasCategories(categories tokenCategory) bool {
	return (t.Category & categories) == categories
}

// checkCategory serves some currently unknown function
func (t token) CheckCategory(fe, fn tokenFlags, cat tokenCategory) bool {
	if t.HasFlags(fe) {
		return t.Category == cat
	} else if t.HasFlags(fn) {
		return t.Category != cat
	}
	return false
}

// validateFlags validates that a token has sane flags and categories set
func (t token) ValidateFlags() error {
	// if has any enclosed flags
	if t.HasFlags(toeknFlagsMaskEnclosed) {
		var success bool
		// if has enclosed flag, error if t.Enclosed is false
		if t.HasFlags(tokenFlagsEnclosed) {
			success = t.Enclosed
			// if has not enclosed flag, error if t.Enclosed is true
		} else {
			success = !t.Enclosed
		}
		// error based on above
		if !success {
			return errors.Trace(ErrorPlaceholder)
		}
	}

	// if has any category flags
	// pass overall validation if has a matching category and normal flag, or if has
	// a normal not flag and does not have the matching category
	// REVIEW: does this logic even make sense; it's borrowed from anitomy
	if t.HasFlags(tokenFlagsMaskCategories) {
		// if has Bracket flag and has Bracket category flag, pass overall validation
		if t.CheckCategory(tokenFlagsBracket, tokenFlagsNotBracket, tokenCategoryBracket) {
			return errors.Trace(nil)
		}
		if t.CheckCategory(tokenFlagsDelimiter, tokenFlagsNotDelimiter, tokenCategoryDelimiter) {
			return errors.Trace(nil)
		}
		if t.CheckCategory(tokenFlagsIdentifier, tokenFlagsNotIdentifier, tokenCategoryIdentifier) {
			return errors.Trace(nil)
		}
		if t.CheckCategory(tokenFlagsUnknown, tokenFlagsNotUnknown, tokenCategoryUnknown) {
			return errors.Trace(nil)
		}
		if t.CheckCategory(tokenFlagsNotValid, tokenFlagsValid, tokenCategoryInvalid) {
			return errors.Trace(nil)
		}
		// has category flags, but category does not match; fail validation
		return errors.Trace(ErrorPlaceholder)
	}

	// else, pass validation
	return errors.Trace(nil)
}
