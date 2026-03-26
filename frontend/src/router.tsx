import { createBrowserRouter } from 'react-router-dom'
import { RootLayout } from './layouts/RootLayout'
import { HomePage } from './pages/HomePage'
import { PlaceholderPage } from './pages/PlaceholderPage'

/** Same as in production behind Go: real paths (`/`, `/admin/games`, …). */
export const routerBasename = '/'

export const routeObjects = [
  {
    path: '/',
    element: <RootLayout />,
    children: [
      { index: true, element: <HomePage /> },
      { path: 'admin/games', element: <PlaceholderPage title="Admin — games" /> },
      { path: '*', element: <PlaceholderPage /> },
    ],
  },
]

export const router = createBrowserRouter(routeObjects, {
  basename: routerBasename,
})
