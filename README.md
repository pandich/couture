![Couture](docs/couture.png)

---

Couture combines multiple log streams from disparate sources into human-friendly output.

### Installation

	make install
	make install-shell-completions # optional

### Usage:

	couture --help

### Configuration

[comment]: <> (TODO config doc)

---

### Examples:

#### Single-line

	couture --highlight --filter=+distincto --filter=+'\"first_name\"\s*:\s*\"B' --filter=+quinoa @@fake

![Couture](docs/example-fake-single-line.gif)

#### Multi-line

	couture --highlight --filter=+distincto --filter=+'\"first_name\"\s*:\s*\"B' --filter=+quinoa --expand --multiline @@fake

![Couture](docs/example-fake-multi-line.gif)
