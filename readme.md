# unsafe-markdown
A Regexp-based Markdown to Html Converter that expects input text to have already been sanitized.

Ideal for use in simple conversions where the output can be verified / is not critical.

### Supported Markdown Syntax:
* Bolding
* Italics
* Strikethrough
* Headers
* Blockquotes
* Images
* Links
* Linebreaks (all kinds specified by Unicode standard)
* Code tag (single line only)

### Known Issues

* No unordered/ordered list support
* No multi-line support
* The repeating nature of the converter means that some nested Markdown syntax cases will be converted  
* Undefined behavior for some edge cases
* Not enough tests!