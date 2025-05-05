// src/utils/auth.guard.js
// Simple auth guard for Vue Router

import axios from 'axios'

export async function requireAuth(to, from, next) {
  console.log('Auth Guard: Checking route:', to.path)
  const publicPages = ['/', '/login', '/auth/microsoft/callback']
  const authRequired = !publicPages.includes(to.path)

  // Only validate the session if the route requires authentication
  if (authRequired) {
    try {
      const response = await axios.get('http://localhost:8080/api/user', {
        withCredentials: true,
      })
      if (response.status === 200) {
        console.log('Auth Guard: Session validated by backend')
        // Allow navigation
        console.log('Auth Guard: Navigation allowed')
        next()
        return
      }
    } catch (error) {
      console.warn('Auth Guard: Session invalid, redirecting to home')
      if (authRequired) {
        console.warn('Auth Guard: No session ID, redirecting to home')
        next({ name: 'home', replace: true })
        return
      }
    }
  } else {
    // Allow navigation for public routes
    console.log('Auth Guard: Navigation allowed for public route')
    next()
  }

  // If logged in and trying to access home or login, redirect to postfeed
  if (sessionID && (to.name === 'home' || to.path === '/login')) {
    console.info('Auth Guard: Session ID found, redirecting to postfeed')
    next({ name: 'postfeed', replace: true })
    return
  }
}
