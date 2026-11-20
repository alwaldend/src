## 1. Search field

- [x] 1.1 Reserve consistent space for the search icon.
- [x] 1.2 Verify search field spacing before and during focus.

## Evidence

Chromium measured an 11px gap between the icon box and text in the homepage
header at 1024px and 1280px, and in the documentation sidebar at 1280px.
The gap held before and during focus. The offline search browser test passed.
