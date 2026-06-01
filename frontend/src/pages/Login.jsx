import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../AuthContext'
import './Login.css'

const URL = 'http://localhost:8000'

export default function Login() {
  const { login } = useAuth()
  const navigate = useNavigate()
  const [form, setForm] = useState({ username: '', password: '' })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleChange = e => {
    setForm(f => ({ ...f, [e.target.name]: e.target.value }))
    setError('')
  }

  const handleSubmit = async e => {
    e.preventDefault()
    if (!form.username || !form.password) {
      setError('Ingresa usuario y contraseña')
      return
    }
    setLoading(true)
    try {
      const res = await fetch(`${URL}/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(form)
      })
      if (!res.ok) { setError('Credenciales inválidas'); return }
      const data = await res.json()
      login({ username: data.username, rol: data.rol }, data.token)
      if (data.rol === 'auditor') navigate('/reportes')
      else navigate('/')
    } catch {
      setError('Error de conexión')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="login-page">
      <div className="login-left">
        <img src="../ramo.jpg" alt="Logo" className="login-logo" />
      </div>
      <div className="login-right">
        <div className="login-card">
          <div className="login-brand">
            <h1>Bloom</h1>
            <p>Ingresa con tu cuenta para continuar</p>
          </div>
          <form onSubmit={handleSubmit}>
            <div className="field">
              <label>Usuario</label>
              <input
                name="username"
                value={form.username}
                onChange={handleChange}
                placeholder="Ej: vendedor1"
                autoComplete="username"
              />
            </div>
            <div className="field">
              <label>Contraseña</label>
              <input
                name="password"
                type="password"
                value={form.password}
                onChange={handleChange}
                placeholder="••••••••"
                autoComplete="current-password"
              />
            </div>
            {error && <div className="login-error">{error}</div>}
            <button
              type="submit"
              className="btn btn-primary login-btn"
              disabled={loading}
            >
              {loading ? 'Ingresando…' : 'Ingresar'}
            </button>
          </form>
        </div>
      </div>
    </div>
  )
}