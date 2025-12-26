---
name: accessibility-design
description: WCAG 2.1 AA compliance patterns, screen reader compatibility, keyboard navigation, and ARIA best practices. Use when implementing accessible interfaces, reviewing UI components, or auditing accessibility compliance. Covers semantic HTML, focus management, color contrast, and assistive technology testing.
---

# Accessibility Standards

## Overview

This skill provides comprehensive guidance for implementing WCAG 2.1 AA compliance across digital products. It establishes patterns for semantic markup, assistive technology compatibility, and inclusive interaction design.

## Core Principles

### 1. POUR Framework

All accessibility work follows the four POUR principles:

**Perceivable** - Information must be presentable in ways users can perceive
- Provide text alternatives for non-text content
- Create content adaptable to different presentations
- Make content distinguishable (color, contrast, audio control)

**Operable** - Interface components must be operable by all users
- Make all functionality keyboard accessible
- Provide sufficient time to read and use content
- Avoid content that causes seizures or physical reactions
- Help users navigate, find content, and determine location

**Understandable** - Information and operation must be understandable
- Make text readable and understandable
- Make pages appear and operate predictably
- Help users avoid and correct mistakes

**Robust** - Content must be robust enough for diverse user agents
- Maximize compatibility with current and future assistive technologies
- Use valid, semantic markup
- Ensure programmatic access to all functionality

### 2. Semantic HTML First

Semantic HTML is the foundation of accessibility. ARIA should enhance, never replace, proper markup.

**Correct approach:**
```html
<button type="submit">Submit Form</button>
```

**Incorrect approach:**
```html
<div role="button" tabindex="0" onclick="submit()">Submit Form</div>
```

Use native HTML elements whenever possible:
- `<button>` for actions
- `<a href>` for navigation
- `<input>`, `<select>`, `<textarea>` for form controls
- `<nav>`, `<main>`, `<aside>`, `<header>`, `<footer>` for landmarks
- `<h1>` through `<h6>` for heading hierarchy

### 3. Progressive Enhancement

Build accessibility in layers:
1. Semantic HTML provides baseline accessibility
2. CSS enhances presentation without breaking structure
3. JavaScript adds interactivity while maintaining keyboard access
4. ARIA fills gaps for complex widgets not covered by HTML

## Implementation Patterns

### Document Structure

Every page requires proper document structure:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Descriptive Page Title - Site Name</title>
</head>
<body>
  <a href="#main-content" class="skip-link">Skip to main content</a>

  <header role="banner">
    <nav aria-label="Main navigation">
      <!-- Navigation content -->
    </nav>
  </header>

  <main id="main-content" role="main">
    <h1>Page Heading</h1>
    <!-- Page content -->
  </main>

  <footer role="contentinfo">
    <!-- Footer content -->
  </footer>
</body>
</html>
```

### Heading Hierarchy

Maintain logical heading order without skipping levels:

```html
<h1>Main Page Title</h1>
  <h2>Section Title</h2>
    <h3>Subsection Title</h3>
    <h3>Another Subsection</h3>
  <h2>Another Section</h2>
    <h3>Subsection</h3>
      <h4>Sub-subsection</h4>
```

### Skip Links

Provide skip links for keyboard users to bypass repetitive content:

```css
.skip-link {
  position: absolute;
  top: -40px;
  left: 0;
  padding: 8px 16px;
  background: #000;
  color: #fff;
  z-index: 100;
}

.skip-link:focus {
  top: 0;
}
```

### Focus Management

#### Visible Focus Indicators

Focus indicators must be visible and meet contrast requirements:

```css
:focus {
  outline: 2px solid #005fcc;
  outline-offset: 2px;
}

/* Enhanced focus for better visibility */
:focus-visible {
  outline: 3px solid #005fcc;
  outline-offset: 3px;
  box-shadow: 0 0 0 6px rgba(0, 95, 204, 0.25);
}

/* Never remove focus indicators entirely */
:focus:not(:focus-visible) {
  outline: 2px solid transparent;
  box-shadow: 0 0 0 2px #005fcc;
}
```

#### Focus Trapping for Modals

```javascript
function trapFocus(element) {
  const focusableElements = element.querySelectorAll(
    'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
  );
  const firstFocusable = focusableElements[0];
  const lastFocusable = focusableElements[focusableElements.length - 1];

  element.addEventListener('keydown', (e) => {
    if (e.key !== 'Tab') return;

    if (e.shiftKey) {
      if (document.activeElement === firstFocusable) {
        lastFocusable.focus();
        e.preventDefault();
      }
    } else {
      if (document.activeElement === lastFocusable) {
        firstFocusable.focus();
        e.preventDefault();
      }
    }
  });
}
```

### Keyboard Navigation

#### Standard Keyboard Patterns

| Key | Action |
|-----|--------|
| Tab | Move focus to next focusable element |
| Shift+Tab | Move focus to previous focusable element |
| Enter/Space | Activate buttons and links |
| Arrow keys | Navigate within components (menus, tabs, radio groups) |
| Escape | Close modals, menus, dropdowns |
| Home/End | Move to first/last item in lists |

#### Tab Order

Ensure logical tab order follows visual order:

```html
<!-- Correct: Tab order matches visual order -->
<form>
  <label for="name">Name</label>
  <input id="name" type="text">

  <label for="email">Email</label>
  <input id="email" type="email">

  <button type="submit">Submit</button>
</form>
```

Never use positive `tabindex` values. Only use:
- `tabindex="0"` to make non-interactive elements focusable
- `tabindex="-1"` to make elements programmatically focusable but not in tab order

### ARIA Patterns

#### When to Use ARIA

Use ARIA only when native HTML cannot achieve the required accessibility:

1. **Custom widgets** - Tabs, accordions, carousels, tree views
2. **Dynamic content** - Live regions for updates
3. **Relationships** - Connecting labels to complex controls
4. **States** - Expanded/collapsed, selected, pressed

#### ARIA Roles

Common roles for custom widgets:

```html
<!-- Tab pattern -->
<div role="tablist" aria-label="Settings tabs">
  <button role="tab" aria-selected="true" aria-controls="panel1" id="tab1">
    General
  </button>
  <button role="tab" aria-selected="false" aria-controls="panel2" id="tab2">
    Security
  </button>
</div>
<div role="tabpanel" id="panel1" aria-labelledby="tab1">
  <!-- Panel content -->
</div>
<div role="tabpanel" id="panel2" aria-labelledby="tab2" hidden>
  <!-- Panel content -->
</div>
```

#### ARIA States and Properties

| Attribute | Purpose | Example |
|-----------|---------|---------|
| `aria-expanded` | Indicates expandable element state | `aria-expanded="false"` |
| `aria-selected` | Indicates selection state | `aria-selected="true"` |
| `aria-pressed` | Indicates toggle button state | `aria-pressed="mixed"` |
| `aria-hidden` | Hides content from assistive tech | `aria-hidden="true"` |
| `aria-live` | Announces dynamic content | `aria-live="polite"` |
| `aria-describedby` | References descriptive content | `aria-describedby="hint"` |
| `aria-labelledby` | References labeling content | `aria-labelledby="heading"` |

#### Live Regions

Announce dynamic content changes:

```html
<!-- Polite: Announces when user is idle -->
<div aria-live="polite" aria-atomic="true">
  Form saved successfully
</div>

<!-- Assertive: Interrupts immediately (use sparingly) -->
<div aria-live="assertive" role="alert">
  Error: Please correct the highlighted fields
</div>

<!-- Status messages -->
<div role="status">
  Loading... 50% complete
</div>
```

### Form Accessibility

#### Labels and Instructions

Every form control requires a label:

```html
<!-- Explicit label association -->
<label for="username">Username</label>
<input type="text" id="username" name="username"
       aria-describedby="username-hint">
<span id="username-hint">Must be 3-20 characters</span>

<!-- Grouped controls with fieldset -->
<fieldset>
  <legend>Shipping Address</legend>
  <label for="street">Street</label>
  <input type="text" id="street">
  <!-- More fields -->
</fieldset>
```

#### Error Handling

Provide clear, actionable error messages:

```html
<label for="email">Email</label>
<input type="email" id="email"
       aria-invalid="true"
       aria-describedby="email-error">
<span id="email-error" role="alert">
  Please enter a valid email address (example: user@domain.com)
</span>
```

Error message requirements:
- Identify the field in error
- Describe what went wrong
- Provide guidance for correction
- Use `aria-invalid` on the field
- Announce errors via live region or `role="alert"`

#### Required Fields

```html
<label for="name">
  Name <span aria-hidden="true">*</span>
  <span class="sr-only">(required)</span>
</label>
<input type="text" id="name" required aria-required="true">
```

### Images and Media

#### Alternative Text

Provide meaningful alt text for informative images:

```html
<!-- Informative image -->
<img src="chart.png" alt="Q3 revenue increased 15% compared to Q2">

<!-- Decorative image -->
<img src="divider.png" alt="" role="presentation">

<!-- Complex image with extended description -->
<figure>
  <img src="flowchart.png" alt="User registration process"
       aria-describedby="flowchart-desc">
  <figcaption id="flowchart-desc">
    The registration process begins with email verification,
    followed by profile creation, and ends with account activation.
  </figcaption>
</figure>
```

#### Video and Audio

```html
<video controls>
  <source src="video.mp4" type="video/mp4">
  <track kind="captions" src="captions.vtt" srclang="en" label="English">
  <track kind="descriptions" src="descriptions.vtt" srclang="en"
         label="Audio descriptions">
</video>
```

### Color and Contrast

#### Minimum Contrast Ratios

| Content Type | Minimum Ratio (AA) | Enhanced Ratio (AAA) |
|--------------|-------------------|---------------------|
| Normal text (< 18pt) | 4.5:1 | 7:1 |
| Large text (>= 18pt or 14pt bold) | 3:1 | 4.5:1 |
| UI components and graphics | 3:1 | N/A |

#### Never Rely on Color Alone

```html
<!-- Bad: Color only indicates error -->
<input style="border-color: red">

<!-- Good: Color plus icon and text -->
<input aria-invalid="true" aria-describedby="error">
<span id="error">
  <svg aria-hidden="true"><!-- Error icon --></svg>
  Invalid email format
</span>
```

### Motion and Animation

#### Respect User Preferences

```css
/* Reduce motion for users who prefer it */
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}
```

#### Auto-playing Content

- Avoid auto-playing video or audio
- Provide controls to pause, stop, or hide moving content
- Limit animations to under 5 seconds or provide stop mechanism

## Testing Methodology

### Automated Testing

Run automated tools to catch common issues:

1. **axe DevTools** - Browser extension for page analysis
2. **WAVE** - Web accessibility evaluation tool
3. **Lighthouse** - Chrome DevTools accessibility audit
4. **pa11y** - Command-line accessibility testing

Automated testing catches approximately 30-40% of issues.

### Manual Testing

#### Keyboard Testing

1. Disconnect or disable mouse
2. Navigate entire page using Tab/Shift+Tab
3. Verify all interactive elements are reachable
4. Verify focus indicators are visible
5. Test all keyboard shortcuts
6. Escape from all modals and menus

#### Screen Reader Testing

Test with multiple screen readers:

| Platform | Screen Reader | Browser |
|----------|---------------|---------|
| Windows | NVDA (free) | Firefox, Chrome |
| Windows | JAWS | Chrome, Edge |
| macOS | VoiceOver | Safari |
| iOS | VoiceOver | Safari |
| Android | TalkBack | Chrome |

Screen reader testing checklist:
- [ ] All images have appropriate alt text
- [ ] Headings convey page structure
- [ ] Links and buttons have descriptive names
- [ ] Form fields have labels
- [ ] Error messages are announced
- [ ] Dynamic content updates are announced
- [ ] Tables have proper headers
- [ ] Custom widgets announce state changes

#### Visual Testing

- Zoom to 200% and verify no horizontal scrolling
- Test with high contrast mode enabled
- Disable images and verify content is still understandable
- Test with browser text size increased

### WCAG Audit

Perform manual audit against WCAG 2.1 AA success criteria. See `checklists/wcag-checklist.md` for the complete checklist.

## Common Patterns

### Modal Dialog

```html
<div role="dialog" aria-modal="true" aria-labelledby="dialog-title">
  <h2 id="dialog-title">Confirm Action</h2>
  <p>Are you sure you want to proceed?</p>
  <button type="button">Cancel</button>
  <button type="button">Confirm</button>
</div>
```

Requirements:
- Focus moves to dialog when opened
- Focus is trapped within dialog
- Escape key closes dialog
- Focus returns to trigger element on close

### Accordion

```html
<div class="accordion">
  <h3>
    <button aria-expanded="true" aria-controls="section1">
      Section 1
    </button>
  </h3>
  <div id="section1" role="region" aria-labelledby="section1-heading">
    <!-- Content -->
  </div>
</div>
```

### Navigation Menu

```html
<nav aria-label="Main">
  <ul>
    <li><a href="/" aria-current="page">Home</a></li>
    <li>
      <button aria-expanded="false" aria-haspopup="true">
        Products
      </button>
      <ul>
        <li><a href="/products/a">Product A</a></li>
        <li><a href="/products/b">Product B</a></li>
      </ul>
    </li>
  </ul>
</nav>
```

### Data Table

```html
<table>
  <caption>Q3 Sales by Region</caption>
  <thead>
    <tr>
      <th scope="col">Region</th>
      <th scope="col">Units Sold</th>
      <th scope="col">Revenue</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <th scope="row">North</th>
      <td>1,234</td>
      <td>$45,678</td>
    </tr>
  </tbody>
</table>
```

## Resources

- [WCAG 2.1 Guidelines](https://www.w3.org/WAI/WCAG21/quickref/)
- [WAI-ARIA Authoring Practices](https://www.w3.org/WAI/ARIA/apg/)
- [Inclusive Components](https://inclusive-components.design/)
- [A11y Project Checklist](https://www.a11yproject.com/checklist/)
