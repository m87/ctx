# ctx-go

![Build](https://github.com/YourGitHubUsername/ctx-go/actions/workflows/ci.yml/badge.svg)
![License](https://img.shields.io/github/license/YourGitHubUsername/ctx-go)
![Issues](https://img.shields.io/github/issues/YourGitHubUsername/ctx-go)
![Last Commit](https://img.shields.io/github/last-commit/YourGitHubUsername/ctx-go)

> **ctx-go** – a modern dashboard for context management, built with Go (backend) and React + TypeScript + Nx + Vite (frontend).

---

## Table of Contents

- [About](#about)
- [Demo](#demo)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Testing](#testing)
- [CI/CD](#cicd)
- [Contributing](#contributing)
- [License](#license)

---

## About

**ctx-go** is a web application for managing work contexts, featuring a fast interface, multiple views, and backend integration.  
The project uses a modern stack: Go (backend), React 19, TypeScript, Nx, Vite, TailwindCSS, Radix UI (shadcn/ui), React Query.

---

## Demo

> *(Add a link to a live demo or screenshots here)*

---

## Features

- Manage contexts (create, switch, browse)
- Timeline and daily summary views
- Fast search and filtering
- Modern, responsive UI (Tailwind, shadcn/ui)
- E2E tests (Cypress)
- Automated CI/CD (GitHub Actions)

---

## Installation

### Backend (Go)

```bash
git clone https://github.com/YourGitHubUsername/ctx-go.git
cd ctx-go
go build -o ctx-go ./cmd/server
./ctx-go
```

By default, the backend runs on [http://localhost:8080](http://localhost:8080).

### Frontend

```bash
cd server/ui/ctx-dashboard
npm install
npm start
# or with Nx:
nx serve ctx-dashboard
```

Frontend runs on [http://localhost:4200](http://localhost:4200).

---

## Testing

### Backend

```bash
go test ./...
```

### Frontend

```bash
npm test
# or
nx test
```

---

## CI/CD

This project uses GitHub Actions for automatic build, test, and lint on every push and pull request.

---

## Contributing

Want to help?  
Report issues, suggest features, or open a pull request!  
See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

---

## License

MIT © [m87](https://github.com/m87)

---

<!--
Suggestions:
- Add a "Screenshots" section with UI images
- Add a "Roadmap" section for planned features
- Add "FAQ" or "Known Issues"
- Add a link to backend API documentation if available
-->