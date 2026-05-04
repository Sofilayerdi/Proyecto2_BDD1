import { NavLink } from 'react-router-dom'
import './Navbar.css'

export default function Navbar() {
  return (
    <nav className="navbar">
      <div className="navbar-brand">
        <span className="brand-icon">❧</span>
        <span className="brand-name">Pétalos <em>&amp;</em> Co.</span>
      </div>
      <div className="navbar-links">
        <NavLink to="/" end className={({ isActive }) => isActive ? 'active' : ''}>
          Inventario
        </NavLink>
        <NavLink to="/reportes" className={({ isActive }) => isActive ? 'active' : ''}>
          Reportes
        </NavLink>
      </div>
    </nav>
  )
}
