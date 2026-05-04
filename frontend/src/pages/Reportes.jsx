import { useState, useEffect } from 'react'
import { api } from '../api'
import { useToast } from '../components/useToast'
import './Reportes.css'

const TABS = [
  { id: 'empleados', label: 'Ventas por empleado' },
  { id: 'productos', label: 'Productos más vendidos' },
  { id: 'inventario', label: 'Stock bajo' },
  { id: 'ventas', label: 'Historial de ventas' },
]

export default function Reportes() {
  const [tab, setTab] = useState('empleados')
  const [data, setData] = useState([])
  const [loading, setLoading] = useState(false)
  const { showToast, Toast } = useToast()

  useEffect(() => {
    setData([])
    setLoading(true)
    const loaders = {
      empleados: api.getVentasPorEmpleado,
      productos: api.getProductosEnRamos,
      inventario: api.getInventarioBajo,
      ventas: api.getVistaVentas,
    }
    loaders[tab]()
      .then(setData)
      .catch(err => showToast(err.message, 'error'))
      .finally(() => setLoading(false))
  }, [tab])

  return (
    <div className="reportes-page">
      <Toast />

      <div className="reportes-header">
        <h1>Reportes</h1>
        <p className="subtitle">Análisis de ventas e inventario</p>
      </div>

      <div className="reporte-tabs">
        {TABS.map(t => (
          <button
            key={t.id}
            className={`tab-btn ${tab === t.id ? 'active' : ''}`}
            onClick={() => setTab(t.id)}
          >
            {t.label}
          </button>
        ))}
      </div>

      <div className="reporte-content">
        {/* Descripción de la query para cada reporte */}
        <div className="query-badge">
          {tab === 'empleados' && <span>Vista SQL · GROUP BY + HAVING + SUM · JOIN empleado ↔ venta</span>}
          {tab === 'productos' && <span>CTE (WITH) · JOIN ramo_producto → venta → producto → proveedor</span>}
          {tab === 'inventario' && <span>Subquery correlacionado EXISTS · Productos bajo el promedio de stock</span>}
          {tab === 'ventas' && <span>Vista SQL · JOIN venta ↔ cliente ↔ empleado · ORDER BY fecha</span>}
        </div>

        {loading ? (
          <div className="spinner">Cargando reporte…</div>
        ) : data.length === 0 ? (
          <div className="empty-reporte">Sin datos disponibles</div>
        ) : (
          <>
            {tab === 'empleados' && <TablaEmpleados data={data} />}
            {tab === 'productos' && <TablaProductos data={data} />}
            {tab === 'inventario' && <TablaInventario data={data} />}
            {tab === 'ventas' && <TablaVentas data={data} />}
          </>
        )}
      </div>
    </div>
  )
}

function TablaEmpleados({ data }) {
  const maxIngresos = Math.max(...data.map(d => d.ingresos))
  return (
    <div>
      <div className="reporte-cards">
        {data.map((e, i) => (
          <div key={i} className="emp-card">
            <div className="emp-avatar">{e.empleado.charAt(0)}</div>
            <div className="emp-info">
              <p className="emp-nombre">{e.empleado}</p>
              <p className="emp-rol">{e.rol}</p>
            </div>
            <div className="emp-stats">
              <div className="stat">
                <span className="stat-val">{e.total_ventas}</span>
                <span className="stat-lbl">ventas</span>
              </div>
              <div className="stat">
                <span className="stat-val" style={{ color: 'var(--moss)' }}>Q{Number(e.ingresos).toFixed(0)}</span>
                <span className="stat-lbl">ingresos</span>
              </div>
            </div>
            <div className="emp-bar-wrap">
              <div className="emp-bar" style={{ width: `${(e.ingresos / maxIngresos) * 100}%` }} />
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}

function TablaProductos({ data }) {
  const max = Math.max(...data.map(d => d.total_vendido))
  return (
    <table className="reporte-table">
      <thead>
        <tr>
          <th>#</th>
          <th>Producto</th>
          <th>Categoría</th>
          <th>Proveedor</th>
          <th>Unidades vendidas</th>
        </tr>
      </thead>
      <tbody>
        {data.map((p, i) => (
          <tr key={i}>
            <td className="rank">#{i + 1}</td>
            <td><strong>{p.producto}</strong></td>
            <td><span className="cat-chip">{p.categoria}</span></td>
            <td className="muted">{p.proveedor}</td>
            <td>
              <div className="bar-cell">
                <div className="mini-bar" style={{ width: `${(p.total_vendido / max) * 100}%` }} />
                <span>{p.total_vendido}</span>
              </div>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}

function TablaInventario({ data }) {
  return (
    <div>
      <p className="reporte-note">
        Productos con stock por debajo del promedio general — requieren reabastecimiento.
      </p>
      <table className="reporte-table">
        <thead>
          <tr>
            <th>Producto</th>
            <th>Categoría</th>
            <th>Proveedor</th>
            <th>Stock actual</th>
            <th>Precio</th>
          </tr>
        </thead>
        <tbody>
          {data.map((p, i) => (
            <tr key={i} className={p.cantidad === 0 ? 'row-danger' : p.cantidad <= 10 ? 'row-warn' : ''}>
              <td><strong>{p.nombre}</strong></td>
              <td><span className="cat-chip">{p.categoria}</span></td>
              <td className="muted">{p.proveedor}</td>
              <td>
                <span className={`stock-badge ${p.cantidad === 0 ? 'out' : p.cantidad <= 10 ? 'low' : ''}`}>
                  {p.cantidad === 0 ? 'Sin stock' : p.cantidad}
                </span>
              </td>
              <td>Q{Number(p.precio).toFixed(2)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function TablaVentas({ data }) {
  return (
    <table className="reporte-table">
      <thead>
        <tr>
          <th>#</th>
          <th>Fecha</th>
          <th>Cliente</th>
          <th>Empleado</th>
          <th>Total</th>
        </tr>
      </thead>
      <tbody>
        {data.map((v, i) => (
          <tr key={i}>
            <td className="muted">{v.id_venta}</td>
            <td>{v.fecha?.split('T')[0] || v.fecha}</td>
            <td><strong>{v.cliente}</strong></td>
            <td className="muted">{v.empleado}</td>
            <td style={{ fontFamily: 'var(--font-display)', fontSize: 16, color: 'var(--moss-dark)' }}>
              Q{Number(v.precio_total).toFixed(2)}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}
