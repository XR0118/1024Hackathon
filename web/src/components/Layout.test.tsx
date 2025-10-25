import { describe, it, expect } from 'vitest'
import { screen } from '@testing-library/react'
import { renderWithRouter } from '@/test/utils'
import Layout from './Layout'

describe('Layout', () => {
  it('renders application title', () => {
    renderWithRouter(<Layout />)
    expect(screen.getByText('部署平台')).toBeInTheDocument()
  })

  it('renders all navigation menu items', () => {
    renderWithRouter(<Layout />)
    
    expect(screen.getByText('仪表板')).toBeInTheDocument()
    expect(screen.getByText('版本管理')).toBeInTheDocument()
    expect(screen.getByText('应用管理')).toBeInTheDocument()
    expect(screen.getByText('环境管理')).toBeInTheDocument()
    expect(screen.getByText('部署管理')).toBeInTheDocument()
  })

  it('renders correct number of nav items', () => {
    renderWithRouter(<Layout />)
    const navItems = screen.getAllByRole('link')
    expect(navItems).toHaveLength(5)
  })

  it('has correct navigation links', () => {
    renderWithRouter(<Layout />)
    
    const dashboardLink = screen.getByRole('link', { name: /仪表板/ })
    expect(dashboardLink).toHaveAttribute('href', '/')
    
    const versionsLink = screen.getByRole('link', { name: /版本管理/ })
    expect(versionsLink).toHaveAttribute('href', '/versions')
    
    const appsLink = screen.getByRole('link', { name: /应用管理/ })
    expect(appsLink).toHaveAttribute('href', '/applications')
    
    const envsLink = screen.getByRole('link', { name: /环境管理/ })
    expect(envsLink).toHaveAttribute('href', '/environments')
    
    const deploymentsLink = screen.getByRole('link', { name: /部署管理/ })
    expect(deploymentsLink).toHaveAttribute('href', '/deployments')
  })
})
