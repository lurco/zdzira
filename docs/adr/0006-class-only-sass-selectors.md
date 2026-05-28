# Class-only selectors in SASS — no tag or ID selectors

All SASS rules target class selectors. Tag selectors (`button`, `input`, `body`) and ID selectors (`#board`) are not used, even when a rule applies to exactly one element type.

**Consequences visible in the code:** Base styles go on `:root` (a pseudo-class) rather than `body`. Nested interactive elements get explicit classes — `.mode-switch__btn` instead of `.mode-switch button`, `.search__input` instead of `.search input`. Pug templates and `board.js` are updated to match.

**Why:** Tag selectors bleed across component boundaries and raise specificity unpredictably. A rule like `.popover button` catches any `<button>` that appears inside a popover, including ones added later by a mixin or a JS fragment. Class selectors are opt-in — a button only gets popover-button styles when the author says so. ID selectors are excluded for the same specificity reason; IDs in HTML are fine for JS hooks (`document.getElementById`), but not as CSS targets.

**Naming convention:** BEM-style double-underscore for element classes scoped to a block (`.new-card-form__textarea`, `.popover__btn`). Strict BEM is not required; the rule is "class selectors only," not "BEM everywhere."

**Trade-off accepted:** More verbose Pug and JS — every interactive element needs an explicit class. Tag selectors are more concise and self-documenting for single-element contexts. We accept the verbosity because specificity bugs are hard to diagnose and the flat selector graph makes the stylesheet easier to audit.
