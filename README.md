# Parse SQL to JSON (via pg_query)

This Go program reads an SQL script and converts it to JSON using [**pg_query**](https://github.com/pganalyze/pg_query_go). It supports:

- **File** input (by specifying a path as a positional argument)
- **STDIN** input (when no file is given)
- **Clipboard** input (via `--use-clipboard` or `-c`)
- **Pretty-printing** and indentation (`--no-pretty` to disable, `--indent` for specifying spaces or tab)
- **Colorized JSON** output (`--no-color` to disable)
- **Output** to a specified file (`--out=FILE`) or STDOUT by default

## Features

1. **SQL** → **JSON** using [github.com/pganalyze/pg_query_go/v6](https://github.com/pganalyze/pg_query_go).
2. **Clipboard** support with [github.com/atotto/clipboard](https://github.com/atotto/clipboard).
3. **Colorized** output using [github.com/TylerBrock/colorjson](https://github.com/TylerBrock/colorjson).
4. **Configurable indentation**:
   - `--indent=tab`: Indent with real tab characters in non-colorized output.
   - `--indent=2` or any integer: Indent with that many spaces per level.
5. **Smart defaults**:
   - If you **do not** specify `--no-pretty`, the JSON is indented.
   - If you **specify** `--out`, the JSON is also indented by default (unless `--no-pretty` is given).
   - If you **do not** specify `--no-color`, JSON is colorized (unless you choose minimal output).

## Installation

```bash
# Clone this repository or copy the main.go file
# Then install dependencies:
go get github.com/pganalyze/pg_query_go/v6
go get github.com/atotto/clipboard
go get github.com/TylerBrock/colorjson

# Build the tool:
go build -o parse_sql_to_json main.go
```

This will produce a binary named `parse_sql_to_json`.

## Usage

Basic usage:

```bash
./parse_sql_to_json [options] [filePath]
```

- `[filePath]` (positional argument) is **optional**. If omitted, the tool reads from `STDIN` by default.
- If `--use-clipboard` / `-c` is set, it will read from the **system clipboard** instead, ignoring file/STDIN.

### Command-line Flags

| Flag                | Short | Default | Description                                                                                              |
|---------------------|-------|---------|----------------------------------------------------------------------------------------------------------|
| `--no-pretty`       | (none)| `false` | If set, **disable** pretty-printing (indentation).                                                      |
| `--no-color`        | (none)| `false` | If set, **disable** colorized output.                                                                   |
| `--out`             | (none)| (empty) | Write output to a file. If not specified, output goes to `STDOUT`.                                      |
| `--indent`          | (none)| `"tab"` | Indentation style for pretty-printing. Accepts `tab` or an integer (e.g., `2`).                         |
| `--use-clipboard`   | `-c`  | `false` | If set, read the SQL script from the system clipboard. Overrides file/STDIN input.                      |

### Examples

1. **Parse from a file**:
   ```bash
   ./parse_sql_to_json my_script.sql
   ```
   - Reads `my_script.sql`, parses to JSON, colorizes, and prints to STDOUT with indentation.

2. **Parse from STDIN**:
   ```bash
   cat my_script.sql | ./parse_sql_to_json
   ```
   - Same outcome as above, but reading from a pipe.

3. **Clipboard**:
   ```bash
   ./parse_sql_to_json --use-clipboard
   ```
   or
   ```bash
   ./parse_sql_to_json -c
   ```
   - Reads whatever SQL content is in your clipboard, parses to JSON, and prints it out.

4. **Disable pretty-printing**:
   ```bash
   ./parse_sql_to_json --no-pretty my_script.sql
   ```
   - Resulting JSON is a single line (unless you also request color, in which case it will still be colorized, but not indented).

5. **Disable color**:
   ```bash
   ./parse_sql_to_json --no-color my_script.sql
   ```
   - Output will be plain JSON without ANSI color codes.

6. **Set custom indentation**:
   ```bash
   ./parse_sql_to_json --indent=4 my_script.sql
   ```
   - Indent with 4 spaces per level.

7. **Write to an output file**:
   ```bash
   ./parse_sql_to_json --out=output.json my_script.sql
   ```
   - Creates (or truncates) `output.json` with the resulting JSON. Pretty-printed by default.

8. **Combine flags**:
   ```bash
   ./parse_sql_to_json --no-color --indent=2 --out=output.json my_script.sql
   ```
   - Writes to `output.json`, using 2 spaces per indent level, with no color codes.

## License

Distributed under the MIT License. See `LICENSE` for more information, if applicable.

## Contributing

Feel free to open an issue or pull request to contribute improvements or fixes. If you have questions about usage, please reach out!
