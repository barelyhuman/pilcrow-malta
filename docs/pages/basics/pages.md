---
title: "Writing pages"
---

# Writing pages

Pages are written in markdown files in the `pages` directory.

## Markdown

Malta supports most standard markdown syntaxes. It also includes basic syntax highlighting for code blocks.

````md
Regular text.

_Italic_

**Bold**

# Headings 1~6

[links](https://example.com)

`code`

```ts
const message = "hello world";
```

- item 1
- item 2

1. item 1
2. item 2
````

## Attributes

Add pages must have a `title` attribute.

```md
---
title: "Malta documentation"
---
```
