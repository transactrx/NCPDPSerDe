# Scripts

## gen_schema_docx.py — NCPDP Schema Documentation Generator

Generates `NCPDP-JSON-Schema-Reference.docx`, a Word reference document built
from the JSON Schemas in `pkg/ncpdp/schemas/`. The document is intended for
upload to Confluence and covers:

- Overview (transaction codes, key concepts)
- Every request/response transaction: structure tables, header fields,
  per-claim segments, and a full JSON layout skeleton
- Every shared segment/type definition (`ncpdp-schemas.json` `$defs`):
  field tables with NCPDP field codes plus a fully expanded JSON layout
- Example payloads from `pkg/ncpdp/schemas/examples/`

### Prerequisites

- Python 3.10+
- `python-docx`:

  ```
  pip install python-docx
  ```

  On machines behind a TLS-intercepting proxy, if pip fails with an SSL
  certificate error:

  ```
  pip install --trusted-host pypi.org --trusted-host files.pythonhosted.org python-docx
  ```

### Regenerating the document

Whenever anything under `pkg/ncpdp/schemas/` changes (new transaction, new
segment, field additions, description/constraint updates), rerun from the
repo root:

```
python scripts/gen_schema_docx.py
```

This rewrites `NCPDP-JSON-Schema-Reference.docx` in the repo root. Optional
arguments override the input directory and output path:

```
python scripts/gen_schema_docx.py <schemas-dir> <output.docx>
```

The script needs no configuration for new schema files — it picks up every
`*.json` in `request/`, `response/`, and `examples/` automatically, and every
definition in `ncpdp-schemas.json` `$defs`. To control where a new
transaction appears in the document, add its file stem to `FILE_ORDER` at the
top of the script; unlisted files are appended at the end.

Note: field tables derive the "NCPDP Code" column from the trailing
parenthesized code in each field's `description` (e.g. `"Cardholder ID (C2)"`),
so keep that convention when adding fields to the schemas.

### Uploading to Confluence

1. Open (or create) the target Confluence page.
2. From the page's **•••** menu choose **Import Word Document**.
3. Upload `NCPDP-JSON-Schema-Reference.docx`. Confluence can optionally split
   the document into child pages by heading level — the document uses
   Heading 1 for major sections and Heading 2 per transaction, so splitting
   at Heading 1 or 2 both work well.
4. Re-importing over an existing page replaces its content; re-run the import
   after each regeneration to keep Confluence in sync.
