import { render, screen } from '@testing-library/react'
import { createMemoryRouter, RouterProvider } from 'react-router-dom'
import { describe, expect, it } from 'vitest'
import { routeObjects } from './router'

describe('router', () => {
  it('renders home', () => {
    const router = createMemoryRouter(routeObjects, {
      initialEntries: ['/?react=true'],
      basename: '/',
    })
    render(<RouterProvider router={router} />)
    expect(screen.getByRole('heading', { name: /cribbly/i })).toBeInTheDocument()
  })

  it('renders room code page', () => {
    const router = createMemoryRouter(routeObjects, {
      initialEntries: ['/room?react=true'],
      basename: '/',
    })
    render(<RouterProvider router={router} />)
    expect(screen.getByText(/enter room code/i)).toBeInTheDocument()
  })
})
