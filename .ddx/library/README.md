# DDx Library

This repository contains the DDx (Document-Driven Development eXperience) library of templates, patterns, prompts, and configurations for AI-assisted development.

## Contents

- **prompts/** - AI prompts and instructions for various development tasks
- **templates/** - Project templates and boilerplates
- **patterns/** - Reusable code patterns and solutions
- **personas/** - AI personality definitions for consistent interactions
- **mcp-servers/** - Model Context Protocol server configurations
- **configs/** - Tool configurations (ESLint, Prettier, TypeScript, etc.)
- **workflows/** - Complete development methodologies (HELIX, etc.)
- **tools/** - Development tool integrations and scripts
- **environments/** - Development environment configurations

## Usage

This library is designed to be used with the [DDx CLI tool](https://github.com/easel/ddx). The CLI automatically syncs this library into your projects using git subtree for bidirectional synchronization.

### Installation via DDx CLI

```bash
# Install DDx CLI
curl -fsSL https://raw.githubusercontent.com/easel/ddx/main/install.sh | bash

# Initialize in your project (automatically pulls this library)
cd your-project
ddx init
```

### Manual Usage

You can also use git subtree directly to include this library in your projects:

```bash
# Add the library to your project
git subtree add --prefix=.ddx https://github.com/easel/ddx-library main --squash

# Update to latest library version
git subtree pull --prefix=.ddx https://github.com/easel/ddx-library main --squash

# Contribute improvements back
git subtree push --prefix=.ddx https://github.com/easel/ddx-library main
```

## Contributing

We welcome contributions! To contribute improvements:

1. **Via DDx CLI**: Make changes to your local `.ddx/` directory and run `ddx contribute`
2. **Via GitHub**: Fork this repository, make changes, and submit a pull request
3. **Via git subtree**: Make changes locally and use `git subtree push` to contribute back

## License

This library is open source software licensed under the MIT License.

## Related Projects

- [DDx CLI](https://github.com/easel/ddx) - The command-line tool for using this library
- [HELIX Workflow](./workflows/helix/) - Six-phase AI-assisted development methodology