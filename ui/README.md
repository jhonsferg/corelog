# CoreLog UI

React, TypeScript, and Vite frontend for CoreLog, structured as a hexagonal architecture so business logic, meaning domain and application, stays independent of React, Axios, Ant Design, or any other framework detail.

## Table of contents

- [Layout](#layout)
- [Container / presentational split](#container--presentational-split)
- [Path aliases](#path-aliases)
- [Stack](#stack)
- [Getting started](#getting-started)
- [Available scripts](#available-scripts)
- [Environment variables](#environment-variables)
- [Styling conventions](#styling-conventions)

## Layout

```
src/
├── core/                          framework-free hexagon core: domain and application
│   ├── auth/
│   │   ├── domain/                     types and business rules such as User and Role
│   │   └── application/
│   │       ├── ports/                    outbound port interface, AuthRepositoryPort
│   │       └── auth.service.ts           inbound port implementation, the use cases
│   ├── tickets/                     same layout: domain/ and application/
│   ├── teams/                       same layout: domain/ and application/
│   └── users/                       search and team-assignment for other users, reuses the auth User type
├── adapters/
│   ├── out/                         driven adapters that implement core ports
│   │   ├── http/                         axios-based repositories plus the shared api-client
│   │   └── storage/                      localStorage-backed token storage
│   └── in/                          driving adapters that call the core
│       ├── state/                        zustand stores bridging UI and application services
│       └── ui/                           pages as containers, presentational components, hooks/
│           ├── auth/                       LoginPage, RegisterPage
│           ├── common/                     ProtectedRoute, AdminRoute
│           ├── layout/                     AppLayout
│           ├── profile/                    ProfilePage, ProfileForm, PasswordForm
│           ├── teams/                      TeamsPage, TeamList, TeamFormModal, TeamMembersPanel
│           └── tickets/                    TicketsListPage, NewTicketPage, TicketDetailPage,
│                                            TicketList, TicketForm, TicketBadge, hooks/
├── app/                            composition root: providers, router, App.tsx
├── shared/                         cross-cutting config, env
└── styles/                         global, non-component CSS
```

**Dependency rule**: `adapters/in/ui` → `adapters/in/state` → `core/*/application` → `core/*/domain`, and `core/*/application` → `core/*/ports` ← `adapters/out/*`. Nothing under `core/` imports React, Axios, or Ant Design. It is plain, framework-free TypeScript that could be unit-tested or reused outside a browser.

## Container / presentational split

Each feature under `adapters/in/ui/<feature>` follows the same split:

- **Container pages**, such as `TicketsListPage`, `NewTicketPage`, and `TicketDetailPage`, own data fetching, navigation, and toasts. They delegate the actual fetch and mutate logic to a custom hook in `<feature>/hooks/`, such as `useTickets` or `useCreateTicket`, which in turn talks to the Zustand store.
- **Presentational components**, such as `TicketList`, `TicketForm`, and `TicketBadge`, are pure: they receive data and callbacks as props and render UI only, with no store, router, or API access inside them.

This keeps every component testable in isolation: render with props and assert the output. It also keeps data-fetching logic in exactly one place per feature.

`ProtectedRoute` and `AdminRoute`, under `adapters/in/ui/common`, are route guards rather than pages: `ProtectedRoute` redirects to `/login` once the session check has run and found no authenticated user, and `AdminRoute` redirects to `/tickets` when the authenticated user's role is not `admin`. Both read `useAuthStore` directly and wrap route subtrees in `app/router.tsx` via nested routes, rather than being rendered by any page.

## Path aliases

Imports that cross a top-level `src/` boundary, such as `app`, `core`, `adapters`, or `shared`, use the `@/` alias instead of relative `../../../` chains:

```ts
import { useAuthStore } from "@/adapters/in/state/auth.store";
import type { Ticket } from "@/core/tickets/domain/ticket";
```

Imports that stay within the same directory or one level up, such as siblings within a feature, stay relative. Examples: a component and its own `.module.css`, or an application service and its own `domain/`:

```ts
import styles from "./TicketList.module.css";
import type { User } from "../domain/user";
```

The alias is configured in two places that must stay in sync. `vite.config.ts` sets it via `resolve.alias`, used by the bundler and dev server. `tsconfig.app.json` sets it via `compilerOptions.paths`, used by the TypeScript language service and `tsc -b`.

## Stack

- **React 19, TypeScript, and Vite**
- **Routing**: React Router, configured in `app/router.tsx` as a data router via `createBrowserRouter`
- **State**: Zustand, in `adapters/in/state`
- **HTTP**: Axios, in `adapters/out/http`
- **Forms**: React Hook Form plus Zod
- **UI components**: Ant Design plus `@ant-design/icons`
- **Styling**: CSS Modules, see [Styling conventions](#styling-conventions) below

TanStack Query was intentionally avoided after a past npm supply-chain compromise of its packages. Zustand, Axios, and React Hook Form with Zod cover the same needs. Server state caching is minimal here since the app is small.

## Getting started

```sh
cp .env.example .env
npm install
npm run dev
```

The API base URL defaults to `http://localhost:8080/api/v1`. Override it via `VITE_API_BASE_URL` in `.env` if the backend runs elsewhere. See the repository root [`README.md`](../README.md) for running the backend with either PostgreSQL or SQLite.

## Available scripts

| Command           | Description                                      |
| ----------------- | ------------------------------------------------ |
| `npm run dev`     | Start the Vite dev server with HMR               |
| `npm run build`   | Type-check via `tsc -b` and build for production |
| `npm run lint`    | Run ESLint                                       |
| `npm run preview` | Serve the production build locally               |

## Environment variables

| Variable            | Default                        | Description                       |
| ------------------- | ------------------------------ | --------------------------------- |
| `VITE_API_BASE_URL` | `http://localhost:8080/api/v1` | Base URL the Axios client targets |

## Styling conventions

Every component owns its own `<Component>.module.css`, imported as `styles` and referenced via `className={styles.foo}`. Never use global CSS or inline `style={{}}` props for anything reusable. `src/styles/global.css` is the only global stylesheet and is limited to resets and root-level tokens. It is loaded once, in `main.tsx`.
