// src/utils/auth.guard.js
// Simple auth guard for Vue Router

export async function requireAuth(to, from, next) {
  console.log('Auth Guard: Checking route:', to.path)
  const publicPages = ['/', '/login', '/auth/microsoft/callback']
  const authRequired = !publicPages.includes(to.path)

  // Only validate the session if the route requires authentication
  if (authRequired) {
    const response = await fetch('http://localhost:8080/api/user', {
      credentials: 'include',
    })
    if (response.ok) {
      console.log('Auth Guard: Session validated by backend')
      // Allow navigation
      console.log('Auth Guard: Navigation allowed')
      next()
      return
    } else {
      console.warn('Auth Guard: Session invalid, redirecting to home')
      // If route requires auth and no session ID, redirect to home
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
