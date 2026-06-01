import { NavLink, useNavigate } from 'react-router-dom'
import { useAuth } from '../AuthContext'
import './Navbar.css'

export default function Navbar() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = async () => {
    await fetch('http://localhost:8000/logout', { method: 'POST' })
    logout()
    navigate('/login')
  }

  return (
    <nav className="navbar">
      <div className="navbar-brand">
        <span className="brand-name">Bloom</span>
      </div>
      <div className="navbar-links">
        {['superadmin', 'gerente', 'vendedor', 'comprador'].includes(user?.rol) && (
          <NavLink to="/" end className={({ isActive }) => isActive ? 'active' : ''}>
            Inventario
          </NavLink>
        )}
        {['superadmin', 'gerente', 'auditor'].includes(user?.rol) && (
          <NavLink to="/reportes" className={({ isActive }) => isActive ? 'active' : ''}>
            Reportes
          </NavLink>
        )}
      </div>
      <div className="navbar-user">
        <span className="navbar-username">{user?.username}</span>
        <span className="navbar-rol">{user?.rol}</span>
        <button className="btn btn-secondary btn-sm" onClick={handleLogout}>
          Cerrar sesión
        </button>
      </div>
    </nav>
  )
}