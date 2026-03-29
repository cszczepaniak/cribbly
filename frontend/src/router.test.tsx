import { fireEvent, render, screen } from '@testing-library/react'
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
    expect(screen.getByRole('heading', { name: /welcome to cribbly/i })).toBeInTheDocument()
  })

  it('shows room code entry when the home page switch is toggled', () => {
    const router = createMemoryRouter(routeObjects, {
      initialEntries: ['/?react=true'],
      basename: '/',
    })
    render(<RouterProvider router={router} />)
    expect(screen.getByRole('heading', { name: /welcome to cribbly/i })).toBeInTheDocument()
    fireEvent.click(screen.getByRole('switch', { name: /show room code entry instead of the home page/i }))
    expect(screen.getByRole('heading', { name: /enter room code/i })).toBeInTheDocument()
  })
})
