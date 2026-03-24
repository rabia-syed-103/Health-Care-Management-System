import { useState } from 'react'
import { addMedicine } from '../../api/medicines'
import PageHeader from '../../components/PageHeader'
import Alert from '../../components/Alert'

const empty = { batch_no: '', name: '', stock: '', expiry_date: '' }

export default function PharmacistAddMed() {
  const [form,    setForm]    = useState(empty)
  const [alert,   setAlert]   = useState({ type: '', message: '' })
  const [loading, setLoading] = useState(false)
  const [result,  setResult]  = useState(null)

  const showAlert = (type, message) => {
    setAlert({ type, message })
    setTimeout(() => setAlert({ type: '', message: '' }), 5000)
  }

  const handleSubmit = async () => {
    if (!form.batch_no || !form.name || !form.stock || !form.expiry_date)
      return showAlert('error', 'All fields are required')
    setLoading(true)
    setResult(null)
    try {
      const res = await addMedicine({ ...form, stock: parseInt(form.stock) })
      setResult(res.data)
      showAlert('success', res.data.action === 'created' ? 'Medicine added successfully' : 'Medicine restocked successfully')
      if (res.data.action === 'created') setForm(empty)
    } catch (err) {
      showAlert('error', err.response?.data?.error || 'Failed to save medicine')
    } finally { setLoading(false) }
  }

  return (
    <div>
      <PageHeader title="Add / Restock Medicine" subtitle="Add a new medicine or restock an existing batch" />

      <div className="max-w-lg">
        <Alert type={alert.type} message={alert.message} onClose={() => setAlert({ type: '', message: '' })} />

        <div className="card space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Batch Number</label>
            <input className="input-field" value={form.batch_no} onChange={e => setForm({...form, batch_no: e.target.value})} placeholder="B001" />
            <p className="text-xs text-gray-400 mt-1">If this batch already exists, stock will be added on top</p>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Medicine Name</label>
            <input className="input-field" value={form.name} onChange={e => setForm({...form, name: e.target.value})} placeholder="Paracetamol" />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Stock Quantity</label>
            <input type="number" min="1" className="input-field" value={form.stock} onChange={e => setForm({...form, stock: e.target.value})} placeholder="100" />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Expiry Date</label>
            <input type="date" className="input-field" value={form.expiry_date} onChange={e => setForm({...form, expiry_date: e.target.value})} />
          </div>
          <button onClick={handleSubmit} disabled={loading} className="btn-primary w-full">
            {loading ? 'Saving...' : 'Save Medicine'}
          </button>
        </div>

        {result && (
          <div className="card mt-4 bg-green-50 border-green-200">
            <h3 className="font-semibold text-green-800 mb-2">
              {result.action === 'created' ? '✅ Medicine Added' : '✅ Medicine Restocked'}
            </h3>
            {result.action === 'created' ? (
              <p className="text-sm text-green-700">New medicine ID: <strong>#{result.id}</strong></p>
            ) : (
              <>
                <p className="text-sm text-green-700">Batch: <strong>{result.batch_no}</strong></p>
                <p className="text-sm text-green-700">New total stock: <strong>{result.new_stock}</strong></p>
              </>
            )}
          </div>
        )}
      </div>
    </div>
  )
}
