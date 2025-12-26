# Accessibility Checklist

WCAG 2.1 Level AA compliance checklist for inclusive design.

## Perceivable

### Text Alternatives (1.1)

- [ ] All images have meaningful alt text describing content or function
- [ ] Decorative images have empty alt (`alt=""`) or are CSS backgrounds
- [ ] Complex images (charts, diagrams) have extended descriptions
- [ ] Icons used as controls have accessible labels
- [ ] CAPTCHAs provide alternative methods

### Time-Based Media (1.2)

- [ ] Videos have captions
- [ ] Videos have audio descriptions for visual-only content
- [ ] Live audio content has captions
- [ ] Pre-recorded audio has transcripts

### Adaptable (1.3)

- [ ] Content structure uses semantic HTML (headings, lists, tables)
- [ ] Heading hierarchy is logical (h1 -> h2 -> h3)
- [ ] Tables have proper headers (`<th>` with scope)
- [ ] Form inputs have associated labels
- [ ] Reading order is logical when CSS is disabled
- [ ] Instructions do not rely solely on sensory characteristics (color, shape)

### Distinguishable (1.4)

- [ ] Color is not the only means of conveying information
- [ ] Text color contrast ratio is at least 4.5:1 (3:1 for large text)
- [ ] UI component contrast ratio is at least 3:1
- [ ] Text can be resized to 200% without loss of content
- [ ] Text is not in images (except logos)
- [ ] Content reflows at 320px width without horizontal scrolling
- [ ] Line height is at least 1.5, paragraph spacing at least 2x font size
- [ ] No content is lost when user overrides text spacing

## Operable

### Keyboard Accessible (2.1)

- [ ] All functionality is available via keyboard
- [ ] No keyboard traps (user can navigate away from all components)
- [ ] Keyboard shortcuts can be turned off or remapped
- [ ] Focus order is logical and intuitive
- [ ] Custom components have appropriate keyboard interactions

### Enough Time (2.2)

- [ ] Time limits can be turned off, adjusted, or extended
- [ ] Auto-updating content can be paused, stopped, or hidden
- [ ] No timing-dependent interactions unless essential
- [ ] Session timeouts warn users and allow extension

### Seizures and Physical Reactions (2.3)

- [ ] No content flashes more than 3 times per second
- [ ] Motion animations can be disabled (respect prefers-reduced-motion)

### Navigable (2.4)

- [ ] Skip link allows bypassing repetitive content
- [ ] Pages have descriptive titles
- [ ] Focus order matches visual order
- [ ] Link purpose is clear from link text (not "click here")
- [ ] Multiple ways to find pages (navigation, search, sitemap)
- [ ] Headings and labels describe content
- [ ] Focus indicator is visible

### Input Modalities (2.5)

- [ ] Touch targets are at least 44x44 CSS pixels
- [ ] Functionality requiring multi-point gestures has alternatives
- [ ] Functionality requiring motion has alternatives
- [ ] Dragging operations have single-pointer alternatives

## Understandable

### Readable (3.1)

- [ ] Page language is set (`<html lang="en">`)
- [ ] Language changes within content are marked
- [ ] Unusual words or jargon are defined
- [ ] Abbreviations are expanded on first use

### Predictable (3.2)

- [ ] Focus does not trigger unexpected context changes
- [ ] Input does not automatically trigger context changes
- [ ] Navigation is consistent across pages
- [ ] Components with same function are identified consistently

### Input Assistance (3.3)

- [ ] Error messages identify the field and describe the error
- [ ] Form fields have labels and instructions
- [ ] Error suggestions are provided when known
- [ ] Users can review, correct, and confirm submissions
- [ ] Help is available for complex inputs

## Robust

### Compatible (4.1)

- [ ] HTML is valid and well-formed
- [ ] Custom components have proper ARIA roles, states, and properties
- [ ] Status messages are announced by screen readers (aria-live)
- [ ] Name, role, and value are programmatically determinable

## Component-Specific Checks

### Forms

- [ ] All inputs have visible labels
- [ ] Required fields are indicated (not by color alone)
- [ ] Error messages are associated with inputs
- [ ] Form validation errors are announced
- [ ] Autocomplete attributes are used appropriately

### Modals and Dialogs

- [ ] Focus moves to modal when opened
- [ ] Focus is trapped within modal while open
- [ ] Focus returns to trigger element when closed
- [ ] Modal can be closed with Escape key
- [ ] Background content is hidden from screen readers (aria-hidden)

### Navigation Menus

- [ ] Current page/section is indicated
- [ ] Dropdowns are keyboard accessible
- [ ] Expanded/collapsed state is announced (aria-expanded)
- [ ] Mobile menu is accessible

### Data Tables

- [ ] Headers are properly associated with cells
- [ ] Complex tables use id/headers attributes
- [ ] Sortable columns indicate sort state
- [ ] Tables are not used for layout

### Carousels and Sliders

- [ ] Pause control is available for auto-advancing
- [ ] Navigation controls are keyboard accessible
- [ ] Current slide is announced to screen readers
- [ ] Slide content is accessible when visible

## Testing Methods

### Automated Testing

- [ ] axe-core or similar tool runs in CI
- [ ] Lighthouse accessibility audit passes
- [ ] HTML validation passes

### Manual Testing

- [ ] Navigate entire interface with keyboard only
- [ ] Test with screen reader (VoiceOver, NVDA)
- [ ] Test with browser zoom at 200%
- [ ] Test with high contrast mode
- [ ] Test with prefers-reduced-motion enabled

## Usage Notes

1. Accessibility is a legal requirement in many jurisdictions
2. Not every item applies to every component - use judgment
3. Automated tools catch only ~30% of issues - manual testing is essential
4. Include users with disabilities in usability testing when possible
5. Document exceptions with rationale and remediation plan
6. Accessibility should be part of the definition of done
