# Senatus Areae Locali Nexus

A Roman Senate themed todo application built with Go, a-h/templ, HTMX, and Bootstrap 5.

## Description

*"Senatus Areae Locali Nexus"* (Senate Local Area Network) is a modern task management system styled after the ancient Roman Senate. Manage your tasks with the gravitas of a Roman senator, complete with Latin-inspired terminology and ancient Roman aesthetics.

## Features

- Create, view, edit, and delete tasks
- Assign priority levels (Maxima, Alta, Media, Humilis, Minima)
- Track task status (Proposed, Approved, In Progress, Completed, Vetoed)
- Filter tasks by status or priority
- Single-page application-like experience with HTMX
- Roman-inspired UI with scrolls, marble, and ancient iconography

## Technologies

- **Backend**: Go (stdlib) with the new ServeMux introduced in Go 1.22
- **Templating**: a-h/templ for type-safe HTML templates
- **Frontend Enhancement**: HTMX for dynamic interactions
- **Styling**: Bootstrap 5 with custom Roman-themed CSS

## Prerequisites

- Go 1.22 or newer
- a-h/templ CLI tool

## Setup

1. Clone the repository:
   ```
   git clone <repository-url>
   cd Senatus-Areae-Locali-Nexus
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Install templ CLI (if not already installed):
   ```
   go install github.com/a-h/templ/cmd/templ@latest
   ```

4. Generate templ components:
   ```
   templ generate
   ```

5. Run the application:
   ```
   go run cmd/server/main.go
   ```

6. Open your browser and navigate to `http://localhost:8080`

## Project Structure

- `/cmd/server`: Server entry point and HTTP handlers
- `/components`: Templ components for UI rendering
- `/models`: Data models and repository implementation
- `/static`: Static assets (CSS, images)

## Usage

- **View Tasks**: Browse the Tabularium (home page) to see all your tasks
- **Create Task**: Click "Nova Propositio" to propose a new task
- **Edit Task**: Select the stylus icon to amend a proposal
- **Delete Task**: Use the fire icon to remove a task from the agenda
- **Filter Tasks**: Use the dropdown menus to filter by Status or Priority
- **Change Status**: On the task detail page, use the action buttons to change the status

## License

This project is licensed under the MIT License - see the LICENSE file for details. 