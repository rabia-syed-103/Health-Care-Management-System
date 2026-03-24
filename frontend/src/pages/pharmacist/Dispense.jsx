import { useState } from 'react'
import { dispenseMedicines } from '../../api/prescriptions'
import { useAuth } from '../../context/AuthContext'
import PageHeader from '../../components/PageHeader'
import Alert from '../../components/Alert'

export default function PharmacistDispense() {
  const { user } = useAuth()
  const [form,    setForm]    = useState({ prescription_id: '', pharmacist_id: user?.user_id || '' })
  const [alert,   setAlert]   = useState({ type: '', message: '' })
  const [loading, setLoading] = useState(false)
  const [result,  setResult]  = useState(null)

  const showAlert = (type, message) => {
    setAlert({ type, message })
    setTimeout(() => setAlert({ type: '', message: '' }), 5000)
  }

  const handleSubmit = async () => {
    if (!form.prescription_id) return showAlert('error', 'Prescription ID is required')
    if (!form.pharmacist_id)   return showAlert('error', 'Pharmacist ID is required')
    setLoading(true)
    setResult(null)
    try {
      const res = await dispenseMedicines({
        prescription_id: parseInt(form.prescription_id),
        pharmacist_id:   parseInt(form.pharmacist_id),
      })
      setResult(res.data)
      showAlert('success', 'Medicines dispensed successfully')
      setForm({ ...form, prescription_id: '' })
    } catch (err) {
      showAlert('error', err.response?.data?.error || 'Dispensing failed — transaction rolled back')
    } finally { setLoading(false) }
  }

  return (
    <div>
      <PageHeader title="Dispense Medicines" subtitle="Dispense medicines from a prescription" />

      <div className="max-w-lg">
        <Alert type={alert.type} message={alert.message} onClose={() => setAlert({ type: '', message: '' })} />

        <div className="card space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Prescription ID</label>
            <input
              type="number"
              className="input-field"
              value={form.prescription_id}
              onChange={e => setForm({...form, prescription_id: e.target.value})}
              placeholder="Enter prescription ID from pending list"
            />
            <p className="text-xs text-gray-400 mt-1">Check Pending Prescriptions page for the ID</p>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Pharmacist ID</label>
            <input
              type="number"
              className="input-field"
              value={form.pharmacist_id}
              onChange={e => setForm({...form, pharmacist_id: e.target.value})}
            />
          </div>
          <button onClick={handleSubmit} disabled={loading} className="btn-primary w-full">
            {loading ? 'Dispensing...' : 'Dispense Medicines'}
          </button>
        </div>

        {result && (
          <div className="card mt-4 bg-green-50 border-green-200">
            <h3 className="font-semibold text-green-800 mb-2">✅ Medicines Dispensed</h3>
            <p className="text-sm text-green-700">Prescription ID: <strong>#{result.prescription_id}</strong></p>
            <p className="text-sm text-green-700">Pharmacist: <strong>{result.pharmacist}</strong></p>
            <p className="text-sm text-green-700">Medicines Dispensed: <strong>{result.medicines_count}</strong></p>
          </div>
        )}
      </div>
    </div>
  )
}
