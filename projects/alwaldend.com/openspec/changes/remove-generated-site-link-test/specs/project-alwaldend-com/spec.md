## MODIFIED Requirements

### Requirement: Resolve links relative to documentation sources

The site SHALL resolve Markdown links and images relative to their source
directory, map packaged documentation pages and resources to published URLs,
and preserve unknown internal destinations rather than silently redirecting
them to GitHub. Resolution behavior SHALL remain defined by the site layout,
not by a generated-output link test.

#### Scenario: Documentation contains an unresolved internal link

- **WHEN** a Markdown destination has no matching packaged page or resource
- **THEN** it remains unresolved in the rendered output rather than being
  silently redirected to GitHub
