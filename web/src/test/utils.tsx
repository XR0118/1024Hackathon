import { ReactElement } from 'react'
import { render, RenderOptions } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'

interface CustomRenderOptions extends Omit<RenderOptions, 'wrapper'> {
  initialEntries?: string[]
}

export function renderWithRouter(
  ui: ReactElement,
  { initialEntries = ['/'], ...options }: CustomRenderOptions = {}
) {
  return render(ui, {
    wrapper: ({ children }) => (
      <BrowserRouter>{children}</BrowserRouter>
    ),
    ...options,
  })
}

export * from '@testing-library/react'
export { userEvent } from '@testing-library/user-event'
