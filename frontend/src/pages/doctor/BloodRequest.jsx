import { useState } from 'react'
import { createBloodRequest } from '../../api/blood'
import { useAuth } from '../../context/AuthContext'
import PageHeader from '../../components/PageHeader'
import Alert from '../../components/Alert'

export default function DoctorBloodRequest() {
  const { user } = useAuth()

  const [form, setForm] = useState({
    doctor_id:      user?.user_id || '',
    patient_email:  '',
    quantity_needed: 1,
  })
  const [alert,   setAlert]   = useState({ type: '', message: '' })
  const [loading, setLoading] = useState(false)
  const [result,  setResult]  = useState(null)

  const showAlert = (type, message) => {
    setAlert({ type, message })
    setTimeout(() => setAlert({ type: '', message: '' }), 5000)
  }

  const handleSubmit = async () => {
    if (!form.patient_email)  return showAlert('error', 'Patient email is required')
    if (!form.doctor_id)      return showAlert('error', 'Doctor ID is required')
    if (form.quantity_needed < 1) return showAlert('error', 'Quantity must be at least 1')

    setLoading(true)
    setResult(null)
    try {
      const res = await createBloodRequest({
        doctor_id:       parseInt(form.doctor_id),
        patient_email:   form.patient_email,
        quantity_needed: parseInt(form.quantity_needed),
      })
      setResult(res.data)
      showAlert('success', 'Blood request created successfully')
      setForm({ ...form, patient_email: '', quantity_needed: 1 })
    } catch (err) {
      showAlert('error', err.response?.data?.error || 'Failed to create blood request')
    } finally { setLoading(false) }
  }

  return (
    <div>
      <PageHeader title="Create Blood Request" subtitle="Request blood units for a patient" />

      <div className="max-w-lg">
        <Alert type={alert.type} message={alert.message} onClose={() => setAlert({ type: '', message: '' })} />

        <div className="card space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Doctor ID</label>
            <input
              type="number"
              className="input-field"
              value={form.doctor_id}
              onChange={e => setForm({...form, doctor_id: e.target.value})}
              placeholder="Your doctor ID"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Patient Email</label>
            <input
              className="input-field"
              value={form.patient_email}
              onChange={e => setForm({...form, patient_email: e.target.value})}
              placeholder="patient@email.com"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Units Needed</label>
            <input
              type="number"
              min="1"
              className="input-field"
              value={form.quantity_needed}
              onChange={e => setForm({...form, quantity_needed: e.target.value})}
            />
          </div>

          <button onClick={handleSubmit} disabled={loading} className="btn-primary w-full">
            {loading ? 'Submitting...' : 'Submit Blood Request'}
          </button>
        </div>

        {result && (
          <div className="card mt-4 bg-green-50 border-green-200">
            <h3 className="font-semibold text-green-800 mb-2">✅ Blood Request Created</h3>
            <p className="text-sm text-green-700">Request ID: <strong>#{result.request_id}</strong></p>
            <p className="text-sm text-green-700">Patient: <strong>{result.patient}</strong></p>
            <p className="text-sm text-green-700">Blood Group: <strong>{result.patient_blood_gr}</strong></p>
            <p className="text-sm text-green-700">Units Needed: <strong>{result.quantity_needed}</strong></p>
            <p className="text-sm text-green-700">Status: <span className="badge-warning">{result.status}</span></p>
          </div>
        )}
      </div>
    </div>
  )
}
