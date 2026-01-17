# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Hot-reloading Development (Recommended)
```bash
air
```
This runs the full build pipeline (templ generate, tailwindcss build, go build) and auto-restarts on file changes.

### Manual Build Steps
```bash
# Generate Go code from .templ files
templ generate

# Build CSS from Tailwind input
tailwindcss -i ./input.css -o ./static/output.css

# Build the binary
go build -o ./tmp/main ./cmd/api/main.go

# Run the application
./tmp/main
```

### Installing Dependencies
```bash
# Go dependencies
go mod download

# Node.js dependencies (TailwindCSS only)
pnpm install
```

## Architecture

This is a Go web application using **Fiber** (Express-inspired web framework) with a clean architecture pattern:

### Layer Structure
- **cmd/api/main.go** - Application entry point: Fiber server setup, middleware, MongoDB connection
- **internal/handler/** - HTTP request handlers that call render functions
- **internal/render/** - Template rendering bridge between Fiber and Templ components
- **internal/db/** - Database connection management (MongoDB)
- **internal/logger/** - Structured logging wrapper around Zap (environment-aware: dev/prod)
- **internal/view/** - Templ HTML templates (layouts/ and pages/)

### Request Flow
1. HTTP request hits Fiber route in `main.go`
2. Handler function in `internal/handler/` processes request
3. Handler calls `render.HTML()` with a Templ component
4. Templ component is rendered and returned as HTML

### Frontend Stack
- **Templ** - Type-safe HTML templating (generates `_templ.go` files from `.templ` files)
- **TailwindCSS** - Utility-first CSS (input.css → static/output.css)
- **HTMX** - Client-side interactivity (served from static/)
- **templui** - UI component library for reusable components

### Environment Configuration
- `.env` file contains `MONGODB_URI` and `APP_ENV` (dev/prod)
- Logger uses `APP_ENV` to switch between Development (debug) and Production (info) modes
- Default port: 3000

## Key Patterns

### Template Rendering
All handlers use the pattern:
```go
return render.HTML(c, layouts.Base("pageTitle", view.ComponentName()))
```
The `layouts.Base()` wrapper includes HTMX script and Tailwind CSS.

### Logging
Use the structured logging helpers:
```go
logger.Info("ServiceName", "message", zap.String("key", "value"))
logger.Error("ServiceName", "error message")
```

### Adding New Pages
1. Create `.templ` file in `internal/view/pages/`
2. Run `templ generate` to generate Go code
3. Add handler in `internal/handler/`
4. Register route in `cmd/api/main.go`

## UI Design System

### Color Palette
```css
/* Primary Brand (Indigo/Blurple) */
--color-primary: #5D5CFF;       /* CTA buttons, headlines */

/* Secondary Colors */
--color-text-dark: #4B5563;     /* Subheadlines, descriptions */
--color-text-nav: #1F2937;      /* Navigation links */
--color-text-body: #374151;     /* Tag text, body text */

/* Backgrounds */
--color-bg-white: #FFFFFF;      /* Main background */
--color-bg-tag: #F3F4F6;        /* Tag/pill backgrounds */
```

### Typography
- **Font Family**: Modern sans-serif (Inter, Roboto, or Open Sans)
- **Headings**: Bold, large for H1
- **Body**: Regular weight, optimized for readability

### Component Patterns

**Buttons (Pill-shaped)**
```html
<button class="rounded-full bg-[#5D5CFF] text-white px-6 py-2">
  Button Text
</button>
```

**Tags/Pills**
```html
<span class="rounded-full bg-[#F3F4F6] text-[#374151] px-3 py-1 text-sm">
  Role Label
</span>
```

**Navigation**
```html
<nav class="flex justify-between items-center">
  <div>Logo</div>
  <div class="flex gap-6 text-[#1F2937]">
    <a>Home</a>
    <a>Solutions</a>
    <!-- etc -->
  </div>
</nav>
```

**Hero Section**
- Centered alignment
- Narrower column for subtext (max-width ~600px)
- Large bold H1 with primary color
- Floating elements positioned absolutely around headline

**Images**
- Rounded corners on all photos
- Consistent border radius
