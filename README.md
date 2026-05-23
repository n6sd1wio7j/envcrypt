# envcrypt

Utility to encrypt and version-control `.env` files using [age](https://github.com/FiloSottile/age) encryption with team key management support.

---

## Installation

```bash
go install github.com/yourusername/envcrypt@latest
```

Or download a prebuilt binary from the [releases page](https://github.com/yourusername/envcrypt/releases).

---

## Usage

**Encrypt a `.env` file:**
```bash
envcrypt encrypt --in .env --out .env.age --recipients keys.txt
```

**Decrypt a `.env` file:**
```bash
envcrypt decrypt --in .env.age --out .env --key ~/.age/key.txt
```

**Add a team member's public key:**
```bash
envcrypt keys add --key "age1ql3z7hjy..." --recipients keys.txt
```

**Re-encrypt for all current recipients after a key change:**
```bash
envcrypt reencrypt --in .env.age --recipients keys.txt
```

The encrypted `.env.age` and `keys.txt` files are safe to commit to version control. Private keys should **never** be committed.

---

## Workflow

1. Add `.env` to `.gitignore`
2. Encrypt with `envcrypt encrypt` and commit `.env.age`
3. Share public keys via `keys.txt` in the repo
4. Each team member decrypts locally using their private key

---

## Requirements

- Go 1.21+
- [age](https://github.com/FiloSottile/age) (`brew install age` or `apt install age`)

---

## License

MIT © [yourusername](https://github.com/yourusername)