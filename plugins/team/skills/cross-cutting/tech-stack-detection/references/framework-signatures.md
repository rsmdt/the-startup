# Reference: Framework Signatures

## Overview

Comprehensive detection patterns for identifying frameworks, libraries, and tools across major ecosystems. Use these signatures to accurately identify project tech stacks.

## Quick Reference

| Ecosystem | Primary Manifest | Lock File | Config Pattern |
|-----------|-----------------|-----------|----------------|
| Node.js | `package.json` | `package-lock.json`, `yarn.lock`, `pnpm-lock.yaml` | `*.config.js/ts/mjs` |
| Python | `pyproject.toml`, `setup.py` | `poetry.lock`, `Pipfile.lock`, `requirements.txt` | `pyproject.toml`, `setup.cfg` |
| Rust | `Cargo.toml` | `Cargo.lock` | `Cargo.toml` |
| Go | `go.mod` | `go.sum` | `go.mod` |
| Ruby | `Gemfile` | `Gemfile.lock` | `config/*.rb` |
| PHP | `composer.json` | `composer.lock` | `*.php` configs |
| .NET | `*.csproj`, `*.sln` | `packages.lock.json` | `appsettings.json` |
| Java | `pom.xml`, `build.gradle` | Various | `application.properties/yml` |

---

## Frontend Frameworks

### React

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"react"` in package.json | HIGH |
| Dependency | `"react-dom"` in package.json | HIGH |
| File extension | `.jsx`, `.tsx` files | MEDIUM |
| Import pattern | `import React from 'react'` | HIGH |
| Config | `babel.config.js` with `@babel/preset-react` | HIGH |

**Common companions**: `react-router-dom`, `redux`, `@tanstack/react-query`, `zustand`

### Vue.js

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"vue"` in package.json | HIGH |
| Config file | `vue.config.js` | HIGH |
| Config file | `vite.config.ts` with `@vitejs/plugin-vue` | HIGH |
| File extension | `.vue` files | HIGH |
| Directory | `src/components/*.vue` | HIGH |

**Version detection**:
- Vue 2: `"vue": "^2.x"`
- Vue 3: `"vue": "^3.x"`, presence of Composition API patterns

### Angular

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Config file | `angular.json` | HIGH |
| Dependency | `"@angular/core"` in package.json | HIGH |
| File pattern | `*.component.ts`, `*.module.ts`, `*.service.ts` | HIGH |
| Directory | `src/app/` structure | MEDIUM |
| CLI | `@angular/cli` in devDependencies | HIGH |

### Svelte

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"svelte"` in package.json | HIGH |
| Config file | `svelte.config.js` | HIGH |
| File extension | `.svelte` files | HIGH |
| Vite plugin | `@sveltejs/vite-plugin-svelte` | HIGH |

### Solid.js

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"solid-js"` in package.json | HIGH |
| Config | `vite.config.ts` with `solid-start` | HIGH |
| File pattern | `.tsx` with `createSignal`, `createEffect` | MEDIUM |

---

## Meta-Frameworks

### Next.js

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"next"` in package.json | HIGH |
| Config file | `next.config.js` or `next.config.mjs` | HIGH |
| Directory | `app/` (App Router) or `pages/` (Pages Router) | HIGH |
| File pattern | `layout.tsx`, `page.tsx` in app/ | HIGH |
| File pattern | `_app.tsx`, `_document.tsx` in pages/ | HIGH |

**Router detection**:
- App Router: `app/` directory with `layout.tsx`
- Pages Router: `pages/` directory with `_app.tsx`

### Nuxt.js

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"nuxt"` in package.json | HIGH |
| Config file | `nuxt.config.ts` or `nuxt.config.js` | HIGH |
| Directory | `pages/`, `components/`, `composables/` | HIGH |
| File pattern | `app.vue` | HIGH |

**Version detection**:
- Nuxt 2: `"nuxt": "^2.x"`
- Nuxt 3: `"nuxt": "^3.x"`, TypeScript-first config

### Remix

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"@remix-run/react"` in package.json | HIGH |
| Config file | `remix.config.js` | HIGH |
| Directory | `app/routes/` | HIGH |
| File pattern | `root.tsx`, route files with `loader`/`action` exports | HIGH |

### SvelteKit

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"@sveltejs/kit"` in package.json | HIGH |
| Config file | `svelte.config.js` with `@sveltejs/kit` | HIGH |
| Directory | `src/routes/` | HIGH |
| File pattern | `+page.svelte`, `+layout.svelte` | HIGH |

### Astro

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"astro"` in package.json | HIGH |
| Config file | `astro.config.mjs` | HIGH |
| Directory | `src/pages/`, `src/components/` | HIGH |
| File extension | `.astro` files | HIGH |

---

## Backend Frameworks (Node.js)

### Express.js

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"express"` in package.json | HIGH |
| Code pattern | `const app = express()` | HIGH |
| Code pattern | `app.get()`, `app.post()`, `app.use()` | HIGH |
| File pattern | `app.js`, `server.js`, `index.js` entry point | MEDIUM |

### NestJS

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"@nestjs/core"` in package.json | HIGH |
| Config file | `nest-cli.json` | HIGH |
| File pattern | `*.controller.ts`, `*.service.ts`, `*.module.ts` | HIGH |
| Directory | `src/` with module structure | MEDIUM |
| Decorator pattern | `@Controller()`, `@Injectable()`, `@Module()` | HIGH |

### Fastify

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"fastify"` in package.json | HIGH |
| Code pattern | `fastify()`, `fastify.register()` | HIGH |

### Hono

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"hono"` in package.json | HIGH |
| Code pattern | `new Hono()`, `app.get()` | HIGH |

### Koa

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"koa"` in package.json | HIGH |
| Code pattern | `new Koa()`, `ctx.body` | HIGH |

---

## Python Frameworks

### Django

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `django` in requirements/pyproject.toml | HIGH |
| File | `manage.py` | HIGH |
| File | `settings.py` or `settings/` directory | HIGH |
| Directory | `<app>/models.py`, `<app>/views.py`, `<app>/urls.py` | HIGH |
| Config | `INSTALLED_APPS`, `MIDDLEWARE` in settings | HIGH |

### FastAPI

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `fastapi` in requirements/pyproject.toml | HIGH |
| Code pattern | `from fastapi import FastAPI` | HIGH |
| Code pattern | `@app.get()`, `@app.post()` decorators | HIGH |
| Companion | `uvicorn` in dependencies | MEDIUM |

### Flask

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `flask` in requirements/pyproject.toml | HIGH |
| Code pattern | `from flask import Flask` | HIGH |
| Code pattern | `app = Flask(__name__)` | HIGH |
| File pattern | `app.py`, `application.py` | MEDIUM |

### Starlette

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `starlette` in requirements/pyproject.toml | HIGH |
| Code pattern | `from starlette.applications import Starlette` | HIGH |

---

## Build Tools and Bundlers

### Vite

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"vite"` in devDependencies | HIGH |
| Config file | `vite.config.ts` or `vite.config.js` | HIGH |
| Script | `"dev": "vite"` in package.json scripts | HIGH |

### Webpack

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"webpack"` in devDependencies | HIGH |
| Config file | `webpack.config.js` | HIGH |
| Directory | `dist/` or `build/` output | MEDIUM |

### esbuild

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"esbuild"` in devDependencies | HIGH |
| Script | Build script using `esbuild` | HIGH |

### Turbopack

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Config | `next.config.js` with `experimental.turbo` | HIGH |
| Script | `next dev --turbo` | HIGH |

### Rollup

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"rollup"` in devDependencies | HIGH |
| Config file | `rollup.config.js` | HIGH |

### Parcel

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"parcel"` in devDependencies | HIGH |
| Script | `"start": "parcel"` | HIGH |

---

## CSS and Styling

### Tailwind CSS

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"tailwindcss"` in devDependencies | HIGH |
| Config file | `tailwind.config.js` or `tailwind.config.ts` | HIGH |
| CSS directive | `@tailwind base/components/utilities` | HIGH |
| Class pattern | `className="flex items-center..."` | MEDIUM |

### Styled Components

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"styled-components"` in dependencies | HIGH |
| Code pattern | ``styled.div`...` `` template literals | HIGH |

### Emotion

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"@emotion/react"` or `"@emotion/styled"` | HIGH |
| Code pattern | `css` prop or `styled` API | HIGH |

### Sass/SCSS

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"sass"` in devDependencies | HIGH |
| File extension | `.scss`, `.sass` files | HIGH |

### CSS Modules

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| File extension | `*.module.css`, `*.module.scss` | HIGH |
| Import pattern | `import styles from './Component.module.css'` | HIGH |

---

## Database and ORM

### Prisma

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"prisma"` or `"@prisma/client"` | HIGH |
| File | `prisma/schema.prisma` | HIGH |
| Directory | `prisma/migrations/` | HIGH |

### Drizzle

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"drizzle-orm"` in dependencies | HIGH |
| Config file | `drizzle.config.ts` | HIGH |
| Directory | `drizzle/` migrations folder | HIGH |

### TypeORM

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"typeorm"` in dependencies | HIGH |
| File | `ormconfig.json` or `data-source.ts` | HIGH |
| Decorator pattern | `@Entity()`, `@Column()` | HIGH |

### Sequelize

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"sequelize"` in dependencies | HIGH |
| Directory | `migrations/`, `seeders/`, `models/` | HIGH |

### Mongoose

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"mongoose"` in dependencies | HIGH |
| Code pattern | `mongoose.Schema`, `mongoose.model` | HIGH |

### SQLAlchemy

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `sqlalchemy` in requirements | HIGH |
| Code pattern | `from sqlalchemy import ...` | HIGH |

---

## Testing Frameworks

### Jest

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"jest"` in devDependencies | HIGH |
| Config file | `jest.config.js` or `jest.config.ts` | HIGH |
| File pattern | `*.test.js`, `*.spec.js` | HIGH |
| Script | `"test": "jest"` | HIGH |

### Vitest

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"vitest"` in devDependencies | HIGH |
| Config file | `vitest.config.ts` | HIGH |
| Script | `"test": "vitest"` | HIGH |

### Playwright

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"@playwright/test"` in devDependencies | HIGH |
| Config file | `playwright.config.ts` | HIGH |
| Directory | `tests/` or `e2e/` | MEDIUM |

### Cypress

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"cypress"` in devDependencies | HIGH |
| Config file | `cypress.config.ts` or `cypress.json` | HIGH |
| Directory | `cypress/` | HIGH |

### pytest

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `pytest` in requirements | HIGH |
| Config | `pytest.ini` or `pyproject.toml` [tool.pytest] | HIGH |
| File pattern | `test_*.py`, `*_test.py` | HIGH |

---

## API and Data Fetching

### tRPC

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"@trpc/server"` or `"@trpc/client"` | HIGH |
| Directory | `server/trpc/` or similar | HIGH |
| Code pattern | `router`, `procedure` | HIGH |

### GraphQL

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"graphql"` in dependencies | HIGH |
| Dependency | `"@apollo/client"` or `"urql"` | HIGH |
| File extension | `.graphql`, `.gql` files | HIGH |
| Directory | `graphql/` or `schema/` | MEDIUM |

### React Query / TanStack Query

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"@tanstack/react-query"` | HIGH |
| Code pattern | `useQuery`, `useMutation` hooks | HIGH |

### SWR

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"swr"` in dependencies | HIGH |
| Code pattern | `useSWR` hook | HIGH |

---

## Monorepo Tools

### Turborepo

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"turbo"` in devDependencies | HIGH |
| Config file | `turbo.json` | HIGH |

### Nx

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Config file | `nx.json` | HIGH |
| Directory | `apps/`, `libs/` structure | HIGH |

### Lerna

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Config file | `lerna.json` | HIGH |
| Dependency | `"lerna"` in devDependencies | HIGH |

### pnpm Workspaces

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Config file | `pnpm-workspace.yaml` | HIGH |
| Lock file | `pnpm-lock.yaml` | HIGH |

---

## Mobile Frameworks

### React Native

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"react-native"` in dependencies | HIGH |
| Config file | `metro.config.js` | HIGH |
| Directory | `ios/`, `android/` | HIGH |
| File | `app.json` with `expo` or RN config | HIGH |

### Expo

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"expo"` in dependencies | HIGH |
| Config file | `app.json` with `expo` key | HIGH |
| Config file | `app.config.js` or `app.config.ts` | HIGH |

### Flutter

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| File | `pubspec.yaml` | HIGH |
| Directory | `lib/`, `android/`, `ios/` | HIGH |
| File extension | `.dart` files | HIGH |

---

## Deployment Platforms

### Vercel

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Config file | `vercel.json` | HIGH |
| Directory | `.vercel/` | HIGH |

### Netlify

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Config file | `netlify.toml` | HIGH |
| Directory | `.netlify/` | HIGH |

### Docker

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| File | `Dockerfile` | HIGH |
| File | `docker-compose.yml` or `compose.yaml` | HIGH |
| File | `.dockerignore` | MEDIUM |

### Kubernetes

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Directory | `k8s/`, `kubernetes/`, `manifests/` | HIGH |
| File pattern | `*.yaml` with `apiVersion`, `kind` | HIGH |
| File | `Helm` charts in `charts/` | HIGH |

---

## State Management

### Redux

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"redux"` or `"@reduxjs/toolkit"` | HIGH |
| Directory | `store/`, `slices/` | MEDIUM |
| Code pattern | `createSlice`, `configureStore` | HIGH |

### Zustand

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"zustand"` in dependencies | HIGH |
| Code pattern | `create()` store pattern | HIGH |

### Jotai

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"jotai"` in dependencies | HIGH |
| Code pattern | `atom()`, `useAtom()` | HIGH |

### Pinia

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"pinia"` in dependencies | HIGH |
| Directory | `stores/` | MEDIUM |
| Code pattern | `defineStore()` | HIGH |

---

## Authentication

### NextAuth.js / Auth.js

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"next-auth"` or `"@auth/core"` | HIGH |
| File | `[...nextauth].ts` or `auth.ts` | HIGH |
| Directory | `app/api/auth/` | HIGH |

### Clerk

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"@clerk/nextjs"` or similar | HIGH |
| Code pattern | `<ClerkProvider>` | HIGH |

### Supabase Auth

| Indicator | Pattern | Confidence |
|-----------|---------|------------|
| Dependency | `"@supabase/supabase-js"` | HIGH |
| Code pattern | `supabase.auth` | HIGH |

---

## External Resources

- [State of JS Survey](https://stateofjs.com) - Annual ecosystem trends
- [npm trends](https://npmtrends.com) - Package popularity comparison
- [BuiltWith](https://builtwith.com) - Technology detection reference
- [Wappalyzer](https://www.wappalyzer.com) - Web technology profiler patterns
