export default function Alert({ type = 'info', message, onClose }) {
  if (!message) return null

  const styles = {
    success: 'bg-green-50 text-green-800 border-green-200',
    error:   'bg-red-50 text-red-800 border-red-200',
    warning: 'bg-yellow-50 text-yellow-800 border-yellow-200',
    info:    'bg-blue-50 text-blue-800 border-blue-200',
  }

  return (
    <div className={`border rounded-lg px-4 py-3 flex items-start justify-between gap-3 mb-4 ${styles[type]}`}>
      <p className="text-sm">{message}</p>
      {onClose && (
        <button onClick={onClose} className="text-lg leading-none opacity-60 hover:opacity-100">
          ×
        </button>
      )}
    </div>
  )
}
