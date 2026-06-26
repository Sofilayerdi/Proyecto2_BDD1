import { useState, useEffect } from 'react'
import { apiFetch } from '../utils/api'
import { useToast } from '../components/useToast'
import './Reportes.css'

const TABS = [
  { id: 'mensuales', label: 'Ventas mensuales' },
  { id: 'productos', label: 'Top productos' },
]

const ENDPOINTS = {
  mensuales: '/reportes/ventas-mensuales',
  productos: '/reportes/top-productos',
}

export default function Reportes() {
  const [tab, setTab] = useState('mensuales')
  const [data, setData] = useState([])
  const [loading, setLoading] = useState(false)
  const { showToast, Toast } = useToast()

  useEffect(() => {
    setData([])
    setLoading(true)
    apiFetch(`${ENDPOINTS[tab]}`)
      .then(res => res.json())
      .then(setData)
      .catch(() => showToast('Error cargando reporte', 'error'))
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
          <button key={t.id} className={`tab-btn ${tab === t.id ? 'active' : ''}`} onClick={() => setTab(t.id)}>
            {t.label}
          </button>
        ))}
      </div>

      <div className="reporte-content">
        {loading ? (
          <div className="spinner">Cargando reporte</div>
        ) : data.length === 0 ? (
          <div className="empty-reporte">Sin datos disponibles</div>
        ) : (
          <>
            {tab === 'mensuales' && <TablaMensuales data={data} />}
            {tab === 'productos' && <TablaProductos data={data} />}
          </>
        )}
      </div>
    </div>
  )
}

function TablaMensuales({ data }) {
  const maxIngresos = Math.max(...data.map(d => d.ingresos))
  const maxVentas   = Math.max(...data.map(d => d.total_ventas))
  const total       = data.reduce((s, d) => s + d.ingresos, 0)
  const totalVentas = data.reduce((s, d) => s + d.total_ventas, 0)

  return (
    <div>
      <div className="summary-cards">
        <div className="summary-card rose">
          <p className="summary-card-label">Total ingresos</p>
          <p className="summary-card-value">Q{Number(total).toFixed(0)}</p>
        </div>
        <div className="summary-card lavender">
          <p className="summary-card-label">Total ventas</p>
          <p className="summary-card-value">{totalVentas}</p>
        </div>
        <div className="summary-card mint">
          <p className="summary-card-label">Meses activos</p>
          <p className="summary-card-value">{data.length}</p>
        </div>
      </div>

      <div className="bar-chart-wrap">
        <p className="bar-chart-title">Ingresos por mes</p>
        <div className="bar-chart">
          {data.map((d, i) => (
            <div key={i} className="bar-col">
              <span className="bar-value">Q{Number(d.ingresos).toFixed(0)}</span>
              <div
                className="bar-fill"
                style={{ height: `${(d.ingresos / maxIngresos) * 140}px` }}
                title={`${d.mes}: Q${Number(d.ingresos).toFixed(2)} · ${d.total_ventas} ventas`}
              />
              <span className="bar-label">{d.mes?.slice(5)}/{d.mes?.slice(2, 4)}</span>
            </div>
          ))}
        </div>
      </div>

      <table className="reporte-table">
        <thead>
          <tr>
            <th>Mes</th>
            <th>Ventas</th>
            <th>Ingresos</th>
            <th>Promedio / venta</th>
          </tr>
        </thead>
        <tbody>
          {[...data].reverse().map((d, i) => (
            <tr key={i}>
              <td><span className="mes-badge">{d.mes}</span></td>
              <td>
                <div className="ventas-bar-cell">
                  <div
                    className="ventas-mini-bar"
                    style={{ width: `${(d.total_ventas / maxVentas) * 80}px` }}
                  />
                  <span>{d.total_ventas}</span>
                </div>
              </td>
              <td className="td-ingresos">Q{Number(d.ingresos).toFixed(2)}</td>
              <td className="muted">
                Q{d.total_ventas > 0 ? (d.ingresos / d.total_ventas).toFixed(2) : '0.00'}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

const CAT_COLORS = {
  flor:    { bg: '#fdf0f4', text: '#c4687e', border: '#f2c4cf' },
  follaje: { bg: '#f0f5f2', text: '#4a8060', border: '#a0c8b0' },
  liston:  { bg: '#fdf5ee', text: '#a06030', border: '#e8c890' },
  papel:   { bg: '#f0f0f8', text: '#5060a0', border: '#b0b8e0' },
}

function rankClass(i) {
  if (i === 0) return 'rank-number gold'
  if (i === 1) return 'rank-number silver'
  return 'rank-number rest'
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
        {data.map((p, i) => {
          const col = CAT_COLORS[p.categoria] || CAT_COLORS.flor
          return (
            <tr key={i}>
              <td>
                <span className={rankClass(i)}>#{i + 1}</span>
              </td>
              <td>
                <strong className="prod-nombre-tabla">{p.producto}</strong>
              </td>
              <td>
                <span
                  className="cat-chip"
                  style={{ background: col.bg, color: col.text, border: `1px solid ${col.border}` }}
                >
                  {p.categoria}
                </span>
              </td>
              <td className="muted">{p.proveedor}</td>
              <td>
                <div className="bar-cell">
                  <div className="mini-bar" style={{ width: `${(p.total_vendido / max) * 140}px` }} />
                  <span className="unidades-count">{p.total_vendido}</span>
                </div>
              </td>
            </tr>
          )
        })}
      </tbody>
    </table>
  )
}