import { useState, useEffect } from 'react'
import { apiFetch } from '../utils/api'
import './ProductoModal.css'

const CATEGORIAS = ['flor', 'follaje', 'liston', 'papel']

export default function ProductoModal({ producto, onClose, onSaved }) {
  const esEdicion = Boolean(producto)
  const [form, setForm] = useState({ nombre: '', categoria: 'flor', id_proveedor: '', cantidad: '', precio: '' })
  const [proveedores, setProveedores] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    apiFetch(`/proveedores`)
      .then(res => res.json())
      .then(setProveedores)
      .catch(() => {})

    if (producto) {
      setForm({
        nombre:       producto.nombre,
        categoria:    producto.categoria,
        id_proveedor: producto.id_proveedor,
        cantidad:     producto.cantidad,
        precio:       producto.precio,
      })
    }
  }, [producto])

  const handleChange = e => {
    setForm(f => ({ ...f, [e.target.name]: e.target.value }))
    setError('')
  }

  const handleSubmit = e => {
    e.preventDefault()
    setLoading(true)
    setError('')

    const payload = {
      ...form,
      id_proveedor: parseInt(form.id_proveedor),
      cantidad:     parseInt(form.cantidad),
      precio:       parseFloat(form.precio),
    }

    const endpoint = esEdicion ? `/productos/${producto.id_producto}` : `/productos`
    const method   = esEdicion ? 'PUT' : 'POST'

    apiFetch(endpoint, { method, headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(payload) })
      .then(res => {
        if (res.ok) {
          onSaved()
        } else {
          return res.text().then(msg => setError(msg))
        }
      })
      .catch(() => setError('Error de conexión'))
      .finally(() => setLoading(false))
  }

  return (
    <div className="modal-overlay" onClick={e => e.target === e.currentTarget && onClose()}>
      <div className="modal">
        <h2>{esEdicion ? 'Editar producto' : 'Nuevo producto'}</h2>

        <form onSubmit={handleSubmit}>
          <div className="field">
            <label>Nombre</label>
            <input name="nombre" value={form.nombre} onChange={handleChange} placeholder="Ej: Rosa roja, Tulipán..." required />
          </div>

          <div className="field">
            <label>Categoría</label>
            <select name="categoria" value={form.categoria} onChange={handleChange}>
              {CATEGORIAS.map(c => (
                <option key={c} value={c}>
                  {c.charAt(0).toUpperCase() + c.slice(1)}
                </option>
              ))}
            </select>
          </div>

          <div className="field">
            <label>Proveedor</label>
            <select name="id_proveedor" value={form.id_proveedor} onChange={handleChange} required>
              <option value="">Seleccionar proveedor</option>
              {proveedores.map(p => (
                <option key={p.id_proveedor} value={p.id_proveedor}>{p.nombre}</option>
              ))}
            </select>
          </div>

          <div className="producto-grid-2">
            <div className="field">
              <label>Cantidad</label>
              <input name="cantidad" type="number" min="0" value={form.cantidad} onChange={handleChange} required />
            </div>
            <div className="field">
              <label>Precio (Q)</label>
              <input name="precio" type="number" min="0" step="0.01" value={form.precio} onChange={handleChange} required />
            </div>
          </div>

          {error && <div className="producto-error">{error}</div>}

          <div className="modal-actions">
            <button type="button" className="btn btn-secondary" onClick={onClose}>Cancelar</button>
            <button type="submit" className="btn btn-primary" disabled={loading}>
              {loading ? 'Guardando…' : esEdicion ? 'Guardar cambios' : 'Crear producto'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}