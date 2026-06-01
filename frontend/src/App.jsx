import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuth } from './AuthContext'
import Navbar from './components/Navbar'
import PrivateRoute from './components/PrivateRoute'
import Login from './pages/Login'
import Productos from './pages/Productos'
import Reportes from './pages/Reportes'
import NoAutorizado from './pages/NoAutorizado'

export default function App() {
  const { user } = useAuth()

  return (
    <>
      {user && <Navbar />}
      <Routes>
        <Route path="/login" element={
          user ? <Navigate to="/" replace /> : <Login />
        } />
        <Route path="/" element={
          <PrivateRoute roles={['superadmin', 'gerente', 'vendedor', 'comprador']}>
            <Productos />
          </PrivateRoute>
        } />
        <Route path="/reportes" element={
          <PrivateRoute roles={['superadmin', 'gerente', 'auditor']}>
            <Reportes />
          </PrivateRoute>
        } />
        <Route path="/no-autorizado" element={<NoAutorizado />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </>
  )
}