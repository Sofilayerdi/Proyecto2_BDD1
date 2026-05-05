import { useState, useEffect, useCallback } from 'react'
import ProductoModal from '../components/ProductoModal'
import CompraModal from '../components/CompraModal'
import { useToast } from '../components/useToast'
import './Productos.css'

const URL = 'http://localhost:8000'

const CATEGORIAS = ['todos', 'flor', 'follaje', 'liston', 'papel']

const CAT_COLORS = {
  flor:    { bg: '#fdf0f4', text: '#8b3a52', border: '#e8b0c0' },
  follaje: { bg: '#f0f5f0', text: '#3d5a3e', border: '#a8c8a8' },
  liston:  { bg: '#f5f3e8', text: '#6b5a2e', border: '#d4c080' },
  papel:   { bg: '#f0f2f5', text: '#3a4a6b', border: '#a0b0d0' },
}

export default function Productos() {
  const [productos, setProductos] = useState([])
  const [filtro, setFiltro] = useState('todos')
  const [loading, setLoading] = useState(true)
  const [productoEditar, setProductoEditar] = useState(null)
  const [showModal, setShowModal] = useState(false)
  const [carrito, setCarrito] = useState([])
  const [showCarrito, setShowCarrito] = useState(false)
  const [showCompra, setShowCompra] = useState(false)
  const { showToast, Toast } = useToast()

  const cargarProductos = useCallback(() => {
    setLoading(true)
    const query = filtro !== 'todos' ? `?categoria=${filtro}&limit=100` : '?limit=100'
    fetch(`${URL}/productos${query}`)
      .then(res => res.json())
      .then(data => setProductos(data))
      .catch(() => showToast('Error cargando productos', 'error'))
      .finally(() => setLoading(false))
  }, [filtro])

  useEffect(() => { cargarProductos() }, [cargarProductos])

  const eliminar = (id) => {
    if (!confirm('¿Eliminar este producto?')) return
    fetch(`${URL}/productos/${id}`, { method: 'DELETE' })
      .then(res => {
        if (res.ok) {
          showToast('Producto eliminado')
          cargarProductos()
        } else {
          return res.text().then(msg => showToast(msg, 'error'))
        }
      })
      .catch(() => showToast('Error al eliminar', 'error'))
  }

  const agregarAlCarrito = (producto) => {
    setCarrito(c => {
      const existe = c.find(i => i.id_producto === producto.id_producto)
      if (existe) {
        if (existe.cantidad >= producto.cantidad) {
          showToast('Stock insuficiente', 'error')
          return c
        }
        return c.map(i => i.id_producto === producto.id_producto ? { ...i, cantidad: i.cantidad + 1 } : i)
      }
      if (producto.cantidad === 0) { showToast('Sin stock disponible', 'error'); return c }
      return [...c, { ...producto, cantidad: 1 }]
    })
    showToast(`${producto.nombre} agregado al ramo`)
  }

  const quitarDelCarrito = (id) => setCarrito(c => c.filter(i => i.id_producto !== id))

  const cambiarCantidad = (id, delta) => {
    setCarrito(c => c.map(i => {
      if (i.id_producto !== id) return i
      const nueva = i.cantidad + delta
      return nueva <= 0 ? null : { ...i, cantidad: nueva }
    }).filter(Boolean))
  }

  const totalCarrito = carrito.reduce((s, i) => s + i.precio * i.cantidad, 0)

  return (
    <div className="productos-page">
      <Toast />

      <div className="productos-header">
        <div>
          <h1>Inventario</h1>
          <p className="subtitle">Gestión de productos y ventas</p>
        </div>
        <button className="btn btn-primary" onClick={() => { setProductoEditar(null); setShowModal(true) }}>
          + Nuevo producto
        </button>
      </div>

      <div className="filtros">
        {CATEGORIAS.map(cat => (
          <button key={cat} className={`filtro-btn ${filtro === cat ? 'active' : ''}`} onClick={() => setFiltro(cat)}>
            {cat === 'todos' ? 'Todos' : cat.charAt(0).toUpperCase() + cat.slice(1)}
          </button>
        ))}
      </div>

      {loading ? (
        <div className="spinner">Cargando…</div>
      ) : productos.length === 0 ? (
        <div className="empty">No hay productos en esta categoría</div>
      ) : (
        <div className="productos-grid">
          {productos.map(p => {
            const col = CAT_COLORS[p.categoria] || CAT_COLORS.flor
            const enCarrito = carrito.find(i => i.id_producto === p.id_producto)
            return (
              <div key={p.id_producto} className="producto-card">
                <div className="card-top">
                  <span className="cat-badge" style={{ background: col.bg, color: col.text, border: `1px solid ${col.border}` }}>
                    {p.categoria}
                  </span>
                  {p.cantidad <= 20 && (
                    <span className="stock-warn">{p.cantidad === 0 ? 'Sin stock' : `⚠ ${p.cantidad} uds`}</span>
                  )}
                </div>
                <div className="card-body">
                  <h3 className="prod-nombre">{p.nombre}</h3>
                  <p className="prod-proveedor">{p.proveedor}</p>
                  <div className="prod-info">
                    <span className="prod-precio">Q{Number(p.precio).toFixed(2)}</span>
                    <span className="prod-stock">{p.cantidad} en stock</span>
                  </div>
                </div>
                <div className="card-actions">
                  <button
                    className={`btn-carrito ${enCarrito ? 'en-carrito' : ''}`}
                    onClick={() => agregarAlCarrito(p)}
                    disabled={p.cantidad === 0}
                  >
                    {enCarrito ? `En ramo (${enCarrito.cantidad})` : '+ Agregar al ramo'}
                  </button>
                  <div className="card-edit-actions">
                    <button className="btn btn-secondary btn-sm" onClick={() => { setProductoEditar(p); setShowModal(true) }}>Editar</button>
                    <button className="btn btn-danger btn-sm" onClick={() => eliminar(p.id_producto)}>Eliminar</button>
                  </div>
                </div>
              </div>
            )
          })}
        </div>
      )}

      {carrito.length > 0 && (
        <div className="carrito-fab">
          <button className="carrito-toggle" onClick={() => setShowCarrito(s => !s)}>
            <span className="carrito-icon">🌿</span>
            <span>{carrito.length} producto{carrito.length > 1 ? 's' : ''}</span>
            <span className="carrito-total">Q{totalCarrito.toFixed(2)}</span>
          </button>

          {showCarrito && (
            <div className="carrito-panel">
              <div className="carrito-header">
                <span>Ramo actual</span>
                <button className="close-btn" onClick={() => setShowCarrito(false)}>✕</button>
              </div>
              <div className="carrito-items">
                {carrito.map(item => (
                  <div key={item.id_producto} className="carrito-item">
                    <div className="ci-info">
                      <span className="ci-nombre">{item.nombre}</span>
                      <span className="ci-precio">Q{(item.precio * item.cantidad).toFixed(2)}</span>
                    </div>
                    <div className="ci-controls">
                      <button onClick={() => cambiarCantidad(item.id_producto, -1)}>−</button>
                      <span>{item.cantidad}</span>
                      <button onClick={() => cambiarCantidad(item.id_producto, +1)}>+</button>
                      <button className="ci-remove" onClick={() => quitarDelCarrito(item.id_producto)}>✕</button>
                    </div>
                  </div>
                ))}
              </div>
              <div className="carrito-footer">
                <div className="carrito-subtotal">
                  <span>Total</span>
                  <span>Q{totalCarrito.toFixed(2)}</span>
                </div>
                <button className="btn btn-gold" style={{ width: '100%' }}
                  onClick={() => { setShowCarrito(false); setShowCompra(true) }}>
                  Comprar
                </button>
              </div>
            </div>
          )}
        </div>
      )}

      {showModal && (
        <ProductoModal
          producto={productoEditar}
          onClose={() => setShowModal(false)}
          onSaved={() => {
            setShowModal(false)
            cargarProductos()
            showToast(productoEditar ? 'Producto actualizado' : 'Producto creado')
          }}
        />
      )}

      {showCompra && (
        <CompraModal
          carrito={carrito}
          onClose={() => setShowCompra(false)}
          onComprado={() => {
            setShowCompra(false)
            setCarrito([])
            cargarProductos()
            showToast('¡Venta realizada con éxito!', 'success')
          }}
        />
      )}
    </div>
  )
}