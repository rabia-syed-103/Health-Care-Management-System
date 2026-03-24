import { createContext, useContext, useState, useEffect } from 'react'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null)
  const [token, setToken] = useState(null)
  const [role, setRole] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const savedToken = localStorage.getItem('token')
    const savedRole  = localStorage.getItem('role')
    const savedUser  = localStorage.getItem('user')
    if (savedToken && savedRole) {
      setToken(savedToken)
      setRole(savedRole)
      setUser(savedUser ? JSON.parse(savedUser) : null)
    }
    setLoading(false)
  }, [])

  const login = (token, role, user) => {
    localStorage.setItem('token', token)
    localStorage.setItem('role', role)
    localStorage.setItem('user', JSON.stringify(user))
    setToken(token)
    setRole(role)
    setUser(user)
  }

  const logout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('role')
    localStorage.removeItem('user')
    setToken(null)
    setRole(null)
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{ user, token, role, login, logout, loading }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  return useContext(AuthContext)
}
