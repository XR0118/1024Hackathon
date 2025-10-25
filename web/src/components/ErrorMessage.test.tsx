import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import ErrorMessage from './ErrorMessage'

describe('ErrorMessage', () => {
  it('renders error message correctly', () => {
    render(
      <ErrorMessage 
        message="Test error message" 
        onDismiss={() => {}} 
      />
    )
    
    expect(screen.getByText('Test error message')).toBeInTheDocument()
  })

  it('calls onDismiss when close button is clicked', async () => {
    const onDismiss = vi.fn()
    const user = userEvent.setup()
    
    render(
      <ErrorMessage 
        message="Test error" 
        onDismiss={onDismiss} 
      />
    )
    
    const closeButton = screen.getByRole('button')
    await user.click(closeButton)
    
    expect(onDismiss).toHaveBeenCalledTimes(1)
  })

  it('displays error alert style', () => {
    const { container } = render(
      <ErrorMessage 
        message="Error" 
        onDismiss={() => {}} 
      />
    )
    
    const alert = container.querySelector('.alert-danger')
    expect(alert).toBeInTheDocument()
  })
})
