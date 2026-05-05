import { Routes, Route } from 'react-router-dom'
import Navbar from './components/Navbar'
import Productos from './pages/Productos'
import Reportes from './pages/Reportes'

export default function App() {
  return (
    <>
      <Navbar />
      <Routes>
        <Route path="/" element={<Productos />} />
        <Route path="/reportes" element={<Reportes />} />
      </Routes>
    </>
  )
}
