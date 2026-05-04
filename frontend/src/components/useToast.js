import { useState, useCallback } from 'react'

export function useToast() {
  const [toast, setToast] = useState({ msg: '', type: '', show: false })

  const showToast = useCallback((msg, type = 'success') => {
    setToast({ msg, type, show: true })
    setTimeout(() => setToast(t => ({ ...t, show: false })), 3000)
  }, [])

  const Toast = () => (
    <div className={`toast ${toast.type} ${toast.show ? 'show' : ''}`}>
      {toast.msg}
    </div>
  )

  return { showToast, Toast }
}
