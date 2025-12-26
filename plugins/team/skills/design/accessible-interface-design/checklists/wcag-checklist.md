# WCAG 2.1 AA Compliance Checklist

This checklist covers all WCAG 2.1 Level A and Level AA success criteria. Use this for comprehensive accessibility audits.

## 1. Perceivable

Information and user interface components must be presentable to users in ways they can perceive.

### 1.1 Text Alternatives

| Criterion | Level | Requirement | Status |
|-----------|-------|-------------|--------|
| 1.1.1 Non-text Content | A | All non-text content has text alternatives | [ ] |

**Testing steps for 1.1.1:**
- [ ] Informative images have descriptive alt text
- [ ] Decorative images have empty alt (`alt=""`) or `role="presentation"`
- [ ] Complex images have extended descriptions
- [ ] Form inputs have accessible names
- [ ] Icons used alone have text alternatives
- [ ] CAPTCHAs have text alternatives and alternative forms

### 1.2 Time-based Media

| Criterion | Level | Requirement | Status |
|-----------|-------|-------------|--------|
| 1.2.1 Audio-only and Video-only | A | Alternatives for prerecorded audio-only and video-only | [ ] |
| 1.2.2 Captions (Prerecorded) | A | Captions for prerecorded audio in synchronized media | [ ] |
| 1.2.3 Audio Description or Media Alternative | A | Audio description or text alternative for video | [ ] |
| 1.2.4 Captions (Live) | AA | Captions for live audio in synchronized media | [ ] |
| 1.2.5 Audio Description (Prerecorded) | AA | Audio description for prerecorded video | [ ] |

**Testing steps for 1.2:**
- [ ] Pre-recorded audio has transcript
- [ ] Pre-recorded video has captions
- [ ] Pre-recorded video has audio descriptions for visual content
- [ ] Live streams have real-time captions
- [ ] Media players have accessible controls

### 1.3 Adaptable

| Criterion | Level | Requirement | Status |
|-----------|-------|-------------|--------|
| 1.3.1 Info and Relationships | A | Information, structure, and relationships are programmatically determinable | [ ] |
| 1.3.2 Meaningful Sequence | A | Reading order is programmatically determinable | [ ] |
| 1.3.3 Sensory Characteristics | A | Instructions do not rely solely on sensory characteristics | [ ] |
| 1.3.4 Orientation | AA | Content not restricted to single display orientation | [ ] |
| 1.3.5 Identify Input Purpose | AA | Input field purpose is programmatically determinable | [ ] |

**Testing steps for 1.3.1:**
- [ ] Headings use proper heading elements (h1-h6)
- [ ] Lists use list elements (ul, ol, dl)
- [ ] Tables have proper headers (th with scope)
- [ ] Form fields have labels
- [ ] Related form fields are grouped (fieldset/legend)
- [ ] Landmarks identify page regions

**Testing steps for 1.3.2:**
- [ ] DOM order matches visual order
- [ ] CSS does not alter meaningful reading sequence
- [ ] Tab order follows logical sequence

**Testing steps for 1.3.3:**
- [ ] Instructions do not rely solely on shape (e.g., "click the round button")
- [ ] Instructions do not rely solely on location (e.g., "menu on the left")
- [ ] Instructions do not rely solely on sound

**Testing steps for 1.3.4:**
- [ ] Content works in both portrait and landscape
- [ ] No JavaScript locks orientation

**Testing steps for 1.3.5:**
- [ ] Input fields use appropriate autocomplete attributes
- [ ] Name, email, phone, address fields have autocomplete values

### 1.4 Distinguishable

| Criterion | Level | Requirement | Status |
|-----------|-------|-------------|--------|
| 1.4.1 Use of Color | A | Color is not the only visual means of conveying information | [ ] |
| 1.4.2 Audio Control | A | Mechanism to pause/stop audio that plays automatically | [ ] |
| 1.4.3 Contrast (Minimum) | AA | Text has contrast ratio of at least 4.5:1 (3:1 for large text) | [ ] |
| 1.4.4 Resize Text | AA | Text can be resized up to 200% without loss of functionality | [ ] |
| 1.4.5 Images of Text | AA | Text is used instead of images of text | [ ] |
| 1.4.10 Reflow | AA | Content reflows without horizontal scrolling at 320px width | [ ] |
| 1.4.11 Non-text Contrast | AA | UI components and graphics have contrast ratio of at least 3:1 | [ ] |
| 1.4.12 Text Spacing | AA | No loss of content with increased text spacing | [ ] |
| 1.4.13 Content on Hover or Focus | AA | Additional content on hover/focus is dismissible, hoverable, persistent | [ ] |

**Testing steps for 1.4.1:**
- [ ] Links are distinguishable from surrounding text beyond color
- [ ] Form errors are indicated by more than color
- [ ] Charts/graphs use patterns or labels in addition to color
- [ ] Required fields are indicated by more than color

**Testing steps for 1.4.2:**
- [ ] No audio plays automatically for more than 3 seconds
- [ ] Mechanism to stop or control audio volume

**Testing steps for 1.4.3:**
- [ ] Normal text (under 18pt): 4.5:1 contrast ratio
- [ ] Large text (18pt+ or 14pt bold): 3:1 contrast ratio
- [ ] Use contrast checker tool to verify all text

**Testing steps for 1.4.4:**
- [ ] Increase browser zoom to 200%
- [ ] All content remains visible and functional
- [ ] No text is clipped or overlapped

**Testing steps for 1.4.5:**
- [ ] Essential text is not presented as images
- [ ] Logos are the only exception

**Testing steps for 1.4.10:**
- [ ] Set viewport to 320px width
- [ ] Content reflows to single column
- [ ] No horizontal scrolling required
- [ ] All content accessible

**Testing steps for 1.4.11:**
- [ ] Form field borders: 3:1 contrast
- [ ] Button borders/backgrounds: 3:1 contrast
- [ ] Focus indicators: 3:1 contrast
- [ ] Icons: 3:1 contrast against background

**Testing steps for 1.4.12:**
Apply these styles and verify no content loss:
```css
* {
  line-height: 1.5 !important;
  letter-spacing: 0.12em !important;
  word-spacing: 0.16em !important;
  p { margin-bottom: 2em !important; }
}
```

**Testing steps for 1.4.13:**
- [ ] Hover/focus content can be dismissed (Escape)
- [ ] User can hover over the additional content
- [ ] Content persists until dismissed or no longer relevant

## 2. Operable

User interface components and navigation must be operable.

### 2.1 Keyboard Accessible

| Criterion | Level | Requirement | Status |
|-----------|-------|-------------|--------|
| 2.1.1 Keyboard | A | All functionality available via keyboard | [ ] |
| 2.1.2 No Keyboard Trap | A | Keyboard focus can always be moved away from any component | [ ] |
| 2.1.4 Character Key Shortcuts | A | Single character key shortcuts can be turned off or remapped | [ ] |

**Testing steps for 2.1.1:**
- [ ] Tab navigates to all interactive elements
- [ ] Enter/Space activates buttons and links
- [ ] Arrow keys work within components
- [ ] No functionality requires mouse-only gestures

**Testing steps for 2.1.2:**
- [ ] Tab can exit all components
- [ ] Modals can be closed with Escape
- [ ] No component traps focus permanently

**Testing steps for 2.1.4:**
- [ ] Single-key shortcuts (if any) can be disabled
- [ ] Or shortcuts require modifier key

### 2.2 Enough Time

| Criterion | Level | Requirement | Status |
|-----------|-------|-------------|--------|
| 2.2.1 Timing Adjustable | A | Time limits can be turned off, adjusted, or extended | [ ] |
| 2.2.2 Pause, Stop, Hide | A | Moving, blinking, scrolling content can be paused | [ ] |

**Testing steps for 2.2.1:**
- [ ] Session timeouts can be extended
- [ ] Users warned before timeout
- [ ] Option to disable timeout where possible

**Testing steps for 2.2.2:**
- [ ] Carousels have pause control
- [ ] Animations can be stopped
- [ ] Auto-updating content can be paused

### 2.3 Seizures and Physical Reactions

| Criterion | Level | Requirement | Status |
|-----------|-------|-------------|--------|
| 2.3.1 Three Flashes or Below | A | No content flashes more than 3 times per second | [ ] |

**Testing steps for 2.3.1:**
- [ ] No content flashes more than 3 times/second
- [ ] Flash area is small (less than 341x256 pixels)

### 2.4 Navigable

| Criterion | Level | Requirement | Status |
|-----------|-------|-------------|--------|
| 2.4.1 Bypass Blocks | A | Mechanism to skip repetitive content | [ ] |
| 2.4.2 Page Titled | A | Pages have descriptive titles | [ ] |
| 2.4.3 Focus Order | A | Focus order preserves meaning and operability | [ ] |
| 2.4.4 Link Purpose (In Context) | A | Link purpose can be determined from link text or context | [ ] |
| 2.4.5 Multiple Ways | AA | Multiple ways to locate pages within a site | [ ] |
| 2.4.6 Headings and Labels | AA | Headings and labels describe topic or purpose | [ ] |
| 2.4.7 Focus Visible | AA | Keyboard focus indicator is visible | [ ] |

**Testing steps for 2.4.1:**
- [ ] Skip link present at top of page
- [ ] Skip link visible on focus
- [ ] Skip link moves focus to main content
- [ ] Landmark regions properly identified

**Testing steps for 2.4.2:**
- [ ] Page title describes page content
- [ ] Page title is unique across site
- [ ] Title format: "Page Name - Site Name"

**Testing steps for 2.4.3:**
- [ ] Tab order follows logical reading order
- [ ] No unexpected focus jumps
- [ ] Focus moves predictably through interactive elements

**Testing steps for 2.4.4:**
- [ ] Link text describes destination
- [ ] No "click here" or "read more" without context
- [ ] Image links have descriptive alt text

**Testing steps for 2.4.5:**
- [ ] Site has navigation menu
- [ ] Site has search functionality or site map
- [ ] Multiple pathways to content exist

**Testing steps for 2.4.6:**
- [ ] Headings describe section content
- [ ] Form labels describe expected input
- [ ] Labels are specific and unique

**Testing steps for 2.4.7:**
- [ ] Focus indicator visible on all elements
- [ ] Focus indicator has sufficient contrast
- [ ] Focus indicator not obscured by other content

### 2.5 Input Modalities

| Criterion | Level | Requirement | Status |
|-----------|-------|-------------|--------|
| 2.5.1 Pointer Gestures | A | Multi-point or path-based gestures have single-point alternatives | [ ] |
| 2.5.2 Pointer Cancellation | A | Down-event does not trigger function; up-event can abort | [ ] |
| 2.5.3 Label in Name | A | Visible label text is included in accessible name | [ ] |
| 2.5.4 Motion Actuation | A | Motion-triggered functions have alternatives and can be disabled | [ ] |

**Testing steps for 2.5.1:**
- [ ] Pinch-to-zoom has button alternatives
- [ ] Swipe has button alternatives
- [ ] Path gestures have alternatives

**Testing steps for 2.5.2:**
- [ ] Actions trigger on mouse up, not mouse down
- [ ] Dragging can be cancelled by releasing outside target

**Testing steps for 2.5.3:**
- [ ] Button visible text matches accessible name
- [ ] Icon buttons have accessible names
- [ ] Links accessible names include visible text

**Testing steps for 2.5.4:**
- [ ] Shake-to-undo can be disabled
- [ ] Tilt features have alternatives

## 3. Understandable

Information and operation of user interface must be understandable.

### 3.1 Readable

| Criterion | Level | Requirement | Status |
|-----------|-------|-------------|--------|
| 3.1.1 Language of Page | A | Default language is programmatically identified | [ ] |
| 3.1.2 Language of Parts | AA | Language of passages or phrases is identified | [ ] |

**Testing steps for 3.1.1:**
- [ ] `<html lang="en">` attribute present
- [ ] Language code matches page content

**Testing steps for 3.1.2:**
- [ ] Foreign language passages have lang attribute
- [ ] Language changes marked inline

### 3.2 Predictable

| Criterion | Level | Requirement | Status |
|-----------|-------|-------------|--------|
| 3.2.1 On Focus | A | Focus does not cause unexpected context change | [ ] |
| 3.2.2 On Input | A | Input does not cause unexpected context change | [ ] |
| 3.2.3 Consistent Navigation | AA | Navigation is consistently ordered | [ ] |
| 3.2.4 Consistent Identification | AA | Components with same function are consistently identified | [ ] |

**Testing steps for 3.2.1:**
- [ ] Focus does not submit forms
- [ ] Focus does not open new windows
- [ ] Focus does not change page content significantly

**Testing steps for 3.2.2:**
- [ ] Selecting option does not submit form
- [ ] Entering data does not cause unexpected navigation
- [ ] Clear indication before context changes

**Testing steps for 3.2.3:**
- [ ] Navigation appears in same location
- [ ] Navigation items in same order
- [ ] Consistent across all pages

**Testing steps for 3.2.4:**
- [ ] Search icon always means search
- [ ] Similar functions have similar labels
- [ ] Icons consistent throughout site

### 3.3 Input Assistance

| Criterion | Level | Requirement | Status |
|-----------|-------|-------------|--------|
| 3.3.1 Error Identification | A | Input errors are identified and described in text | [ ] |
| 3.3.2 Labels or Instructions | A | Labels or instructions provided for user input | [ ] |
| 3.3.3 Error Suggestion | AA | Suggestions provided when input errors are detected | [ ] |
| 3.3.4 Error Prevention (Legal, Financial, Data) | AA | Reversible, verified, or confirmed for significant submissions | [ ] |

**Testing steps for 3.3.1:**
- [ ] Error messages identify the field in error
- [ ] Error messages describe the error
- [ ] Errors announced to screen readers

**Testing steps for 3.3.2:**
- [ ] All form fields have visible labels
- [ ] Required fields clearly indicated
- [ ] Format requirements stated (e.g., date format)

**Testing steps for 3.3.3:**
- [ ] Errors suggest valid input
- [ ] Examples provided where helpful
- [ ] Clear guidance for correction

**Testing steps for 3.3.4:**
For legal, financial, or data submissions:
- [ ] Data can be reviewed before submission
- [ ] Submission can be reversed, or
- [ ] Confirmation step before final submission

## 4. Robust

Content must be robust enough to be interpreted by a wide variety of user agents.

### 4.1 Compatible

| Criterion | Level | Requirement | Status |
|-----------|-------|-------------|--------|
| 4.1.1 Parsing | A | (Obsolete in WCAG 2.2 but relevant for 2.1) Markup is valid | [ ] |
| 4.1.2 Name, Role, Value | A | All UI components have accessible names, roles, values | [ ] |
| 4.1.3 Status Messages | AA | Status messages are announced without focus change | [ ] |

**Testing steps for 4.1.1:**
- [ ] HTML validates without errors
- [ ] No duplicate IDs
- [ ] Elements properly nested
- [ ] Opening/closing tags matched

**Testing steps for 4.1.2:**
- [ ] All interactive elements have accessible names
- [ ] Custom widgets have appropriate ARIA roles
- [ ] State changes are programmatically communicated
- [ ] Values are programmatically determinable

**Testing steps for 4.1.3:**
- [ ] Success messages use `role="status"` or `aria-live="polite"`
- [ ] Error messages use `role="alert"` or `aria-live="assertive"`
- [ ] Loading indicators announced
- [ ] Search results count announced

## Quick Reference Testing Checklist

### Keyboard Testing (Disconnect Mouse)
- [ ] Tab through entire page
- [ ] All interactive elements reachable
- [ ] Focus visible at all times
- [ ] No keyboard traps
- [ ] Modal Escape key works
- [ ] Dropdown arrow keys work

### Screen Reader Testing
- [ ] Page title announced
- [ ] Landmarks identified
- [ ] Headings convey structure
- [ ] Images have alt text
- [ ] Links are descriptive
- [ ] Forms labels announced
- [ ] Errors announced
- [ ] Dynamic content announced

### Visual Testing
- [ ] 200% zoom works
- [ ] Color contrast passes
- [ ] No color-only meaning
- [ ] Focus indicators visible
- [ ] Text is real text

### Automated Testing
- [ ] axe DevTools: 0 issues
- [ ] WAVE: 0 errors
- [ ] Lighthouse: 90+ accessibility score

## Severity Classification

When logging issues, classify severity:

| Severity | Definition | Example |
|----------|------------|---------|
| Critical | Blocks users completely | No keyboard access, no form labels |
| Major | Significant barrier | Low contrast, missing alt text |
| Minor | Inconvenience | Minor contrast issue, verbose alt text |
| Advisory | Best practice | Could be improved but not required |

## Sign-off

| Item | Auditor | Date | Status |
|------|---------|------|--------|
| Level A criteria | | | [ ] Pass [ ] Fail |
| Level AA criteria | | | [ ] Pass [ ] Fail |
| Keyboard testing | | | [ ] Pass [ ] Fail |
| Screen reader testing | | | [ ] Pass [ ] Fail |
| Visual testing | | | [ ] Pass [ ] Fail |
| Automated testing | | | [ ] Pass [ ] Fail |

**Overall WCAG 2.1 AA Compliance:** [ ] Compliant [ ] Non-compliant

**Notes:**
