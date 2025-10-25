# Frontend Testing

## Overview

This directory contains testing setup, utilities, and mock data for frontend testing with Vitest and React Testing Library.

## Setup

Testing dependencies have been installed:
- `vitest` - Test runner (Vite-native)
- `@vitest/ui` - UI for test visualization
- `@testing-library/react` - React component testing utilities
- `@testing-library/jest-dom` - DOM matchers
- `@testing-library/user-event` - User interaction simulation
- `jsdom` - DOM environment for tests

## Directory Structure

```
test/
├── mocks/
│   ├── api.ts        # Mock API data for all entities
│   └── server.ts     # Mock API server helpers
├── utils.tsx         # Testing utilities and custom renderers
├── setup.ts          # Test environment setup
└── README.md         # This file
```

## Running Tests

```bash
npm test              # Run tests in watch mode
npm test -- --run     # Run tests once
npm run test:ui       # Run tests with UI
npm run test:coverage # Run tests with coverage report
```

## Mock API Data

Mock data is available for all API entities in `mocks/api.ts`:

- `mockVersions` - Version data with git info and applications
- `mockApplications` - Application data with versions and nodes
- `mockEnvironments` - Environment data (k8s/physical)
- `mockDeployments` - Deployment data with various statuses
- `mockDeploymentDetail` - Detailed deployment with steps and logs
- `mockDashboardStats` - Dashboard statistics
- `mockDeploymentTrends` - Deployment trend data

## Using Mock Data in Tests

```typescript
import { mockApplications, mockDeployments } from '@/test/mocks/api'

it('displays applications', () => {
  // Use mock data in your tests
  expect(mockApplications[0].name).toBe('api-service')
})
```

## Using Testing Utilities

```typescript
import { renderWithRouter, screen, userEvent } from '@/test/utils'

it('navigates on click', async () => {
  const user = userEvent.setup()
  renderWithRouter(<MyComponent />)
  
  const button = screen.getByRole('button')
  await user.click(button)
  
  expect(screen.getByText('Success')).toBeInTheDocument()
})
```

## Writing New Tests

1. Create test file next to the component: `ComponentName.test.tsx`
2. Import testing utilities from `@/test/utils`
3. Import mock data from `@/test/mocks/api` as needed
4. Write tests using `describe`, `it`, and `expect`

Example:
```typescript
import { describe, it, expect } from 'vitest'
import { renderWithRouter, screen } from '@/test/utils'
import { mockApplications } from '@/test/mocks/api'
import MyComponent from './MyComponent'

describe('MyComponent', () => {
  it('renders application name', () => {
    renderWithRouter(<MyComponent app={mockApplications[0]} />)
    expect(screen.getByText('api-service')).toBeInTheDocument()
  })
})
```

## Best Practices

1. **Use descriptive test names** - Clearly describe what is being tested
2. **Test user behavior** - Test what users see and do, not implementation details
3. **Use mock data** - Leverage existing mock data instead of creating new data in tests
4. **Keep tests isolated** - Each test should be independent
5. **Clean up after tests** - Use `afterEach(cleanup)` (already configured in setup.ts)
6. **Test accessibility** - Use `getByRole` and other accessibility queries when possible

## CI/CD Integration

Tests should be run as part of the CI/CD pipeline:

```bash
npm test -- --run
```

This will run all tests once and exit with appropriate status code for CI systems.
