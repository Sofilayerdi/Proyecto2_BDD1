import { Navigate } from 'react-router-dom'
import { useAuth } from '../AuthContext'

export default function PrivateRoute({ children, roles }) {
  const { user } = useAuth()

  if (!user) return <Navigate to="/login" replace />

  if (roles && !roles.includes(user.rol)) {
    return <Navigate to="/no-autorizado" replace />
  }

  return children
}