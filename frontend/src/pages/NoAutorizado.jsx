export default function NoAutorizado() {
  return (
    <div style={{
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      justifyContent: 'center',
      minHeight: '60vh',
      gap: '12px'
    }}>
      <h2 style={{ fontSize: '32px', color: 'var(--rose-dark)', fontStyle: 'italic' }}>
        Acceso denegado
      </h2>
      <p style={{ color: 'var(--plum-light)' }}>
        No tienes permisos para ver esta página
      </p>
    </div>
  )
}