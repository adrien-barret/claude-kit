---
description: Frontend rules applied when working on UI components, pages, and styles
globs: ["src/components/**", "src/pages/**", "src/views/**", "src/styles/**", "**/*.css", "**/*.scss", "**/*.tsx", "**/*.jsx"]
---

## Component Design

- Keep components small and focused on a single responsibility; split if exceeding ~150 lines
- Separate presentation (UI) from logic (hooks, services, stores)
- Handle loading, error, and empty states in every data-fetching component

## Accessibility (WCAG AA)

- Use semantic HTML elements (`<nav>`, `<main>`, `<button>`, `<label>`) instead of generic `<div>`
- Add ARIA attributes only when native semantics are insufficient
- Ensure all interactive elements are keyboard-navigable with visible focus indicators
- Maintain minimum 4.5:1 color contrast ratio for text; 3:1 for large text and UI components
- Provide descriptive `alt` text for images; use `alt=""` for decorative images

## Responsive Design and Performance

- Use mobile-first approach: base styles for small screens, `min-width` media queries for larger
- Lazy-load images and below-the-fold components to improve initial page load
- Use code splitting (dynamic `import()`) for routes and heavy components
- Use appropriate image formats (WebP with fallback) and responsive `srcset` attributes

## Styling

- Avoid inline styles; use the project's styling solution consistently
- Define design tokens (colors, spacing, typography) as variables; do not hardcode values

## State and Forms

- Keep state as local as possible; lift to parent or global store only when multiple components need it
- Implement client-side form validation with accessible error messages tied via `aria-describedby`
- Sanitize any user-provided content rendered in the DOM to prevent XSS
