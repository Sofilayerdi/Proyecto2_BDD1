import { useState, useEffect } from 'react'
import './CompraModal.css'

const URL = 'http://localhost:8000'

export default function CompraModal({ carrito, onClose, onComprado }) {
  const [clientes, setClientes] = useState([])
  const [empleados, setEmpleados] = useState([])
  const [form, setForm] = useState({
    id_cliente: '',
    id_empleado: '',
    fecha: new Date().toISOString().split('T')[0]
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    fetch(`${URL}/clientes`)
      .then(res => res.json())
      .then(setClientes)
      .catch(() => {})

    fetch(`${URL}/empleados`)
      .then(res => res.json())
      .then(setEmpleados)
      .catch(() => {})
  }, [])

  const total = carrito.reduce((sum, item) => sum + item.precio * item.cantidad, 0)

  const handleChange = e => {
    setForm(f => ({ ...f, [e.target.name]: e.target.value }))
    setError('')
  }

  const handleComprar = () => {
    if (!form.id_cliente || !form.id_empleado || !form.fecha) {
      setError('Todos los campos son requeridos')
      return
    }
    setLoading(true)
    setError('')

    fetch(`${URL}/ramos`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        productos: carrito.map(i => ({ id_producto: i.id_producto, cantidad: i.cantidad }))
      })
    })
      .then(res => {
        if (!res.ok) return res.text().then(msg => { throw new Error(msg) })
        return res.json()
      })
      .then(ramo => fetch(`${URL}/ventas`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          id_cliente:  parseInt(form.id_cliente),
          id_empleado: parseInt(form.id_empleado),
          fecha:       form.fecha,
          ramos:       [ramo.id_ramo],
        })
      }))
      .then(res => {
        if (!res.ok) return res.text().then(msg => { throw new Error(msg) })
        onComprado()
      })
      .catch(err => setError(err.message))
      .finally(() => setLoading(false))
  }

  return (
    <div className="modal-overlay" onClick={e => e.target === e.currentTarget && onClose()}>
      <div className="modal">
        <h2>Confirmar compra</h2>

        <div className="compra-resumen">
          <p className="compra-resumen-title">🌿 Resumen del ramo</p>
          {carrito.map(item => (
            <div key={item.id_producto} className="compra-resumen-item">
              <span>{item.nombre} × {item.cantidad}</span>
              <span>Q{(item.precio * item.cantidad).toFixed(2)}</span>
            </div>
          ))}
          <div className="compra-total">
            <span className="compra-total-label">Total</span>
            <span className="compra-total-amount">Q{total.toFixed(2)}</span>
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

        {error && <div className="compra-error">{error}</div>}

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