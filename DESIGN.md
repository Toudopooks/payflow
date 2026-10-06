# PayFlow Interface System

PayFlow uses a restrained operations-console language: compact navigation, strong information hierarchy, quiet surfaces, and a single green accent. The visual reference is Linear's density and surface discipline, adapted to an original light-first fintech interface. No Linear colors, assets, copy, or proprietary type are used.

## Direction

- Design variance: 3. Predictable grids support financial review.
- Motion intensity: 2. Motion is limited to feedback and navigation state.
- Visual density: 6. Tables and identifiers remain readable without feeling sparse.
- Tone: factual, calm, and explicit about simulation behavior.

## Tokens

- Canvas: `#f5f7f8`
- Surface: `#ffffff`
- Text: `#20262e`
- Muted text: `#63707e`
- Border: `#d9dfe3`
- Accent: `#08765f`
- Radius: 8px for controls and 10px for containers
- Type: Geist when available, then Segoe UI and system sans-serif

Dark mode follows system preference using the same semantic hierarchy. Status colors always include a text label and never rely on color alone.

## Interaction Rules

- Controls are at least 42px high and 44px on mobile.
- Every interactive element has a visible `:focus-visible` outline.
- Animations only change opacity or transform and stop under reduced motion.
- Tables scroll horizontally on narrow screens instead of clipping data.
- Financial values use tabular numerals and locale-aware formatting.
- Forms place labels above fields and errors below them.
- Loading skeletons match the final row shape.

