import { useState, useEffect } from 'react'
import { api } from '../api'

export default function CompraModal({ carrito, onClose, onComprado }) {
  const [clientes, setClientes] = useState([])
  const [empleados, setEmpleados] = useState([])
  const [form, setForm] = useState({ id_cliente: '', id_empleado: '', fecha: new Date().toISOString().split('T')[0] })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    api.getClientes().then(setClientes).catch(() => {})
    api.getEmpleados().then(setEmpleados).catch(() => {})
  }, [])

  const total = carrito.reduce((sum, item) => sum + item.precio * item.cantidad, 0)

  const handleChange = e => {
    setForm(f => ({ ...f, [e.target.name]: e.target.value }))
    setError('')
  }

  const handleComprar = async () => {
    if (!form.id_cliente || !form.id_empleado || !form.fecha) {
      setError('Todos los campos son requeridos')
      return
    }
    setLoading(true)
    setError('')
    try {
      // 1. Crear el ramo con los productos del carrito
      const ramo = await api.crearRamo({
        productos: carrito.map(i => ({ id_producto: i.id_producto, cantidad: i.cantidad }))
      })
      // 2. Crear la venta con ese ramo
      await api.crearVenta({
        id_cliente: parseInt(form.id_cliente),
        id_empleado: parseInt(form.id_empleado),
        fecha: form.fecha,
        ramos: [ramo.id_ramo],
      })
      onComprado()
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="modal-overlay" onClick={e => e.target === e.currentTarget && onClose()}>
      <div className="modal">
        <h2>Confirmar compra</h2>

        {/* Resumen del ramo */}
        <div style={{ background: 'var(--cream)', borderRadius: 'var(--radius-sm)', padding: '14px 16px', marginBottom: 20 }}>
          <p style={{ fontFamily: 'var(--font-display)', fontSize: 13, color: 'var(--bark-mid)', marginBottom: 8, textTransform: 'uppercase', letterSpacing: '0.6px' }}>Resumen del ramo</p>
          {carrito.map(item => (
            <div key={item.id_producto} style={{ display: 'flex', justifyContent: 'space-between', fontSize: 13, padding: '3px 0', borderBottom: '1px solid var(--cream-dark)' }}>
              <span>{item.nombre} × {item.cantidad}</span>
              <span style={{ color: 'var(--moss)' }}>Q{(item.precio * item.cantidad).toFixed(2)}</span>
            </div>
          ))}
          <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: 10, fontWeight: 500, fontSize: 15 }}>
            <span>Total</span>
            <span style={{ color: 'var(--moss-dark)', fontFamily: 'var(--font-display)', fontSize: 18 }}>Q{total.toFixed(2)}</span>
          </div>
        </div>

        <div className="field">
          <label>Cliente</label>
          <select name="id_cliente" value={form.id_cliente} onChange={handleChange}>
            <option value="">Seleccionar cliente</option>
            {clientes.map(c => (
              <option key={c.id_cliente} value={c.id_cliente}>{c.nombre}</option>
            ))}
          </select>
        </div>

        <div className="field">
          <label>Empleado que atiende</label>
          <select name="id_empleado" value={form.id_empleado} onChange={handleChange}>
            <option value="">Seleccionar empleado</option>
            {empleados.map(e => (
              <option key={e.id_empleado} value={e.id_empleado}>{e.nombre} — {e.rol}</option>
            ))}
          </select>
        </div>

        <div className="field">
          <label>Fecha</label>
          <input type="date" name="fecha" value={form.fecha} onChange={handleChange} />
        </div>

        {error && <p style={{ color: 'var(--error)', fontSize: 13, marginBottom: 8 }}>{error}</p>}

        <div className="modal-actions">
          <button className="btn btn-secondary" onClick={onClose}>Cancelar</button>
          <button className="btn btn-gold" onClick={handleComprar} disabled={loading}>
            {loading ? 'Procesando…' : `Comprar — Q${total.toFixed(2)}`}
          </button>
        </div>
      </div>
    </div>
  )
}
