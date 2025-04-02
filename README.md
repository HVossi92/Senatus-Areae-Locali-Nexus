# Senatus - LAN Party Organizer

Senatus is a web application designed to streamline the organization of LAN parties.
It helps groups coordinate game selections, scheduling, and food preferences through a democratic voting system,
significantly reducing the overhead typically involved in organizing gaming events.

## Features

- **Event Creation & Management**: Create and manage LAN party events with specific time slots
- **Democratic Game Selection**: Vote on which games to play during different time slots
- **Food Coordination**: Organize food preferences and orders for the event
- **Responsive Design**: Works seamlessly on both desktop and mobile devices

## Tech Stack

### Frontend

- HTML with HTMX for dynamic interactions
- Bootstrap for styling

### Backend

- Go 1.22+ with standard library's net/http package
- [a-h/templ](https://github.com/a-h/templ) for type-safe HTML templating
- SQLite for data persistence
- [sqlc](https://sqlc.dev/) for type-safe SQL queries

## Development Setup

### Prerequisites

- Go 1.22 or newer
- [air](https://github.com/cosmtrek/air) for live reloading
- [sqlc](https://sqlc.dev/) for SQL code generation

### Local Development

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/senatus.git
   cd senatus
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

4. Start the development server:
   ```bash
   air
   ```
   This will:
   - Watch for file changes
   - Automatically regenerate templ templates
   - Rebuild and restart the server

### Database Changes

When modifying SQL schemas or queries:

```bash
make sqlc
```

## Deployment

### Building for Production

```bash
make build
```

This creates a compressed archive containing the Linux binary and necessary assets.

### Server Deployment

1. Transfer the build to your server:

   ```bash
   scp senatus.zip user@your-server:/path/to/deployment
   ```

2. On the server:
   ```bash
   unzip senatus.zip
   cd senatus
   nohup ./main &
   ```

Logs will be available in `nohup.out`

## TODOs

- [ ] Handle date changes and general date management
- [ ] Implement date-based pagination
- [ ] Hide past events when following events are in progress
- [ ] Add confirmation dialog for delete operations
- [ ] Live updates using HTMX for a smooth, interactive experiencedd confirmation dialog for delete operations