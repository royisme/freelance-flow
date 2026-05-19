# Tally - FreelanceFlow

A desktop application for Canadian freelancers to manage clients, projects, timesheets, and invoices. Built with CRA tax compliance in mind.

## Features

- **Client Management**: Track client details and project history
- **Projects**: Manage project lifecycle from start to completion
- **Timesheets**: Log work hours with billable/non-billable tracking
- **Invoices**: Generate professional invoices with HST/GST support
- **Dashboard**: At-a-glance business metrics and summaries

## Tech Stack

| Component | Technology |
|-----------|------------|
| Backend | Go (Wails v2) |
| Frontend | Vue 3, TypeScript |
| Styling | Tailwind CSS |
| Runtime | Bun |
| Database | SQLite |
| State Management | Pinia |

## Prerequisites

- **Go**: 1.21+ ([Install](https://go.dev/doc/install))
- **Wails CLI**: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- **Bun**: [Install](https://bun.sh/)
- **Xcode Command Line Tools**: `xcode-select --install` (macOS only)

## Quick Start

```bash
# Clone and navigate
git clone https://github.com/yourusername/freelance-flow.git
cd freelance-flow

# Install frontend dependencies
cd frontend && bun install && cd ..

# Run in development mode (hot-reload enabled)
wails dev
```

The app will open at `http://localhost:34115`.

## Build

```bash
# Build for your current platform
wails build

# Build for Windows
wails build -platform windows/amd64

# Build for macOS
wails build -platform darwin/universal
```

Executables are generated in `build/bin/`.

## Canadian Freelancer Setup

### Registering for HST/GST with CRA

1. **When to register**: If your annual revenue exceeds $30,000, you must register for GST/HST.
2. **How to register**: Visit [CRA My Business Account](https://www.canada.ca/en/revenue-agency.html) or call 1-800-959-5525.
3. **After registration**: You'll receive a GST/HST account number (format: RT-XXXX-XXXXXXX).

### Configuring Tax Settings in Tally

1. Go to **Settings** > **Tax Configuration**
2. Enter your GST/HST rate (5% GST, 13% HST Ontario, or your provincial rate)
3. Add your CRA Business Number as the tax ID on invoices
4. Select whether to include tax in or exclude from quoted prices

### Invoice Formatting for Canadian Standards

Include these elements on all invoices:

- Your business name, address, and phone number
- Your GST/HST registration number (RT-XXXX-XXXXXXX)
- Client's complete billing address
- Invoice number (sequential)
- Invoice date and due date
- Description of services
- Hours, rates, and amounts (subtotal, tax, total)
- Payment terms

## Project Structure

```
freelance-flow/
├── frontend/           # Vue 3 application
├── internal/           # Go backend (services, models, database)
├── main.go             # Application entry point
└── wails.json          # Wails configuration
```

## Development

```bash
# Run with hot-reload
wails dev

# Pre-commit hooks (recommended)
brew install pre-commit
pre-commit install
```

## License

MIT